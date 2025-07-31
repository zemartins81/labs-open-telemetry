package main

import (
	"log"

	"github.com/go-delve/delve/pkg/config"
	"github.com/jesseduffield/go-git/v5/plumbing/transport/server"
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
