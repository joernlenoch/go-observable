# Go-Observable

Highly reduced go observable system based on the spirit of ReactiveX. 

#### Code snippets

```go
// Create a new observer
o := observable.New()

// Register a listener method.
sub := o.Subscribe(func(e interface{}) {
  log.Println("Hello " e) // Emits 'Hello Tom'
})

// Unsubscribe
sub.Unsubscribe()

// Register a listener method for the whol event.
o.SubscribeEvent(func(evt Event) {
  log.Print("Data", evt.Data())
  evt.Unsubscribe()
})

// Emit an event.
o.Next("Tom")

// Emit an event.
o.Next(ClientDisconnectedEvent {
  IP: "10.2.0.2",
})

// Close the observable
o.Complete()
```

## Description

Implements the minimal observer pattern. Aims to mimic the reactive pattern in the future.

## Roadmap
 - [x] Subscriptions
 - [x] Middlewares
 - [ ] More RX-like support
 - [ ] Full test coverage

## Install

```go
go get -u github.com/joernlenoch/go-observable
```

or

```go
import "github.com/joernlenoch/go-observable"
```

## Examples

### Emit events

Example of event listener and emitter.

```go
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

### Middleware

```go
// Use middlewares
o.Pipe(
  // Only allow value 10
  Filter(func(i interface{}) bool {
    return i == 10
  }),

  // Map the value on 2*x
  Map(func(i interface{}) interface{} {
    num := i.(int)
    return num * 2
  }),

  // Unsubscribe after the first event
  First(),
).SubscribeEvent(func(i interface{}) {
  evt := i.(int)
  log.Print("Recieved", evt)
})
```