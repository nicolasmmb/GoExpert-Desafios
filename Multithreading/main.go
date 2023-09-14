package main

import (
	"encoding/json"
	"fmt"
	"io"
	md "local/models"
	"log"
	"net/http"
	"os"
	"time"
)

func fetchFromAPI[C md.RequestType](url string, ch chan md.RequestDTO) {
	var data C
	var err error

	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(url)

	if err != nil {
		log.Default().Println("Erro ao consultar API: ", url, " - Erro:", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Default().Println("Erro ao consultar API: ", url, " - Status:", resp.Status)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Default().Println("Erro ao ler resposta da API: ", url, " - Erro:", err)
		return
	}
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Default().Println("Erro ao converter resposta da API: ", url, " - Erro:", err)
		return
	}

	ch <- md.NewRequestDTO(url, data.GetCode(), data.GetCity(), data.GetState())

}

func main() {
	// read from command line arguments or use default

	var cep string
	if len(os.Args) < 2 {
		cep = "12460-000"
	} else {
		cep = os.Args[1]
	}

	ch := make(chan md.RequestDTO, 2)

	urlViaCEP := fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep)
	urlAPICEP := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)

	go fetchFromAPI[md.RequestViaCEP](urlViaCEP, ch)
	go fetchFromAPI[md.RequestAPICEP](urlAPICEP, ch)

	select {
	case result := <-ch:
		println(result.String())

	case <-time.After(1 * time.Second):
		fmt.Println("Tempo limite excedido para ambas as APIs.")
	}
}
