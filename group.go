package hchannel

import (
	"sync"
)

type HandleAction func(*Context)

type actionHandler struct {
	id     int32
	handle HandleAction
	owner  IGroup
}

func newAction(id int32, handle HandleAction, owner IGroup) *actionHandler {
	return &actionHandler{id: id, handle: handle, owner: owner}
}

func (a *actionHandler) do(ctx *Context) {
	if a.owner != nil {
		a.owner.doBefore(ctx)
		if !ctx.isAbort() {
			if a.handle != nil {
				a.handle(ctx)
			}
			a.owner.doAfter(ctx)
		}
	}
}

type IGroup interface {
	Group() IGroup
	AddBeforeMiddleWare(handlers ...HandleAction) IGroup
	AddAfterMiddleWare(handlers ...HandleAction) IGroup
	Register(id int32, handler func(ctx *Context)) IGroup

	doBefore(ctx *Context)
	doAfter(ctx *Context)
}

type Group struct {
	root           *group
	lock           sync.RWMutex
	pool           *ContextPool
	actionHandlers map[int32]*actionHandler
}

func NewGroup() *Group {
	g := &Group{actionHandlers: map[int32]*actionHandler{}, pool: NewContextPool()}
	g.root = newGroup(nil, g)
	return g
}

func (g *Group) SetContextPool(pool *ContextPool) *Group {
	g.pool = pool
	return g
}

func (g *Group) Root() IGroup {
	return g.root
}

func (g *Group) Do(pid int32, values Values) {
	ctx := g.pool.Get(pid).SetValues(values)
	defer g.pool.Put(ctx)

	g.do(ctx)
}

func (g *Group) do(ctx *Context) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	if action := g.actionHandlers[ctx.id]; action != nil {
		action.do(ctx)
	}
}

func (g *Group) register(action *actionHandler) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.actionHandlers[action.id] = action
}

type group struct {
	lock   sync.RWMutex
	parent IGroup
	root   *Group

	beforeMiddleWare []HandleAction
	afterMiddleWare  []HandleAction
}

func newGroup(parent IGroup, root *Group) *group {
	return &group{parent: parent, root: root}
}

func (g *group) Group() IGroup {
	return newGroup(g, g.root)
}

func (g *group) Register(id int32, handler func(ctx *Context)) IGroup {
	g.root.register(newAction(id, handler, g))
	return g
}

func (g *group) AddBeforeMiddleWare(handlers ...HandleAction) IGroup {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.beforeMiddleWare == nil {
		g.beforeMiddleWare = make([]HandleAction, 0)
	}
	g.beforeMiddleWare = append(g.beforeMiddleWare, handlers...)
	return g
}

func (g *group) AddAfterMiddleWare(handlers ...HandleAction) IGroup {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.afterMiddleWare == nil {
		g.afterMiddleWare = make([]HandleAction, 0)
	}
	g.afterMiddleWare = append(g.afterMiddleWare, handlers...)
	return g
}

func (g *group) doBefore(ctx *Context) {
	if g.parent != nil {
		g.parent.doBefore(ctx)
	}
	g.doHandlers(ctx, g.beforeMiddleWare)
}

func (g *group) doAfter(ctx *Context) {
	g.doHandlers(ctx, g.afterMiddleWare)
	if g.parent != nil {
		g.parent.doAfter(ctx)
	}
}

func (g *group) doHandlers(ctx *Context, handlers []HandleAction) {
	if ctx.isAbort() {
		return
	}
	for _, handler := range handlers {
		if ctx.isAbort() {
			return
		} else {
			handler(ctx)
		}
	}
}
