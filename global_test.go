package eventbus

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGlobalSubscribeAndPublish(t *testing.T) {
	a := false
	Subscribe("test1", func(name string, time time.Time, args ...interface{}) {
		a = true
	})

	Publish("test1")
	WaitAll()

	assert.True(t, a)
}

func TestGlobalUnsubscribe(t *testing.T) {
	a := false
	fna := func(name string, time time.Time, args ...interface{}) {
		a = true
	}
	Subscribe("test1", fna)
	Unsubscribe("test1", fna)

	Publish("test1")
	WaitAll()

	assert.False(t, a)
}

func TestGlobalUnsubscribeAll(t *testing.T) {
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
	Subscribe("test1", fna)
	Subscribe("test1", fnb)
	Subscribe("test1", fnc)
	UnsubscribeAll("test1")

	Publish("test1")
	WaitAll()

	assert.False(t, a)
	assert.False(t, b)
	assert.False(t, c)
}
