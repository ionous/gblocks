package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/stretchr/testify/require"
	"testing"
)

type MutationElControl struct {
	PreviousStatement, NextStatement interface{}
}
type MutationAltControl struct {
	PreviousStatement, NextStatement interface{}
}

// for testing.
func newMutatorLikeWorkspace() *Workspace {
	obj := js.Global.Get("Blockly").Get("Workspace").New()
	obj.Set("isMutator", true)
	return initWorkspace(obj)
}

func TestMutationDecompose(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		require.False(t, ws.IsMutator)
		//
		t.Log("new block")
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		ctx := ws.Context(b.Id)
		require.NotNil(t, ctx)
		require.False(t, ctx.IsValid())

		t.Log("data by id")
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant = append(d.Mutant, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		t.Log("decomposing")
		mutationui := newMutatorLikeWorkspace()
		containerBlock, e := b.decompose(ws, mutationui)
		require.NoError(t, e)
		//
		t.Log("reducing")
		require.NotNil(t, containerBlock, "reduced")
		mutationString := reduceInputs(containerBlock)
		//
		t.Log("matching", mutationString)
		require.Equal(t, []string{
			"MUTANT", "mutation_el_control", "mutation_alt_control", "mutation_alt_control",
		}, mutationString)
	})
}

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
func xTestMutationSaveConnections(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		t.Log("new block")
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		//
		t.Log("data by id")
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant = append(d.Mutant, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		t.Log("decomposing")
		mutationui := newMutatorLikeWorkspace()
		containerBlock, e := b.decompose(ws, mutationui)
		require.NoError(t, e, "decomposing")
		b.saveConnections(ws, containerBlock)
		var connections []*Connection
		//
		for mi, mcount := 0, containerBlock.NumInputs(); mi < mcount; mi++ {
			firstInput := containerBlock.Input(mi)
			if c := firstInput.Connection(); c != nil {
				connections = append(connections, c)
				for itemBlock := c.TargetBlock(); itemBlock != nil; {
					// next block in the mutation ui
					if c := itemBlock.NextConnection(); c != nil {
						connections = append(connections, c)
						itemBlock = c.TargetBlock()
					} else {
						break
					}
				}
			}
		}
		// 1+ the number of bocks [ b/c of the trailing edge ]
		require.Len(t, connections, 4)
	})
}

func TestMutationCompose(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		// ShapeTest
		// - Input  *ShapeTest
		// - Mutant []interface{} `mutation:"TestMutation"`
		// - Field  string

	})
}
