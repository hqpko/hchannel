package hchannel

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestActionChannel(t *testing.T) {
	count := int32(0)
	handlerAddCount := func(ctx *Context) { atomic.AddInt32(&count, 1) }

	channel := NewActionChannel()
	rootGroup := channel.Root()
	rootGroup.AddBeforeMiddleWare(handlerAddCount)
	rootGroup.AddAfterMiddleWare(handlerAddCount)
	rootGroup.Register(1, handlerAddCount)

	subGroup := rootGroup.Group()
	{
		subGroup.AddBeforeMiddleWare(handlerAddCount)
		subGroup.Register(2, handlerAddCount)

		subGroup2 := subGroup.Group()
		{
			subGroup2.AddBeforeMiddleWare(handlerAddCount)
			subGroup2.Register(3, handlerAddCount)
		}
	}
	channel.Start()

	channel.Input(1, nil)
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&count) != 3 {
		t.Errorf("action group do fail")
	}

	count = 0
	channel.Input(2, nil)
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&count) != 4 {
		t.Errorf("action group do fail")
	}

	count = 0
	channel.Input(3, nil)
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&count) != 5 {
		t.Errorf("action group do fail")
	}
}
