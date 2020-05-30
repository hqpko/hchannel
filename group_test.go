package hchannel

import (
	"testing"
)

func TestGroup(t *testing.T) {
	count := 0
	handlerAddCount := func(ctx *Context) { count++ }

	group := NewGroup()
	root := group.Root()
	root.AddBeforeMiddleWare(handlerAddCount)
	root.AddAfterMiddleWare(handlerAddCount)
	root.Register(1, handlerAddCount)

	subGroup := root.Group()
	{
		subGroup.AddBeforeMiddleWare(handlerAddCount)
		subGroup.Register(2, handlerAddCount)

		subGroup2 := subGroup.Group()
		{
			subGroup2.AddBeforeMiddleWare(handlerAddCount)
			subGroup2.Register(3, handlerAddCount)
		}
	}

	group.Do(1, nil)
	if count != 3 {
		t.Errorf("action group do fail")
	}

	count = 0
	group.Do(2, nil)
	if count != 4 {
		t.Errorf("action group do fail")
	}

	count = 0
	group.Do(3, nil)
	if count != 5 {
		t.Errorf("action group do fail")
	}
}
