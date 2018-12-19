package gblocks

import (
	"github.com/ionous/errutil"
	r "reflect"
	"strings"
)

// Create XML to record the number of user-defined input blocks.
func (b *Block) mutationToDom(ws *Workspace) *DomElement {
	var out mutationOutput
	blockData := ws.BlockData(b)
	//
	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
		in := b.Input(i)
		if m, ok := blockData.Mutation(in); ok {
			els := m.Elements()
			if n := els.Len(); n > 0 {
				elements := make([]int, n)
				nameToIndex := make(map[r.Type]int)
				var typeNames []string

				for i := 0; i < n; i++ {
					// note: we have to deref the element into its actual value before asking for its type.
					iface := els.Index(i)
					ptr := iface.Elem()
					el := ptr.Elem()
					t := el.Type()
					typeIndex := len(nameToIndex)
					if existingIndex, ok := nameToIndex[t]; ok {
						typeIndex = existingIndex
					} else {
						nameToIndex[t] = typeIndex
						typeName := toTypeName(t)
						typeNames = append(typeNames, typeName)
					}
					elements[i] = typeIndex
				}
				el := out.newMutation()
				el.SetAttribute("name", in.Name)
				el.SetAttribute("types", typeNames)
				el.SetAttribute("elements", elements)
			}
		}
	}
	return out.outer
}

func (b *Block) domToMutation(ws *Workspace, dom *DomElement) (err error) {
	data := ws.BlockData(b)
	//
	kids := dom.Children()
	for i, cnt := 0, kids.Num(); i < cnt; i++ {
		el := kids.Index(i)
		attrName := el.GetAttribute("name").String()
		if els, ok := data.Elements(b.GetInput(attrName)); ok {
			err = errutil.New("unknown field", attrName)
		} else {
			out := els.Slice(0, 0)
			types := strings.Split(el.GetAttribute("types").Value.String(), ",")
			val := el.GetAttribute("elements").Value
			for i := 0; i < val.Length(); i++ {
				index := val.Index(i).Int()
				typeName := types[index]
				if v, e := ws.reg.New(typeName); e != nil {
					err = errutil.Append(err, e)
				} else {
					out = r.Append(out, v)
				}
			}
			els.Set(out)
		}
	}
	return
}

type mutationOutput struct {
	outer *DomElement
}

type mutationGroup struct {
	*DomElement
}

func (o *mutationOutput) newMutation() *mutationGroup {
	if o.outer == nil {
		o.outer = NewDomElement("mutation")
	}
	child := NewDomElement("input")
	o.outer.AppendChild(child)
	return &mutationGroup{&DomElement{Object: child.Object}}
}
