package observable

import (
	"github.com/labstack/gommon/log"
	"github.com/rs/xid"
	"sync"
)

type EventListener func(Event)
type Listener func(interface{})

type Observable struct {
	middlewares []Middleware
	listeners   map[string]interface{}
	mutex       sync.Mutex
}

func (o *Observable) Subscribe(l Listener) Subscription {
	return o.subscribe(l)
}

func (o *Observable) SubscribeEvent(l EventListener) Subscription {
	return o.subscribe(l)
}

func (o *Observable) subscribe(l interface{}) Subscription {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	id := xid.New().String()

	if o.listeners == nil {
		o.listeners = map[string]interface{}{}
	}

	o.listeners[id] = l

	return &subscription{
		obs: o,
		id:  id,
	}
}

func (o *Observable) Next(evt interface{}) {
	o.next(&event{
		data: evt,
	})
}

func (o *Observable) next(evt Event) {

	o.mutex.Lock()
	defer o.mutex.Unlock()

	for _, middleware := range o.middlewares {
		result := middleware.Pipe(evt)
		if result.Abort {
			return
		}

		if result.Unsubscribe {
			evt.Unsubscribe()
		}

		if result.Event != nil {
			evt = result.Event
		}
	}

	for subID, listener := range o.listeners {
		switch fn := listener.(type) {
		case Listener:
			fn(evt.Data())
		case EventListener:
			fn(evt.withSubscription(&subscription{
				id:  subID,
				obs: o,
			}))
		default:
			log.Errorf("unknown callback type: %v", listener)
		}
	}
}

func (o *Observable) Pipe(middlewares ...Middleware) *Observable {
	next := &Observable{
		middlewares: middlewares,
	}

	// Add the current observable as source
	o.SubscribeEvent(func(evt Event) {
		next.next(evt)
	})

	return next
}

func (o *Observable) unsubscribe(id string) {
	// o.mutex.Lock()
	// defer o.mutex.Unlock()
	delete(o.listeners, id)
}

func (o *Observable) Complete() {
	o.listeners = nil
}
