package event

import (
	"fmt"
	"iter"
	"maps"
	"sync"
	"time"
	"wb_l2/18/internal/model"
)

var ErrorEventNotFound = fmt.Errorf("Event is not found")

type EventRepository struct {
	events        map[int]*model.Event
	autoincrement int

	mu sync.Mutex
}

func NewEventRepositoryInMemory() *EventRepository {
	return &EventRepository{
		events:        make(map[int]*model.Event),
		autoincrement: 1,
	}
}

func (r *EventRepository) Create(event *model.Event) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.autoincrement
	event.ID = id
	r.events[id] = event

	r.autoincrement++

	return id, nil
}

func (r *EventRepository) ListForDay(userID int, date time.Time) ([]*model.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]*model.Event, 0)
	next, stop := iter.Pull(maps.Values(r.events))

	for {
		event, ok := next()
		if !ok {
			stop()
			break
		}

		if event.UserID != userID {
			continue
		}

		if !event.Date.Equal(date) {
			continue
		}

		res = append(res, event)
	}

	return res, nil
}

func (r *EventRepository) ListForWeek(userID int, starting time.Time) ([]*model.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]*model.Event, 0)
	next, stop := iter.Pull(maps.Values(r.events))

	ending := starting.AddDate(0, 0, 7)

	for {
		event, ok := next()
		if !ok {
			stop()
			break
		}

		if event.UserID != userID {
			continue
		}

		if event.Date.Compare(starting) < 0 || event.Date.Compare(ending) > 0 {
			continue
		}

		res = append(res, event)
	}

	return res, nil
}

func (r *EventRepository) ListForMonth(userID int, starting time.Time) ([]*model.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]*model.Event, 0)
	next, stop := iter.Pull(maps.Values(r.events))

	ending := starting.AddDate(0, 1, 0)

	for {
		event, ok := next()
		if !ok {
			stop()
			break
		}

		if event.UserID != userID {
			continue
		}

		if event.Date.Compare(starting) < 0 || event.Date.Compare(ending) > 0 {
			continue
		}

		res = append(res, event)
	}

	return res, nil
}

func (r *EventRepository) Update(ID int, event *model.Event) (*model.Event, error) {
	if _, ok := r.events[ID]; !ok {
		return nil, ErrorEventNotFound
	}

	if event.Name != "" {
		r.events[ID].Name = event.Name
	}

	if !event.Date.IsZero() {
		r.events[ID].Date = event.Date
	}

	return r.events[ID], nil
}

func (r *EventRepository) Delete(ID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.events[ID]; !ok {
		return ErrorEventNotFound
	}

	delete(r.events, ID)
	return nil
}
