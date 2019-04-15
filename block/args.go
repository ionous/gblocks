package block

import (
	"strconv"
	"strings"
)

// Args - accumulate the inputs and fields of blocks.
type Args struct {
	msgs []string
	list []Dict
}

func NewArgs(msg string, list ...Dict) Args {
	return Args{[]string{msg}, list}
}

func (args *Args) Message() string {
	return strings.Join(args.msgs, " ")
}

func (args *Args) List() []Dict {
	return args.list
}

func (args *Args) AddArg(desc Dict) {
	args.list = append(args.list, desc)
	args.msgs = append(args.msgs, "%"+strconv.Itoa(len(args.list)))
}
