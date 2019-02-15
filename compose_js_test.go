package gblocks

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/named"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func addTestAtoms(ws *Workspace, reg *Registry, t *testing.T) (ret *Block, err error) {
	if b, e := ws.NewBlock((*ShapeTest)(nil)); e != nil {
		err = e
	} else if in, index := b.InputByName("MUTANT"); index < 0 {
		err = errutil.New("addTestAtoms missing input")
	} else if m := in.Mutation(); m == nil {
		err = errutil.New("addTestAtoms missing mutation")
	} else {
		for _, atomType := range []named.Type{"atom_test", "atom_alt_test", "atom_test"} {
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
		b, e := addTestAtoms(ws, reg, t)
		require.NoError(t, e)
		//
		t.Log("decomposing")
		mui := NewBlankWorkspace(true)
		mui.idGen = &orderedGenerator{name: "mui"}
		muiContainer, e := b.decompose(reg, mui)
		require.NoError(t, e, "created mui container")
		//
		t.Log("reducing")
		require.NotNil(t, muiContainer, "reduced")
		mutationString := reduceInputs(muiContainer)
		//
		t.Log("matching", mutationString)
		require.Equal(t, []string{
			"MUTANT", "mui$test_mutation$atom_test", "mui$test_mutation$atom_alt_test", "mui$test_mutation$atom_test",
		}, mutationString)
	})
}

func reduceInputs(block *Block) (ret []string) {
	for i, cnt := 0, block.NumInputs(); i < cnt; i++ {
		in := block.Input(i)
		// inputs created for fields dont have names
		if n := in.Name.String(); len(n) > 0 {
			if n == "undefined" {
				panic("field")
			}
			ret = append(ret, n)
		}
		fields := in.Fields()
		for i, cnt := 0, fields.Length(); i < cnt; i++ {
			field := fields.Field(i)
			if n := field.Name(); len(n) > 0 {
				ret = append(ret, n)
			}
		}

		in.visitStack(func(b *Block) (keepGoing bool) {
			ret = append(ret, b.Type.String())
			return true
		})
	}
	return
}

func (cs *Connections) blocks() (ret []string) {
	for i, cnt := 0, cs.Length(); i < cnt; i++ {
		if c := cs.Connection(i); c != nil {
			if tgt := c.GetSourceBlock(); tgt != nil {
				ret = append(ret, tgt.Id)
			}
		}
	}
	return
}

type listed struct {
	id      string // id of the mui block
	targets []string
}

func listConnections(b *Block) (ret []listed) {
	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
		in := b.Input(i)
		in.visitStack(func(nextBlock *Block) (keepGoing bool) {
			if cs := nextBlock.CachedConnections(); cs != nil {
				targets := cs.blocks()
				ret = append(ret, listed{nextBlock.Id, targets})
			}
			return true
		})
	}
	return
}

// re/create the workspace blocks from the mutation dialog ui
func TestMutationCompose(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		// create mutation blocks
		mui := NewBlankWorkspace(true)
		mui.idGen = &orderedGenerator{name: "mui"}
		muiContainer, err := mui.NewBlock(named.SpecialType("mui_container", "shape_test"))
		require.NoError(t, err)

		var muiBlocks [3](*Block)
		src := [3]named.Type{
			named.SpecialType("mui", "test_mutation", "atom_test"),
			named.SpecialType("mui", "test_mutation", "atom_alt_test"),
			named.SpecialType("mui", "test_mutation", "atom_alt_test"),
		}

		t.Log("building blocks")
		for i := 0; i < len(src); i++ {
			b, e := mui.NewBlock(src[i])
			require.NoError(t, e)
			require.NotNilf(t, b, "new muiBlocks %d", i)
			muiBlocks[i] = b
		}

		//
		t.Log("connecting a->b->c")
		muiContainer.Input(0).Connection().Connect(muiBlocks[0].PreviousConnection())
		muiBlocks[0].NextConnection().Connect(muiBlocks[1].PreviousConnection())
		muiBlocks[1].NextConnection().Connect(muiBlocks[2].PreviousConnection())

		b, err := ws.NewBlock("shape_test")
		require.NoError(t, err)

		if e := b.compose(reg, muiContainer); e != nil {
			t.Fatal(e)
		} else {
			// test the composed block
			composed := reduceInputs(b)
			str := strings.Join(composed, ",")
			require.Equal(t, "INPUT,MUTANT,MUTANT/1/ATOM_INPUT,MUTANT/2/ATOM_FIELD,MUTANT/3/ATOM_FIELD,FIELD", str)
		}
	})
}

// new a block with data. run a minimal check of connections.
// save connections requires compose
func TestMutationConnections(t *testing.T) {
	testShape(t, func(ws *Workspace, reg *Registry) {
		b, e := addTestAtoms(ws, reg, t)
		require.NoError(t, e)
		//
		in, where := b.InputByName("MUTANT/1/ATOM_INPUT")
		require.NotEqual(t, -1, where)
		require.NotNil(t, in)

		// connect the first input
		target, e := ws.NewBlock("shape_test")
		require.NoError(t, e)
		in.Connection().Connect(target.OutputConnection())

		// decompose to create a mui
		mui := NewBlankWorkspace(true)
		mui.idGen = &orderedGenerator{name: "mui"}
		muiContainer, e := b.decompose(reg, mui)
		if e != nil {
			t.Fatal(e, "created mui container")
		} else {
			e := b.saveConnections(muiContainer)
			require.NoError(t, e, "initial save")
			//
			targets := listConnections(muiContainer)
			t.Log("initial targets:", pretty.Sprint(targets))
			require.Len(t, targets, 3)
			requires := require.NotEmptyf
			for i, tgt := range targets {
				requires(t, tgt.targets, "initial target %d", i)
				requires = require.Emptyf
			}
		}

		// disconnect block
		firstIn := muiContainer.Input(0) // MUTANT
		c := firstIn.Connection()
		require.NotNil(t, c, "first statement connection")
		muiBlock := c.TargetBlock()
		nextNext := muiBlock.GetNextBlock().GetNextBlock().NextConnection() // remember this one

		require.Equal(t,
			[]string{"MUTANT", "mui$test_mutation$atom_test", "mui$test_mutation$atom_alt_test", "mui$test_mutation$atom_test"},
			reduceInputs(muiContainer), "before unplug")
		muiBlock.Unplug(true)
		require.Equal(t,
			[]string{"MUTANT", "mui$test_mutation$atom_alt_test", "mui$test_mutation$atom_test"},
			reduceInputs(muiContainer), "after unplug")

		// removing the mui block hasnt changed the atom's number of inputs
		// removing the first block should act as if the first
		// compose *might+ happend on block change before save conncetions
		if cs := muiBlock.CachedConnections(); cs == nil {
			t.Fatal("no cached connections")
		} else {
			require.NotNil(t, cs.Connection(0))
		}

		// check connections after disconnect
		if e := b.compose(reg, muiContainer); e != nil {
			t.Fatal("recompose after disconnect:", e)
		} else {
			e := b.saveConnections(muiContainer)
			require.NoError(t, e, "disconnected save")

			targets := listConnections(muiContainer)
			t.Log("disconnected targets:", pretty.Sprint(targets))
			require.Len(t, targets, 2)
			for i, tgt := range targets {
				require.Emptyf(t, tgt.targets, "disconnected target %d", i)
			}
		}

		// connect block at the end
		if cs := muiBlock.CachedConnections(); cs == nil {
			t.Fatal("no chaced connections")
		} else {
			require.NotNil(t, cs.Connection(0), "preconnected")
			nextNext.Connect(muiBlock.PreviousConnection())

			require.Equal(t,
				[]string{"MUTANT", "mui$test_mutation$atom_alt_test", "mui$test_mutation$atom_test", "mui$test_mutation$atom_test"},
				reduceInputs(muiContainer), "after reconnect")
			require.Equal(t, 1, cs.Length(), "reconnected length")
			require.NotNil(t, cs.Connection(0), "reconnected")
		}

		// check connections
		if e := b.compose(reg, muiContainer); e != nil {
			t.Fatal("recompose after reconnect:", e)
		} else {
			e := b.saveConnections(muiContainer)
			require.NoError(t, e, "reconnected save")

			// t.Log("reconnected", muiBlock.Id, muiBlock.connections.blocks())
			targets := listConnections(muiContainer)
			t.Log("reconnected targets:", pretty.Sprint(targets))
			require.Len(t, targets, 3)
			require.Empty(t, targets[0].targets, "reconnected target 0")
			require.Empty(t, targets[1].targets, "reconnected target 1")
			require.NotEmpty(t, targets[2].targets, "reconnected target 2")
		}
	})
}
