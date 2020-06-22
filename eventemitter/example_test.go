package eventemitter

import (
	"fmt"
)

func ExampleEventEmitter() {
	// Construct a new EventEmitter instance
	emitter := New()

	emitter.On("hello", func() {
		fmt.Println("Hello World")
	})

	emitter.On("hello", func() {
		fmt.Println("Hello Hello World")
	})

	// Wait until all handlers have finished
	<-emitter.Emit("hello")
	// Output:
	// Hello World
	// Hello Hello World
}

func ExampleEventEmitter_Emit() {
	emitter := New()

	emitter.On("hello", func(name string) {
		fmt.Printf("Hello World %s\n", name)
	})

	<-emitter.Emit("hello", "John")
	// Output:
	// Hello World John
}

