package main

import (
	"context"
	"log"

	"github.com/thiagohmm/allLog/configuration"
	"github.com/thiagohmm/allLog/internal/service"
)

func LoadConfig() (*configuration.Conf, error) {
	// Carrega as configurações do arquivo .env
	cfg, err := configuration.LoadConfig("../../.env")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}
	return cfg, err
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	log.Printf("Configuração carregada com Sucesso")

	// Inicializa o serviço de mensagens
	messageService := service.MessageService{}

	messageService.ListenToQueue(context.Background(), cfg.ENV_RABBITMQ, "logs")

}
