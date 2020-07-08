# An EventEmitter for Go

**Build Status:** ![Build Status by goci.me](http://goci.me/project/image/github.com/CHH/eventemitter)

## Install

With `go get`:

    % go get github.com/CHH/eventemitter

## Usage

_For more information please also see the [Package Docs](http://go.pkgdoc.org/github.com/CHH/eventemitter)._

A new EventEmitter is created by the `New` function.

```go
import "github.com/CHH/eventemitter"

func main() {
    emitter := eventemitter.New()
}
```

A listener is of type `func (event *ee.Event)`
Listeners can be bound to event names with the `On` method:

```go
emitter.On("foo", func(name string) {
    fmt.Printf("Hello World %s", name)
})
```

An event can be triggered by calling the `Emit` method:

```go
<- emitter.Emit("foo", "John")
```

When `Emit` is called, each registered listener is called in
its own Goroutine. They all share a common channel, which is
returned by the `Emit` function.

```go
var c chan interface{} 

c = emitter.Emit("foo", "John")
```

This channel can be used to wait until all listeners have finished, by using the
`<-` operator without variable:

```go
<- emitter.Emit("foo", "John")
```

Each listener yields a reference to `eventemitter.Response` on the
channel once it has finished:

```go
c := emitter.Emit("foo", "John")

for resp := <- c {
    // Do something
}
```

### Embedding

EventEmitters can also be embedded in other types. When embedding you've
to call the `Init` function on the EventEmitter, so the memory is
correctly allocated:

```go
type Server struct {
    eventemitter.EventEmitter
}

func NewServer() *Server {
    s := new(Server)

    // Allocates the EventEmitter's memory.
    s.EventEmitter.Init()

    // All functions of the EventEmitter are available:
    s.On("foo", func() {
        
    })
}
```

## License

EventEmitter is distributed under the Terms of the MIT License. See
the bundled file `LICENSE.txt` for more information.

