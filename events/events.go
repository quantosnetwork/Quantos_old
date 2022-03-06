package events

import (
	"github.com/google/uuid"
	"github.com/minio/blake2b-simd"
	"go.uber.org/zap"
	"lukechampine.com/frand"
)

type Event struct {
	EventID      string `json:"event_id"`
	Name         string `json:"event_name"`
	Handler      EventHandler
	subscribable bool
	observable   bool
	once         bool
	hasState     bool
	subscribers  map[string]bool
}

func (e Event) Observe(event *Event) {
	//TODO implement me
	panic("implement me")
}

func (e Event) Dispatch() {
	//TODO implement me
	panic("implement me")
}

func (e Event) Subscribe(subscriber string) {
	//TODO implement me
	panic("implement me")
}

func (e Event) Unsubscribe(subscriber string) {
	//TODO implement me
	panic("implement me")
}

func (e Event) Log() *zap.Logger {
	z := zap.L()
	return z
}

func (e Event) Observable() bool {
	return e.observable
}

func (e Event) Subscribable() bool {
	return e.subscribable
}

func (e Event) Once() bool {
	return e.once
}

func (e Event) HasState() bool {
	return e.hasState
}

type EventHandler = func(eventName string) bool

type Events struct {
	List map[*Event]bool
}

type EventInterface interface {
	Create(name string, handler EventHandler, canSubscribe, canObserve, onlyOnce, hasState bool) *Event
	Observe(*Event)
	Dispatch()
	Subscribe(subscriber string)
	Unsubscribe(subscriber string)
	Log() *zap.Logger
	Observable() bool
	Subscribable() bool
	Once() bool
	HasState() bool
}

var EventsList Events

func New(name string, handler EventHandler, canSubscribe, canObserve, onlyOnce, hasState bool) EventInterface {
	var newEvent Event
	e := newEvent.Create(name, handler, canSubscribe, canObserve, onlyOnce, hasState)
	EventsList.List[e] = true
	return e

}

func (e Event) Create(name string, handler EventHandler, canSubscribe, canObserve, onlyOnce, hasState bool) *Event {
	event := &Event{}
	buf := make([]byte, 32)
	frand.Read(buf)
	h := blake2b.New256()
	space, _ := uuid.NewUUID()
	data := []byte("events_" + name)
	version := 4

	id := uuid.NewHash(h, space, data, version)
	event.EventID = id.String()
	event.Name = name
	event.Handler = handler
	event.subscribable = canSubscribe
	event.observable = canObserve
	event.once = onlyOnce
	event.hasState = hasState
	return event
}
