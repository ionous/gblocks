package inspect

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
	"strconv"
	"strings"
)

// describes the inputs and fields of blocks
type ArgumentBuilder struct {
	parent block.Item // empty if a top level block
	enums  EnumPairs
	deps   DependencyPool
	msgs   []string
	list   []block.Dict
}

func NewArgs(parent block.Item, enums EnumPairs, deps DependencyPool) *ArgumentBuilder {
	return &ArgumentBuilder{parent: parent, enums: enums, deps: deps}
}

func (a *ArgumentBuilder) Len() int {
	return len(a.msgs)
}

func (a *ArgumentBuilder) Message() string {
	return strings.Join(a.msgs, " ")
}

func (a *ArgumentBuilder) List() []block.Dict {
	return a.list
}

// send the current argument to the list of all args
func (a *ArgumentBuilder) addArg(argDesc block.Dict) {
	a.list = append(a.list, argDesc)
	a.msgs = append(a.msgs, "%"+strconv.Itoa(len(a.list)))
}

func (a *ArgumentBuilder) AddItem(it *Item) (err error) {
	// if the field has decoration; add a placeholder label for it.
	if decoration, ok := it.Options[option.Decor].(string); ok {
		// fix? validate dc is a valid decoration?
		a.addArg(block.Dict{
			option.Name:  block.ItemFromString(block.ItemDecor).Push(it.Name),
			option.Type:  block.LabelField,
			option.Class: decoration,
			option.Text:  "",
		})
	}

	// we dont add args ( other than decorations ) for next/prev links
	if cls := it.Class; !cls.IsLink() {
		itemDesc := it.Options.Copy()
		if e := it.ItemDesc(a.parent, itemDesc, a.enums.GetPairs(it.Type)); e != nil {
			err = e
		} else if cls.Connects() {
			if types, ok := a.deps.GetConstraints(it.Type); ok {
				itemDesc.Insert(option.Check, types)
			}
		}
		a.addArg(itemDesc)
	}
	return
}
