package eventbus

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()

	bus := New()
	assert.NotNil(t, bus)
}

func TestSubscribeAndPublish(t *testing.T) {
	t.Parallel()
	bus := New()

	a := false
	b := false
	bus.Subscribe("test1", func(name string, time time.Time, args ...interface{}) {
		a = true
	})
	bus.Subscribe("test1", func(name string, time time.Time, args ...interface{}) {
		b = true
	})
	bus.Subscribe("test2", func(name string, time time.Time, args ...interface{}) {
		a = false
	})
	bus.Subscribe("test3", func(name string, time time.Time, args ...interface{}) {
		b = args[0].(bool)
	})

	bus.Publish("test1")
	bus.WaitAll()

	assert.True(t, a)
	assert.True(t, b)

	bus.Publish("test2")
	bus.WaitAll()

	assert.False(t, a)
	assert.True(t, b)

	bus.Publish("test3", false)
	bus.WaitAll()

	assert.False(t, a)
	assert.False(t, b)

	bus.Publish("test3", true)
	bus.WaitAll()

	assert.False(t, a)
	assert.True(t, b)
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()
	bus := New()

	a := false
	b := false
	c := false
	fna := func(name string, time time.Time, args ...interface{}) {
		a = true
	}
	fnb := func(name string, time time.Time, args ...interface{}) {
		b = true
	}
	fnc := func(name string, time time.Time, args ...interface{}) {
		c = true
	}
	bus.Subscribe("test1", fna)
	bus.Subscribe("test1", fnb)
	bus.Subscribe("test1", fnc)
	bus.Unsubscribe("test1", fnb)
	bus.Unsubscribe("test2", fnc)

	bus.Publish("test1")
	bus.WaitAll()

	assert.True(t, a)
	assert.False(t, b)
	assert.True(t, c)
}

func TestUnsubscribeAll(t *testing.T) {
	t.Parallel()
	bus := New()

	a := false
	b := false
	c := false
	fna := func(name string, time time.Time, args ...interface{}) {
		a = true
	}
	fnb := func(name string, time time.Time, args ...interface{}) {
		b = true
	}
	fnc := func(name string, time time.Time, args ...interface{}) {
		c = true
	}
	bus.Subscribe("test1", fna)
	bus.Subscribe("test1", fnb)
	bus.Subscribe("test1", fnc)
	bus.UnsubscribeAll("test1")

	bus.Publish("test1")
	bus.WaitAll()

	assert.False(t, a)
	assert.False(t, b)
	assert.False(t, c)
}
