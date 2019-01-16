package gblocks

import (
	r "reflect"
	"strconv"
	"strings"
)

// Context - storage for additional per block data
type Context struct {
	block *Block  // block for which the context was created
	elem  r.Value // data pointer
}

func (ctx *Context) String() string {
	return ctx.block.Type.String() + ":" + ctx.block.Id
}

func (ctx *Context) IsValid() bool {
	return ctx.elem.IsValid()
}

func (ctx *Context) FieldForInput(in InputName) (ret r.Value) {
	// input name format: INPUT_NAME/AtomIndex/INPUT_NAME
	if name := in.FieldPath(); len(name) > 0 {
		subs := strings.Split(name, "/")
		field := ctx.Elem().FieldByName(subs[0])
		if len(subs) == 1 {
			ret = field
		} else if len(subs) == 3 {
			if i, e := strconv.Atoi(subs[1]); e == nil {
				atom := unpack(field.Index(i))
				ret = atom.FieldByName(subs[2])
			}
		}
	}
	return
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
