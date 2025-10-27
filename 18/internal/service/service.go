package service

import (
	"fmt"
	"wb_l2/18/internal/repository"
)

var InvalidQuery = fmt.Errorf("Invalid query provided")

type Service struct {
	Event *EventService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Event: NewEventService(repo),
	}
}
