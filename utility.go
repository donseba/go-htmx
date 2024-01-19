package htmx

const (
	// notificationSuccess is the success notification type
	notificationSuccess notificationType = "success"
	// notificationInfo is the info notification type
	notificationInfo notificationType = "info"
	// notificationWarning is the warning notification type
	notificationWarning notificationType = "warning"
	// notificationError is the error notification type
	notificationError notificationType = "error"

	notificationKeyLevel   = "level"
	notificationKeyMessage = "message"
)

type notificationType string

func (n *notificationType) String() string {
	return string(*n)
}

func (h *Handler) notifyObject(nt notificationType, message string, vars ...map[string]string) {
	details := map[string]string{
		notificationKeyLevel:   nt.String(),
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

	t := NewTrigger().AddEventObject(DefaultNotificationKey, details)

	h.TriggerWithObject(t)
}

func (h *Handler) TriggerSuccess(message string, vars ...map[string]string) {
	h.notifyObject(notificationSuccess, message, vars...)
}

func (h *Handler) TriggerInfo(message string, vars ...map[string]string) {
	h.notifyObject(notificationInfo, message, vars...)
}

func (h *Handler) TriggerWarning(message string, vars ...map[string]string) {
	h.notifyObject(notificationWarning, message, vars...)
}

func (h *Handler) TriggerError(message string, vars ...map[string]string) {
	h.notifyObject(notificationError, message, vars...)
}

func (h *Handler) TriggerCustom(custom, message string, vars ...map[string]string) {
	h.notifyObject(notificationType(custom), message, vars...)
}
