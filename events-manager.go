package core

import (
	"errors"
	"fmt"
)

type EventsManager struct {
	eventsJobsList map[string][]EventJob
	firedEvents    []*Event
}

var manager *EventsManager
var rqc *Context

func NewEventsManager() *EventsManager {
	manager = &EventsManager{
		eventsJobsList: map[string][]EventJob{},
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
	if disableEvents {
		return nil
	}
	if e.Name == "" {
		return errors.New("event name is empty")
	}
	_, exists := m.eventsJobsList[e.Name]
	if !exists {
		return errors.New(fmt.Sprintf("event %v is not registered", e.Name))
	}

	m.firedEvents = append(m.firedEvents, e)
	return nil
}

func (m *EventsManager) Register(eName string, job EventJob) {
	if disableEvents {
		return
	}
	if eName == "" {
		panic("event name is empty")
	}
	_, exists := m.eventsJobsList[eName]
	if !exists {
		m.eventsJobsList[eName] = []EventJob{job}
		return
	}

	for key, jobs := range m.eventsJobsList {
		if key == eName {
			jobs = append(jobs, job)
			m.eventsJobsList[key] = jobs
		}
	}
}

func (m *EventsManager) processFiredEvents() {
	if disableEvents {
		return
	}
	for _, event := range m.firedEvents {
		m.executeEventJobs(event)
	}
	m.firedEvents = []*Event{}
}

func (m *EventsManager) executeEventJobs(event *Event) {
	for key, jobs := range m.eventsJobsList {
		if key == event.Name {
			for _, job := range jobs {
				job(event, rqc)
			}
		}
	}
}
