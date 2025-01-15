package service

import (
    "log"
    "project/internal/entity"
    "project/internal/repository"
)

type MessageService struct {
    Repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
    return &MessageService{Repo: repo}
}

func (s *MessageService) ProcessMessage(content string) {
    message := entity.Message{Content: content}
    err := s.Repo.SaveMessage(message)
    if err != nil {
        log.Printf("Error saving message: %v", err)
    } else {
        log.Println("Message saved successfully")
    }
}
