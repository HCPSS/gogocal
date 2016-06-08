package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
	"gopkg.in/redis.v3"
)

// GoGoCal is the representation of our application.
type GoGoCal struct {
	// The repository of events.
	rc *redis.Client

	// The repository of Calendars
	cs *calendar.Service

	// The logger
	l *log.Logger
}

// SetCalendarService sets the calendar service.
func (ggc *GoGoCal) SetCalendarService(cs *calendar.Service) {
	ggc.cs = cs
}

// SetRedisClient sets the redis client.
func (ggc *GoGoCal) SetRedisClient(rc *redis.Client) {
	ggc.rc = rc
}

// SetLogger sets the logger.
func (ggc *GoGoCal) SetLogger(l *log.Logger) {
	ggc.l = l
}

// MarkAsFailed marks a keys as having failed to process.
func (ggc *GoGoCal) MarkAsFailed(key string) {
	failKey := NewKey("event", "status", "failed").String()
	ggc.rc.SAdd(failKey, key)
}

// ProcessEvent processes an event.
func (ggc *GoGoCal) ProcessEvent(key string, log chan string, ec chan error) {
	log <- fmt.Sprintf("Processing %s", key)

	// Fetch the event from Redis
	e, err := ggc.rc.HGetAllMap(key).Result()
	if err != nil {
		log <- fmt.Sprintf("Failed to fetch event %s: %s", key, err.Error())
		ec <- err
		return
	}

	// Convert the Json encoded event into a calendar.Event.
	event := new(calendar.Event)
	err = json.Unmarshal([]byte(e["event"]), event)
	if err != nil {
		log <- fmt.Sprintf("Failed to create event %s: %s", key, err.Error())
		ec <- err
		return
	}

	// Send the event to google calendar.
	if event.Id == "" {
		// New event
		event, err = ggc.cs.Events.Insert(e["calendar"], event).Do()
	} else {
		// Update event
		event, err = ggc.cs.Events.Update(e["calendar"], event.Id, event).Do()
	}
	if err != nil {
		log <- fmt.Sprintf(
			"Failed to send event %s to google: %s", key, err.Error(),
		)
		ec <- err
		return
	}

	// Convert the response from google back into a serialized string.
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log <- fmt.Sprintf("Failed to serialize event %s: %s", key, err.Error())
		ec <- err
		return
	}

	// Mark the event as processed
	pipe := ggc.rc.Pipeline()
	defer pipe.Close()

	pipe.HSet(key, "event", string(eventJSON))
	pipe.SAdd(NewKey("event", "status", "processed").String(), key)
	_, err = pipe.Exec()
	if err != nil {
		log <- fmt.Sprintf(
			"Failed to mark event %s as processed: %s", key, err.Error(),
		)
		ec <- err
		return
	}

	log <- fmt.Sprintf("%s processed", key)
	ec <- nil
}

// Run runs the application
func (ggc *GoGoCal) Run() {
	toProcKey := NewKey("event", "status", "to-process").String()

	logs := make(chan string)
	defer close(logs)

	go func() {
		// Log whatever comes in on the logs chanel
		for log := range logs {
			ggc.l.Print(log)
		}
	}()

	for {
		time.Sleep(time.Second)

		// Find out if there are any events to process.
		eventKey, err := ggc.rc.SPop(toProcKey).Result()
		if err != nil {
			// No keys found. Nothing to process.
			continue
		}

		// We have an event key.
		go func(key string) {
			e := make(chan error)
			defer close(e)

			go ggc.ProcessEvent(key, logs, e)

			select {
			case err := <-e:
				if err != nil {
					// Got an error
					ggc.MarkAsFailed(key)
				}
			case <-time.After(2 * time.Minute):
				// Timeout
				ggc.MarkAsFailed(key)
			}
		}(eventKey)
	}
}
