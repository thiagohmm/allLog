package entity







type MessageRepository interface {
    SaveMessage(message Message) error
}