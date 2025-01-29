package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Message struct {
	table  string
	fields map[string]interface{}
}

type MessageRepositoryDB struct {
	DB *sql.DB
}

func (m *MessageRepositoryDB) SaveMessage(ctx context.Context, message *Message) error {
	query := fmt.Sprintf("INSERT INTO %s (", message.table)
	var fields []string
	var placeholders []string
	var values []interface{}
	for field, value := range message.fields {
		fields = append(fields, field)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}
	query += strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"
	_, err := m.DB.Exec(query, values...)
	return err
}
