package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	fmt.Println("Fields:", message.Fields)
	fmt.Println("Values:", message.Values)
	fmt.Println("Mensagem inteira", message)

	placeholders := make([]string, len(message.Fields))
	for i := range message.Fields {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		message.Tabela,
		quoteIdentifiers(message.Fields),
		strings.Join(placeholders, ", "),
	)

	// Exibe a query para debug
	fmt.Println("Query:", query)

	_, err := m.DB.ExecContext(ctx, query, message.Values...)
	return err
}
