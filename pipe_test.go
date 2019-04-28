package observable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterMiddleware_Pipe_Filter(t *testing.T) {

	obs := Observable{}
	var called bool

	obs.Pipe(
		Filter(func(i interface{}) bool {
			return i != nil
		}),
	).Subscribe(func(data interface{}) {
		called = data.(bool)
	})

	obs.Next(true)
	assert.True(t, called, "listener was not called")
}

func TestFilterMiddleware_Pipe_FilterNot(t *testing.T) {
	obs := Observable{}

	obs.Pipe(
		Filter(func(i interface{}) bool {
			return i != nil
		}),
	).Subscribe(func(i interface{}) {
		assert.Fail(t, "listener must not be called")
	})

	obs.Next(nil)
}

func TestFirstMiddleware_Pipe(t *testing.T) {
	obs := Observable{}
	var numCalled int

	obs.Pipe(
		First(),
	).Subscribe(func(i interface{}) {
		numCalled++
	})

	obs.Next(true)
	obs.Next(false)

	assert.Equal(t, 1, numCalled, "must unsubscribe after first call")
}

func TestMapMiddleware_Pipe(t *testing.T) {
	obs := Observable{}

	obs.Pipe(
		Map(func(i interface{}) interface{} {
			return true
		}),
	).Subscribe(func(i interface{}) {
		assert.True(t, i.(bool), "must be transformed to 'true'")
	})

	obs.Next(nil)
}

func TestMiddlewares(t *testing.T) {
	obs := Observable{}

	numCalled := 0

	obs.Pipe(
		Filter(func(i interface{}) bool {
			return i == 10
		}),
		Map(func(i interface{}) interface{} {
			num := i.(int)
			return num * 2
		}),
		First(),
	).Subscribe(func(i interface{}) {
		numCalled++
		assert.Equal(t, 20, i)
	})

	obs.Next(1)   // Filtered
	obs.Next(nil) // Filtered
	obs.Next(10)  // OK
	obs.Next(10)  // Unsubscribed
	assert.Equal(t, 1, numCalled, "must unsubscribe after first event")
}

func TestTakeUntilMiddleware_Pipe(t *testing.T) {
	obs := &Observable{}
	numCalled := 0
	trigger := &Observable{}

	obs.Pipe(
		TakeUntil(trigger),
	).Subscribe(func(i interface{}) {
		numCalled++
	})

	obs.Next(nil)
	obs.Next(nil)
	trigger.Next(nil)
	obs.Next(nil)

	assert.Equal(t, 2, numCalled, "must unsubscribe when trigger is called")
}
