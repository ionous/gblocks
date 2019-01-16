package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	r "reflect"
)

type ToolboxPosition int

const (
	ToolboxAtTop ToolboxPosition = iota
	ToolboxAtBottom
	ToolboxAtLeft
	ToolboxAtRight
)

type Workspace struct {
	*js.Object
	Id               string          `js:"id"`
	Options          *js.Object      `js:"options"`
	Rtl              bool            `js:"RTL"`
	HorizontalLayout bool            `js:"horizontalLayout"`
	ToolboxPosition  ToolboxPosition `js:"toolboxPosition"`
	Rendered         bool            `js:"rendered"`
	IsClearing       bool            `js:"isClearing"`

	// workspace svg
	IsMutator bool `js:"isMutator"`

	// custom fields
	context map[string]*Context // blockId-> context
}

var TheWorkspace *Workspace

func NewBlankWorkspace() (ret *Workspace) {
	if blockly := js.Global.Get("Blockly"); blockly.Bool() {
		obj := blockly.Get("Workspace").New()
		ws := initWorkspace(obj)
		TheWorkspace = ws
		ret = ws
	}
	return
}

func NewWorkspace(elementId, mediaPath string, tools interface{}) (ret *Workspace) {
	// note: toolbox can be an xml string containing the toolbox
	// grrr --- init calls:
	// - Blockly.VerticalFlyout.Blockly.Flyout.show
	// - Blockly.Events.Create
	// - Object.Blockly.Xml.blockToDom
	// - Blockly.BlockSvg.mutationToDom
	// -> registerBlock's mutationToDom; and TheWorkspace is nil.
	if blockly := js.Global.Get("Blockly"); blockly.Bool() {
		obj := blockly.Call("inject", "blockly", map[string]interface{}{
			"media":   mediaPath,
			"toolbox": tools,
		})
		ws := initWorkspace(obj)
		TheWorkspace = ws
		ret = ws
	}
	return
}

func initWorkspace(obj *js.Object) *Workspace {
	ws := &Workspace{Object: obj}
	if !ws.IsMutator {
		ws.AddChangeListener(ws.mirror)
		ws.context = make(map[string]*Context)
	}
	return ws
}

// GetDataById custom function to get go-lang mirror
func (ws *Workspace) GetDataById(blockId string) (ret interface{}) {
	if ctx := ws.Context(blockId); ctx != nil {
		ret = ctx.Ptr().Interface()
	}
	return
}

// returns pointer to element; returns nil if blockId refers to a mutation block
func (ws *Workspace) Context(blockId string) (ret *Context) {
	if !ws.IsMutator {
		if ctx, ok := ws.context[blockId]; ok {
			ret = ctx
		} else if b := ws.GetBlockById(blockId); b != nil {
			ctx := &Context{block: b}
			ws.context[blockId] = ctx
			ret = ctx
		}
	}
	return
}

func (ws *Workspace) Dispose() {
	ws.Call("dispose")
}

// func (ws* Workspace) addTopBlock (block) {
//  return ws.Call("addTopBlock")
// }
// func (ws* Workspace) removeTopBlock (block) {
//  return ws.Call("removeTopBlock")
// }
// func (ws* Workspace) getTopBlocks (ordered) {
//  return ws.Call("getTopBlocks")
// }
// func (ws* Workspace) addTopComment (comment) {
//  return ws.Call("addTopComment")
// }
// func (ws* Workspace) removeTopComment (comment) {
//  return ws.Call("removeTopComment")
// }
// func (ws* Workspace) getTopComments (ordered) {
//  return ws.Call("getTopComments")
// }
// func (ws* Workspace) getAllBlocks (ordered) {
//  return ws.Call("getAllBlocks")
// }
// func (ws* Workspace) clear () {
//  return ws.Call("clear")
// }
// func (ws* Workspace) renameVariableById (id, newName) {
//  return ws.Call("renameVariableById")
// }
// func (ws* Workspace) createVariable (name, opt_type, opt_id) {
//  return ws.Call("createVariable")
// }
// func (ws* Workspace) getVariableUsesById (id) {
//  return ws.Call("getVariableUsesById")
// }
// func (ws* Workspace) deleteVariableById (id) {
//  return ws.Call("deleteVariableById")
// }
// func (ws* Workspace) deleteVariableInternal_ (variable, uses) {
//  return ws.Call("deleteVariableInternal_")
// }
// func (ws* Workspace) getVariable (name, opt_type) {
//  return ws.Call("getVariable")
// }
// func (ws* Workspace) getVariableById (id) {
//  return ws.Call("getVariableById")
// }
// func (ws* Workspace) getVariablesOfType (type) {
//  return ws.Call("getVariablesOfType")
// }
// func (ws* Workspace) getVariableTypes () {
//  return ws.Call("getVariableTypes")
// }
// func (ws* Workspace) getAllVariables () {
//  return ws.Call("getAllVariables")
// }
// func (ws* Workspace) getWidth () {
//  return ws.Call("getWidth")
// }

// where t is either a TypeName, string, or pointer to type.
func (ws *Workspace) NewBlock(t interface{}) (*Block, error) {
	return ws.NewBlockWithId(t, "")
}

func (ws *Workspace) NewBlockWithId(t interface{}, opt_id string) (ret *Block, err error) {
	var prototypeName TypeName
	switch t := t.(type) {
	case TypeName:
		prototypeName = t
	case string:
		prototypeName = TypeName(t)
	case r.Type:
		prototypeName = toTypeName(t)
	default:
		prototypeName = toTypeName(r.TypeOf(t).Elem())
	}
	// pattern for handling thrown errors
	defer func() {
		if e := recover(); e != nil {
			if e, ok := e.(*js.Error); ok {
				err = e
			} else {
				panic(e)
			}
		}
	}()
	if obj := ws.Call("newBlock", prototypeName, opt_id); obj != nil {
		ret = &Block{Object: obj}
	}
	return
}

// func (ws* Workspace) remainingCapacity () {
//  return ws.Call("remainingCapacity")
// }
// func (ws* Workspace) undo (redo) {
//  return ws.Call("undo")
// }
// func (ws* Workspace) clearUndo () {
//  return ws.Call("clearUndo")
// }
func (ws *Workspace) AddChangeListener(fn func(evt interface{})) *js.Object {
	wrappedFn := js.MakeFunc(func(self *js.Object, args []*js.Object) interface{} {
		fn(decodeEvent(args[0]))
		return nil
	})
	ws.Call("addChangeListener", wrappedFn)
	return wrappedFn
}

func (ws *Workspace) RemoveChangeListener(wrappedFn *js.Object) {
	ws.Call("removeChangeListener", wrappedFn)
}

// func (ws* Workspace) fireChangeListener (event) {
//  return ws.Call("fireChangeListener")
// }

// GetBlockById lookup ( and wrap ) a blockly block for use with go apis.
func (ws *Workspace) GetBlockById(id string) (ret *Block) {
	if obj := ws.Call("getBlockById", id); obj.Bool() {
		ret = &Block{Object: obj}
	}
	return
}

// func (ws* Workspace) getCommentById (id) {
//  return ws.Call("getCommentById")
// }
// func (ws* Workspace) allInputsFilled = function(
// func (ws* Workspace) getPotentialVariableMap () {
//  return ws.Call("getPotentialVariableMap")
// }
// func (ws* Workspace) createPotentialVariableMap () {
//  return ws.Call("createPotentialVariableMap")
// }
// func (ws* Workspace) getVariableMap () {
//  return ws.Call("getVariableMap")
// }
func (ws *Workspace) Clear() {
	ws.Call("clear")
}
func (ws *Workspace) ClearUndo() {
	ws.Call("clearUndo")
}

// func (ws* Workspace) addChangeListener= function() {
// func (ws* Workspace) removeChangeListener= function() {

// listen to changes in the workspace, reflect them into the go-data.
func (ws *Workspace) mirror(evt interface{}) {
	switch evt := evt.(type) {
	case *BlockDelete:
		// ids is an array of js strings
		for i := 0; i < evt.Ids.Length(); i++ {
			key := evt.Ids.Index(i).String()
			delete(ws.context, key)
		}

	case *BlockChange:
		if evt.Element == "field" {
			if ctx := ws.Context(evt.BlockId); ctx != nil {
				name := InputName(evt.Name)
				dst := ctx.FieldForInput(name)

				switch v := evt.NewValue; dst.Kind() {
				case r.Bool:
					var v bool = v.Bool()
					dst.Set(r.ValueOf(v))
				case r.Int, r.Int8, r.Int16, r.Int32:
					var v int = v.Int()
					dst.Set(r.ValueOf(v).Convert(dst.Type()))
				case r.Int64:
					var v int64 = v.Int64()
					dst.Set(r.ValueOf(v))
				case r.Uint, r.Uint8, r.Uint16, r.Uint32:
					var v uint64 = v.Uint64()
					dst.Set(r.ValueOf(v).Convert(dst.Type()))
				case r.Uint64:
					var v uint64 = v.Uint64()
					dst.Set(r.ValueOf(v))
				case r.Float32:
					var v float64 = v.Float()
					dst.Set(r.ValueOf(float32(v)))
				case r.Float64:
					var v float64 = v.Float()
					dst.Set(r.ValueOf(v))
				case r.String:
					var v string = v.String()
					dst.Set(r.ValueOf(v))
				default:
					e := errutil.New("unknown destination in block change", dst.Kind())
					panic(e.Error())
				}
			}
		}

	case *BlockMove:
		if ctx := ws.Context(evt.BlockId); ctx != nil {
			// disconnect the block from the parent; and the parent from the block
			if pid := evt.PrevParentId(); len(pid) > 0 {
				oldParent := ws.Context(pid)

				in := evt.PrevInputName()
				if len(in) == 0 {
					in = NextInput
					// fix up the block's previous input to point to nothing
					if prev := ctx.FieldForInput(PreviousInput); !prev.IsValid() {
						panic(oldParent.String() + " missing previous statement")
					} else {
						prev.Set(r.Zero(prev.Type()))
					}
				}
				// fix up the parent's input to point to nothing
				dst := oldParent.FieldForInput(in)
				dst.Set(r.Zero(dst.Type()))
			}

			// connect the block to the parent; and the parent to the block
			if pid := evt.NextParentId(); len(pid) > 0 {
				newParent := ws.Context(pid)
				in := evt.NextInputName()
				// a blank input means a vertical (next/prev) connection
				if len(in) == 0 {
					in = NextInput
					// fix up the block's previous to point to the parent
					if prev := ctx.FieldForInput(PreviousInput); !prev.IsValid() {
						panic("missing previous statement")
					} else {
						prev.Set(newParent.Ptr())
					}
				}
				// fix up the parent's input to point to this block
				if dst := newParent.FieldForInput(in); !dst.IsValid() {
					panic(newParent.String() + " missing field " + in.String())
				} else {
					valPtr := ctx.Ptr()
					dst.Set(valPtr)
				}
			}
		}
	default:
		// pass
	}
}
