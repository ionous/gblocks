package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type MutationElStart struct {
	NextStatement interface{}
}

type MutationElControl struct {
	PreviousStatement, NextStatement interface{}
}
type MutationAltControl struct {
	PreviousStatement, NextStatement interface{}
}

func TestComposeDecompose(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		t.Log("new block")
		b, e := ws.NewBlock((*ShapeTest)(nil))
		require.NoError(t, e)
		// _, foundMutation := ws.Context(b.Id).FieldByName("MUTANT")
		// require.True(t, foundMutation, "found mutation")
		//
		t.Log("data by id")
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant = append(d.Mutant, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		t.Log("decomposing")
		containerBlock, e := b.decompose(ws)
		require.NoError(t, e)
		//
		t.Log("reducing")
		require.NotNil(t, containerBlock, "reduced")
		mutationString := reduceBlocks(containerBlock)
		//
		t.Log("matching")
		require.Equal(t, []string{
			"MUTANT", "mutation_el_start", "mutation_el_control", "mutation_alt_control", "mutation_alt_control",
		}, mutationString)
	})
}

func reduceInputs(block *Block) (ret []string) {
	for i, cnt := 0, block.NumInputs(); i < cnt; i++ {
		in := block.Input(i)
		ret = append(ret, in.Name.String())
		block := in.Connection.TargetBlock()
		ret = append(ret, reduceBlocks(block)...)
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

// new a block with data.
// 	you have to build containerBlock .
// 	check connections.
func TestComposeSave(t *testing.T) {
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
		containerBlock, e := b.decompose(ws)
		require.NoError(t, e, "decomposing")
		b.saveConnections(ws, containerBlock)
		var connections []*Connection
		//
		for mi, mcount := 0, containerBlock.NumInputs(); mi < mcount; mi++ {
			firstInput := containerBlock.Input(mi)
			if c := firstInput.Connection; c != nil {
				connections = append(connections, c)
				for itemBlock := c.TargetBlock(); itemBlock != nil; {
					// next block in the mutation ui
					if c := itemBlock.NextConnection; c != nil {
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
