package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"servico-a/internal/models"
)

// ServiceBClient é responsável pela comunicação com o Serviço B
type ServiceBClient struct {
	baseURL string
	client  *http.Client
}

// NewServiceBClient cria uma nova instância do cliente do Serviço B
func NewServiceBClient(baseURL string) *ServiceBClient {
	return &ServiceBClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ForwardCEPRequest encaminha a requisição de CEP para o Serviço B
func (s *ServiceBClient) ForwardCEPRequest(cepReq models.CEPRequest) (*models.ServiceBResponse, error) {
	// Converte para JSON
	jsonData, err := json.Marshal(cepReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar JSON: %w", err)
	}

	// Cria a requisição para o Serviço B
	url := s.baseURL + "/temperature"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Faz a requisição
	log.Printf("Encaminhando CEP %s para Serviço B: %s", cepReq.CEP, url)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para Serviço B: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do Serviço B: %w", err)
	}

	log.Printf("Resposta do Serviço B - Status: %d, Body: %s", resp.StatusCode, string(body))

	return &models.ServiceBResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
