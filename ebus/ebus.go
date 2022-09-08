package ebus

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"
)

type contextKey int

const (
	keyPub contextKey = iota
)

// Publisher is the interface that wraps the basic Publish method.
//
// Publish publishes event. It will be distributed to subscribed handler.
type Publisher interface {
	Publish(ctx context.Context, e Event)
}

// Subscriber is the interface that wraps the basic Subscribe method.
//
// Subscribe subscribes for event. Every events will be delivered to the Handler.
type Subscriber interface {
	Subscribe(Handler)
}

// Handler is the interface that wraps th basic Handle method.
//
// Handle will be invoked by pass an Event if the event occurred.
type Handler interface {
	Handle(ctx context.Context, e Event)
}

// HandlerFunc is the function adapter of Handler.
type HandlerFunc func(ctx context.Context, e Event)

// Handle handles the event. It invoke the f(e).
func (f HandlerFunc) Handle(ctx context.Context, e Event) {
	f(ctx, e)
}

// Event holds the event information.
type Event struct {
	Name        string      `json:"eventName"`
	Data        interface{} `json:"data"`
	OccuredTime time.Time   `json:"occuredTime"`
}

// JSONString return the JSON format of the Event structs
func (e Event) JSONString() (res string, err error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return
	}
	res = string(bytes)
	return
}

// Bus is event bus implementation. It center of the event nerve system.
type Bus struct {
	mu       sync.RWMutex
	handlers []Handler
}

// NewEbus will initialize the ebus. If the handler specified, it will create a new Ebus with the specified handler
func NewEbus(handlers ...Handler) *Bus {
	bus := new(Bus)
	if len(handlers) > 0 {
		bus.handlers = handlers
	}
	return bus
}

// Publish publishes the event.
func (b *Bus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	for _, h := range b.handlers {
		h.Handle(ctx, e)
	}
	b.mu.RUnlock()
}

// Subscribe subscribes h to receive event.
func (b *Bus) Subscribe(h Handler) {
	b.mu.Lock()
	b.handlers = append(b.handlers, h)
	b.mu.Unlock()
}

// ContextWithPublisher build new Context from existing parent with Publisher inside.
func ContextWithPublisher(parent context.Context, pub Publisher) context.Context {
	return context.WithValue(parent, keyPub, pub)
}

// PublisherFromContext get Publisher from the ctx.
func PublisherFromContext(ctx context.Context) Publisher {
	pub, ok := ctx.Value(keyPub).(Publisher)
	if !ok {
		return nil
	}

	return pub
}

// NamedEvent creates event using name inferred from the eventData type name.
func NamedEvent(eventData interface{}) Event {
	name := reflect.TypeOf(eventData).Name()
	return Event{
		Name:        name,
		Data:        eventData,
		OccuredTime: time.Now(),
	}
}

// PublishNamedEvent publishes the event with event name inferred from the eventData type name.
func PublishNamedEvent(ctx context.Context, eventData interface{}) {
	event := NamedEvent(eventData)
	Publish(ctx, event)
}

// Publish event.
func Publish(ctx context.Context, event Event) {
	pub := PublisherFromContext(ctx)
	if pub == nil {
		return
	}

	pub.Publish(ctx, event)
}
