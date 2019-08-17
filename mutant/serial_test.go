package mutant_test

//

//
//	"github.com/ionous/gblocks/block"
//	"github.com/ionous/gblocks/dom"
//	"github.com/ionous/gblocks/mock"
//	"github.com/ionous/gblocks/mutant"
//	"github.com/kr/pretty"
//	"github.com/stretchr/testify/require"

//
// func TestSerialization(t *testing.T) {
// 	data := dom.BlockMutation{dom.Mutations{
// 		&dom.Mutation{
// 			Input: "M1",
// 			Atoms: dom.Atoms{
// 				[]string{"A"},
// 			},
// 		},
// 		&dom.Mutation{¥
// 			Input: "M3",
// 			Atoms: dom.Atoms{
// 				[]string{"B", "C"},
// 			},
// 		},
// 	}}
// 	expandedInputs := []string{
// 		"M1:input_dummy",
// 		/**/ block.Scope("a", "M1", "0", "TERM:input_value"),
// 		"M2:input_dummy",
// 		"M3:input_dummy",
// 		/**/ block.Scope("a", "M3", "0", "NUM:field_number"),
// 		/**/ block.Scope("a", "M3", "1", "TEXT:field_input"),
// 		/**/ block.Scope("a", "M3", "1", "STATE:input_statement"),
// 	}
// 	expectedOutput := `` +
// 		`<mutation>` +
// 		/**/ `<pin name="M1">` +
// 		/* */ `<atom type="A"></atom>` +
// 		/**/ `</pin>` +
// 		/**/ `<pin name="M3">` +
// 		/* */ `<atom type="B"></atom>` +
// 		  `<atom type="C"></atom>` +
// 		/**/ `</pin>` +
// 		`</mutation>`

// 	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
// 	arch := mock.NewMutations(common.inputs, common.quarks)

// 	db := &mock.MockDatabase{common.atomProducts}
// 	blocks := mutant.NewMutatedBlocks()

// 	mb := blocks.CreateMutatedBlock(b, arch, db)
// 	e := mb.LoadMutation(&data)
// 	require.NoError(t, e)
// 	expanded := listInputs(b)
// 	require.Equal(t, expandedInputs, expanded)

// 	// save the in-memory data, makes sure it matches the original input
// 	serial := mb.SaveMutation()
// 	if v := pretty.Diff(data, serial); len(v) != 0 {
// 		t.Log(pretty.Sprint(serial))
// 		t.Fatal(v)
// 	}
// 	if str, e := serial.MarshalMutation(); e != nil {
// 		t.Fatal(e)
// 	} else {
// 		require.Equal(t, expectedOutput, str)
// 	}
// }
