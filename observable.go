package observable

import (
	"sync"
)

type Listener func(interface{})

type Observable struct {
	quit      chan bool
	events    chan interface{}
	listeners []Listener
	mutex     *sync.Mutex
	Verbose   bool
}

func (o *Observable) Open() {
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
func (o *Observable) Close() error {
	// Close event eventLoop
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

func (o *Observable) eventLoop() {
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
func (o *Observable) handleEvent(event interface{}) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	for _, listener := range o.listeners {
		go listener(event)
	}
}

// NewObservable creates a new observable
func New() *Observable {
	obs := Observable{}
	obs.Open()
	return &obs
}
