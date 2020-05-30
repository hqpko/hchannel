package hchannel

import (
	"sync/atomic"
	"testing"
)

func TestChannel(t *testing.T) {
	total := 0
	count := 16
	c := NewChannel(16, func(i interface{}) {
		total += i.(int)
	}).Run()
	for i := 0; i < count; i++ {
		if ok := c.Input(1); !ok {
			t.Errorf("input fail")
		}
	}
	c.Close()
	if total != count {
		t.Errorf("channel multi error")
	}
}

func TestChannelMulti(t *testing.T) {
	total := int64(0)
	count := int64(16)
	c := NewChannelMulti(16, 4, func(i interface{}) {
		atomic.AddInt64(&total, i.(int64))
	}).Run()
	for i := int64(0); i < count; i++ {
		if ok := c.Input(int64(1)); !ok {
			t.Errorf("input fail")
		}
	}
	c.Close()
	if total != count {
		t.Errorf("channel multi error")
	}
}
