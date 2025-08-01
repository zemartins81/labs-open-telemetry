package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"servico-b/internal/models"
	"servico-b/internal/services"
	"servico-b/internal/validators"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TemperatureHandler é responsável por lidar com requisições de temperatura
type TemperatureHandler struct {
	temperatureService *services.TemperatureService
}

// NewTemperatureHandler cria uma nova instância do handler de temperatura
func NewTemperatureHandler(temperatureService *services.TemperatureService) *TemperatureHandler {
	return &TemperatureHandler{
		temperatureService: temperatureService,
	}
}

// HandleTemperature processa requisições de temperatura por CEP
func (h *TemperatureHandler) HandleTemperature(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("servico-b")
	ctx, span := tracer.Start(r.Context(), "HandleTemperature")
	defer span.End()

	// Configura headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		span.SetAttributes(attribute.String("http.method", r.Method))
		w.WriteHeader(http.StatusOK)
		return
	}

	// Apenas aceita método POST
	if r.Method != http.MethodPost {
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("error", "method not allowed"),
		)
		span.SetStatus(codes.Error, "method not allowed")
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	span.SetAttributes(attribute.String("http.method", r.Method))

	// Lê o body da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Erro ao ler body da requisição: %v", err)
		span.SetStatus(codes.Error, "failed to read request body")
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	// Parse do JSON
	var cepReq models.CEPRequest
	if err := json.Unmarshal(body, &cepReq); err != nil {
		log.Printf("Erro ao fazer parse do JSON: %v", err)
		span.SetStatus(codes.Error, "invalid json format")
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid json format")
		return
	}

	span.SetAttributes(attribute.String("cep.received", cepReq.CEP))

	// Valida o CEP
	if !validators.ValidateCEP(cepReq.CEP) {
		log.Printf("CEP inválido recebido: %s", cepReq.CEP)
		span.SetAttributes(attribute.Bool("cep.valid", false))
		span.SetStatus(codes.Error, "invalid zipcode")
		h.writeErrorResponse(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	span.SetAttributes(attribute.Bool("cep.valid", true))
	log.Printf("Processando CEP: %s", cepReq.CEP)

	// Busca temperatura pelo CEP
	temperatureInfo, err := h.temperatureService.GetTemperatureByCEP(ctx, cepReq.CEP)
	if err != nil {
		log.Printf("Erro ao buscar temperatura para CEP %s: %v", cepReq.CEP, err)

		// Verifica se é erro de CEP não encontrado
		if strings.Contains(err.Error(), "CEP não encontrado") {
			span.SetStatus(codes.Error, "zipcode not found")
			h.writeErrorResponse(w, http.StatusNotFound, "can not find zipcode")
			return
		}

		// Outros erros são considerados erro interno
		span.SetStatus(codes.Error, "internal server error")
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Monta resposta de sucesso
	response := models.TemperatureResponse{
		City:  temperatureInfo.City,
		TempC: temperatureInfo.TempC,
		TempF: temperatureInfo.TempF,
		TempK: temperatureInfo.TempK,
	}

	span.SetAttributes(
		attribute.String("city.name", response.City),
		attribute.Float64("temperature.celsius", response.TempC),
		attribute.Float64("temperature.fahrenheit", response.TempF),
		attribute.Float64("temperature.kelvin", response.TempK),
	)

	log.Printf("Resposta enviada para CEP %s: %s - %.1f°C",
		cepReq.CEP, response.City, response.TempC)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse escreve uma resposta de erro padronizada
func (h *TemperatureHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Message: message})
}
