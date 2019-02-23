package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/block"
)

type EventId string

// Event --
// on the blockly side, events are raised with Blockly.Events.fire.
// there is no node "hierarchy" -- all listeners get all events.
type Event struct {
	*js.Object           // Blockly.Events.Abstract;
	Type        EventId  `js:"type"`
	BlockId     block.Id `js:"blockId"`
	WorkspaceId string   `js:"workspaceId"`
	Group       string   `js:"group"`
	RecordUndo  bool     `js:"recordUndo"`
}

func decodeEvent(obj *js.Object) (ret interface{}) {
	evt := &Event{Object: obj}
	switch evt.Type {
	case "create":
		ret = &BlockCreate{Event: evt}
	case "delete":
		ret = &BlockDelete{Event: evt}
	case "change":
		ret = &BlockChange{Event: evt}
	case "move":
		ret = &BlockMove{Event: evt}
	case "var_create":
		ret = &VarCreate{Event: evt}
	case "var_delete":
		ret = &VarDelete{Event: evt}
	case "var_rename":
		ret = &VarRename{Event: evt}
	case "ui":
		ret = &UiChange{Event: evt}
	case "comment_create":
		ret = &CommentChange{Event: evt}
	case "comment_delete":
		ret = &CommentDelete{Event: evt}
	case "comment_change":
		ret = &CommentChange{Event: evt}
	case "comment_move":
		ret = &CommentMove{Event: evt}
	default:
		ret = evt
	}
	return
}

// toJson
// fromJson(json js.M) {}
// isNull()
// run(foward bool)
// getEventWorkspace_()

// triggered after the block has been added to the workspace; the block object can be found via the event's blockId
type BlockCreate struct {
	*Event
	// xml	*js.Object	An XML tree defining the new block and any connected child blocks.
	Ids *js.Object `js:"ids"` // 	An array containing the UUIDs of the new block and any connected child blocks.
}

type BlockDelete struct {
	*Event
	// oldXml	*js.Object	An XML tree defining the deleted block and any connected child blocks.
	Ids *js.Object `js:"ids"` // An array containing the UUIDs of the deleted block and any connected child blocks.
}

// BlockChange -- when the status of a block has changed.
// ( note, changes to inputs are reported by BlockMove. )
type BlockChange struct {
	*Event
	// Element events include:
	//   collapsed - .collapsed status changed ( via setCollapsed )
	//   comment - comment text changed ( via setCommentText )
	//   disabled - disabled status changed ( via setDisabled )
	//   field - when checkbox, colour, dropdown, text input, variable, etc. ( via setValue )
	//   inline - inputs inline changed ( via setInputsInline )
	//   mutation - "a procedure definition changes its parameters". or workspaceChanged, -- via Blockly.Procedures.mutatateCallers.
	Element string `js:"element"`
	// Name of element affected
	Name     string     `js:"name"`
	OldValue *js.Object `js:"oldValue"`
	NewValue *js.Object `js:"newValue"`
}

// BlockMove -- event when a block has been dragged/dropped into a new slot.
// for why this uses methods instead of properties see https://github.com/gopherjs/gopherjs/issues/617
type BlockMove struct {
	*Event
	// NewCoordinate *js.Object `js:"newCoordinate"` // X and Y coordinates if it is a top level block. Undefined if it has a parent.
	// OldCoordinate *js.Object `js:"oldCoordinate"` // X and Y coordinates if it was a top level block. Undefined if it had a parent.
}

// NextParentId - UUID of new parent block. Empty if it is a top level block.
func (evt *BlockMove) NextParentId() (ret string) {
	if p := evt.Get("newParentId"); p.Bool() {
		ret = p.String()
	}
	return
}

// PrevParentId - UUID of old parent block. Empty if it was a top level block.
func (evt *BlockMove) PrevParentId() (ret string) {
	if p := evt.Get("oldParentId"); p.Bool() {
		ret = p.String()
	}
	return
}

// NextInputName - Input in new parent ( if any ). Empty if it's the parent's next block.
func (evt *BlockMove) NextInputName() (ret block.Item) {
	if p := evt.Get("newInputName"); p.Bool() {
		ret = block.Item(p.String())
	}
	return
}

// PrevInputName - Input in old parent ( if any ). Empty if it's the parent's next block.
func (evt *BlockMove) PrevInputName() (ret block.Item) {
	if p := evt.Get("oldInputName"); p.Bool() {
		ret = block.Item(p.String())
	}
	return
}

type CommentCreate struct {
	*Event
}
type CommentDelete struct {
	*Event
}
type CommentChange struct {
	*Event
}
type CommentMove struct {
	*Event
}

type VarCreate struct {
	*Event
	VarId   string `js:"varId"`
	VarName string `js:"varName"`
	VarType string `js:"varType"`
}

type VarDelete struct {
	*Event
	VarId   string `js:"varId"`
	VarName string `js:"varName"`
	VarType string `js:"varType"`
}

type VarRename struct {
	*Event
	VarId   string `js:"varId"`
	OldName string `js:"oldName"`
	NewName string `js:"newName"`
}

// ex. warningOpen ( showing or hding the warning bubble ), mutatorOpen, commentOpen,  (block) click, (block) selected.
type UiChange struct {
	*Event
	Element  string     `js:"element"`
	OldValue *js.Object `js:"oldValue"`
	NewValue *js.Object `js:"newValue"`
}
