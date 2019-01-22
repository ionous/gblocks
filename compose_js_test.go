package gblocks

import (
	// "github.com/gopherjs/gopherjs/js"
	// "github.com/kr/pretty"
	"github.com/ionous/errutil"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type MutationElControl struct {
	PreviousStatement, NextStatement interface{}
}
type MutationAltControl struct {
	PreviousStatement, NextStatement interface{}
}

func buildMutation(ws *Workspace, reg *Registry, t *testing.T) (ret *Block, err error) {
	if b, e := ws.NewBlock((*ShapeTest)(nil)); e != nil {
		err = e
	} else if in, index := b.InputByName("MUTANT"); index < 0 {
		err = errutil.New("missing input")
	} else if m := in.Mutation(); m == nil {
		err = errutil.New("missing mutation")
	} else {
		for _, atomType := range []TypeName{"atom_test", "atom_alt_test", "atom_test"} {
			if numInputs, e := m.addAtom(reg, atomType); e != nil {
				err = e
				break
			} else if numInputs != 1 {
				err = errutil.New(atomType, "generated unexpected inputs", numInputs)
				break
			}
		}
		if err == nil {
			ret = b
		}
	}
	return
}

func TestMutationDecompose(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		b, e := buildMutation(ws, reg, t)
		require.NoError(t, e)
		//
		t.Log("decomposing")
		mui := NewBlankWorkspace(true)
		muiContainer, e := b.decompose(reg, mui)
		require.NoError(t, e, "created mui container")
		//
		t.Log("reducing")
		require.NotNil(t, muiContainer, "reduced")
		mutationString := reduceInputs(muiContainer)
		//
		t.Log("matching", mutationString)
		require.Equal(t, []string{
			"MUTANT", "mutation_el_control", "mutation_alt_control", "mutation_el_control",
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

func listConnections(b *Block) (ret []string) {
	for mi, mcount := 0, b.NumInputs(); mi < mcount; mi++ {
		firstInput := b.Input(mi)
		if muiConnection := firstInput.Connection(); muiConnection != nil {
			for muiBlock := muiConnection.TargetBlock(); muiBlock != nil; {
				//
				cs := muiBlock.connections
				// for i, cnt := 0, cs.Length(); i < cnt; i++ {
				var id string
				if c := cs.Connection(0); c != nil {
					if tgt := c.TargetBlock(); tgt != nil {
						id = tgt.Id
					}
				}
				ret = append(ret, id)
				// }
				if muiConnection := muiBlock.NextConnection(); muiConnection != nil {
					muiBlock = muiConnection.TargetBlock()
				} else {
					break
				}
			}
		}
	}
	return
}

func listBlocks(b *Block) (ret []string) {
	for mi, mcount := 0, b.NumInputs(); mi < mcount; mi++ {
		firstInput := b.Input(mi)
		if muiConnection := firstInput.Connection(); muiConnection != nil {
			for muiBlock := muiConnection.TargetBlock(); muiBlock != nil; {

				var id string
				if c := muiBlock.connections.Connection(0); c != nil {
					if t := c.TargetBlock(); t != nil {
						id = t.Id
					}
				}
				ret = append(ret, muiBlock.Id+"----"+id)
				//
				if muiConnection := muiBlock.NextConnection(); muiConnection != nil {
					muiBlock = muiConnection.TargetBlock()
				} else {
					break
				}
			}
		}
	}
	return
}

// new a block with data. run a minimal check of connections.
// func TestMutationConnections(t *testing.T) {
// 	testShape(t, func(ws *Workspace, reg *Registry) {
// 		b, e := buildMutation(ws, reg, t)
// 		require.NoError(t, e)
// 		//
// 		in, where := b.InputByName("MUTANT/0/SUB_INPUT")
// 		require.NotEqual(t, -1, where)
// 		require.NotNil(t, in)
// 		// connect the first input
// 		target, e := ws.NewBlock("shape_test")
// 		require.NoError(t, e)
// 		in.Connection().Connect(target.OutputConnection())

// 		// decompose to create a mui
// 		mui := NewBlankWorkspace(true)
// 		muiContainer, e := b.decompose(reg, mui)
// 		require.NoError(t, e, "created mui container")

// 		// save connections
// 		{
// 			e := b.saveConnections(muiContainer)
// 			require.NoError(t, e, "initial save")
// 			blocks := listBlocks(muiContainer)
// 			t.Log("xxx\n", strings.Join(blocks, "\n"))
// 			//
// 			targets := listConnections(muiContainer)
// 			require.Len(t, targets, 3)
// 			t.Log("initial", strings.Join(targets, ","))
// 			require.NotEmpty(t, targets[0], "initial connected")
// 			require.Empty(t, targets[2], "initial not connected")
// 		}

// 		// disconnect block
// 		firstIn := muiContainer.Input(0) // MUTANT
// 		c := firstIn.Connection()
// 		require.NotNil(t, c, "first statement connection")
// 		muiBlock := c.TargetBlock()
// 		// nextNext := muiBlock.GetNextBlock().NextConnection()
// 		muiBlock.Unplug(true)

// 		// removing the mui block hasnt changed the atom's number of inputs
// 		// removing the first block should act as if the first
// 		// compose *might+ happend on block change before save conncetions

// 		// check connections
// 		if e := muiBlock.compose(reg, muiContainer); e != nil {
// 			t.Fatal(e)
// 		} else {
// 			e := b.saveConnections(muiContainer)
// 			require.NoError(t, e, "disconnected save")
// 			blocks := listBlocks(muiContainer)
// 			t.Log("xxx\n", strings.Join(blocks, "\n"))

// 			targets := listConnections(muiContainer)
// 			t.Log("disconnected", strings.Join(targets, ","))
// 			require.Empty(t, targets[0], "disconnected connected")
// 			require.NotEmpty(t, targets[2], "disconnected not connected")
// 		}

// 		// connect block at the end
// 		// nextNext.Connect(muiBlock.PreviousConnection())

// 		// // check connections
// 		// {
// 		// 	e := b.saveConnections(muiContainer)
// 		// 	require.NoError(t, e, "reconnected save")
// 		// 	targets := listConnections(muiContainer)
// 		// 	require.Len(t, targets, 3)
// 		// 	t.Log("reconnected", strings.Join(targets, ","))
// 		// 	require.Empty(t, targets[0], "reconnected connected")
// 		// 	require.NotEmpty(t, targets[2], "reconnected not connected")
// 		// }
// 	})
// }

// re/create the workspace blocks from the mutation dialog ui
func TestMutationCompose(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		// create mutation blocks
		mui := NewBlankWorkspace(true)
		container, err := mui.NewBlock("shape_test$mutation")
		require.NoError(t, err)

		var block [3](*Block)
		src := [3]interface{}{
			(*MutationElControl)(nil),
			(*MutationAltControl)(nil),
			(*MutationAltControl)(nil),
		}

		t.Log("building blocks")
		for i := 0; i < len(src); i++ {
			b, e := mui.NewBlock(src[i])
			require.NoError(t, e)
			require.NotNilf(t, b, "new block %d", i)
			block[i] = b
		}

		//
		t.Log("connecting a->b->c")
		container.Input(0).Connection().Connect(block[0].PreviousConnection())
		block[0].NextConnection().Connect(block[1].PreviousConnection())
		block[1].NextConnection().Connect(block[2].PreviousConnection())

		b, err := ws.NewBlock("shape_test")
		require.NoError(t, err)

		if e := b.compose(reg, container); e != nil {
			t.Fatal(e)
		} else {
			// test the composed block
			composed := reduceInputs(b)
			str := strings.Join(composed, ",")
			require.Equal(t, "INPUT,MUTANT,MUTANT/0/SUB_INPUT,,,", str)
		}
	})
}
