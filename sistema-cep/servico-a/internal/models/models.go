package models

type CEPRequest struct {
	CEP string
}

type ErrorResponse struct {
	Message string
}

type TemperatureResponse struct {
	City  string
	TempC float64
	TempF float64
	TempK float64
}

type ServiceBResponse struct {
	StatusCode int
	Body       []byte
}
