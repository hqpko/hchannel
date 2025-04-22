package hchannel

import (
	"math"
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
	return NewChannelMulti(chanSize, 1, handler)
}

func NewChannelMulti(chanSize, goroutineCount int, handler func(interface{})) *Channel {
	return &Channel{c: make(chan interface{}, chanSize), gc: goroutineCount, h: handler, t: time.NewTimer(math.MaxInt64)}
}

func (c *Channel) Run() *Channel {
	for i := 0; i < c.gc; i++ {
		go c.run()
	}
	return c
}

func (c *Channel) Reset(d time.Duration) {
	c.stopTimer()
	c.t.Reset(d)
}

func (c *Channel) stopTimer() {
	if !c.t.Stop() {
		select {
		case <-c.t.C:
		default:
		}
	}
}

func (c *Channel) run() {
	defer c.wg.Done()

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
