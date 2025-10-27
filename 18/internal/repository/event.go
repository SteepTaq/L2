package repository

import (
	"time"
	"wb_l2/18/internal/model"
)

type eventRepository interface {
	Create(event *model.Event) (int, error)
	ListForDay(userID int, date time.Time) ([]*model.Event, error)
	ListForWeek(userID int, starting time.Time) ([]*model.Event, error)
	ListForMonth(userID int, starting time.Time) ([]*model.Event, error)
	Update(ID int, event *model.Event) (*model.Event, error)
	Delete(ID int) error
}
