package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/thiagohmm/allLog/internal/entity"
)

type DTOIN struct {
	Table  string        `json:"tabela"`
	Fields []string      `json:"fields"`
	Values []interface{} `json:"values"`
}

type LogUseCase struct {
	Repo entity.MessageRepository
}

func NewLogUseCase(repo entity.MessageRepository) *LogUseCase {
	if repo == nil {
		log.Fatal("repository cannot be nil")
	}
	return &LogUseCase{
		Repo: repo,
	}
}

func (l *LogUseCase) UsecaseSaveLog(ctx context.Context, dto DTOIN) error {
	if ctx == nil {
		return errors.New("context is nil")
	}
	if l.Repo == nil {
		return errors.New("repository is nil")
	}

	message := &entity.Message{
		Tabela: dto.Table,
		Fields: dto.Fields,
		Values: dto.Values,
	}
	fmt.Println("Processing message for table:", message.Tabela)

	return l.Repo.SaveMessage(ctx, *message)
}
