package main

import (
	"log"
	"servico-a/internal/config"
	"servico-a/internal/handlers"
	"servico-a/internal/server"
	"servico-a/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	serviceBClient := services.NewServiceBClient(cfg.ServiceBURL)

	cepHandler := handlers.NewCEPHandler(serviceBClient)

	srv := server.NewServer(cfg.Port, cepHandler)

	log.Printf("Serviço Iniciado na Porta %s", cfg.Port)

	if err := srv.Start(); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
