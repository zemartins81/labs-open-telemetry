package server

import (
	"net/http"

	"servico-b/internal/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	mux := s.setupRoutes()

	// Inicia o servidor com instrumentação OpenTelemetry
	handler := otelhttp.NewHandler(mux, "servico-b")
	return http.ListenAndServe(":"+s.port, handler)
}

// setupRoutes configura as rotas da aplicação
func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/temperature", s.temperatureHandler.HandleTemperature)
	mux.HandleFunc("/health", s.healthCheck)
	return mux
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
