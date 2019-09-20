package toolbox

import (
	r "reflect" // for inspecting go-lang ptrValues

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
	"github.com/ionous/gblocks/tin"
)

// see also: https://developers.google.com/blockly/guides/configure/web/toolbox
type Builder struct {
	blocks    dom.Blocks // target of xml generation
	shadowing Shadowing
	gen       domGenerator
}

// build a toolbox detecting and (optionally) registering new blocks to the Maker.
func NewBlocks(s Shadowing, events Events, names UniqueNames) *Builder {
	return &Builder{shadowing: s, gen: domGenerator{events, names}}
}

func (l *Builder) Blocks() dom.BlockList {
	return dom.BlockList{l.blocks}
}

func (l *Builder) AddTerm(ptr interface{}) *Builder {
	ptrVal := r.ValueOf(ptr)
	if _, ok := ptrVal.Elem().Type().FieldByName(block.NextStatement); ok {
		e := errutil.New("terms shouldnt have next statements", ptrVal.String())
		panic(e)
	}
	return l.addPtrVal(ptrVal, tin.TermBlock)
}

func (l *Builder) AddStatement(ptr interface{}) *Builder {
	ptrVal := r.ValueOf(ptr)
	return l.addPtrVal(ptrVal, tin.MidBlock)
}

func (l *Builder) AddTopStatement(ptr interface{}) *Builder {
	ptrVal := r.ValueOf(ptr)
	return l.addPtrVal(ptrVal, tin.TopBlock)
}

// add a list of concrete structs;
// signature has to be interface so that lists can be of concrete types.
func (l *Builder) AddBlocks(blocks interface{}) *Builder {
	slice := r.ValueOf(blocks)
	for i, cnt := 0, slice.Len(); i < cnt; i++ {
		v := slice.Index(i)
		var structType r.Type
		var ptrVal r.Value
		if v.Kind() == r.Struct {
			ptrVal, structType = v.Addr(), v.Type()
		} else if v.Kind() == r.Ptr {
			ptrVal, structType = v, v.Type().Elem()
		} else if v.Kind() == r.Interface {
			v := v.Elem()
			ptrVal, structType = v, v.Type().Elem()
		} else {
			e := errutil.New("v is unknown", v.Kind().String())
			panic(e)
		}

		model := tin.TermBlock
		if _, ok := structType.FieldByName(block.NextStatement); ok {
			model = tin.MidBlock
		}
		l.addPtrVal(ptrVal, model)
	}
	return l
}

func (l *Builder) addPtrVal(ptrVal r.Value, model tin.Model) *Builder {
	if l.gen.events != nil {
		if t, e := model.TypeInfo(ptrVal.Type()); e != nil {
			l.gen.events.OnError(e)
		} else {
			l.gen.events.OnBlock(t)
		}
	}
	b := l.gen.genBlock(ptrVal.Elem(), l.shadowing)
	l.blocks = append(l.blocks, b)
	return l
}
