package gblocks

import (
	r "reflect"
	"strconv"
	"strings"
)

type argDesc struct {
	Name string
	Type r.Type
	Path string
	Dict
}

func (a *argDesc) InputName() InputName {
	return InputName(a.Name)
}

func (a *argDesc) TypeName() TypeName {
	return toTypeName(a.Type)
}

func (a *argDesc) String() string {
	return strings.Join([]string{a.Name, a.Type.String()}, ":")
}

func makeArg(f r.StructField, path string) argDesc {
	options := parseTags(string(f.Tag))
	name := path + pascalToCaps(f.Name)
	options.Insert(opt_name, name)
	return argDesc{name, f.Type, path, options}
}

type argsOut struct {
	msgs      []string // helper for building message strings
	list      []Dict
	mutations []*mutationInput
}

// send the current argument to the list of all args
func (a *argsOut) addArg(argDesc argDesc) {
	a.list = append(a.list, argDesc.Dict)
	a.msgs = append(a.msgs, "%"+strconv.Itoa(len(a.list)))
}

func (a *argsOut) addMutation(argDesc argDesc, mui *mutationInput) {
	a.addArg(argDesc)
	a.mutations = append(a.mutations, mui)
}

func (a *argsOut) message() string {
	return strings.Join(a.msgs, " ")
}
