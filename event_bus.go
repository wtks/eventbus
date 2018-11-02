package eventbus

import (
	"github.com/gobwas/glob"
	"reflect"
	"sync"
	"time"
)

type Handler func(event string, time time.Time, args ...interface{})

type EventBus struct {
	handlers map[string][]Handler
	compiled map[string]glob.Glob
	cached   map[string][]Handler
	m        sync.RWMutex
	wg       sync.WaitGroup
	cacheM   sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
		compiled: make(map[string]glob.Glob),
		cached:   make(map[string][]Handler),
	}
}

func (bus *EventBus) Subscribe(event string, fn Handler) error {
	bus.m.Lock()
	defer bus.m.Unlock()
	if _, ok := bus.compiled[event]; !ok {
		g, err := glob.Compile(event, ':')
		if err != nil {
			return err
		}
		bus.compiled[event] = g
	}
	bus.handlers[event] = append(bus.handlers[event], fn)
	bus.forgetCache()
	return nil
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
	bus.forgetCache()
	return false
}

func (bus *EventBus) UnsubscribeAll(event string) {
	bus.m.Lock()
	defer bus.m.Unlock()

	bus.handlers[event] = make([]Handler, 0)
	bus.forgetCache()
}

func (bus *EventBus) Publish(event string, args ...interface{}) {
	bus.m.RLock()
	defer bus.m.RUnlock()

	now := time.Now()
	if hs := bus.findHandlers(event); len(hs) > 0 {
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

func (bus *EventBus) findHandlers(topic string) []Handler {
	bus.cacheM.RLock()
	if hs, ok := bus.cached[topic]; ok {
		bus.cacheM.RUnlock()
		return hs
	}
	bus.cacheM.RUnlock()

	var result []Handler
	for k, matcher := range bus.compiled {
		if matcher.Match(topic) {
			result = append(result, bus.handlers[k]...)
		}
	}

	bus.cacheM.Lock()
	bus.cached[topic] = result
	bus.cacheM.Unlock()
	return result
}

func (bus *EventBus) forgetCache() {
	bus.cacheM.Lock()
	bus.cached = map[string][]Handler{}
	bus.cacheM.Unlock()
}
