package gblocks

// Prepare the mutator's dialog with this block's components.
func decompose(ws *Workspace, m Mutation) *Block {
	// ask for the type used for the head block
	mutationType := m.MutationForType(nil)
	// turn that into a block
	containerBlock := ws.NewBlock(mutationType)
	containerBlock.InitSvg()
	connection := firstConnection(containerBlock)
	els := m.Elements()
	for i, cnt := 0, els.Len(); i < cnt; i++ {
		iface := els.Index(i)
		ptr := iface.Elem()
		el := ptr.Elem()
		// alt:
		mutationType := m.MutationForType(el.Type())
		if block := ws.NewBlock(mutationType); block == nil {
			// FIX: return error
			break
		} else {
			block.InitSvg()
			connection.Connect(block.PreviousConnection)
			connection = block.NextConnection
		}
	}

	return containerBlock
}
