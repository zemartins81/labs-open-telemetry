package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"servico-b/internal/models"
)

// ViaCEPService é responsável pela comunicação com a API ViaCEP
type ViaCEPService struct {
	baseURL string
	client  *http.Client
}

// NewViaCEPService cria uma nova instância do serviço ViaCEP
func NewViaCEPService() *ViaCEPService {
	return &ViaCEPService{
		baseURL: "https://viacep.com.br/ws",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetLocationByCEP busca informações de localização pelo CEP
func (v *ViaCEPService) GetLocationByCEP(cep string) (*models.LocationInfo, error) {
	url := fmt.Sprintf("%s/%s/json/", v.baseURL, cep)

	log.Printf("Buscando informações do CEP %s na ViaCEP: %s", cep, url)

	resp, err := v.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP retornou status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da ViaCEP: %w", err)
	}

	var viaCEPResp models.ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResp); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da resposta ViaCEP: %w", err)
	}

	// Verifica se o CEP foi encontrado
	if viaCEPResp.Erro {
		log.Printf("CEP %s não encontrado na ViaCEP", cep)
		return nil, fmt.Errorf("CEP não encontrado")
	}

	// Verifica se os campos essenciais estão presentes
	if viaCEPResp.Localidade == "" {
		log.Printf("CEP %s retornou dados incompletos da ViaCEP", cep)
		return nil, fmt.Errorf("dados de localização incompletos")
	}

	location := &models.LocationInfo{
		City:  viaCEPResp.Localidade,
		State: viaCEPResp.UF,
	}

	log.Printf("CEP %s encontrado: %s/%s", cep, location.City, location.State)

	return location, nil
}
