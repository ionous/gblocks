package mutant

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

type muiBuilder struct {
	mins      *BlockMutations // turns atoms into blocks
	wsBlockId string
	container block.Shape   // container to fill
	inputs    MutableInputs // atoms which will create blocks to fill the container
}

// Build mui from existing workspace blocks.
// aka decompose
func (l *muiBuilder) fillContainer() (err error) {
	l.container.InitSvg() // from blockly examples
	for i, cnt := 0, l.container.NumInputs(); i < cnt; i++ {
		muiInput := l.container.Input(i)
		if e := l.fillInput(muiInput); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// create blocks to fill the passed input
func (l *muiBuilder) fillInput(muiInput block.Input) (err error) {
	inputName := muiInput.InputName()
	if min, ok := l.mins.GetMutation(inputName); !ok {
		err = errutil.New("input not mutable", inputName)
	} else if atoms, ok := l.inputs[inputName]; ok {
		stack := muiInput.Connection()
		for i, atom := range atoms {
			if b, e := l.createBlock(min, inputName, atom, i); e != nil {
				err = errutil.Append(err, e)
			} else {
				// link the new block into the stack
				stack.Connect(b.PreviousConnection())
				stack = b.NextConnection()
			}
		}
	}
	return
}

// create a mui block to represent the named quark
// atomNum is zeroIndexed
func (l *muiBuilder) createBlock(min Mutation, inputName, atom string, atomNum int) (ret block.Shape, err error) {
	if q, ok := FindQuark(min, atom); !ok {
		err = errutil.New("couldnt find atom", min, atom)
	} else {
		mui := l.container.BlockWorkspace()
		zeroIndexed := strconv.Itoa(atomNum)
		muiBlockId := block.Scope(l.wsBlockId, inputName, zeroIndexed)
		if muiBlock, e := mui.NewBlockWithId(muiBlockId, q.BlockType()); e != nil {
			err = e
		} else {
			// r/o first block?
			if f, ok := min.FirstBlock(); ok && f.Name() == q.Name() {
				muiBlock.SetFlag(block.Movable, false)
				muiBlock.SetFlag(block.Editable, false)
			}
			muiBlock.InitSvg()
			ret = muiBlock
		}
	}
	return
}
