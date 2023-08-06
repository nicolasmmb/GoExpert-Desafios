package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	resp, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		fmt.Printf("Error while sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error while reading response body: %v\n", err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Error while unmarshalling JSON: %v\n", err)
		return
	}

	if bid, ok := data["bid"]; ok {
		dollarValue := fmt.Sprintf("%v", bid)
		err = saveToFile("cotacao.txt", "Dólar: "+dollarValue)
		if err != nil {
			fmt.Printf("Error while saving to file: %v\n", err)
			return
		}
		fmt.Println("Dólar:", dollarValue)
	} else {
		fmt.Println("Error: bid not found in the response.")
	}
}

func saveToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
