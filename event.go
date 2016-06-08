package main

import (
	"encoding/json"
	"fmt"

	"google.golang.org/api/calendar/v3"
)

// Event is the structural representation of an event.
type Event struct {
	Source        string
	SourceID      int
	CalendarEvent *calendar.Event
	Calendar      string
}

// ToKey returns the key associated with the Event.
func (e *Event) ToKey() *Key {
	return NewKey(
		"event",
		"id",
		fmt.Sprintf("%s/%d", e.Source, e.SourceID),
	)
}

// EventList is just a list of Event types, we just want to add some helper
// methods to it.
type EventList []*Event

// Keys gets a slice of keys
func (el EventList) Keys() KeyList {
	var keys KeyList

	for _, e := range el {
		key := e.ToKey()
		keys = append(keys, key)
	}

	return keys
}

// ExcludeKeys returns an EventList of events excluding the ones with the
// provided keys.
func (el EventList) ExcludeKeys(keys KeyList) (list EventList) {
	for _, key := range keys {
		for _, event := range el {
			if KeysEqual(key, event.ToKey()) {
				list = append(list, event)
			}
		}
	}

	return
}

// ToKeyValue converts the EventList to a slice with alternating keys and JSON
// encoded values. Seems like a strange format, but thats how redis.MSet likes
// it.
func (el EventList) ToKeyValue() ([]string, error) {
	var m []string

	for _, event := range el {
		key := event.ToKey()
		value, e := json.Marshal(event)
		if e != nil {
			return nil, e
		}

		m = append(m, key.String(), string(value))
	}

	return m, nil
}
