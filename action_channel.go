package hchannel

type ActionChannel struct {
	group             *Group
	mainActionChannel *Channel
}

func NewActionChannel() *ActionChannel {
	return NewActionChannelWithOption(1<<10, 1)
}

func NewActionChannelWithOption(ActionChannelSize, goroutineCount int) *ActionChannel {
	c := &ActionChannel{group: NewGroup().SetContextPool(NewContextPool())}
	c.mainActionChannel = NewChannelMulti(ActionChannelSize, goroutineCount, c.doAction)
	return c
}

func (c *ActionChannel) SetContextPool(pool *ContextPool) *ActionChannel {
	c.group.SetContextPool(pool)
	return c
}

func (c *ActionChannel) Start() *ActionChannel {
	c.mainActionChannel.Run()
	return c
}

func (c *ActionChannel) Root() IGroup {
	return c.group.Root()
}

func (c *ActionChannel) Input(pid int32, values Values) bool {
	return c.mainActionChannel.Input(c.getContext(pid).SetValues(values))
}

func (c *ActionChannel) doAction(i interface{}) {
	if ctx, ok := i.(*Context); ok {
		c.group.do(ctx)
		c.group.pool.Put(ctx)
	}
}

func (c *ActionChannel) getContext(id int32) *Context {
	return c.group.pool.Get(id)
}

func (c *ActionChannel) Stop() {
	c.mainActionChannel.Close()
}
