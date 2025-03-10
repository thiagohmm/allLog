package main

import (
	"context"
	"fmt"
	"log"

	"github.com/thiagohmm/allLog/configuration"

	"github.com/thiagohmm/allLog/internal/database"
	"github.com/thiagohmm/allLog/internal/repository"
	"github.com/thiagohmm/allLog/internal/service"
	"github.com/thiagohmm/allLog/internal/usecase"
)

func LoadConfig() (*configuration.Conf, error) {
	// Carrega as configurações do arquivo .env
	cfg, err := configuration.LoadConfig("../../.env")
	fmt.Println(cfg)
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

	db, err := database.ConectarBanco(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializa o serviço de mensagens
	messageService := service.MessageService{}
	messageService.UseCase = usecase.NewLogUseCase(repository.NewMessageRepository(db))

	messageService.ListenToQueue(context.Background(), cfg.ENV_RABBITMQ, "logs")

}
