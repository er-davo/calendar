package service

import (
	"calendar/internal/models"
	"time"
)

// Service interface for event management
type Service interface {
	CreateEvent(event *models.Event)
	GetEvents(id int64) (*models.Event, error)
	UpdateEvent(e *models.Event) error
	DeleteEvent(id int64) error
	GetSince(userID int64, t time.Time, d models.Date) *models.Events
}

type eventService struct {
	events models.Events
}

// NewService creates a new eventService
func NewService() Service {
	return &eventService{
		events: models.Events{},
	}
}

// CreateEvent creates a new event
func (es *eventService) CreateEvent(event *models.Event) {
	es.events.Add(event)
}

// GetEvents returns a specific event
func (es *eventService) GetEvents(id int64) (*models.Event, error) {
	return es.events.Get(id)
}

// UpdateEvent updates an existing event
func (es *eventService) UpdateEvent(e *models.Event) error {
	return es.events.Update(e)
}

// DeleteEvent deletes an event
func (es *eventService) DeleteEvent(id int64) error {
	return es.events.Remove(id)
}

// GetSince returns events since a given time
func (es *eventService) GetSince(userID int64, t time.Time, d models.Date) *models.Events {
	return es.events.GetSince(userID, t, d)
}
