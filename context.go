package gblocks

import (
	r "reflect"
)

// Context - storage for additional per block data
type Context struct {
	block *Block  // block for which the context was created
	elem  r.Value // data pointer
}

func (ctx *Context) IsValid() bool {
	return ctx.elem.IsValid()
}

func (ctx *Context) Elem() r.Value {
	if !ctx.elem.IsValid() {
		if x, e := TheRegistry.NewData(ctx.block.Type); e == nil {
			ctx.elem = x.Elem()
		} else {
			panic(e)
		}
	}
	return ctx.elem
}

func (ctx *Context) Ptr() r.Value {
	return ctx.Elem().Addr()
}
