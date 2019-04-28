package observable

type Subscription interface {
	Unsubscribe()
	subscriptionID() string
}

type subscription struct {
	obs *Observable
	id  string
}

func (s *subscription) subscriptionID() string {
	return s.id
}

func (s *subscription) Unsubscribe() {
	s.obs.unsubscribe(s.id)
}

var _ Subscription = (*subscription)(nil)
