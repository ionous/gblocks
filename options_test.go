package gblocks

import (
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/tin"
	"github.com/stretchr/testify/require"
)

// TestBlockOptions - dependency pool generates constraints
func TestBlockOptions(t *testing.T) {
	type Description struct {
		Msg          block.Option `message0:"The description of %1 is"`
		InputsInline block.Option `inputsInline:"true"`
		Input        interface{}
	}
	if ti, e := tin.TermBlock.PtrInfo((*Description)(nil)); e != nil {
		t.Fatal(e)
	} else {
		// when there are no types to match we should get "offlimits"
		var m Maker
		desc, e := m.makeDescByType(ti, nil, nil)
		require.NoError(t, e)
		expected := block.Dict{
			"message0": "The description of %1 is",
			"type":     "description",
			"tooltip":  "description",
			// MOD: output is now produced for all TermBlocks
			// ( tin.TermBlock.PtrInfo above )
			"output":       []string{},
			"inputsInline": true,
			"args0": []block.Dict{
				block.Dict{
					"name":  "INPUT",
					"type":  "input_value",
					"check": []string{},
				},
			},
		}
		require.Equal(t, expected, desc)
	}
}
