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
// 	Else *ControlsIfMutator "mutaiton*ControlsIfMutator"

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

// // func firstClause(containerBlock *Block) (ret *Block) {
// // 	if c := firstConnection(containerBlock); c != nil {
// // 		ret = c.TargetBlock()
// // 	}
// // 	return
// // }
