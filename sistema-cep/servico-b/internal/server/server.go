package server

import (
	"net/http"

	"servico-b/internal/handlers"
)

// Server representa o servidor HTTP
type Server struct {
	port               string
	temperatureHandler *handlers.TemperatureHandler
}

// NewServer cria uma nova instância do servidor
func NewServer(port string, temperatureHandler *handlers.TemperatureHandler) *Server {
	return &Server{
		port:               port,
		temperatureHandler: temperatureHandler,
	}
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	// Configura as rotas
	s.setupRoutes()

	// Inicia o servidor
	return http.ListenAndServe(":"+s.port, nil)
}

// setupRoutes configura as rotas da aplicação
func (s *Server) setupRoutes() {
	http.HandleFunc("/temperature", s.temperatureHandler.HandleTemperature)
	http.HandleFunc("/health", s.healthCheck)
}

// healthCheck endpoint para verificação de saúde da aplicação
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "servico-b"}`))
}
