package eventbus

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
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
	var c int32 = 0
	var d int32 = 0
	var e int32 = 0
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
	bus.Subscribe("*", func(name string, time time.Time, args ...interface{}) {
		atomic.AddInt32(&c, 1)
	})
	bus.Subscribe("test:*", func(event string, time time.Time, args ...interface{}) {
		atomic.AddInt32(&d, 1)
	})
	bus.Subscribe("test:**", func(event string, time time.Time, args ...interface{}) {
		atomic.AddInt32(&e, 1)
	})

	bus.Publish("test1")
	bus.WaitAll()

	assert.True(t, a)
	assert.True(t, b)
	assert.EqualValues(t, 1, c)

	bus.Publish("test2")
	bus.WaitAll()

	assert.False(t, a)
	assert.True(t, b)
	assert.EqualValues(t, 2, c)

	bus.Publish("test3", false)
	bus.WaitAll()

	assert.False(t, a)
	assert.False(t, b)
	assert.EqualValues(t, 3, c)

	bus.Publish("test3", true)
	bus.WaitAll()

	assert.False(t, a)
	assert.True(t, b)
	assert.EqualValues(t, 4, c)

	bus.Publish("test:created")
	bus.Publish("test:deleted")
	bus.Publish("test:updated")
	bus.Publish("test:a:b")
	bus.WaitAll()

	assert.EqualValues(t, 3, d)
	assert.EqualValues(t, 4, e)
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
