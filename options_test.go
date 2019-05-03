package gblocks

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/tin"
	"github.com/stretchr/testify/require"
)

// TestBlockOptions - dependency pool generates constaints
func TestBlockOptions(t *testing.T) {
	type Description struct {
		Msg          block.Option `message0:"The description of %1 is"`
		InputsInline block.Option `inputsInline:"true"`
		Input        interface{}
	}
	if ti, e := tin.TermBlock.PtrInfo((*Description)(nil)); e != nil {
		t.Fatal(e)
	} else {
		// when all types match, we should get an unlimited next statement
		var m Maker
		desc, e := m.makeDescByType(ti, nil, nil)
		require.NoError(t, e)
		expected := block.Dict{
			"message0":     "The description of %1 is",
			"type":         "description",
			"tooltip":      "description",
			"inputsInline": true,
			"args0": []block.Dict{
				block.Dict{
					"name": "INPUT",
					"type": "input_value",
				},
			},
		}
		require.Equal(t, expected, desc)
	}
}
