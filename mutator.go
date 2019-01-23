package gblocks

import "github.com/gopherjs/gopherjs/js"

type Mutator struct {
	*js.Object
}

func NewMutator(quarkNames []TypeName) (ret *Mutator) {
	if blockly := js.Global.Get("Blockly"); blockly.Bool() {
		obj := blockly.Get("Mutator").New(quarkNames)
		ret = &Mutator{Object: obj}
	}
	return
}

/**
 * Reconnect an block to a mutated input.
 * @return {boolean} True iff a reconnection was made, false otherwise.
 */
func reconnect(block *Block, i int, tgtConnection *Connection) (okay bool) {
	if tgtConnection != nil {
		// ensure the block hasnt been disposed.
		if parentBlock := tgtConnection.GetSourceBlock(); parentBlock != nil && parentBlock.hasWorkspace() {
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
