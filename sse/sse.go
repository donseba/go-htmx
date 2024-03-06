package sse

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Client interface {
	ID() string
	Chan() chan Message
}

// Message defines the interface for  messages.
type Message interface {
	String() string // Represent the message as a string for transmission.
}

type TextMessage struct {
	ClientID string
	Event    string
	Data     string
}

func NewMessage(data string) *TextMessage {
	return &TextMessage{
		Data: data,
	}
}

func (m *TextMessage) String() string {
	if m.Event != "" {
		return fmt.Sprintf(`event: %s
data: %s

`, m.Event, m.Data)
	}

	return fmt.Sprintf("data: %s\n\n", m.Data)
}

func (m *TextMessage) WithClientID(id string) Message {
	m.ClientID = id
	return m
}

func (m *TextMessage) WithEvent(event string) Message {
	m.Event = event
	return m
}

type Manager struct {
	clients        sync.Map // Concurrent map for client management
	broadcast      chan Message
	workerPoolSize int
}

// NewManager initializes and returns a new Manager instance.
func NewManager(workerPoolSize int) *Manager {
	manager := &Manager{
		broadcast:      make(chan Message),
		workerPoolSize: workerPoolSize,
	}
	manager.startWorkers()
	manager.startKeepAlive()
	return manager
}

// startWorkers starts worker goroutines for message broadcasting.
func (manager *Manager) startWorkers() {
	for i := 0; i < manager.workerPoolSize; i++ {
		go func() {
			for message := range manager.broadcast {
				manager.clients.Range(func(key, value any) bool {
					client, ok := value.(Client)
					if !ok {
						return true // Continue iteration
					}
					select {
					case client.Chan() <- message:
					default:
						manager.clients.Delete(key) // Remove client if channel is full/closed
					}
					return true // Continue iteration
				})
			}
		}()
	}
}

// startKeepAlive starts a goroutine to send keep-alive messages to connected clients.
func (manager *Manager) startKeepAlive() {
	ticker := time.NewTicker(30 * time.Second) // Send a keep-alive every 30 seconds
	go func() {
		for range ticker.C {
			manager.clients.Range(func(_, value any) bool {
				client, ok := value.(Client)
				if !ok {
					return true // Skip if the type assertion fails
				}
				select {
				case client.Chan() <- NewMessage(":\n"): // Send  comment as keep-alive
				default: // If the channel is blocked, assume the client is disconnected
					manager.clients.Delete(client.ID())
				}
				return true
			})
		}
	}()
}

// Send broadcasts a message to all connected clients.
func (manager *Manager) Send(message Message) {
	manager.broadcast <- message
}

// Handle manages an  connection for a given client.
func (manager *Manager) Handle(w http.ResponseWriter, cl Client) {
	manager.register(cl)
	defer manager.unregister(cl.ID())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for msg := range cl.Chan() {
		_, _ = fmt.Fprint(w, msg.String())
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

// register adds a client to the manager.
func (manager *Manager) register(client Client) {
	manager.clients.Store(client.ID(), client)
}

// unregister removes a client from the manager.
func (manager *Manager) unregister(clientID string) {
	if _, ok := manager.clients.Load(clientID); ok {
		manager.clients.Delete(clientID)
	}
}
