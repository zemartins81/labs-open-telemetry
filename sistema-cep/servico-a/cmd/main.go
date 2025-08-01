package main

import (
	"log"
	"servico-a/internal/config"
	"servico-a/internal/handlers"
	"servico-a/internal/server"
	"servico-a/internal/services"
	"servico-a/internal/telemetry"
)

func main() {
	cfg := config.LoadConfig()

	shutdown, err := telemetry.InitTracer("servico-a")
	if err != nil {
		log.Fatal("Erro ao inicializar telemetria: ", err)
	}
	defer shutdown()

	serviceBClient := services.NewServiceBClient(cfg.ServiceBURL)

	cepHandler := handlers.NewCEPHandler(serviceBClient)

	srv := server.NewServer(cfg.Port, cepHandler)

	log.Printf("Serviço A iniciado na porta %s com tracing habilitado", cfg.Port)

	if err := srv.Start(); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
