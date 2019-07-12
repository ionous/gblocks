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
	// format args are one indexed
	args.list = append(args.list, desc)
	oneIndexed := strconv.Itoa(len(args.list))
	args.msgs = append(args.msgs, "%"+oneIndexed)
}
