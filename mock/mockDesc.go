package mock

import (
	"strings"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
)

func MakeDesc(blockType string, inputs []string) block.Dict {
	var args block.Args
	for _, in := range inputs {
		parts := strings.Split(in, ":")
		inputName, inputType := parts[0], parts[1]
		args.AddArg(block.Dict{
			option.Name: inputName,
			option.Type: inputType,
		})
	}
	return block.Dict{
		option.Type:       blockType,
		option.Message(0): args.Message(),
		option.Args(0):    args.List(),
	}
}
