package observable

type Event interface {
	Subscription
	Data() interface{}

	withSubscription(i *subscription) Event
	withData(mapFunc interface{}) Event
}

type event struct {
	*subscription
	data interface{}
}

func (e event) withData(data interface{}) Event {
	e.data = data
	return e
}

func (e event) withSubscription(s *subscription) Event {
	e.subscription = s
	return e
}

func (e event) Data() interface{} {
	return e.data
}

var _ Event = (*event)(nil)
