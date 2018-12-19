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

	// custom fields
	reg      *Registry
	dataPtrs ptrMap // stores pointer to values
}

type ptrMap map[string]r.Value

func NewBlankWorkspace(reg *Registry) *Workspace {
	//workspace = new Blockly.Workspace();
	obj := js.Global.Get("Blockly").Get("Workspace").New()
	return initWorkspace(obj, reg)
}

func NewWorkspace(elementId, mediaPath string, reg *Registry) *Workspace {
	div := js.Global.Get("document").Call("getElementById", elementId)
	// note: toolbox can be an xml string containing the toolbox
	obj := js.Global.Get("Blockly").Call("inject", "  blockly", map[string]interface{}{
		"media":   mediaPath,
		"toolbox": div,
	})
	return initWorkspace(obj, reg)
}

func initWorkspace(obj *js.Object, reg *Registry) *Workspace {
	ws := &Workspace{Object: obj}
	ws.AddChangeListener(ws.mirror)
	ws.reg = reg
	ws.dataPtrs = make(ptrMap)
	return ws
}

// GetDataById custom function to get go-lang mirror
func (ws *Workspace) GetDataById(id string) (ret interface{}) {
	if val := ws.dataPointerById(id); val.IsValid() {
		ret = val.Interface()
	}
	return
}

// returns pointer to element
func (ws *Workspace) dataPointerById(id string) (ret r.Value) {
	if val, ok := ws.dataPtrs[id]; ok {
		ret = val
	} else if b := ws.GetBlockById(id); b != nil {
		if t, ok := ws.reg.types[b.Type]; !ok {
			e := errutil.New("unknown type", b.Object)
			panic(e.Error())
		} else {
			val := r.New(t)
			ws.dataPtrs[id] = val
			ret = val
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
func (ws *Workspace) NewBlock(t interface{}) *Block {
	return ws.NewBlockWithId(t, "")
}

func (ws *Workspace) NewBlockWithId(t interface{}, opt_id string) (ret *Block) {
	var prototypeName string
	switch t := t.(type) {
	case r.Type:
		prototypeName = toTypeName(t)
	case string:
		prototypeName = t
	default:
		prototypeName = toTypeName(r.TypeOf(t).Elem())
	}
	if _, ok := ws.reg.types[prototypeName]; ok {
		if obj := ws.Call("newBlock", prototypeName, opt_id); obj != nil {
			ret = &Block{Object: obj}
		}
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
	//println("mirroring", r.TypeOf(evt).Elem().Name(), evt)

	switch evt := evt.(type) {
	case *BlockCreate:
		if b := ws.GetBlockById(evt.BlockId); b == nil {
			panic(evt.BlockId)
		} else if valPtr, e := ws.reg.New(b.Type); e != nil {
			panic(e)
		} else {
			ws.dataPtrs[evt.BlockId] = valPtr
		}

	case *BlockDelete:
		// ids is an array of js strings
		for i := 0; i < evt.Ids.Length(); i++ {
			obj := evt.Ids.Index(i)
			delete(ws.dataPtrs, obj.String())
		}

	case *BlockChange:
		//println("block change", evt.Object)
		if evt.Element == "field" {
			valPtr := ws.dataPointerById(evt.BlockId)
			dst := valPtr.Elem().FieldByName(underscoreToPascal(evt.Name))

			switch v := evt.NewValue; dst.Kind() {
			case r.Bool:
				var v bool = v.Bool()
				dst.Set(r.ValueOf(v))
			case r.Int:
				var v int = v.Int()
				dst.Set(r.ValueOf(v))
			case r.Int8, r.Int16, r.Int32:
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

	case *BlockMove:
		valPtr := ws.dataPointerById(evt.BlockId)

		// disconnect the block from the parent; and the parent from the block
		if pid := evt.OldParentId(); len(pid) > 0 {
			oldParent := ws.dataPointerById(pid)
			in := evt.OldInputName()
			if len(in) != 0 {
				in = underscoreToPascal(in)
			} else {
				in = nextStatement
				// fix up the block's previous input to point to nothing
				if prev := valPtr.Elem().FieldByName(previousStatement); !prev.IsValid() {
					panic("missing previous statement")
				} else {
					prev.Set(r.Zero(prev.Type()))
				}
			}
			// fix up the parent's input to point to nothing
			dst := oldParent.Elem().FieldByName(in)
			dst.Set(r.Zero(dst.Type()))
		}

		// connect the block to the parent; and the parent to the block
		if pid := evt.NewParentId(); len(pid) > 0 {
			newParent := ws.dataPointerById(pid)
			in := evt.NewInputName()
			// a blank input means a vertical (next/prev) connection
			if len(in) != 0 {
				in = underscoreToPascal(in)
			} else {
				in = nextStatement
				// fix up the block's previous to point to the parent
				if prev := valPtr.Elem().FieldByName(previousStatement); !prev.IsValid() {
					panic("missing previous statement")
				} else {
					prev.Set(newParent)
				}
			}
			// fix up the parent's input to point to this block
			if dst := newParent.Elem().FieldByName(in); !dst.IsValid() {
				panic("missing field")
			} else {
				dst.Set(valPtr)
			}
		}
	default:
		// pass
	}
}
