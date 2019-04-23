# Go-Observable

Highly reduced go observable and event listener based on the work of "github.com/yaacov/observer".

#### Code snippets

``` go
// Create a new observer
o := observable.New()

// Register a listener method.
o.Subscribe(func(e interface{}) {
  log.Println("Hello " e) // Emits 'Hello Tom'
})

// Emit an event.
o.Next("Tom")

// Close the observable
o.Close()
```

## Description

Implements the minimal observer pattern using the golang [channels](https://gobyexample.com/channels). Aims to
mimic the reactive pattern in the future.

## Install

```
go get -u github.com/joernlenoch/go-observable
```

or

```
import "github.com/joernlenoch/go-observable"
```

## Examples

### Emit events

Example of event listener and emitter.

``` go
// Open an observer and start running
o := observable.New()
defer o.Close()

// Add a listener that logs events
o.Subscribe(func(e interface{}) {
  log.Printf("Received: %s.\n", e.(string))
})

// This event will be picked by the listener
o.Next("Hello")
```
