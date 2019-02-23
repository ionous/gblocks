package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/block"
	r "reflect"
)

type ToolboxPosition int

const (
	ToolboxAtTop ToolboxPosition = iota
	ToolboxAtBottom
	ToolboxAtLeft
	ToolboxAtRight
)

type IdGenerator interface {
	NewId() string
}

// Workspace - a container for Blockly blocks.
// The mutation popups, and the main editing space are examples of separate workspaces.
// ( The toolbox uses the main workspace. )
type Workspace struct {
	*js.Object
	Id               string          `js:"id"`
	Options          *js.Object      `js:"options"`
	Rtl              bool            `js:"RTL"`
	HorizontalLayout bool            `js:"horizontalLayout"`
	ToolboxPosition  ToolboxPosition `js:"toolboxPosition"`
	Rendered         bool            `js:"rendered"`
	IsClearing       bool            `js:"isClearing"`
	IsMutator        bool            `js:"isMutator"` // from workspacesvg

	idGen IdGenerator
}

func NewBlankWorkspace(isMutator bool, optGen IdGenerator) (ret *Workspace) {
	if blockly := getBlockly(); blockly != nil {
		obj := blockly.Get("Workspace").New()
		ret = &Workspace{Object: obj}
		ret.IsMutator = isMutator
		ret.idGen = optGen
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
	// -> registerType's mutationToDom; and TheWorkspace is nil.
	if blockly := getBlockly(); blockly != nil {
		obj := blockly.Call("inject", "blockly", block.Dict{
			"media":   mediaPath,
			"toolbox": tools,
		})
		ret = &Workspace{Object: obj}
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
// func (ws* Workspace) createVariable (name, OptType, OptId) {
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
// func (ws* Workspace) getVariable (name, OptType) {
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

// where t is either a block.Type, string, or pointer to type.
func (ws *Workspace) NewBlock(t interface{}) (*Block, error) {
	var id string
	if ws.idGen != nil {
		id = ws.idGen.NewId()
	}
	return ws.NewBlockWithId(t, id)
}

func (ws *Workspace) NewBlockWithId(t interface{}, OptId string) (ret *Block, err error) {
	var prototypeName block.Type
	switch t := t.(type) {
	case block.Type:
		prototypeName = t
	case string:
		prototypeName = block.Type(t)
	case r.Type:
		prototypeName = block.TypeFromStruct(t.Elem())
	default:
		prototypeName = block.TypeFromStruct(r.TypeOf(t).Elem())
	}
	// pattern for handling thrown errors
	defer func() {
		if e := recover(); e != nil {
			if jserror, ok := e.(*js.Error); ok {
				err = jserror
			} else {
				panic(e)
			}
		}
	}()
	if obj := ws.Call("newBlock", prototypeName, OptId); obj != nil {
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
func (ws *Workspace) GetBlockById(id block.Id) (ret *Block) {
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
