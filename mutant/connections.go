package mutant

import (
	"strings"

	"github.com/ionous/gblocks/block"
)

// map pending mui block inputs to workspace connections
// so rearranging mui blocks can rearrange the workspace.
type SavedConnections map[string]block.Connection

func (cs SavedConnections) saveInput(in block.Input) {
	if c := in.Connection(); c != nil && c.IsConnected() {
		name := in.InputName()
		cs[name] = c.TargetConnection()
	}
}

// generates a mapping of atom inputs to workspace targets.
func SaveConnections(main, mui block.Inputs) SavedConnections {
	out := make(SavedConnections)
	// for each input in the mutation ui
	for i, cnt := 0, mui.NumInputs(); i < cnt; i++ {
		// visit every mui block connected to this mui container input
		muiInput := mui.Input(i)
		block.VisitStack(muiInput, func(muiBlock block.Shape) bool {
			// name helper:
			muiBlockId := muiBlock.BlockId()
			// inputs names created from existing data ( via atomParser.parseAtom ):
			//    input= "a, wsBlockId, INPUT, atomNum, FIELD"
			//
			// mui block ids created from existing data ( via muiBuilder.createBlock ):
			//    muiBlockId= "wsBlockId, INPUT, atomNum"
			//
			// if the user placed the mui block themselves, the id will be random:
			//    muiBlockId= "auto-generated"
			//
			// re/creating workspace inputs from the mui ( via muiParser.createAtomsAt ) yields:
			// 	  input= "a, <muiBlockId>, FIELD"
			//
			// when we save/restore an input here, we use: "a, <muiBlockId>, FIELD":
			//    save slot= "a, wsBlockId, INPUT, atomNum, FIELD" or
			//               "a, auto-generated, FIELD"
			//
			// note: the blockId is used to differentiate blocks when multiple popups are open
			//
			atomPrefix := block.Scope("a", muiBlockId)
			for i, cnt := 0, main.NumInputs(); i < cnt; i++ {
				in := main.Input(i)
				name := in.InputName()
				if strings.HasPrefix(name, atomPrefix) {
					out.saveInput(in)
				}
			}
			return true
		})

	}
	return out
}

func RestoreConnections(block block.Shape, savedConnections SavedConnections) {
	// we expect names like: "a$ muiBlockId $ FIELD"
	for atomInput, oldDst := range savedConnections {
		Reconnect(oldDst, block, atomInput)
	}
	return
}

// ported from blockly
func Reconnect(connectionChild block.Connection, block block.Shape, inputName string) (okay bool) {
	if isConnectionValid(connectionChild) {
		if in, dex := block.InputByName(inputName); dex >= 0 {
			if connectionParent := in.Connection(); connectionParent != nil {
				currentParent := connectionChild.TargetBlock()
				if currentParent == nil || currentParent == block {
					if connectionParent.TargetConnection() != connectionChild {
						if connectionParent.IsConnected() {
							connectionParent.Disconnect()
						}
						connectionParent.Connect(connectionChild)
						okay = true
					}
				}
			}
		}
	}
	return
}

func isConnectionValid(c block.Connection) (okay bool) {
	if c != nil {
		connectionOwner := c.SourceBlock()
		if connectionOwner != nil && connectionOwner.HasWorkspace() {
			okay = true
		}
	}
	return
}
