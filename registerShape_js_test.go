package gblocks

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/blockly"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	r "reflect"
	"strconv"
	"testing"
)

// TestShapeRegistration using a shape contain mutations ( "MutableBlock" )
func TestShapeRegistration(t *testing.T) {
	expectedDesc := block.Dict{
		"type":     block.Type("mutable_block"),
		"message0": "%1 %2 %3",
		"output":   block.Type("mutable_block"),
		"args0": []block.Dict{
			{
				"name":  "INPUT",
				"type":  "input_value",
				"check": []block.Type{"mutable_block"},
			},
			{
				"name":     "MUTANT",
				"type":     "input_dummy",
				"mutation": block.Type("test_mutation"),
			},
			{
				"name": "FIELD",
				"type": "field_input",
				"text": "field",
			},
		},
	}
	var reg Registry
	blockDesc, e := reg.testRegister(r.TypeOf((*MutableBlock)(nil)))
	require.NoError(t, e)
	if v := pretty.Diff(blockDesc, expectedDesc); len(v) != 0 {
		t.Log(pretty.Sprint(blockDesc))
		t.Fatal(v)
	}
}

// helper that registers a reflected type rather than a pointer.
func (reg *Registry) testRegister(ptrType r.Type) (ret block.Dict, err error) {
	blockDesc := make(block.Dict)
	typeName := block.TypeFromStruct(ptrType.Elem())
	if e := reg.registerType(typeName, ptrType, blockDesc); e != nil {
		err = e
	} else {
		ret = blockDesc
	}
	return
}
