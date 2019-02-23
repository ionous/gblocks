package mutation

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/blockly"
)

// mapping of workspace block id to custom mutation data.
// alt: extend blockly's block and input type.
// that's a bit harder to do well via gopher.
type blockDataMap map[block.Id]*blockData

// per-block, map each input name to its mutation data
type blockData struct {
	inputMap    inputMap
	connections connectionMap
}

type inputMap map[block.Item]*inputData
type connectionMap map[block.Id]connections
type connections []*blockly.Connection

func (bd *blockData) saveConnections(muiBlock *blockly.Block, cs connections) {
	bd.connections[muiBlock.Id] = cs
}

func (bd *blockData) getConnections(muiBlock *blockly.Block) connections {
	return bd.connections[muiBlock.Id]
}

func (cs connections) appendInput(in *blockly.Input) connections {
	var target *blockly.Connection
	if c := in.Connection(); c != nil {
		target = c.TargetConnection()
	}
	return append(cs, target)
}
