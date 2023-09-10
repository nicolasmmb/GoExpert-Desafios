package main

import (
	"fmt"
	"net/http"
	"time"
)

func fetchFromAPI(apiURL string, cep string, ch chan<- string) {
	start := time.Now()
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(apiURL + cep + ".json")
	if err != nil {
		ch <- fmt.Sprintf("Erro na API %s: %v", apiURL, err)
		return
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)
	ch <- fmt.Sprintf("API %s: Tempo de resposta: %v", apiURL, elapsed)
}

func main() {
	cep := "12460000" // Substitua pelo CEP desejado
	ch := make(chan string, 2)

	go fetchFromAPI("https://cdn.apicep.com/file/apicep/", cep, ch)
	go fetchFromAPI("http://viacep.com.br/ws/", cep, ch)

	select {
	case result := <-ch:
		fmt.Println(result)
		<-ch // Descartar a outra resposta
	case <-time.After(1 * time.Second):
		fmt.Println("Tempo limite excedido para ambas as APIs.")
	}
}
