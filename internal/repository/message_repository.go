package repository





type MessageRepositoryDB struct {
  DB *sql.DB
}

func (m *MessageRepositoryDB) SaveMessage(message *entity.Message) error {
  query := fmt.Sprintf("INSERT INTO %s (", message.Table)
  var fields []string
  var placeholders []string
  var values []interface{}
  for field, value := range message.Values {
    fields = append(fields, field)
    placeholders = append(placeholders, "?")
    values = append(values, value)
  }
  query += strings.Join(fields, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"
  _, err := m.DB.Exec(query, values...)
  return err
}