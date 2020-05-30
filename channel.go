package hchannel

import "sync"

type Channel struct {
	wg sync.WaitGroup
	c  chan interface{}
	gc int // goroutine count
	h  func(interface{})
}

func NewChannel(chanSize int, handler func(interface{})) *Channel {
	return NewChannelMulti(chanSize, 1, handler)
}

func NewChannelMulti(chanSize, goroutineCount int, handler func(interface{})) *Channel {
	return &Channel{c: make(chan interface{}, chanSize), gc: goroutineCount, h: handler}
}

func (c *Channel) Run() *Channel {
	for i := 0; i < c.gc; i++ {
		go c.run()
	}
	return c
}

func (c *Channel) run() {
	for i := range c.c {
		c.h(i)
	}
	c.wg.Done()
}

func (c *Channel) Input(i interface{}) bool {
	select {
	case c.c <- i:
		return true
	default:
		return false
	}
}

// mustInput 避免使用 mustInput，在 close 时会引起 panic
func (c *Channel) mustInput(i interface{}) {
	c.c <- i
}

// Close 阻塞等待所有已存入的消息处理完毕
func (c *Channel) Close() {
	c.wg.Add(c.gc)
	close(c.c)
	c.wg.Wait()
}
