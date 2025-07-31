package config

import "os"

// Config representa a configuração da aplicação
type Config struct {
	Port           string
	WeatherAPIKey  string
	ViaCEPURL      string
	WeatherAPIURL  string
}

// LoadConfig carrega as configurações a partir das variáveis de ambiente
func LoadConfig() *Config {
	return &Config{
		Port:           getEnv("PORT", "8081"),
		WeatherAPIKey:  getEnv("WEATHER_API_KEY", ""),
		ViaCEPURL:      getEnv("VIACEP_URL", "https://viacep.com.br/ws"),
		WeatherAPIURL:  getEnv("WEATHER_API_URL", "http://api.weatherapi.com/v1"),
	}
}

// getEnv retorna o valor da variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
