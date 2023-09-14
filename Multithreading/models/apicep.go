package models

type RequestAPICEP struct {
	Code  string `json:"code"`
	City  string `json:"city"`
	State string `json:"state"`
}

func (r RequestAPICEP) GetCode() string {
	return r.Code
}

func (r RequestAPICEP) GetCity() string {
	return r.City
}

func (r RequestAPICEP) GetState() string {
	return r.State
}
