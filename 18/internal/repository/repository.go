package repository

import (
	"fmt"
	inmemory "wb_l2/18/internal/repository/inmemory/event"
)

type StorageType int

const (
	InMemory StorageType = iota
)

var storageTypeName = map[StorageType]string{
	InMemory: "in-memory",
}

func (st StorageType) String() string {
	return storageTypeName[st]
}

type Repository struct {
	Event eventRepository
}

func NewRepository(storageType StorageType) *Repository {
	switch storageType {
	case InMemory:
		return &Repository{
			Event: inmemory.NewEventRepositoryInMemory(),
		}
	default:
		panic(fmt.Errorf("Unknown repository storage type: %s", storageType))
	}
}
