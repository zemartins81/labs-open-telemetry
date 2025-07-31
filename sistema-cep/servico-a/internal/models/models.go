package models

// CEPRequest representa a estrutura da requisição de CEP
type CEPRequest struct {
	CEP string `json:"cep"`
}

// ErrorResponse representa a estrutura de resposta de erro
type ErrorResponse struct {
	Message string `json:"message"`
}

// TemperatureResponse representa a resposta com dados de temperatura
type TemperatureResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ServiceBResponse representa a resposta do Serviço B
type ServiceBResponse struct {
	StatusCode int
	Body       []byte
}
