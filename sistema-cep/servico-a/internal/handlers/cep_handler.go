package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"servico-a/internal/models"
	"servico-a/internal/services"
	"servico-a/internal/validators"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CEPHandler é responsável por lidar com requisições de CEP
type CEPHandler struct {
	serviceBClient *services.ServiceBClient
}

// NewCEPHandler cria uma nova instância do handler de CEP
func NewCEPHandler(serviceBClient *services.ServiceBClient) *CEPHandler {
	return &CEPHandler{
		serviceBClient: serviceBClient,
	}
}

// HandleCEP processa requisições de CEP
func (h *CEPHandler) HandleCEP(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("servico-a")
	ctx, span := tracer.Start(r.Context(), "HandleCEP")
	defer span.End()

	// Configura headers CORS e Content-Type
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
	log.Printf("CEP válido recebido: %s", cepReq.CEP)

	// Encaminha para o Serviço B
	response, err := h.serviceBClient.ForwardCEPRequest(ctx, cepReq)
	if err != nil {
		log.Printf("Erro ao comunicar com Serviço B: %v", err)
		span.SetStatus(codes.Error, "failed to communicate with service B")
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	span.SetAttributes(
		attribute.Int("response.status_code", response.StatusCode),
		attribute.Int("response.body_size", len(response.Body)),
	)

	// Retorna a resposta do Serviço B
	w.WriteHeader(response.StatusCode)
	w.Write(response.Body)
}

// writeErrorResponse escreve uma resposta de erro padronizada
func (h *CEPHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Message: message})
}
