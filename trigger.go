package htmx

import (
	"encoding/json"
	"strings"
)

type eventContent struct {
	event string
	data  any
}

type Trigger struct {
	triggers   []eventContent
	onlySimple bool
}

// NewTrigger returns a new Trigger set
func NewTrigger() *Trigger {
	return &Trigger{
		triggers:   make([]eventContent, 0),
		onlySimple: true,
	}
}

// Add adds a trigger to the Trigger set
func (t *Trigger) add(trigger eventContent) *Trigger {
	t.triggers = append(t.triggers, trigger)

	return t
}

func (t *Trigger) AddEvent(event string) *Trigger {
	return t.add(eventContent{event: event, data: ""})
}

func (t *Trigger) AddEventDetailed(event, message string) *Trigger {
	t.onlySimple = false

	return t.add(eventContent{event: event, data: message})
}

func (t *Trigger) AddEventObject(event string, details map[string]string) *Trigger {
	t.onlySimple = false

	return t.add(eventContent{event: event, data: details})
}

// String returns the string representation of the Trigger set
func (t *Trigger) String() string {
	if t.onlySimple {
		data := make([]string, len(t.triggers))

		for i, trigger := range t.triggers {
			data[i] = trigger.event
		}

		return strings.Join(data, ", ")
	}

	triggerMap := make(map[string]any)
	for _, tr := range t.triggers {
		triggerMap[tr.event] = tr.data
	}
	data, _ := json.Marshal(triggerMap)
	return string(data)
}
