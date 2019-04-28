package observable

type MiddlewareResult struct {
	Event       Event
	Abort       bool
	Unsubscribe bool
}

type Middleware interface {
	Pipe(evt Event) MiddlewareResult
}

// Unsubscribing

type FirstMiddleware struct {
}

func (f *FirstMiddleware) Pipe(evt Event) MiddlewareResult {
	return MiddlewareResult{
		Unsubscribe: true,
	}
}

func First() Middleware {
	return &FirstMiddleware{}
}

// Mapping

type MapFunc func(data interface{}) interface{}

type MapMiddleware struct {
	mapFunc MapFunc
}

func (f *MapMiddleware) Pipe(evt Event) MiddlewareResult {

	return MiddlewareResult{
		Event: evt.withData(f.mapFunc(evt.Data())),
	}
}

func Map(fn MapFunc) Middleware {
	return &MapMiddleware{
		mapFunc: fn,
	}
}

// Filtering

type FilterFunc func(data interface{}) bool

type FilterMiddleware struct {
	FilterFunc FilterFunc
}

func (f *FilterMiddleware) Pipe(evt Event) MiddlewareResult {
	return MiddlewareResult{
		Event: evt,
		Abort: f.FilterFunc(evt.Data()) == false,
	}
}

func Filter(fn FilterFunc) Middleware {
	return &FilterMiddleware{
		FilterFunc: fn,
	}
}

// TakeUntil

type TakeUntilFunc func(data interface{}) bool

type TakeUntilMiddleware struct {
	delete bool
}

func (m *TakeUntilMiddleware) Pipe(evt Event) MiddlewareResult {
	return MiddlewareResult{
		Abort:       m.delete,
		Unsubscribe: m.delete,
	}
}

func TakeUntil(trigger *Observable) Middleware {
	m := &TakeUntilMiddleware{}

	// Subscribe to the first event and unsubscribe from
	// the event as soon as it triggers
	trigger.Pipe(First()).Subscribe(func(i interface{}) {
		m.delete = true
	})

	return m
}
