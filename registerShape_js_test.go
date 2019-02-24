package gblocks

// import (
// 	"github.com/ionous/gblocks/block"
// 	"github.com/ionous/gblocks/blockly"
// .	"github.com/ionous/gblocks/gtest"
// 	"github.com/ionous/gblocks/inspect"
// 	"github.com/kr/pretty"
// 	"github.com/stretchr/testify/require"
// 	r "reflect"
// 	"strconv"
// 	"testing"
// )

// // TestShapeRegistration using a shape contain mutations ( "MutableBlock" )
// func TestRegisterShape(t *testing.T) {
// 	expectedDesc := block.Dict{
// 		"type":     block.Type("mutable_block"),
// 		"message0": "%1 %2 %3",
// 		"output":   block.Type("mutable_block"),
// 		"args0": []block.Dict{
// 			{
// 				"name":  "INPUT",
// 				"type":  "input_value",
// 				"check": []block.Type{"mutable_block"},
// 			},
// 			{
// 				"name":     "MUTANT",
// 				"type":     "input_dummy",
// 				"mutation": block.Type("test_mutation"),
// 			},
// 			{
// 				"name": "FIELD",
// 				"type": "field_input",
// 				"text": "field",
// 			},
// 		},
// 	}
// 	var tp inspect.TypePool
// 	blockDesc, e := tp.BuildDesc(r.TypeOf((*MutableBlock)(nil)), nil)
// 	require.NoError(t, e)
// 	if v := pretty.Diff(blockDesc, expectedDesc); len(v) != 0 {
// 		t.Log(pretty.Sprint(blockDesc))
// 		t.Fatal(v)
// 	}
// }
