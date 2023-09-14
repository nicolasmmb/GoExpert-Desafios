package models

import "fmt"

type RequestType interface {
	RequestAPICEP | RequestViaCEP
	GetCode() string
	GetCity() string
	GetState() string
}

type RequestDTO struct {
	URL   string
	Code  string
	City  string
	State string
}

func (r *RequestDTO) String() string {
	return fmt.Sprintf("URL: %s\nCEP: %s - Estado: %s - Cidade: %s", r.URL, r.Code, r.State, r.City)
}

func NewRequestDTO(url string, code string, city string, state string) RequestDTO {
	return RequestDTO{url, code, city, state}
}
