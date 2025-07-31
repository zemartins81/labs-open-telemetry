package main

import (
	"log"

	"servico-b/internal/config"
	"servico-b/internal/handlers"
	"servico-b/internal/server"
	"servico-b/internal/services"
)

func main() {
	// Carrega configuração
	cfg := config.LoadConfig()

	// Valida configuração crítica
	if cfg.WeatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY é obrigatória. Obtenha uma chave em https://www.weatherapi.com/")
	}

	// Inicializa serviços
	viaCEPService := services.NewViaCEPService()
	weatherService := services.NewWeatherService(cfg.WeatherAPIKey)
	temperatureService := services.NewTemperatureService(viaCEPService, weatherService)

	// Inicializa handlers
	temperatureHandler := handlers.NewTemperatureHandler(temperatureService)

	// Inicializa servidor
	srv := server.NewServer(cfg.Port, temperatureHandler)

	log.Printf("Serviço B iniciado na porta %s", cfg.Port)
	log.Printf("Usando WeatherAPI com chave: %s...", cfg.WeatherAPIKey[:min(len(cfg.WeatherAPIKey), 8)])

	// Inicia o servidor
	if err := srv.Start(); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
