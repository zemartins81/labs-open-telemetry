package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"servico-b/internal/models"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
			Timeout:   10 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// GetLocationByCEP busca informações de localização pelo CEP
func (v *ViaCEPService) GetLocationByCEP(ctx context.Context, cep string) (*models.LocationInfo, error) {
	tracer := otel.Tracer("servico-b")
	ctx, span := tracer.Start(ctx, "ViaCEP.GetLocationByCEP")
	defer span.End()

	url := fmt.Sprintf("%s/%s/json/", v.baseURL, cep)

	span.SetAttributes(
		attribute.String("viacep.cep", cep),
		attribute.String("viacep.url", url),
	)

	log.Printf("Buscando informações do CEP %s na ViaCEP: %s", cep, url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.SetStatus(codes.Error, "failed to create request")
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	resp, err := v.client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "failed to make request")
		return nil, fmt.Errorf("erro ao fazer requisição para ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("viacep.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "non-200 status code")
		return nil, fmt.Errorf("ViaCEP retornou status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(codes.Error, "failed to read response")
		return nil, fmt.Errorf("erro ao ler resposta da ViaCEP: %w", err)
	}

	var viaCEPResp models.ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResp); err != nil {
		span.SetStatus(codes.Error, "failed to parse response")
		return nil, fmt.Errorf("erro ao fazer parse da resposta ViaCEP: %w", err)
	}

	// Verifica se o CEP foi encontrado
	if viaCEPResp.Erro.Bool() {
		log.Printf("CEP %s não encontrado na ViaCEP", cep)
		span.SetAttributes(attribute.Bool("viacep.found", false))
		span.SetStatus(codes.Error, "CEP not found")
		return nil, fmt.Errorf("CEP não encontrado")
	}

	// Verifica se os campos essenciais estão presentes
	if viaCEPResp.Localidade == "" {
		log.Printf("CEP %s retornou dados incompletos da ViaCEP", cep)
		span.SetStatus(codes.Error, "incomplete location data")
		return nil, fmt.Errorf("dados de localização incompletos")
	}

	location := &models.LocationInfo{
		City:  viaCEPResp.Localidade,
		State: viaCEPResp.UF,
	}

	span.SetAttributes(
		attribute.Bool("viacep.found", true),
		attribute.String("viacep.city", location.City),
		attribute.String("viacep.state", location.State),
	)

	log.Printf("CEP %s encontrado: %s/%s", cep, location.City, location.State)

	return location, nil
}
