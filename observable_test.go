package observable_test

import (
	"github.com/joernlenoch/go-observable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObservable_Complete(t *testing.T) {
	obs := observable.Observable{}

	obs.Subscribe(func(i interface{}) {
		assert.Fail(t, "must not be called after completion")
	})

	obs.Complete()

	obs.Next(nil)
}
