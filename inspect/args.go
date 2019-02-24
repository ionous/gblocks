package inspect

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
	r "reflect"
	"strconv"
	"strings"
)

// ArgumentBuilder - accumulate the inputs and fields of blocks.
type ArgumentBuilder struct {
	enums EnumPairs
	deps  DependencyPool
	msgs  []string
	list  []block.Dict
	block block.Dict
}

func NewArgs(block block.Dict, enums EnumPairs, deps DependencyPool) *ArgumentBuilder {
	return &ArgumentBuilder{block: block, enums: enums, deps: deps}
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
func (a *ArgumentBuilder) AddDesc(argDesc block.Dict) {
	a.list = append(a.list, argDesc)
	a.msgs = append(a.msgs, "%"+strconv.Itoa(len(a.list)))
}

var linkOptions = map[Class]string{NextLink: option.Next, PrevLink: option.Prev}

func (a *ArgumentBuilder) AddMembers(parent block.Item, rtype r.Type) (err error) {
	VisitItems(rtype, func(item *Item, e error) bool {
		if e != nil {
			err = errutil.Append(err, e)
		} else if e := a.AddItem(parent, item); e != nil {
			err = errutil.Append(err, e)
		} else if a.block != nil {
			if opt, ok := linkOptions[item.Class]; ok {
				if types, ok := a.deps.GetConstraints(item.Type); !ok {
					err = errutil.Append(err, errutil.New("link has no matching types", item))
				} else {
					a.block.Insert(opt, types)
				}
			}
		}
		return true // keepGoing
	})
	return
}

func (a *ArgumentBuilder) AddItem(parent block.Item, it *Item) (err error) {
	// if the field has decoration; add a placeholder label for it.
	if decoration, ok := it.Options[option.Decor].(string); ok {
		// fix? validate dc is a valid decoration?
		a.AddDesc(block.Dict{
			option.Name:  block.ItemFromString(block.ItemDecor).Push(it.Name),
			option.Type:  block.LabelField,
			option.Class: decoration,
			option.Text:  "",
		})
	}

	// we dont add args ( other than decorations ) for next/prev links
	if cls := it.Class; !cls.IsLink() {
		itemDesc := it.Options.Copy()
		if e := it.ItemDesc(parent, itemDesc, a.enums.GetPairs(it.Type)); e != nil {
			err = e
		} else if cls.Connects() {
			if types, ok := a.deps.GetConstraints(it.Type); ok {
				itemDesc.Insert(option.Check, types)
			}
		}
		a.AddDesc(itemDesc)
	}
	return
}
