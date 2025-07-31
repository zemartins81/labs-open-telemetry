package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"servico-a/internal/models"
	"servico-a/internal/services"
	"servico-a/internal/validators"
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
	// Configura headers CORS e Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Apenas aceita método POST
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Lê o body da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Erro ao ler body da requisição: %v", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	// Parse do JSON
	var cepReq models.CEPRequest
	if err := json.Unmarshal(body, &cepReq); err != nil {
		log.Printf("Erro ao fazer parse do JSON: %v", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid json format")
		return
	}

	// Valida o CEP
	if !validators.ValidateCEP(cepReq.CEP) {
		log.Printf("CEP inválido recebido: %s", cepReq.CEP)
		h.writeErrorResponse(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	log.Printf("CEP válido recebido: %s", cepReq.CEP)

	// Encaminha para o Serviço B
	response, err := h.serviceBClient.ForwardCEPRequest(cepReq)
	if err != nil {
		log.Printf("Erro ao comunicar com Serviço B: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Retorna a resposta do Serviço B
	w.WriteHeader(response.StatusCode)
	w.Write(response.Body)
}

// writeErrorResponse escreve uma resposta de erro padronizada
func (h *CEPHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Message: message})
}
