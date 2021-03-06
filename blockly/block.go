package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/block"
)

type Block struct {
	*js.Object                    // Blockly.Block
	Id                 string     `js:"id"`
	Type               string     `js:"type"`
	outputConnection   *js.Object `js:"outputConnection"`
	nextConnection     *js.Object `js:"nextConnection"`
	previousConnection *js.Object `js:"previousConnection"`

	inputList    *js.Object `js:"inputList"`
	Tooltip      string     `js:"tooltip"`
	ContextMenu  bool       `js:"contextMenu"`
	Comment      string     `js:"comment"`
	IsInFlyout   bool       `js:"isInFlyout"`
	IsInMutator  bool       `js:"isInMutator"`
	Rtl          bool       `js:"RTL"`
	InputsInline bool       `js:"inputsInline"`
	workspace    *js.Object `js:"workspace"`
}

func jsConnection(obj *js.Object) (ret block.Connection) {
	if obj != nil && obj.Bool() {
		ret = &Connection{Object: obj}
	}
	return
}

func (b *Block) BlockId() string {
	return b.Id
}

func (b *Block) BlockType() string {
	return b.Type
}

// left side puzzle connector
func (b *Block) OutputConnection() block.Connection {
	return jsConnection(b.outputConnection)
}

// connection to a piece in the following line
func (b *Block) NextConnection() block.Connection {
	return jsConnection(b.nextConnection)
}

// connection to a piece in the preceeding line
func (b *Block) PreviousConnection() block.Connection {
	return jsConnection(b.previousConnection)
}

// feels like this should have been asynchronous, hidden
// so that its called after initialiation automatically rather than the caller/creator deciding when its a good time. [ it gets littered after NewBlock() and is illegal in "headless" builds ]
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

// func (b *Block) GetNextBlock() (ret *Block) {
// 	if next := b.NextConnection(); next != nil {
// 		ret = next.TargetBlock()
// 	}
// 	return
// }

//func (b* Block)getPreviousBlock  ()  { b.Call("getPreviousBlock") }

// Return the connection on the first statement input
// func (b *Block) GetFirstStatementConnection() (ret block.Connection) {
// 	return jsConnection(b.Call("getFirstStatementConnection"))
// }

func (b *Block) GetRootBlock() (ret *Block) {
	if obj := b.Call("getRootBlock"); obj.Bool() {
		ret = &Block{Object: obj}
	}
	return
}

//func (b* Block)getChildren  (ordered)  { b.Call("getChildren") }
//func (b* Block)setParent  (newParent *Block)  { b.Call("setParent") }
//func (b* Block)getDescendants  (ordered)  { b.Call("getDescendants") }

func (b *Block) GetFlag(flag block.Flag) bool {
	var name string
	switch flag {
	case block.InputsInline:
		name = "getInputsInline"
	default:
		name = "is" + flag.String()
	}
	return b.Call(name).Bool()
}

func (b *Block) SetFlag(flag block.Flag, state bool) {
	var name string
	switch flag {
	case block.Enabled:
		name, state = "setDisabled", !state
	default:
		name = "set" + flag.String()
	}
	b.Call(name, state)
}

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

//func (b* Block)setPreviousStatement  (newBoolean, optCheck)  { b.Call("setPreviousStatement") }
//func (b* Block)setNextStatement  (newBoolean, optCheck)  { b.Call("setNextStatement") }
//func (b* Block)setOutput  (newBoolean, optCheck)  { b.Call("setOutput") }
func (b *Block) SetInputsInline(yes bool) (err error) {
	b.Call("setInputsInline", yes)
	return
}

func (b *Block) GetInputsInline() bool {
	return b.Call("getInputsInline").Bool()
}

//func (b* Block)getInheritedDisabled  ()  { b.Call("getInheritedDisabled") }

//func (b* Block)toString  (optMaxLength, optEmptyToken)  { b.Call("toString") }

// AppendValueInput for blocks with output.
func (b *Block) AppendValueInput(name string) (ret block.Input) {
	return b.appendInput(InputValueType, name)
}

// AppendStatementInput for blocks with previous statements.
// statements give a c-shape; they are slices
func (b *Block) AppendStatementInput(name string) (ret block.Input) {
	return b.appendInput(NextStatementType, name)
}

// AppendDummyInput for standalone fields.
func (b *Block) AppendDummyInput(name string) (ret block.Input) {
	return b.appendInput(DummyInputType, name)
}

func (b *Block) appendInput(inputType InputType, name string) (ret block.Input) {
	newInput := b.Call("appendInput_", inputType, name)
	return &Input{Object: newInput}
}

func (b *Block) JsonInit(opt block.Dict) {
	b.Call("jsonInit", opt)
}

//func (b* Block)mixin  (mixinObj, optDisableCheck)  { b.Call("mixin") }
func (b *Block) Interpolate(msg string, args []block.Dict) {
	b.Call("interpolate_", msg, args)
}

func (b *Block) HasWorkspace() bool {
	return b.workspace != nil && b.workspace.Bool()
}

func (b *Block) BlockWorkspace() block.Workspace {
	return &Workspace{Object: b.workspace}
}

//func (b* Block)moveInputBefore  (name, refName)  { b.Call("moveInputBefore") }
//func (b* Block)moveNumberedIxpnputBefore  (

func (b *Block) RemoveInput(name string) {
	noExceptionIfMissing := false // default in blockly raises exception
	b.Call("removeInput", name, noExceptionIfMissing)
}

func (b *Block) NumInputs() int {
	return b.inputList.Length()
}

// Inputs are zeroIndexed
func (b *Block) Input(i int) block.Input {
	if cnt := b.inputList.Length(); i < 0 || i >= cnt {
		println("out of range", i, "of", cnt)
		panic("out of range")
	}
	in := b.inputList.Index(i)
	return &Input{Object: in}
}

func (b *Block) SetInput(i int, in block.Input) {
	b.inputList.SetIndex(i, in.(*Input).Object)
}

func (b *Block) InputByName(str string) (retInput block.Input, retIndex int) {
	found := false
	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
		if in := b.Input(i); in.InputName() == str {
			retInput, retIndex = in, i
			found = true
			break
		}
	}
	if !found {
		retIndex = -1
	}
	return
}

//func (b* Block)getInputTargetBlock  (name)  { b.Call("getInputTargetBlock") }
//func (b* Block)getCommentText  ()  { b.Call("getCommentText") }
//func (b* Block)setCommentText  (text)  { b.Call("setCommentText") }
//func (b* Block)setWarningText  (_text, _optId)  { b.Call("setWarningText") }

// SetMutator - blockly api to display a button which pops up a dialog to customize this block's inputs.
// func (b *Block) SetMutator(mutator *Mutator) {
// 	b.Call("setMutator", mutator.Object)
// }

//func (b* Block)getRelativeToSurfaceXY  ()  { b.Call("getRelativeToSurfaceXY") }
//func (b* Block)moveBy  (dx, dy)  { b.Call("moveBy") }
//func (b* Block)allInputsFilled  (optShadowBlocksAreFilled)  { b.Call("allInputsFilled") }
//func (b* Block)toDevString  ()  { b.Call("toDevString") }
