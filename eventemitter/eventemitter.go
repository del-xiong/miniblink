package eventemitter

import (
	"reflect"
)

type Response struct {
	// Name of the Event
	EventName string

	// Slice of all the handler's return values
	Ret []interface{}
}

type EventEmitter struct {
	events map[string][]reflect.Value
}

func New() *EventEmitter {
	e := new(EventEmitter)
	e.Init()

	return e
}

// Allocates the EventEmitters memory. Has to be called when
// embedding an EventEmitter in another Type.
func (self *EventEmitter) Init() {
	self.events = make(map[string][]reflect.Value)
}

func (self *EventEmitter) Listeners(event string) []reflect.Value {
	return self.events[event]
}

// Alias to AddListener.
func (self *EventEmitter) On(event string, listener interface{}) {
	self.AddListener(event, listener)
}

// AddListener adds an event listener on the given event name.
func (self *EventEmitter) AddListener(event string, listener interface{}) {
	// Check if the event exists, otherwise initialize the list
	// of handlers for this event.
	if _, exists := self.events[event]; !exists {
		self.events[event] = []reflect.Value{}
	}

	if l, ok := listener.(reflect.Value); ok {
		self.events[event] = append(self.events[event], l)
	} else {
		l := reflect.ValueOf(listener)
		self.events[event] = append(self.events[event], l)
	}
}

// Removes all listeners from the given event.
func (self *EventEmitter) RemoveListeners(event string) {
	delete(self.events, event)
}

// Emits the given event. Puts all arguments following the event name
// into the Event's `Argv` member. Returns a channel if listeners were
// called, nil otherwise.
func (self *EventEmitter) Emit(event string, argv ...interface{}) <-chan *Response {
	listeners, exists := self.events[event]

	if !exists {
		return nil
	}

	var callArgv []reflect.Value
	c := make(chan *Response)

	for _, a := range argv {
		callArgv = append(callArgv, reflect.ValueOf(a))
	}

	for _, listener := range listeners {
		go func(listener reflect.Value) {
			retVals := listener.Call(callArgv)
			var ret []interface{}

			for _, r := range retVals {
				ret = append(ret, r.Interface())
			}

			c <- &Response{EventName: event, Ret: ret}
		}(listener)
	}

	return c
}
