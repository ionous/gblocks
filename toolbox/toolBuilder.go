package toolbox

import (
	r "reflect" // for inspecting go-lang ptrValues

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
func NewBlocks(s Shadowing, ids Ids, events Events) *Builder {
	return &Builder{shadowing: s, gen: domGenerator{ids, events}}
}

func (l *Builder) Blocks() dom.BlockList {
	return dom.BlockList{l.blocks}
}

func (l *Builder) AddTerm(ptr interface{}) *Builder {
	ptrVal := r.ValueOf(ptr)
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
