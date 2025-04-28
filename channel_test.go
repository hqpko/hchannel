package hchannel

import (
	"sync/atomic"
	"testing"
	"time"
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

func TestChannelTimer(t *testing.T) {
	total := 0
	count := 16
	c := NewChannelTimer(16, time.Second, func(i interface{}) {
		if _, ok := i.(time.Time); ok {
			total += 1
		} else {
			total += i.(int)
		}
	}).Run()
	for i := 0; i < count; i++ {
		if ok := c.Input(1); !ok {
			t.Errorf("input fail")
		}
	}
	c.Reset(time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	c.Close()
	if total != count+1 {
		t.Errorf("channel multi error")
	}
}

func TestChannelMulti(t *testing.T) {
	total := int64(0)
	count := int64(16)
	c := NewChannelMulti(16, 4, 0, func(i interface{}) {
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
