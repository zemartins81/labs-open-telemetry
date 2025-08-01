package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"servico-b/internal/models"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// WeatherService é responsável pela comunicação com a API WeatherAPI
type WeatherService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewWeatherService cria uma nova instância do serviço Weather
func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{
		apiKey:  apiKey,
		baseURL: "http://api.weatherapi.com/v1",
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// GetTemperatureByLocation busca informações de temperatura pela localização
func (w *WeatherService) GetTemperatureByLocation(ctx context.Context, location *models.LocationInfo) (*models.TemperatureInfo, error) {
	tracer := otel.Tracer("servico-b")
	ctx, span := tracer.Start(ctx, "WeatherAPI.GetTemperatureByLocation")
	defer span.End()

	// Constrói a query de localização - usa apenas a cidade para evitar problemas de encoding
	query := location.City

	span.SetAttributes(
		attribute.String("weather.query", query),
		attribute.String("weather.city", location.City),
		attribute.String("weather.state", location.State),
	)

	// Constrói a URL com parâmetros
	apiURL := fmt.Sprintf("%s/current.json", w.baseURL)
	params := url.Values{}
	params.Add("key", w.apiKey)
	params.Add("q", query) // url.Values.Add já faz o encoding automático
	params.Add("aqi", "no") // Não precisamos de dados de qualidade do ar

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())
	span.SetAttributes(attribute.String("weather.url", apiURL))

	log.Printf("Buscando temperatura para %s na WeatherAPI", query)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		span.SetStatus(codes.Error, "failed to create request")
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "failed to make request")
		return nil, fmt.Errorf("erro ao fazer requisição para WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("weather.status_code", resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(codes.Error, "failed to read response")
		return nil, fmt.Errorf("erro ao ler resposta da WeatherAPI: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("WeatherAPI retornou erro %d: %s", resp.StatusCode, string(body))

		// Tenta fazer parse da mensagem de erro
		var errorResp map[string]interface{}
		if json.Unmarshal(body, &errorResp) == nil {
			if errorMsg, ok := errorResp["error"]; ok {
				span.SetStatus(codes.Error, fmt.Sprintf("WeatherAPI error: %v", errorMsg))
				return nil, fmt.Errorf("erro da WeatherAPI: %v", errorMsg)
			}
		}

		span.SetStatus(codes.Error, "non-200 status code")
		return nil, fmt.Errorf("WeatherAPI retornou status %d", resp.StatusCode)
	}

	var weatherResp models.WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		span.SetStatus(codes.Error, "failed to parse response")
		return nil, fmt.Errorf("erro ao fazer parse da resposta WeatherAPI: %w", err)
	}

	// Calcula temperatura em Kelvin (K = C + 273.15)
	tempK := weatherResp.Current.TempC + 273.15

	tempInfo := &models.TemperatureInfo{
		City:  weatherResp.Location.Name,
		TempC: weatherResp.Current.TempC,
		TempF: weatherResp.Current.TempF,
		TempK: tempK,
	}

	span.SetAttributes(
		attribute.String("weather.result_city", tempInfo.City),
		attribute.Float64("weather.temp_c", tempInfo.TempC),
		attribute.Float64("weather.temp_f", tempInfo.TempF),
		attribute.Float64("weather.temp_k", tempInfo.TempK),
	)

	log.Printf("Temperatura obtida para %s: %.1f°C, %.1f°F, %.1f K",
		tempInfo.City, tempInfo.TempC, tempInfo.TempF, tempInfo.TempK)

	return tempInfo, nil
}
