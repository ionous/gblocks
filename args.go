package gblocks

import (
	"strconv"
	"strings"
)

type argsOut struct {
	msgs      []string // helper for building message strings
	list      []Dict
	mutations []*mutationDesc
}

// send the current argument to the list of all args
func (a *argsOut) addArg(argDesc Dict) {
	a.list = append(a.list, argDesc)
	a.msgs = append(a.msgs, "%"+strconv.Itoa(len(a.list)))
}

func (a *argsOut) addMutation(mui *mutationDesc) {
	a.mutations = append(a.mutations, mui)
}

func (a *argsOut) message() string {
	return strings.Join(a.msgs, " ")
}
