package mutant_test

//
//import (
//	"sort"
//	"testing"
//
//	"github.com/ionous/gblocks/block"
//	"github.com/ionous/gblocks/mock"
//	"github.com/ionous/gblocks/mutant"
//	"github.com/stretchr/testify/require"
//)
//
//var common = struct {
//	inputs, quarks []string
//	inputAtoms     map[string][]mutant.AtomizedInput
//	muiContainer   block.Dict
//	atomProducts   map[string][]mock.MockAtom
//	// workspace inputs after expanding the mutable input atoms
//	expandedInputs []string
//}{
//	inputs: []string{"M1:input_dummy", "M2:input_dummy", "M3:input_dummy"},
//	quarks: []string{"A", "B", "C"},
//	inputAtoms: map[string][]mutant.AtomizedInput{
//		"M1": []mutant.AtomizedInput{
//			{"A", 1},
//			{"A", 1},
//		},
//		"M2": []mutant.AtomizedInput{
//			{"B", 1},
//		},
//		"M3": []mutant.AtomizedInput{
//			{"A", 1},
//			{"B", 1},
//			{"C", 2}, // it should really be just one;
//			// but mock doesnt collapse the fields into inputs
//			// the same way blockly does
//		},
//	},
//	muiContainer: block.Dict{
//		"type":     block.Scope("mui", "test"),
//		"message0": "%1 %2 %3",
//		"args0": []block.Dict{{
//			"name":  "M1",
//			"type":  "input_statement",
//			"check": "A",
//		}, {
//			"type":  "input_statement",
//			"name":  "M2",
//			"check": "A",
//		}, {
//			"name":  "M3",
//			"type":  "input_statement",
//			"check": "A",
//		}},
//	},
//	atomProducts: map[string][]mock.MockAtom{
//		"A": []mock.MockAtom{{Name: "TERM", Type: block.ValueInput}},
//		"B": []mock.MockAtom{{Name: "NUM", Type: block.NumberField}},
//		"C": []mock.MockAtom{
//			{Name: "TEXT", Type: block.TextField},
//			{Name: "STATE", Type: block.StatementInput},
//		},
//	},
//	expandedInputs: []string{
//		// a, blockId, INPUT, atomNum, FIELD:type
//		"M1:input_dummy",
//		/*A*/ block.Scope("a", "M1", "0", "TERM:input_value"),
//		/*A*/ block.Scope("a", "M1", "1", "TERM:input_value"),
//		"M2:input_dummy",
//		/*B*/ block.Scope("a", "M2", "0", "NUM:field_number"),
//		"M3:input_dummy",
//		/*A*/ block.Scope("a", "M3", "0", "TERM:input_value"),
//		/*B*/ block.Scope("a", "M3", "1", "NUM:field_number"),
//		/*C*/ block.Scope("a", "M3", "2", "TEXT:field_input"),
//		/*C*/ block.Scope("a", "M3", "2", "STATE:input_statement"),
//	},
//}
//
//func TestCreateMockBlock(t *testing.T) {
//	// create a workspace block where we will be writing data to
//	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
//	require.Equal(t, b.NumInputs(), 3)
//	require.Equal(t, common.inputs, listInputs(b))
//}
//
//func TestDescribeContainer(t *testing.T) {
//	arch := mock.NewMutations(common.inputs, common.quarks)
//	muiContainer := arch.DescribeContainer(mutant.ContainerName("test"))
//	require.Equal(t, common.muiContainer, muiContainer)
//}
//
//func TestPreregister(t *testing.T) {
//	var reg mock.Registry
//	arch := mock.NewMutations(common.inputs, common.quarks)
//	require.NoError(t, arch.Preregister("test", &reg))
//	var keys []string
//	for k, _ := range reg.Blocks {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//	// we expect to see unique mutations for each input
//	// ( named after its input )
//	// and we expect to see unique quarks for each mutation
//	// ( also the mui container which can hold the quarks )
//	expected := []string{
//		block.Scope("mui", "M1", "A"),
//		block.Scope("mui", "M1", "B"),
//		block.Scope("mui", "M1", "C"),
//		block.Scope("mui", "M2", "A"),
//		block.Scope("mui", "M2", "B"),
//		block.Scope("mui", "M2", "C"),
//		block.Scope("mui", "M3", "A"),
//		block.Scope("mui", "M3", "B"),
//		block.Scope("mui", "M3", "C"),
//		block.Scope("mui", "test"),
//	}
//	require.Equal(t, expected, keys)
//}
//
//func TestCreateMui(t *testing.T) {
//	var reg mock.Registry
//	muispace := reg.NewMockSpace()
//	blocks := mutant.NewMutatedBlocks()
//	arch := mock.NewMutations(common.inputs, common.quarks)
//	require.NoError(t, arch.Preregister("mockType", &reg))
//	wsblock := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
//	mb := blocks.CreateMutatedBlock(wsblock, arch, nil)
//	for k, v := range common.inputAtoms {
//		mi, ok := mb.GetMutatedInput(k)
//		require.True(t, ok)
//		mi.Atoms = v
//	}
//	c, e := mb.CreateMui(muispace)
//	// note, the block type for mui blocks comes from MockQuark.BlockType
//	// "mui, mutation, atom"
//	require.NoError(t, e, "CreateMui")
//	ids := listStack(c)
//	// MOD: mui block creation used to specify particular ids
//	// now it just uses whatever blockly/mock chooses. ( NewBlock vs. NewBlockWithId )
//	expectedIds := []string{
//		block.Scope("mui", "M1", "A#0:mui", "M1", "A"),
//		/**/ block.Scope("mui", "M1", "A#1:mui", "M1", "A"),
//		block.Scope("mui", "M2", "B#0:mui", "M2", "B"),
//		block.Scope("mui", "M3", "A#0:mui", "M3", "A"),
//		/**/ block.Scope("mui", "M3", "B#0:mui", "M3", "B"),
//		/**/ block.Scope("mui", "M3", "C#0:mui", "M3", "C"),
//	}
//	require.Equal(t, expectedIds, ids, "expectedIds")
//}
//
//// list all blocks connected to inputs of b
//// each mock block returns the format "id:type"
//func listStack(b block.Shape) (ret []string) {
//	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
//		in := b.Input(i)
//		block.VisitStack(in, func(b block.Shape) bool {
//			ret = append(ret, b.(interface{ String() string }).String())
//			return true // keepGoing
//		})
//	}
//	return
//}
//
//// list every named input of b
//// note: with the mock implementation, every field gets promoted to a connectionless input.
//func listInputs(b block.Shape) (ret []string) {
//	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
//		in := b.Input(i)
//		ret = append(ret, in.(interface{ String() string }).String())
//	}
//	return
//}
//
//func TestCreateFromMui(t *testing.T) {
//	var reg mock.Registry
//	// create a workspace block where we will be writing data to
//	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
//	// first create a container and give it some fake atom data.
//	muispace := reg.NewMockSpace()
//	// have to create the mui in order to fill from it.
//	// (  a container with inputs containing stacks of atoms )
//	arch := mock.NewMutations(common.inputs, common.quarks)
//	require.NoError(t, arch.Preregister("mockType", &reg))
//	blocks := mutant.NewMutatedBlocks()
//	db := &mock.MockDatabase{common.atomProducts}
//	mb := blocks.CreateMutatedBlock(b, arch, db)
//	for k, v := range common.inputAtoms {
//		mi, ok := mb.GetMutatedInput(k)
//		require.True(t, ok)
//		mi.Atoms = v
//	}
//
//	c, e := mb.CreateMui(muispace)
//	require.NoError(t, e)
//	// now, fill the b with the blocks from the container
//	e = mb.CreateFromMui(c)
//	require.NoError(t, e)
//	for k, v := range common.inputAtoms {
//		mi, ok := mb.GetMutatedInput(k)
//		require.True(t, ok)
//		require.Equal(t, v, mi.Atoms)
//	}
//	//
//	expandedInputs := listInputs(b)
//	require.Equal(t, common.expandedInputs, expandedInputs, "expandedInputs")
//}
//
//// start a block with atoms then remove those atoms and ensure the block is "empty"
//func TestRemoveAtoms(t *testing.T) {
//	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.expandedInputs))
//	expandedInputs := listInputs(b)
//	require.Equal(t, common.expandedInputs, expandedInputs)
//	mutant.RemoveAtoms(b)
//	clippedInputs := listInputs(b)
//	require.Equal(t, common.inputs, clippedInputs)
//}
//
