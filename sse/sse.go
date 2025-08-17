package sse

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type (
	// Listener defines the interface for the receiving end.
	Listener interface {
		ID() string
		Chan() chan Envelope
	}

	// Envelope defines the interface for content that can be broadcast to clients.
	Envelope interface {
		String() string // Represent the envelope contents as a string for transmission.
	}

	// Manager defines the interface for managing clients and broadcasting messages.
	Manager interface {
		Send(message Envelope)
		Handle(w http.ResponseWriter, r *http.Request, cl Listener)
		Clients() []string
	}

	History interface {
		Add(message Envelope) // Add adds a message to the history.
		Send(c Listener)      // Send sends the history to a client.
	}
)

type Client struct {
	id string
	ch chan Envelope
}

func NewClient(id string) Listener {
	return &Client{
		id: id,
		ch: make(chan Envelope, 50),
	}
}

func (c *Client) ID() string          { return c.id }
func (c *Client) Chan() chan Envelope { return c.ch }

// Message represents a simple message implementation.
type Message struct {
	Event string
	Time  time.Time
	Data  string
}

// NewMessage returns a new message instance.
func NewMessage(data string) *Message {
	return &Message{
		Data: data,
		Time: time.Now(),
	}
}

// String returns the message as a string.
func (m *Message) String() string {
	sb := strings.Builder{}

	if m.Event != "" {
		sb.WriteString(fmt.Sprintf("event: %s\n", m.Event))
	}
	sb.WriteString(fmt.Sprintf("data: %v\n\n", m.Data))

	return sb.String()
}

// WithEvent sets the event name for the message.
func (m *Message) WithEvent(event string) Envelope {
	m.Event = event
	return m
}

// broadcastManager manages the clients and broadcasts messages to them.
type broadcastManager struct {
	clients        sync.Map
	broadcast      chan Envelope
	workerPoolSize int
	messageHistory *history
}

// NewManager initializes and returns a new Manager instance.
func NewManager(workerPoolSize int) Manager {
	manager := &broadcastManager{
		broadcast:      make(chan Envelope),
		workerPoolSize: workerPoolSize,
		messageHistory: newHistory(10),
	}

	manager.startWorkers()

	return manager
}

// Send broadcasts a message to all connected clients.
func (manager *broadcastManager) Send(message Envelope) {
	manager.broadcast <- message
}

// Handle sets up a new client and handles the connection.
func (manager *broadcastManager) Handle(w http.ResponseWriter, r *http.Request, cl Listener) {
	manager.register(cl)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Send history to the newly connected client
	manager.messageHistory.Send(cl)

	for {
		select {
		case msg, ok := <-cl.Chan():
			if !ok {
				// If the channel is closed, return from the function
				return
			}
			_, err := fmt.Fprint(w, msg.String())
			if err != nil {
				// If an error occurs (e.g., client has disconnected), return from the function
				return
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

		case <-r.Context().Done():
			manager.unregister(cl.ID())
			close(cl.Chan())
			return
		}
	}
}

// Clients method to list connected client IDs
func (manager *broadcastManager) Clients() []string {
	var clients []string
	manager.clients.Range(func(key, value any) bool {
		id, ok := key.(string)
		if ok {
			clients = append(clients, id)
		}
		return true
	})
	return clients
}

// startWorkers starts worker goroutines for message broadcasting.
func (manager *broadcastManager) startWorkers() {
	for i := 0; i < manager.workerPoolSize; i++ {
		go func() {
			for message := range manager.broadcast {
				manager.clients.Range(func(key, value any) bool {
					client, ok := value.(Listener)
					if !ok {
						return true // Continue iteration
					}
					select {
					case client.Chan() <- message:
						manager.messageHistory.Add(message)
					default:
						// If the client's channel is full, drop the message
					}
					return true // Continue iteration
				})
			}
		}()
	}
}

// register adds a client to the manager.
func (manager *broadcastManager) register(client Listener) {
	manager.clients.Store(client.ID(), client)
}

// unregister removes a client from the manager.
func (manager *broadcastManager) unregister(clientID string) {
	manager.clients.Delete(clientID)
}

type history struct {
	messages []Envelope
	maxSize  int // Maximum number of messages to retain
}

func newHistory(maxSize int) *history {
	return &history{
		messages: []Envelope{},
		maxSize:  maxSize,
	}
}

func (h *history) Add(message Envelope) {
	h.messages = append(h.messages, message)
	// Ensure history does not exceed maxSize
	if len(h.messages) > h.maxSize {
		// Remove the oldest messages to fit the maxSize
		h.messages = h.messages[len(h.messages)-h.maxSize:]
	}
}

func (h *history) Send(c Listener) {
	for _, msg := range h.messages {
		c.Chan() <- msg
	}
}
