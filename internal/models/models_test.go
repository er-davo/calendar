package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventValidate(t *testing.T) {
	e := &Event{
		UserID:      1,
		Title:       "Test",
		Description: "Desc",
		Date:        DateOnly{t: time.Now()},
	}
	assert.NoError(t, e.Validate())

	e2 := &Event{}
	assert.Error(t, e2.Validate())
}

func TestEventsAddAndGet(t *testing.T) {
	es := &Events{}

	e1 := &Event{UserID: 1, Title: "E1", Description: "D1", Date: DateOnly{t: time.Now()}}
	es.Add(e1)

	assert.Equal(t, int64(1), e1.ID)

	got, err := es.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, e1, got)
}

func TestEventsRemove(t *testing.T) {
	es := &Events{}
	e1 := &Event{UserID: 1, Title: "E1", Description: "D1", Date: DateOnly{t: time.Now()}}
	es.Add(e1)

	err := es.Remove(e1.ID)
	assert.NoError(t, err)
	assert.Empty(t, es.Events)

	err = es.Remove(42)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestEventsUpdate(t *testing.T) {
	es := &Events{}
	e1 := &Event{UserID: 1, Title: "E1", Description: "D1", Date: DateOnly{t: time.Now()}}
	es.Add(e1)

	e1.Title = "Updated"
	err := es.Update(e1)
	assert.NoError(t, err)

	got, _ := es.Get(e1.ID)
	assert.Equal(t, "Updated", got.Title)

	e2 := &Event{ID: 42, UserID: 1, Title: "X"}
	assert.ErrorIs(t, es.Update(e2), ErrNotFound)
}

func TestEventsGetSince(t *testing.T) {
	es := &Events{}

	day := time.Date(2025, 10, 4, 0, 0, 0, 0, time.UTC)
	e1 := &Event{UserID: 1, Title: "Today", Description: "D1", Date: DateOnly{t: day}}
	e2 := &Event{UserID: 1, Title: "Next Day", Description: "D2", Date: DateOnly{t: day.AddDate(0, 0, 1)}}
	e3 := &Event{UserID: 2, Title: "Other User", Description: "D3", Date: DateOnly{t: day}}

	es.Add(e1)
	es.Add(e2)
	es.Add(e3)

	// Должно найти только e1
	resultDay := es.GetSince(1, day, Day)
	assert.Len(t, resultDay.Events, 1)
	assert.Equal(t, "Today", resultDay.Events[0].Title)

	// Должно найти e1 и e2 в одной неделе
	resultWeek := es.GetSince(1, day, Week)
	assert.Len(t, resultWeek.Events, 2)

	// По месяцу тоже найдутся оба
	resultMonth := es.GetSince(1, day, Month)
	assert.Len(t, resultMonth.Events, 2)
}

func TestDateOnlyMarshalUnmarshal(t *testing.T) {
	d := DateOnly{t: time.Date(2025, 10, 4, 0, 0, 0, 0, time.UTC)}

	// Marshal
	data, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-10-04"`, string(data))

	// Unmarshal
	var d2 DateOnly
	err = json.Unmarshal([]byte(`"2025-10-04"`), &d2)
	assert.NoError(t, err)
	assert.Equal(t, d.Time().Year(), d2.Time().Year())
	assert.Equal(t, d.Time().Month(), d2.Time().Month())
	assert.Equal(t, d.Time().Day(), d2.Time().Day())
}
