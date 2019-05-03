package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/jsdom"
)

// implements block.Project
type Globals struct {
	Blockly    *Blockly
	Extensions *Extensions
	Xml        *Xml
}

func Project() *Globals {
	p := &theProject
	if !p.init() {
		panic("couldnt initialize blockly")
	}
	return p
}

var theProject Globals

func (p *Globals) init() (okay bool) {
	if p.Blockly == nil {
		if obj := js.Global.Get("Blockly"); obj.Bool() {
			p.Blockly = &Blockly{Object: obj}
			if obj := p.Blockly.Get("Extensions"); obj.Bool() {
				p.Extensions = &Extensions{Object: obj}
			}
			if obj := p.Blockly.Get("Xml"); obj.Bool() {
				p.Xml = &Xml{Object: obj}
			}
			okay = true
		}
	}
	return
}

func (p *Globals) IsBlockRegistered(blockType string) (ret bool) {
	if blocks := p.Blockly.Get("Blocks"); blocks.Bool() {
		ret = blocks.Get(blockType).Bool()
	}
	return
}

func (p *Globals) RegisterBlock(blockType string, desc block.Dict) (err error) {
	p.Blockly.Call("defineBlocksWithJsonArray", []block.Dict{desc})
	return
}

func (p *Globals) RegisterMutator(name string, m block.Mutator) (err error) {
	mixin := Mixin{
		"mutationToDom": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
			main := &Block{Object: obj}
			if dom, e := m.MutationToDom(main); e != nil {
				panic(e)
			} else {
				temp := js.Global.Get("document").Call("createElement", "template")
				temp.Set("innerHTML", dom)
				ret = temp.Get("content").Get("firstChild")
			}
			return
		}),
		"domToMutation": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
			main := &Block{Object: obj}
			xmlElement := &jsdom.Element{Object: parms[0]}
			dom := xmlElement.OuterHTML()
			if e := m.DomToMutation(main, dom); e != nil {
				panic(e)
			}
			return
		}),
		"decompose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
			main := &Block{Object: obj}
			popup := &Workspace{Object: parms[0]}
			if mui, e := m.Decompose(main, popup); e != nil {
				panic(e)
			} else {
				ret = mui.(*Block).Object
			}
			return
		}),
		"compose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
			main := &Block{Object: obj}
			mui := &Block{Object: parms[0]}
			if e := m.Compose(main, mui); e != nil {
				panic(e)
			}
			return
		}),
		"saveConnections": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
			main := &Block{Object: obj}
			mui := &Block{Object: parms[0]}
			if e := m.SaveConnections(main, mui); e != nil {
				panic(e)
			}
			return
		})}
	post := js.MakeFunc(func(this *js.Object, args []*js.Object) (ret interface{}) {
		main := &Block{Object: this}
		if e := m.PostMixin(main); e != nil {
			panic(e)
		}
		return
	})
	p.Extensions.RegisterMutator(name, mixin, post, m.Quarks())
	return
}

// note: toolbox can be an xml string containing the toolbox
func (p *Globals) NewWorkspace(elementId, mediaPath string, tools interface{}) *Workspace {
	// warning:
	// storing the project pointer or using a Globals workspace pointer is difficult;
	// the pointer can only be recorded after the workspace has finished creating
	//
	// inject() calls:
	// - Blockly.VerticalFlyout.Blockly.Flyout.show
	// - Blockly.Events.Create
	// - Object.Blockly.Xml.blockToDom
	// - Blockly.BlockSvg.mutationToDom <- uses the js workspace pointer.
	//
	obj := p.Blockly.Call("inject",
		elementId,
		block.Dict{
			"media":   mediaPath,
			"toolbox": tools,
		})
	return &Workspace{Object: obj}
}

func (p *Globals) NewBlankWorkspace(isMutator bool) (ret *Workspace) {
	obj := p.Blockly.Get("Workspace").New()
	ret = &Workspace{Object: obj}
	ret.IsMutator = isMutator
	return
}
