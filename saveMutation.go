package gblocks

import (
	"github.com/ionous/errutil"
	r "reflect"
	"strings"
)

// Create XML to record the number of user-defined input blocks.
func (b *Block) mutationToDom(ws *Workspace) *DomElement {
	out := NewDomElement("mutation")
	if cnt := b.NumInputs(); cnt > 0 {
		ctx := ws.Context(b.Id)
		for i := 0; i < cnt; i++ {
			in := b.Input(i)
			if m := in.Mutation(); m != nil {
				if slice := ctx.Elem().FieldByName(in.Name.FieldName()); slice.IsValid() {
					appendMutation(in.Name.String(), slice, out)
				}
			}
		}
	}
	return out
}

func appendMutation(name string, slice r.Value, out *DomElement) {
	if cnt := slice.Len(); cnt > 0 {
		elements := make([]int, cnt)
		nameToIndex := make(map[r.Type]int)
		var typeNames []TypeName
		for i := 0; i < cnt; i++ {
			// note: we have to deref the element into its actual value before asking for its type.
			iface := slice.Index(i)
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
		el := out.AppendChild(NewDomElement("input"))
		el.SetAttribute("name", name)
		el.SetAttribute("types", typeNames)
		el.SetAttribute("elements", elements)
	}
}

func (b *Block) domToMutation(ws *Workspace, dom *DomElement) (err error) {
	ctx := ws.Context(b.Id)
	kids := dom.Children()
	for i, cnt := 0, kids.Num(); i < cnt; i++ {
		el := kids.Index(i)
		attrName := InputName(el.GetAttribute("name").String())
		if els := ctx.Elem().FieldByName(attrName.FieldName()); !els.IsValid() {
			err = errutil.New("unknown field", attrName)
		} else {
			out := els.Slice(0, 0)
			types := strings.Split(el.GetAttribute("types").Value.String(), ",")
			val := el.GetAttribute("elements").Value
			for i := 0; i < val.Length(); i++ {
				index := val.Index(i).Int()
				typeName := TypeName(types[index])
				if v, e := TheRegistry.NewData(typeName); e != nil {
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
