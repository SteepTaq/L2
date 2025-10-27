package model

import (
	"encoding/json"
	"fmt"
	"time"
	"wb_l2/18/pkg/date"
)

var InvalidFormat = fmt.Errorf("Invalid body format")

type Event struct {
	ID   int
	Name string
	Date time.Time

	UserID int
}

type EventOut struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`

	UserID int `json:"user_id"`
}

func EventFromBody(body []byte) (*Event, error) {
	var eventParse EventOut
	if err := json.Unmarshal(body, &eventParse); err != nil {
		return nil, InvalidFormat
	}
	event := new(Event)

	if eventParse.Name == "" {
		return nil, InvalidFormat
	}
	event.Name = eventParse.Name

	if eventParse.Date == "" {
		return nil, InvalidFormat
	}

	time, err := date.TimeFromString(eventParse.Date)
	if err != nil {
		return nil, InvalidFormat
	}
	event.Date = time

	if eventParse.UserID <= 0 {
		return nil, InvalidFormat
	}
	event.UserID = eventParse.UserID

	return event, nil
}

func (e *Event) FormatDate() *EventOut {
	time := date.StringFromTime(e.Date)

	return &EventOut{
		ID:     e.ID,
		Name:   e.Name,
		Date:   time,
		UserID: e.UserID,
	}
}
