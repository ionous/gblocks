package gblocks

import (
	r "reflect"
)

// Context - storage for additional per block data
type Context struct {
	ws    *Workspace // owner of the context
	block *Block     // block for which the context was created
	elem  r.Value    // data pointer
}

func (ctx *Context) IsValid() bool {
	return ctx.elem.IsValid()
}

func (ctx *Context) Elem() r.Value {
	if !ctx.elem.IsValid() {
		if x, e := ctx.ws.reg.New(ctx.block.Type); e == nil {
			ctx.elem = x.Elem()
		}
	}
	return ctx.elem
}

func (ctx *Context) Ptr() r.Value {
	return ctx.Elem().Addr()
}
