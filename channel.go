package hchannel

import (
	"sync"
	"time"
)

type Channel struct {
	wg sync.WaitGroup
	c  chan interface{}
	gc int // goroutine count
	h  func(interface{})
	t  *time.Timer
}

func NewChannel(chanSize int, handler func(interface{})) *Channel {
	return NewChannelMulti(chanSize, 1, 0, handler)
}

func NewChannelTimer(chanSize int, timerDuration time.Duration, handler func(interface{})) *Channel {
	return NewChannelMulti(chanSize, 1, timerDuration, handler)
}

func NewChannelMulti(chanSize, goroutineCount int, timerDuration time.Duration, handler func(interface{})) *Channel {
	c := &Channel{c: make(chan interface{}, chanSize), gc: goroutineCount, h: handler}
	if timerDuration > 0 {
		c.t = time.NewTimer(timerDuration)
	}
	return c
}

func (c *Channel) Run() *Channel {
	for i := 0; i < c.gc; i++ {
		go c.run()
	}
	return c
}

func (c *Channel) Reset(d time.Duration) {
	if c.t == nil {
		panic("channel timer is nil")
	}
	c.stopTimer()
	c.t.Reset(d)
}

func (c *Channel) stopTimer() {
	if c.t == nil {
		return
	}
	if !c.t.Stop() {
		select {
		case <-c.t.C:
		default:
		}
	}
}

func (c *Channel) run() {
	defer c.wg.Done()
	if c.t == nil {
		for i := range c.c {
			c.h(i)
		}
	} else {
		c.runTimer()
	}
}

func (c *Channel) runTimer() {
	for {
		select {
		case i, ok := <-c.c:
			if ok {
				c.h(i)
			} else {
				return
			}
		case t := <-c.t.C:
			c.h(t)
		}
	}
}

func (c *Channel) Input(i interface{}) bool {
	select {
	case c.c <- i:
		return true
	default:
		return false
	}
}

// MustInput 避免使用 MustInput，在 close 时会引起 panic
func (c *Channel) MustInput(i interface{}) {
	c.c <- i
}

// Close 阻塞等待所有已存入的消息处理完毕
func (c *Channel) Close() {
	c.stopTimer()
	c.wg.Add(c.gc)
	close(c.c)
	c.wg.Wait()
	c.h = nil
	c.t = nil
}
