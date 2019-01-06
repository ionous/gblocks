package gblocks

/**
 * Reconnect an block to a mutated input.
 * @return {boolean} True iff a reconnection was made, false otherwise.
 */
func reconnect(block *Block, i int, tgtConnection *Connection) (okay bool) {
	if tgtConnection != nil {
		// ensure the block hasnt been disposed.
		if parentBlock := tgtConnection.GetSourceBlock(); parentBlock != nil && parentBlock.hasWorkspace() {
			src := block.Input(i).Connection()
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
	return
}
