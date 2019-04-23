package observable

import (
	"sync"
)

type Listener func(interface{})

type Observable interface {
	Subscribe(l Listener)
	Next(event interface{})
	Close()
}

type observable struct {
	quit      chan bool
	events    chan interface{}
	listeners []Listener
	mutex     *sync.Mutex
	Verbose   bool
}

func (o *observable) Open() {
	// Check for mutex
	if o.mutex == nil {
		o.mutex = &sync.Mutex{}
	}

	if o.events != nil {
		return
	}

	o.quit = make(chan bool)
	o.events = make(chan interface{})

	// Start the event eventLoop
	go o.eventLoop()
}

// Close the observer channles,
// it will return an error if close fails.
func (o *observable) Close() {
	// Close event eventLoop
	if o.events != nil {
		// Send a quit signal.
		o.quit <- true

		// Close channels.
		close(o.quit)
		close(o.events)
	}
}

func (o *observable) Subscribe(l Listener) {

	if o.mutex == nil {
		o.mutex = &sync.Mutex{}
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.listeners = append(o.listeners, l)
}

func (o *observable) Next(event interface{}) {
	o.events <- event
}

func (o *observable) eventLoop() {
	for {
		select {
		case event := <-o.events:
			o.handleEvent(event)
		case <-o.quit:
			return
		}
	}
}

// handleEvent sends an event to all listeners
func (o *observable) handleEvent(event interface{}) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	for _, listener := range o.listeners {
		go listener(event)
	}
}

var _ Observable = (*observable)(nil)

// NewObservable creates a new observable
func New() Observable {
	obs := &observable{}
	obs.Open()
	return obs
}
