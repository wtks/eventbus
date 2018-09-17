package eventbus

import (
	"reflect"
	"sync"
	"time"
)

type Handler func(event string, time time.Time, args ...interface{})

type EventBus struct {
	handlers map[string][]Handler
	m        sync.RWMutex
	wg       sync.WaitGroup
}

func New() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
		m:        sync.RWMutex{},
		wg:       sync.WaitGroup{},
	}
}

func (bus *EventBus) Subscribe(event string, fn Handler) {
	bus.m.Lock()
	defer bus.m.Unlock()

	bus.handlers[event] = append(bus.handlers[event], fn)
}

func (bus *EventBus) Unsubscribe(event string, fn Handler) bool {
	bus.m.Lock()
	defer bus.m.Unlock()

	target := reflect.ValueOf(fn)
	if hs, ok := bus.handlers[event]; ok {
		for i, h := range hs {
			if reflect.ValueOf(h) == target {
				l := len(bus.handlers[event])
				copy(bus.handlers[event][i:], bus.handlers[event][i+1:])
				bus.handlers[event][l-1] = nil
				bus.handlers[event] = bus.handlers[event][:l-1]
				return true
			}
		}
	}
	return false
}

func (bus *EventBus) UnsubscribeAll(event string) {
	bus.m.Lock()
	defer bus.m.Unlock()

	bus.handlers[event] = make([]Handler, 0)
}

func (bus *EventBus) Publish(event string, args ...interface{}) {
	bus.m.RLock()
	defer bus.m.RUnlock()

	now := time.Now()
	if hs, ok := bus.handlers[event]; ok && len(hs) > 0 {
		for _, h := range hs {
			bus.wg.Add(1)
			go bus.publish(h, event, now, args...)
		}
	}
}

func (bus *EventBus) publish(fn Handler, event string, time time.Time, args ...interface{}) {
	defer bus.wg.Done()
	fn(event, time, args...)
}

func (bus *EventBus) WaitAll() {
	bus.wg.Wait()
}
