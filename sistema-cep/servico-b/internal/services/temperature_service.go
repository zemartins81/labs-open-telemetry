package services

import (
	"fmt"

	"servico-b/internal/models"
)

// TemperatureService orquestra a busca de localização e temperatura
type TemperatureService struct {
	viaCEPService  *ViaCEPService
	weatherService *WeatherService
}

// NewTemperatureService cria uma nova instância do serviço de temperatura
func NewTemperatureService(viaCEPService *ViaCEPService, weatherService *WeatherService) *TemperatureService {
	return &TemperatureService{
		viaCEPService:  viaCEPService,
		weatherService: weatherService,
	}
}

// GetTemperatureByCEP busca a temperatura a partir de um CEP
func (t *TemperatureService) GetTemperatureByCEP(cep string) (*models.TemperatureInfo, error) {
	// 1. Busca informações de localização pelo CEP
	location, err := t.viaCEPService.GetLocationByCEP(cep)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar localização: %w", err)
	}

	// 2. Busca informações de temperatura pela localização
	temperature, err := t.weatherService.GetTemperatureByLocation(location)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar temperatura: %w", err)
	}

	return temperature, nil
}
