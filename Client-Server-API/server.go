package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		log.Fatalf("Error while opening the database: %v", err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid FLOAT, timestamp DATETIME)`)
	if err != nil {
		log.Fatalf("Error while creating table: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Make a request to the currency API
		ctxAPI, cancelAPI := context.WithTimeout(ctx, 200*time.Millisecond)
		defer cancelAPI()

		req, err := http.NewRequestWithContext(ctxAPI, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {
			log.Printf("Error while creating request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error while sending request to API: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error while reading response body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data map[string]map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("Error while unmarshalling JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if _, ok := data["USDBRL"]; !ok {
			log.Println("Error: 'USDBRL' not found in the response.")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		bid, ok := data["USDBRL"]["bid"]
		if !ok {
			log.Println("Error: 'bid' not found in the response.")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Save the bid value to the database
		ctxDB, cancelDB := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancelDB()

		stmt, err := db.PrepareContext(ctxDB, "INSERT INTO cotacoes (bid, timestamp) VALUES (?, ?)")
		if err != nil {
			log.Printf("Error while preparing the statement: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(bid, time.Now())
		if err != nil {
			log.Printf("Error while executing the statement: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{"bid": bid}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error while marshalling JSON response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
