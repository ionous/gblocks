package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	// "github.com/kr/pretty"
	// "github.com/stretchr/testify/require"
	// "strings"
	// "testing"
)

type MutationElControl struct {
	PreviousStatement, NextStatement interface{}
}
type MutationAltControl struct {
	PreviousStatement, NextStatement interface{}
}

// for testing.
func newMutatorLikeWorkspace() (ret *Workspace) {
	if blockly := js.Global.Get("Blockly"); blockly.Bool() {
		obj := blockly.Get("Workspace").New()
		ret := &Workspace{Object: obj}
		ret.IsMutator = true
	}
	return
}

// func xTestMutationDecompose(t *testing.T) {
// 	testShape(t, func(ws *Workspace, reg*Registry) {
// 		require.False(t, ws.IsMutator, "is mutator")
// 		//
// 		t.Log("new block")
// 		b, e := ws.NewBlock((*ShapeTest)(nil))
// 		require.NoError(t, e)

// 		t.Log("data by id")
// 		d := ws.GetDataById(b.Id).(*ShapeTest)
// 		d.Mutant = append(d.Mutant, &AtomTest{}, &AtomAltTest{}, &AtomAltTest{})
// 		//
// 		t.Log("decomposing")
// 		mui := newMutatorLikeWorkspace()
// 		containerBlock, e := b.decompose(ws, mui)
// 		require.NoError(t, e)
// 		//
// 		t.Log("reducing")
// 		require.NotNil(t, containerBlock, "reduced")
// 		mutationString := reduceInputs(containerBlock)
// 		//
// 		t.Log("matching", mutationString)
// 		require.Equal(t, []string{
// 			"MUTANT", "mutation_el_control", "mutation_alt_control", "mutation_alt_control",
// 		}, mutationString)
// 	})
// }

func reduceInputs(block *Block) (ret []string) {
	for i, cnt := 0, block.NumInputs(); i < cnt; i++ {
		in := block.Input(i)
		ret = append(ret, in.Name.String())
		if c := in.Connection(); c != nil {
			block := c.TargetBlock()
			ret = append(ret, reduceBlocks(block)...)
		}
	}
	return
}

func reduceBlocks(block *Block) (ret []string) {
	for i := 0; block != nil && i < 100; i++ {
		ret = append(ret, block.Type.String())
		block = block.GetNextBlock()
	}
	return
}

// new a block with data. check connections.
// func xTestMutationSaveConnections(t *testing.T) {
// 	testShape(t, func(ws *Workspace, reg*Registry) {
// 		t.Log("new block")
// 		b, e := ws.NewBlock((*ShapeTest)(nil))
// 		require.NoError(t, e)
// 		//
// 		t.Log("data by id")
// 		d := ws.GetDataById(b.Id).(*ShapeTest)
// 		d.Mutant = append(d.Mutant, &AtomTest{}, &AtomAltTest{}, &AtomAltTest{})
// 		//
// 		t.Log("decomposing")
// 		mui := newMutatorLikeWorkspace()
// 		containerBlock, e := b.decompose(ws, mui)
// 		require.NoError(t, e, "decomposing")
// 		b.saveConnections(containerBlock)
// 		var connections []*Connection
// 		//
// 		for mi, mcount := 0, containerBlock.NumInputs(); mi < mcount; mi++ {
// 			firstInput := containerBlock.Input(mi)
// 			if c := firstInput.Connection(); c != nil {
// 				connections = append(connections, c)
// 				for itemBlock := c.TargetBlock(); itemBlock != nil; {
// 					// next block in the mutation ui
// 					if c := itemBlock.NextConnection(); c != nil {
// 						connections = append(connections, c)
// 						itemBlock = c.TargetBlock()
// 					} else {
// 						break
// 					}
// 				}
// 			}
// 		}
// 		// 1+ the number of bocks [ b/c of the trailing edge ]
// 		require.Len(t, connections, 4)
// 	})
// }

// re/create the workspace blocks from the mutation dialog ui
// func xTestMutationCompose(t *testing.T) {
// 	testShape(t, func(ws *Workspace, reg*Registry) {
// 		// create mutation blocks
// 		mui := newMutatorLikeWorkspace()
// 		container, err := mui.NewBlock("shape_test$mutation")
// 		require.NoError(t, err)

// 		var block [3](*Block)
// 		src := [3]interface{}{
// 			(*MutationElControl)(nil),
// 			(*MutationAltControl)(nil),
// 			(*MutationAltControl)(nil),
// 		}

// 		t.Log("building blocks")
// 		for i := 0; i < len(src); i++ {
// 			b, e := mui.NewBlock(src[i])
// 			require.NoError(t, e)
// 			require.NotNilf(t, b, "new block %d", i)
// 			block[i] = b
// 		}

// 		//
// 		t.Log("connecting a->b->c")
// 		container.Input(0).Connection().Connect(block[0].PreviousConnection())
// 		block[0].NextConnection().Connect(block[1].PreviousConnection())
// 		block[1].NextConnection().Connect(block[2].PreviousConnection())

// 		b, err := ws.NewBlock("shape_test")
// 		require.NoError(t, err)

// 		b.compose(ws, container)

// 		// test the composed block
// 		composed := reduceInputs(b)
// 		str := strings.Join(composed, ",")
// 		require.Equal(t, "INPUT,MUTANT,MUTANT/0/SUB_INPUT,,,", str)
// 		d := ws.GetDataById(b.Id)

// 		// test the generated data
// 		expected := &ShapeTest{
// 			Mutant: []interface{}{
// 				&AtomTest{},
// 				&AtomAltTest{},
// 				&AtomAltTest{},
// 			},
// 		}
// 		v := pretty.Diff(d, expected)
// 		if len(v) != 0 {
// 			t.Fatal(v)
// 			t.Log(v)
// 		}
// 	})
// }
