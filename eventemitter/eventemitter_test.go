package eventemitter

import (
	"fmt"
	"testing"
)

// Struct for testing Embedding of EventEmitters
type Server struct {
	EventEmitter
}

func TestEmbedding(t *testing.T) {
	s := new(Server)

	// Don't forget to allocate the memory when
	// used as sub type.
	s.EventEmitter.Init()

	s.On("recv", func(msg string) string {
		return msg
	})

	resp := <-s.Emit("recv", "Hello World")

	expected := "Hello World"

	if res := resp.Ret[0].(string); res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}
}

func ExampleEmitReturnsEventOnChan() {
	emitter := New()

	emitter.On("hello", func(name string) string {
		return "Hello World " + name
	})

	e := <-emitter.Emit("hello", "John")

	fmt.Println(e.EventName)
	// Output:
	// hello
}

func BenchmarkEmit(b *testing.B) {
	b.StopTimer()
	emitter := New()
	nListeners := 100

	for i := 0; i < nListeners; i++ {
		emitter.On("hello", func(name string) string {
			return "Hello World " + name
		})
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		<-emitter.Emit("hello", "John")
	}
}
