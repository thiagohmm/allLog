package usecase

import (
	"context"

	"github.com/thiagohmm/allLog/internal/entity"
)

type DTOIN struct {
	Table  string                 `json:"table"`
	Fields map[string]interface{} `json:"fields"`
	Values map[string]interface{} `json:"values"`
}

type LogUseCase struct {
	Repo entity.MessageRepository
}

func NewLogUseCase(repo entity.MessageRepository) *LogUseCase {
	return &LogUseCase{
		Repo: repo,
	}
}

func (l *LogUseCase) SaveLog(ctx context.Context, dto DTOIN) error {
	message := entity.Message{
		Table:  dto.Table,
		Fields: dto.Fields,
		Values: dto.Values,
	}
	err := l.Repo.SaveMessage(message)
	if err != nil {
		return err
	}
	return nil
}
