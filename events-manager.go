package core

import (
	"errors"
	"fmt"
)

type EventsManager struct {
	eventsJobsList map[string]JobsAndPayloadBag
}

type JobsAndPayloadBag struct {
	eventJobs []EventJob
	payload   map[string]interface{}
}

var manager *EventsManager
var rqc *Context

func NewEventsManager() *EventsManager {
	manager = &EventsManager{
		eventsJobsList: map[string]JobsAndPayloadBag{},
	}

	return manager
}

func ResolveEventsManager() *EventsManager {
	return manager
}

func (m *EventsManager) setContext(requestContext *Context) *EventsManager {
	rqc = requestContext

	return m
}

func (m *EventsManager) Fire(e *Event) error {
	if e.Name == "" {
		return errors.New("event name is empty")
	}
	_, exists := m.eventsJobsList[e.Name]
	if !exists {
		return errors.New(fmt.Sprintf("event %v is not registered", e.Name))
	}

	for eName, bag := range m.eventsJobsList {
		if eName == e.Name {
			bag.payload = e.Payload
		}
	}
	return nil
}

func (m *EventsManager) Register(eName string, job EventJob) {
	if eName == "" {
		panic("event name is empty")
	}
	_, exists := m.eventsJobsList[eName]
	if !exists {
		m.eventsJobsList[eName] = JobsAndPayloadBag{
			eventJobs: []EventJob{job},
			payload:   nil,
		}
		return
	}

	for key, bag := range m.eventsJobsList {
		if key == eName {
			bag.eventJobs = append(bag.eventJobs, job)
			m.eventsJobsList[eName] = bag
		}
	}
}

func (m *EventsManager) executeEventsJobs() {
	for eventName, bag := range m.eventsJobsList {
		processEventsJobsExecution(eventName, bag)
	}
}

func processEventsJobsExecution(eventName string, bag JobsAndPayloadBag) {
	for _, job := range bag.eventJobs {
		event := &Event{
			Name:    eventName,
			Payload: bag.payload,
		}
		job(event, rqc)
	}
}
