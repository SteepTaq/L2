package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
	"wb_l2/18/internal/model"
	"wb_l2/18/internal/repository"
	"wb_l2/18/pkg/date"
)

type EventService struct {
	repo *repository.Repository
}

func NewEventService(repo *repository.Repository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(body []byte) (int, error) {
	event, err := model.EventFromBody(body)
	if err != nil {
		return 0, err
	}

	return s.repo.Event.Create(event)
}

type ListFor int

const (
	Day ListFor = iota
	Week
	Month
)

var listForName = map[ListFor]string{
	Day:   "day",
	Week:  "week",
	Month: "month",
}

func (st ListFor) String() string {
	return listForName[st]
}

func (s *EventService) List(query url.Values, by ListFor) ([]*model.EventOut, error) {
	userIDStr := query.Get("user_id")
	if userIDStr == "" {
		return []*model.EventOut{}, InvalidQuery
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return []*model.EventOut{}, InvalidQuery
	}

	dateStr := query.Get("date")
	if dateStr == "" {
		return []*model.EventOut{}, InvalidQuery
	}

	time, err := date.TimeFromString(dateStr)
	if err != nil {
		return []*model.EventOut{}, InvalidQuery
	}

	var res []*model.Event

	switch by {
	case Day:
		res, err = s.repo.Event.ListForDay(userID, time)
	case Week:
		res, err = s.repo.Event.ListForWeek(userID, time)
	case Month:
		res, err = s.repo.Event.ListForMonth(userID, time)
	default:
		panic(fmt.Errorf("Unknown event list type: %s", listForName))
	}

	prettify := make([]*model.EventOut, 0, len(res))
	for _, event := range res {
		prettify = append(prettify, event.FormatDate())
	}

	return prettify, err
}

func (s *EventService) Update(body []byte) (*model.EventOut, error) {
	var eventParse model.EventOut
	if err := json.Unmarshal(body, &eventParse); err != nil {
		return nil, model.InvalidFormat
	}

	var time time.Time
	if eventParse.Date != "" {
		parsed, err := date.TimeFromString(eventParse.Date)
		if err != nil {
			return nil, model.InvalidFormat
		}

		time = parsed
	}

	event := &model.Event{
		Name: eventParse.Name,
		Date: time,
	}

	event, err := s.repo.Event.Update(eventParse.ID, event)
	if err != nil {
		return nil, err
	}

	return event.FormatDate(), nil
}

func (s *EventService) Delete(body []byte) error {
	var eventParse model.EventOut
	if err := json.Unmarshal(body, &eventParse); err != nil {
		return model.InvalidFormat
	}

	return s.repo.Event.Delete(eventParse.ID)
}
