package mutation

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/blockly"
)

type savedConnections struct {
	connections connections
	numInputs   int
}

type savedMutation struct {
	itemName   block.Item
	savedAtoms []savedConnections
}

// re-connect those inputs
func (mb *mutableBlock) reconnect(savedInputs []savedMutation) (err error) {
	for _, savedInput := range savedInputs {
		itemName := savedInput.itemName
		if _, inputIndex := mb.block.InputByName(itemName); inputIndex < 0 {
			err = errutil.Append(err, errutil.New("no input block", itemName))
		} else {
			for _, savedAtom := range savedInput.savedAtoms {
				for i, c := range savedAtom.connections {
					reconnect(mb.block, inputIndex+i+1, c)
				}
				inputIndex += savedAtom.numInputs
			}
		}
	}
	return
}

/**
 * Reconnect an block to a mutated input.
 * @return {boolean} True iff a reconnection was made, false otherwise.
 */
func reconnect(block *blockly.Block, i int, tgtConnection *blockly.Connection) (okay bool) {
	if tgtConnection != nil {
		// ensure the block hasnt been disposed.
		if parentBlock := tgtConnection.GetSourceBlock(); parentBlock != nil && parentBlock.HasWorkspace() {
			if src := block.Input(i).Connection(); src != nil {
				targetBlock := tgtConnection.TargetBlock()
				if ((targetBlock == nil) || (targetBlock == block)) && src.TargetConnection() != tgtConnection {
					if src.IsConnected() {
						// There's already something connected here.  Get rid of it.
						src.Disconnect()
					}
					src.Connect(tgtConnection)
					okay = true
				}
			}
		}
	}
	return
}
