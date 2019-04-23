package observable

import (
	"sync"
	"time"
)

type Listener func(interface{})

type Observable struct {
	quit           chan bool
	events         chan interface{}
	listeners      []Listener
	mutex          *sync.Mutex
	Verbose        bool
}

func (o *Observable) Open() {
	// Check for mutex
	if o.mutex == nil {
		o.mutex = &sync.Mutex{}
	}

	if o.events != nil {
		return
	}

	// Create the observer channels.
	o.quit = make(chan bool)
	o.events = make(chan interface{})

	// Run the observer.
	o.eventLoop()
}

// Close the observer channles,
// it will return an error if close fails.
func (o *Observable) Close() error {
	// Close event loop
	if o.events != nil {
		// Send a quit signal.
		o.quit <- true

		// Close channels.
		close(o.quit)
		close(o.events)
	}

	return nil
}

func (o *Observable) Subscribe(l Listener) {

	if o.mutex == nil {
		o.mutex = &sync.Mutex{}
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.listeners = append(o.listeners, l)
}

func (o *Observable) Next(event interface{}) {
	o.events <- event
}

func (o *Observable) Emit(event interface{}) {
	o.Next(event)
}

// handleEvent handle an event.
func (o *Observable) handleEvent(event interface{}) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	for _, listener := range o.listeners {
		go listener(event)
	}
}

// eventLoop runs the event loop.
func (o *Observable) eventLoop() {
	// Run observer.
	go func() {
		for {
			select {
			case event := <-o.events:
				o.handleEvent(event, nil)
			case <-o.quit:
				return
			}
		}
	}()
}

// NewObservable creates a new observable
func New() *Observable {
	obs := Observable{}
	obs.Open()
	return &obs
}