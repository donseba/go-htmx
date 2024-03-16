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

// add adds a trigger to the Trigger set
func (t *Trigger) add(trigger eventContent) *Trigger {
	t.triggers = append(t.triggers, trigger)

	return t
}

func (t *Trigger) AddEvent(event string) *Trigger {
	return t.add(eventContent{event: event, data: ""})
}

// AddEventDetailed adds a trigger to the Trigger set
func (t *Trigger) AddEventDetailed(event, message string) *Trigger {
	t.onlySimple = false

	return t.add(eventContent{event: event, data: message})
}

// AddEventObject adds a trigger to the Trigger set
func (t *Trigger) AddEventObject(event string, details map[string]any) *Trigger {
	t.onlySimple = false

	return t.add(eventContent{event: event, data: details})
}

func (t *Trigger) AddSuccess(message string, vars ...map[string]any) {
	t.addNotifyObject(notificationSuccess, message, vars...)
}

func (t *Trigger) AddInfo(message string, vars ...map[string]any) {
	t.addNotifyObject(notificationInfo, message, vars...)
}

func (t *Trigger) AddWarning(message string, vars ...map[string]any) {
	t.addNotifyObject(notificationWarning, message, vars...)
}

func (t *Trigger) AddError(message string, vars ...map[string]any) {
	t.addNotifyObject(notificationError, message, vars...)
}

func (t *Trigger) addNotifyObject(nt notificationType, message string, vars ...map[string]any) *Trigger {
	details := map[string]any{
		notificationKeyLevel:   nt,
		notificationKeyMessage: message,
	}

	if len(vars) > 0 {
		for _, m := range vars {
			for k, v := range m {
				if k == notificationKeyLevel || k == notificationKeyMessage {
					k = "_" + k
				}
				details[k] = v
			}
		}
	}

	return t.AddEventObject(DefaultNotificationKey, details)
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

const (
	// notificationSuccess is the success notification type
	notificationSuccess notificationType = "success"
	// notificationInfo is the info notification type
	notificationInfo notificationType = "info"
	// notificationWarning is the warning notification type
	notificationWarning notificationType = "warning"
	// notificationError is the error notification type
	notificationError notificationType = "error"
	// notificationKeyLevel is the notification level key
	notificationKeyLevel = "level"
	// notificationKeyMessage is the notification message key
	notificationKeyMessage = "message"
)

type notificationType string

func (n *notificationType) String() string {
	return string(*n)
}

func (h *Handler) notifyObject(nt notificationType, message string, vars ...map[string]any) {
	t := NewTrigger().addNotifyObject(nt, message, vars...)
	h.TriggerWithObject(t)
}

func (h *Handler) TriggerSuccess(message string, vars ...map[string]any) {
	h.notifyObject(notificationSuccess, message, vars...)
}

func (h *Handler) TriggerInfo(message string, vars ...map[string]any) {
	h.notifyObject(notificationInfo, message, vars...)
}

func (h *Handler) TriggerWarning(message string, vars ...map[string]any) {
	h.notifyObject(notificationWarning, message, vars...)
}

func (h *Handler) TriggerError(message string, vars ...map[string]any) {
	h.notifyObject(notificationError, message, vars...)
}

func (h *Handler) TriggerCustom(custom, message string, vars ...map[string]any) {
	h.notifyObject(notificationType(custom), message, vars...)
}
