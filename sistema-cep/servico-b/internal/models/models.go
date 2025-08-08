package models

import (
	"encoding/json"
	"strconv"
)

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

// ViaCEPResponse representa a resposta da API ViaCEP
type ViaCEPResponse struct {
	CEP         string     `json:"cep"`
	Logradouro  string     `json:"logradouro"`
	Complemento string     `json:"complemento"`
	Bairro      string     `json:"bairro"`
	Localidade  string     `json:"localidade"`
	UF          string     `json:"uf"`
	IBGE        string     `json:"ibge"`
	GIA         string     `json:"gia"`
	DDD         string     `json:"ddd"`
	SIAFI       string     `json:"siafi"`
	Erro        BoolString `json:"erro,omitempty"`
}

type BoolString bool

// UnmarshalJSON implementa o unmarshaler customizado
func (b *BoolString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Se não conseguir fazer unmarshal como string, tenta como bool
		var boolVal bool
		if err := json.Unmarshal(data, &boolVal); err != nil {
			return err
		}
		*b = BoolString(boolVal)
		return nil
	}

	// Converte string para bool
	switch s {
	case "true", "1", "True", "TRUE":
		*b = true
	case "false", "0", "False", "FALSE", "":
		*b = false
	default:
		// Para qualquer outro valor, tenta converter
		if boolVal, err := strconv.ParseBool(s); err == nil {
			*b = BoolString(boolVal)
		} else {
			*b = false
		}
	}
	return nil
}

// Bool converte BoolString para bool nativo
func (b BoolString) Bool() bool {
	return bool(b)
}

// WeatherAPIResponse representa a resposta da API WeatherAPI
type WeatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

// LocationInfo representa informações de localização processadas
type LocationInfo struct {
	City  string
	State string
}

// TemperatureInfo representa informações de temperatura processadas
type TemperatureInfo struct {
	City  string
	TempC float64
	TempF float64
	TempK float64
}
