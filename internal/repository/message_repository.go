package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/thiagohmm/allLog/internal/entity"
)

type MessageRepositoryDB struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepositoryDB {
	if db == nil {
		panic("database connection is nil") // Prevent invalid instantiation
	}
	return &MessageRepositoryDB{DB: db}
}

func quoteIdentifiers(fields []string) string {
	quoted := make([]string, len(fields))
	for i, field := range fields {
		quoted[i] = fmt.Sprintf("\"%s\"", field)
	}
	return strings.Join(quoted, ", ")
}

func (m *MessageRepositoryDB) SaveMessage(ctx context.Context, message entity.Message) error {
	fmt.Println("Salvando mensagem na tabela:", message.Tabela)

	// Converte valores de data para time.Time se os campos forem DATARECEBIMENTO ou DATAPROCESSAMENTO
	for i, field := range message.Fields {
		if field == "DATARECEBIMENTO" || field == "DATAPROCESSAMENTO" {
			if dateStr, ok := message.Values[i].(string); ok {
				t, err := time.Parse(time.RFC3339, dateStr)
				if err != nil {
					return fmt.Errorf("failed to parse date for field %s: %w", field, err)
				}
				message.Values[i] = t
			}
		}
	}

	// Prepara os placeholders no formato ":1, :2, ..." compatível com Oracle
	placeholders := make([]string, len(message.Fields))
	for i := range message.Fields {
		placeholders[i] = fmt.Sprintf(":%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		message.Tabela,
		quoteIdentifiers(message.Fields),
		strings.Join(placeholders, ", "),
	)

	// Exibe a query para debug
	fmt.Println("Data:", time.Now().Format("2006-01-02 15:04:05"), "- Query:", query)
	fmt.Println("Values:", message.Values)

	_, err := m.DB.ExecContext(ctx, query, message.Values...)
	return err
}
