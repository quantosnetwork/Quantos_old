package events

import (
	"Quantos/events/set"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"sync"
	"time"
)

// Listener is the function type to run on events.
type Listener func(interface{})

type Subscribers map[string]Listener

type Observer struct {
	quit           chan bool
	events         chan *Event
	watcher        *fsnotify.Watcher
	watchPatterns  set.Set
	watchDirs      set.Set
	listeners      []Listener
	mutex          *sync.Mutex
	bufferEvents   []*Event
	bufferDuration time.Duration
	Verbose        bool
}

// Open the observer channles and run the event loop,
// it will return an error if event loop already running.
func (o *Observer) Open() error {
	// Check for mutex
	if o.mutex == nil {
		o.mutex = &sync.Mutex{}
	}

	if o.events != nil {
		return fmt.Errorf("Observer already inititated.")
	}

	// Create the observer channels.
	o.quit = make(chan bool)
	o.events = make(chan *Event)

	// Run the observer.
	return o.eventLoop()
}

// Close the observer channles,
// it will return an error if close fails.
func (o *Observer) Close() error {
	// Close event loop
	if o.events != nil {
		// Send a quit signal.
		o.quit <- true

		// Close channels.
		close(o.quit)
		close(o.events)
	}

	// Close file watcher.
	if o.watcher != nil {
		o.watcher.Close()
	}

	return nil
}

// AddListener adds a listener function to run on event,
// the listener function will recive the event object as argument.
func (o *Observer) AddListener(l Listener) {
	// Check for mutex
	if o.mutex == nil {
		o.mutex = &sync.Mutex{}
	}

	// Lock:
	// 1. operations on array listeners
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.listeners = append(o.listeners, l)
}

// Emit an event, and event can be of any type, when event is triggered all
// listeners will be called using the event object.
func (o *Observer) Emit(event *Event) {
	o.events <- event
}

// eventLoop runs the event loop.
func (o *Observer) eventLoop() error {
	// Run observer.
	go func() {
		for {
			select {
			case event := <-o.events:
				o.handleEvent(event)
			case <-o.quit:
				return
			}
		}
	}()

	return nil
}

func (o *Observer) handleEvent(event *Event) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	// If we do not buffer events, just send this event now.
	if o.bufferDuration == 0 {
		o.sendEvent(event)
		return
	}
	o.bufferEvents = append(o.bufferEvents, event)
	// If this is the first event, set a timeout function.
	if len(o.bufferEvents) == 1 {
		time.AfterFunc(o.bufferDuration, func() {
			// Lock:
			// 1. operations on listeners array (sendEvent).
			// 2. operations on bufferEvents array.
			o.mutex.Lock()
			defer o.mutex.Unlock()

			// Send all events in event buffer.
			for i := range o.bufferEvents {
				o.sendEvent(o.bufferEvents[i])
			}

			// Reset events buffer.
			o.bufferEvents = make([]*Event, 0)
		})
	}
}

var listeners Subscribers

func (o *Observer) sendEvent(event *Event) {
	if event.Subscribable() {
		subscribers := event.subscribers
		for s := range subscribers {
			go listeners[s](event)
		}
	}
}
