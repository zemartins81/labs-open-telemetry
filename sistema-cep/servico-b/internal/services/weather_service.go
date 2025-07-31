package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"servico-b/internal/models"
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
			Timeout: 10 * time.Second,
		},
	}
}

// GetTemperatureByLocation busca informações de temperatura pela localização
func (w *WeatherService) GetTemperatureByLocation(location *models.LocationInfo) (*models.TemperatureInfo, error) {
	// Constrói a query de localização
	query := fmt.Sprintf("%s,%s", location.City, location.State)

	// Constrói a URL com parâmetros
	apiURL := fmt.Sprintf("%s/current.json", w.baseURL)
	params := url.Values{}
	params.Add("key", w.apiKey)
	params.Add("q", query)
	params.Add("aqi", "no") // Não precisamos de dados de qualidade do ar

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	log.Printf("Buscando temperatura para %s na WeatherAPI", query)

	resp, err := w.client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da WeatherAPI: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("WeatherAPI retornou erro %d: %s", resp.StatusCode, string(body))

		// Tenta fazer parse da mensagem de erro
		var errorResp map[string]interface{}
		if json.Unmarshal(body, &errorResp) == nil {
			if errorMsg, ok := errorResp["error"]; ok {
				return nil, fmt.Errorf("erro da WeatherAPI: %v", errorMsg)
			}
		}

		return nil, fmt.Errorf("WeatherAPI retornou status %d", resp.StatusCode)
	}

	var weatherResp models.WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
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

	log.Printf("Temperatura obtida para %s: %.1f°C, %.1f°F, %.1f K",
		tempInfo.City, tempInfo.TempC, tempInfo.TempF, tempInfo.TempK)

	return tempInfo, nil
}
