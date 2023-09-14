package models

type RequestViaCEP struct {
	Code  string `json:"cep"`
	City  string `json:"localidade"`
	State string `json:"uf"`
}

func (r RequestViaCEP) GetCode() string {
	return r.Code
}

func (r RequestViaCEP) GetCity() string {
	return r.City
}

func (r RequestViaCEP) GetState() string {
	return r.State
}
