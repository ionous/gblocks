package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type TypeName string

func (n TypeName) StructName() string {
	return underscoreToPascal(n.String())
}

func (n TypeName) String() (ret string) {
	if len(n) > 0 {
		ret = string(n)
	}
	return
}

type Block struct {
	*js.Object                    // Blockly.Block
	Id                 string     `js:"id"`
	Type               TypeName   `js:"type"`
	outputConnection   *js.Object `js:"outputConnection"`
	nextConnection     *js.Object `js:"nextConnection"`
	previousConnection *js.Object `js:"previousConnection"`

	inputList    *js.Object `js:"inputList"`
	Disabled     bool       `js:"disabled"`
	Tooltip      string     `js:"tooltip"`
	ContextMenu  bool       `js:"contextMenu"`
	Comment      string     `js:"comment"`
	IsInFlyout   bool       `js:"isInFlyout"`
	IsInMutator  bool       `js:"isInMutator"`
	Rtl          bool       `js:"RTL"`
	InputsInline bool       `js:"inputsInline"`

	// note: this workspace pointer has limited value;
	// it doesnt point to the Workspace object containing "Context" data
	workspace_ *js.Object `js:"workspace"`
}

func jsConnection(obj *js.Object) (ret *Connection) {
	if obj != nil && obj.Bool() {
		ret = &Connection{Object: obj}
	}
	return
}

// left side puzzle connector
func (b *Block) OutputConnection() *Connection {
	return jsConnection(b.outputConnection)
}

// connection to a piece in the following line
func (b *Block) NextConnection() *Connection {
	return jsConnection(b.nextConnection)
}

// connection to a piece in the preceeding line
func (b *Block) PreviousConnection() *Connection {
	return jsConnection(b.previousConnection)
}

// feels like this should have been asynchronous, hidden
// so that its called after initialiation automatically rather than the caller/creator deciding when its a good time. [ it gets littered after newBlock() and is illegal in "headless" builds ]
func (b *Block) InitSvg() {
	// missing in headless version
	if b.Object.Get("initSvg").Bool() {
		b.Call("initSvg")
	}
}

// Dispose removes the block from the workspace.
// To prevent child blocks from *also* being disposed, Unplug() the block first.
func (b *Block) Dispose() {
	b.Call("dispose")
}

//func (b* Block)initModel  ()  { b.Call("initModel") }
func (b *Block) Unplug(healStack bool) {
	b.Call("unplug", healStack)
}

//func (b* Block)lastConnectionInStack  ()  { b.Call("lastConnectionInStack") }
func (b *Block) GetParent() (ret *Block) {
	if obj := b.Call("getParent"); obj.Bool() {
		ret = &Block{Object: obj}
	}
	return
}

//func (b* Block)getInputWithBlock  (block)  { b.Call("getInputWithBlock") }
//func (b* Block)getSurroundParent  ()  { b.Call("getSurroundParent") }

func (b *Block) GetNextBlock() (ret *Block) {
	if next := b.NextConnection(); next != nil {
		ret = next.TargetBlock()
	}
	return
}

//func (b* Block)getPreviousBlock  ()  { b.Call("getPreviousBlock") }

// Return the connection on the first statement input
func (b *Block) GetFirstStatementConnection() (ret *Connection) {
	if obj := b.Call("getFirstStatementConnection"); obj.Bool() {
		ret = &Connection{Object: obj}
	}
	return
}

func (b *Block) GetRootBlock() (ret *Block) {
	if obj := b.Call("getRootBlock"); obj.Bool() {
		ret = &Block{Object: obj}
	}
	return
}

//func (b* Block)getChildren  (ordered)  { b.Call("getChildren") }
//func (b* Block)setParent  (newParent *Block)  { b.Call("setParent") }
//func (b* Block)getDescendants  (ordered)  { b.Call("getDescendants") }
func (b *Block) IsDeletable() bool {
	return b.Call("isDeletable").Bool()
}

//func (b* Block)setDeletable  (deletable)  { b.Call("setDeletable") }
func (b *Block) IsMovable() bool {
	return b.Call("isMovable").Bool()
}

//func (b* Block)setMovable  (movable)  { b.Call("setMovable") }
func (b *Block) IsShadow() bool {
	return b.Call("isShadow").Bool()
}

//func (b* Block)setShadow  (shadow)  { b.Call("setShadow") }
func (b *Block) IsInsertionMarker() bool {
	return b.Call("isInsertionMarker").Bool()
}

//func (b* Block)setInsertionMarker  (insertionMarker)  { b.Call("setInsertionMarker") }
func (b *Block) IsEditable() bool {
	return b.Call("isEditable").Bool()
}

//func (b* Block)setEditable  (editable)  { b.Call("setEditable") }
//func (b* Block)setConnectionsHidden  (hidden)  { b.Call("setConnectionsHidden") }
//func (b* Block)getMatchingConnection  (otherBlock, conn)  { b.Call("getMatchingConnection") }
func (b *Block) SetHelpUrl(url string) {
	//FIX -- var localizedValue = Blockly.utils.replaceMessageReferences(rawValue);
	b.Call("setHelpUrl", url)
}

func (b *Block) SetTooltip(text string) {
	b.Call("setTooltip", text)
}

// GetColour of the block as an #RRGGBB string.
func (b *Block) GetColour() string {
	return b.Call("getColour").String()
}

// GetHue as 0-360 HSV value
func (b *Block) getHue() int {
	return b.Call("getHue").Int()
}

func (b *Block) SetColour(colour string) {
	b.Call("setColour", colour)
}

//func (b* Block)setOnChange  (onchangeFn)  { b.Call("setOnChange") }
func (b *Block) GetField(name string) (ret *Field) {
	if obj := b.Call("getField", name); obj.Bool() {
		ret = &Field{Object: obj}
	}
	return
}

//func (b* Block)getVars  ()  { b.Call("getVars") }
//func (b* Block)getVarModels  ()  { b.Call("getVarModels") }
//func (b* Block)updateVarName  (variable)  { b.Call("updateVarName") }
//func (b* Block)renameVarById  (oldId, newId)  { b.Call("renameVarById") }

//func (b* Block)getFieldValue  (name)  { b.Call("getFieldValue") }
//func (b* Block)setFieldValue  (newValue, name)  { b.Call("setFieldValue") }

//func (b* Block)setPreviousStatement  (newBoolean, opt_check)  { b.Call("setPreviousStatement") }
//func (b* Block)setNextStatement  (newBoolean, opt_check)  { b.Call("setNextStatement") }
//func (b* Block)setOutput  (newBoolean, opt_check)  { b.Call("setOutput") }
func (b *Block) SetInputsInline(yes bool) (err error) {
	b.Call("setInputsInline", yes)
	return
}

func (b *Block) GetInputsInline() bool {
	return b.Call("getInputsInline").Bool()
}

//func (b* Block)setDisabled  (disabled)  { b.Call("setDisabled") }
//func (b* Block)getInheritedDisabled  ()  { b.Call("getInheritedDisabled") }
func (b *Block) IsCollapsed() bool {
	return b.Call("isCollapsed").Bool()
}

//func (b* Block)setCollapsed  (collapsed)  { b.Call("setCollapsed") }
//func (b* Block)toString  (opt_maxLength, opt_emptyToken)  { b.Call("toString") }

// AppendValueInput for blocks with output.
func (b *Block) AppendValueInput(name InputName) (ret *Input) {
	return b.appendInput(InputValue, name)
}

// AppendStatementInput for blocks with previous statements.
// statements give a c-shape; they are slices
func (b *Block) AppendStatementInput(name InputName) (ret *Input) {
	return b.appendInput(NextStatement, name)
}

// AppendDummyInput for standalone fields.
func (b *Block) AppendDummyInput(name InputName) (ret *Input) {
	return b.appendInput(NextStatement, name)
}

func (b *Block) appendInput(inputType InputType, name InputName) (ret *Input) {
	newInput := b.Call("appendInput_", inputType, name)
	return &Input{Object: newInput}
}

func (b *Block) JsonInit(opt Options) (err error) {
	b.Call("jsonInit", opt)
	return
}

//func (b* Block)mixin  (mixinObj, opt_disableCheck)  { b.Call("mixin") }
func (b *Block) interpolate(msg string, args []Options) {
	b.Call("interpolate_", msg, args)
}

func (b *Block) hasWorkspace() bool {
	return b.workspace_ != nil && b.workspace_ != js.Undefined
}

//func (b* Block)moveInputBefore  (name, refName)  { b.Call("moveInputBefore") }
//func (b* Block)moveNumberedIxpnputBefore  (

func (b *Block) RemoveInput(name InputName) {
	noExceptionIfMissing := false // default in blockly raises exception
	b.Call("removeInput", name, noExceptionIfMissing)
}

func (b *Block) NumInputs() int {
	return b.inputList.Length()
}

func (b *Block) Input(i int) *Input {
	if cnt := b.inputList.Length(); i < 0 || i >= cnt {
		println(i, "of", cnt)
		panic("out of range")
	}
	in := b.inputList.Index(i)
	return &Input{Object: in}
}

func (b *Block) setInput(i int, in *Input) {
	b.inputList.SetIndex(i, in.Object)
}

func (b *Block) InputByName(str InputName) (retInput *Input, retIndex int) {
	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
		if in := b.Input(i); in.Name == str {
			retInput, retIndex = in, i
			break
		}
	}
	return
}

func (b *Block) MutationType() TypeName {
	return b.Type + "$mutation"
}

//func (b* Block)getInputTargetBlock  (name)  { b.Call("getInputTargetBlock") }
//func (b* Block)getCommentText  ()  { b.Call("getCommentText") }
//func (b* Block)setCommentText  (text)  { b.Call("setCommentText") }
//func (b* Block)setWarningText  (_text, _opt_id)  { b.Call("setWarningText") }
//func (b* Block)setMutator  (_mutator)  { b.Call("setMutator") }
//func (b* Block)getRelativeToSurfaceXY  ()  { b.Call("getRelativeToSurfaceXY") }
//func (b* Block)moveBy  (dx, dy)  { b.Call("moveBy") }
//func (b* Block)allInputsFilled  (opt_shadowBlocksAreFilled)  { b.Call("allInputsFilled") }
//func (b* Block)toDevString  ()  { b.Call("toDevString") }
