package models

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var vld = validator.New()

// ErrNotFound is returned when a resource is not found
var ErrNotFound = errors.New("not found")

// Date is a type for date
type Date int

const (
	// Day is a type for day
	Day Date = iota
	// Week is a type for week
	Week
	// Month is a type for month
	Month
)

// DateOnly is a struct for date without time
type DateOnly struct {
	t time.Time
}

// Time returns time
func (d *DateOnly) Time() time.Time {
	return d.t
}

// UnmarshalJSON unmarshals date
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.t = t
	return nil
}

// MarshalJSON marshals date
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time().Format("2006-01-02") + `"`), nil
}

// Event is a type for event
type Event struct {
	ID          int64
	UserID      int64    `json:"user_id" validate:"required"`
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Date        DateOnly `json:"date" validate:"required"`
}

// Validate validates event
func (e *Event) Validate() error {
	return vld.Struct(e)
}

// Events is a type for events
type Events struct {
	Events []*Event
}

// Add adds event to events
func (es *Events) Add(e *Event) {
	if e != nil {
		if len(es.Events) == 0 {
			e.ID = 1
		} else {
			e.ID = es.Events[len(es.Events)-1].ID + 1
		}
		es.Events = append(es.Events, e)
	}
}

// Remove removes event from events
func (es *Events) Remove(id int64) error {
	for i, e := range es.Events {
		if e.ID == id {
			es.Events = append(es.Events[:i], es.Events[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

// Get returns event by id
func (es *Events) Get(id int64) (*Event, error) {
	for _, e := range es.Events {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, ErrNotFound
}

// Update updates event
func (es *Events) Update(e *Event) error {
	for i, event := range es.Events {
		if event.ID == e.ID {
			es.Events[i] = e
			return nil
		}
	}
	return ErrNotFound
}

// GetSince returns events since date
func (es *Events) GetSince(userID int64, t time.Time, d Date) *Events {
	result := new(Events)

	var same func(t1, t2 time.Time) bool

	switch d {
	case Day:
		same = sameDay
	case Week:
		same = sameWeek
	case Month:
		same = sameMonth
	default:
		return result
	}

	t = normalizeDay(t)

	for _, e := range es.Events {
		if userID == e.UserID && same(normalizeDay(e.Date.Time()), t) {
			result.Add(e)
		}
	}

	return result
}

// normalizeDay cuts time to DD-MM-YYYY
func normalizeDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func sameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

func sameWeek(t1, t2 time.Time) bool {
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	return y1 == y2 && w1 == w2
}

func sameMonth(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}
