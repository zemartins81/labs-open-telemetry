package services

import (
	"context"
	"fmt"

	"servico-b/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
func (t *TemperatureService) GetTemperatureByCEP(ctx context.Context, cep string) (*models.TemperatureInfo, error) {
	tracer := otel.Tracer("servico-b")
	ctx, span := tracer.Start(ctx, "GetTemperatureByCEP")
	defer span.End()

	span.SetAttributes(attribute.String("cep", cep))

	// 1. Busca informações de localização pelo CEP
	location, err := t.viaCEPService.GetLocationByCEP(ctx, cep)
	if err != nil {
		span.SetStatus(codes.Error, "failed to get location")
		return nil, fmt.Errorf("erro ao buscar localização: %w", err)
	}

	span.SetAttributes(attribute.String("location.city", location.City))

	// 2. Busca informações de temperatura pela localização
	temperature, err := t.weatherService.GetTemperatureByLocation(ctx, location)
	if err != nil {
		span.SetStatus(codes.Error, "failed to get temperature")
		return nil, fmt.Errorf("erro ao buscar temperatura: %w", err)
	}

	span.SetAttributes(
		attribute.String("result.city", temperature.City),
		attribute.Float64("result.temp_c", temperature.TempC),
	)

	return temperature, nil
}
