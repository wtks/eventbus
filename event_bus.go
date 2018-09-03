package eventbus

import (
	"reflect"
	"sync"
	"time"
)

type Handler func(name string, time time.Time, args ...interface{})

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

func (bus *EventBus) Subscribe(name string, fn Handler) {
	bus.m.Lock()
	defer bus.m.Unlock()

	bus.handlers[name] = append(bus.handlers[name], fn)
}

func (bus *EventBus) Unsubscribe(name string, fn Handler) bool {
	bus.m.Lock()
	defer bus.m.Unlock()

	target := reflect.ValueOf(fn)
	if hs, ok := bus.handlers[name]; ok {
		for i, h := range hs {
			if reflect.ValueOf(h) == target {
				l := len(bus.handlers[name])
				copy(bus.handlers[name][i:], bus.handlers[name][i+1:])
				bus.handlers[name][l-1] = nil
				bus.handlers[name] = bus.handlers[name][:l-1]
				return true
			}
		}
	}
	return false
}

func (bus *EventBus) UnsubscribeAll(name string) {
	bus.m.Lock()
	defer bus.m.Unlock()

	bus.handlers[name] = make([]Handler, 0)
}

func (bus *EventBus) Publish(name string, args ...interface{}) {
	bus.m.RLock()
	defer bus.m.RUnlock()

	now := time.Now()
	if hs, ok := bus.handlers[name]; ok && len(hs) > 0 {
		for _, h := range hs {
			bus.wg.Add(1)
			go bus.publish(h, name, now, args...)
		}
	}
}

func (bus *EventBus) publish(fn Handler, name string, time time.Time, args ...interface{}) {
	defer bus.wg.Done()
	fn(name, time, args...)
}

func (bus *EventBus) WaitAll() {
	bus.wg.Wait()
}
