package gblocks

// import (
// 	"github.com/gopherjs/gopherjs/js"
// 	"github.com/stretchr/testify/require"
// 	"testing"
// )

// // matches 'output'
// type Boolean interface{}

// func (*True) Output() Boolean {
// 	return Boolean(nil)
// }

// // ControlsIf block
// type ControlsIf struct {
// 	If   Boolean
// 	Do   []interface{}
// 	Else *ControlsIfMutator

// 	PreviousStatement, NextStatement interface{}
// }

// // ControlsIfIf mutation
// type ControlsIfIf struct {
// 	NextStatement interface{}
// }

// // ControlsIfElseIf mutation
// type ControlsIfElseIf struct {
// 	PreviousStatement, NextStatement interface{}
// }

// // ControlsIfElse mutation
// type ControlsIfElse struct {
// 	PreviousStatement interface{}
// }

// type IfElseData struct {
// 	ElseIf Boolean
// 	Do     []interface{}
// }

// type ElseData struct {
// 	Do []interface{}
// }

// type ControlsIfMutator struct {
// 	Data []interface{}
// }

// func (m *ControlsIfMutator) NumElements() int {
// 	return len(m.Data)
// }

// func (m *ControlsIfMutator) Element(i int) interface{} {
// 	return m.Data[i]
// }

// // MutationForType - given the passed data type; what block type is needed
// func (m *ControlsIfMutator) MutationForType(dataType interface{}) (ret interface{}) {
// 	switch dataType.(type) {
// 	case nil:
// 		ret = (*ControlsIfIf)(nil)
// 	case *IfElseData:
// 		ret = (*ControlsIfElseIf)(nil)
// 	case *ElseData:
// 		ret = (*ControlsIfElse)(nil)
// 	}
// 	return
// }

// // DomToMutation -- reserve data for this node.

// // re/create the workspace blocks from the mutation dialog ui
// func (b *Block) Compose(containerBlock *Block) {
// 	restore := make(map[string]*Connection)

// 	// clauseBlock: walk down the stack of mutation ui elements.
// 	if clauseBlock := firstClause(containerBlock); clauseBlock != nil {
// 		for i, cnt := 0, len(b.subBlocks); i < cnt; i++ {
// 			// prepare to walk the existing "sub blocks" before they are destroyed.
// 			subBlock := b.subBlocks[i]
// 			var nextInput int
// 			if i+1 < cnt {
// 				nextInput = b.subBlocks[i+1].minInput
// 			} else {
// 				nextInput = len(b.inputList)
// 			}

// 			// get the workspace connections that were assigned to each mutation block during saveConnections
// 			blockConnections := b.connectionMap[clauseBlock.Id]
// 			for _, c := range blockConnections {
// 				for i := subBlock.minInput; i < nextInput; i++ {
// 					if in := b.Input(i); in != nil {
// 						restore[in.Name] = c
// 					}
// 				}
// 			}

// 			// move to the next mutation ui element
// 			if next := clauseBlock.GetNextBlock(); next != nil {
// 				clauseBlock = next
// 			} else {
// 				break
// 			}
// 		}
// 	}

// 	// rebuild the workspace ui
// 	b.updateShape()

// 	// Reconnect any child blocks.
// 	// i *believe* the events will get triggered and reconnect the data automatically.
// 	for inputName, outputConnection := range restore {
// 		reconnect(outputConnection, b, inputName)
// 	}
// }

// // "into" each mutation block, store links to the workspace's connected blocks.
// // ( so that reordering the mutations can re-order the connections )
// func (b *Block) saveConnections(containerBlock *Block) {
// 	// get the users's first placed block from the mutation dialog.
// 	// either its the one after the "head block", or its the one inside a c-shape statement input
// 	clauseBlock := firstClause(containerBlock)
// 	// mutation block to sub-block ( a set of inputs ) from this block.
// 	connectionMap := make(map[string]subBlock)
// 	// inputRange represents the *second* index of each slice;
// 	// inputRange-1 the first index
// 	for inputRange := 1; clauseBlock != nil && inputRange < len(b.inRanges); inputRange++ {
// 		// input ranges saved during composition.
// 		firstInput, lastInput := b.inputRanges[inputRange-1], b.inputRanges[inputRange]
// 		if numConnections := lastInput - firstInput + 1; numConnections > 0 {
// 			connections := make([]*Connection, numConnections)
// 			for i := firstInput; i < lastInput; i++ {
// 				in := b.Input(i)
// 				connections[i] = in.Connection.TargetConnection
// 			}
// 			connectionMap[clauseBlock.Id] = subBlock{
// 				data[inputRange-1], // maybe something like this...
// 				connections,
// 			}
// 		}
// 		clauseBlock = clauseBlock.GetNextBlock()
// 	}
// }
