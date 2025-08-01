package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"servico-a/internal/models"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
			Timeout:   30 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// ForwardCEPRequest encaminha a requisição de CEP para o Serviço B
func (s *ServiceBClient) ForwardCEPRequest(ctx context.Context, cepReq models.CEPRequest) (*models.ServiceBResponse, error) {
	tracer := otel.Tracer("servico-a")
	ctx, span := tracer.Start(ctx, "ForwardCEPRequest")
	defer span.End()

	span.SetAttributes(
		attribute.String("service.name", "servico-b"),
		attribute.String("cep.value", cepReq.CEP),
	)

	// Converte para JSON
	jsonData, err := json.Marshal(cepReq)
	if err != nil {
		span.SetStatus(codes.Error, "failed to marshal JSON")
		return nil, fmt.Errorf("erro ao serializar JSON: %w", err)
	}

	// Cria a requisição para o Serviço B
	url := s.baseURL + "/temperature"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		span.SetStatus(codes.Error, "failed to create request")
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	span.SetAttributes(
		attribute.String("http.method", "POST"),
		attribute.String("http.url", url),
	)

	// Faz a requisição
	log.Printf("Encaminhando CEP %s para Serviço B: %s", cepReq.CEP, url)
	resp, err := s.client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "failed to make request")
		return nil, fmt.Errorf("erro ao fazer requisição para Serviço B: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(codes.Error, "failed to read response")
		return nil, fmt.Errorf("erro ao ler resposta do Serviço B: %w", err)
	}

	span.SetAttributes(attribute.Int("response.body_size", len(body)))
	log.Printf("Resposta do Serviço B - Status: %d, Body: %s", resp.StatusCode, string(body))

	return &models.ServiceBResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
