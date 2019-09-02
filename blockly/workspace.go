package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/block"
)

type ToolboxPosition int

const (
	ToolboxAtTop ToolboxPosition = iota
	ToolboxAtBottom
	ToolboxAtLeft
	ToolboxAtRight
)

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
	onDel            []block.OnDelete
}

func (ws *Workspace) Dispose() {
	ws.Call("dispose")
}

func (ws *Workspace) OnDelete(onDel block.OnDelete) {
	if firstTime := len(ws.onDel) == 0; firstTime {
		ws.AddChangeListener(ws.onEvent)
	}
	ws.onDel = append(ws.onDel, onDel)
}

// note: all of the events we get are only for this workspace
// there are no global events, and none of the events carry information about other workspaces
// basically: events alone cannot learn about the mutation workspace(s).
func (ws *Workspace) onEvent(evt interface{}) {
	switch evt := evt.(type) {
	case *BlockDelete:
		// ids is an array of js strings
		for i := 0; i < evt.Ids.Length(); i++ {
			id := evt.Ids.Index(i).String()
			for _, ondel := range ws.onDel {
				ondel.OnDelete(ws.Id, id)
			}
		}
	}
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

func (ws *Workspace) WorkspaceId() string {
	return ws.Id
}
func (ws *Workspace) NewBlock(blockType string) (block.Shape, error) {
	return ws.NewBlockWithId("", blockType)
}

func (ws *Workspace) NewBlockWithId(blockId string, blockType string) (ret block.Shape, err error) {
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
	if obj := ws.Call("newBlock", blockType, blockId); obj.Bool() {
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
