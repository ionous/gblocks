package gblocks

import (
	"github.com/ionous/errutil"
	"strings"
)

// Serialize mutations returning XML describing the number and type of user-defined inputs.
// note: this is used by blockly during block creation, updates to generate block.xml.
// this function cannot return nil or blockly will throw an error.
func (b *Block) mutationToDom() *XmlElement {
	out := NewXmlElement("mutation")
	for i, cnt := 0, b.NumInputs(); i < cnt; i++ {
		in := b.Input(i)
		if m := in.Mutation(); m != nil {
			if child := m.mutationToDom(); child != nil {
				out.AppendChild(child)
			}
		}
	}
	return out
}

// shared with toolbox creation
func (m *InputMutation) mutationToDom() (ret *XmlElement) {
	in := m.Input()
	if cnt := m.NumAtoms(); cnt > 0 {
		out := NewXmlElement("atoms", Attrs{"name": in.Name.String()})
		for i := 0; i < cnt; i++ {
			atom := m.Atom(i)
			out.AppendChild(NewXmlElement("atom", Attrs{"type": atom.Type.String()}))
		}
		ret = out
	}
	return
}

// Deserialize mutations by expanding XML into atoms.
// returns the total number of inputs added
func (b *Block) domToMutation(reg *Registry, dom *XmlElement) (ret int, err error) {
	// we are "reloading" the mutations; remove all dynamic inputs
	b.removeAtoms()

	kids := dom.Children()
	for i, cnt := 0, kids.Num(); i < cnt; i++ {
		if el := kids.Index(i); !strings.EqualFold(el.TagName, "atoms") {
			err = errutil.Append(err, errutil.New("mutation has unexpected child", el.TagName))
		} else {
			inputName := InputName(el.GetAttribute("name").String())
			if in, index := b.InputByName(inputName); index < 0 {
				err = errutil.New("unknown input", inputName)
			} else if m := in.Mutation(); m == nil {
				err = errutil.New("input is not a mutation", inputName)
			} else {
				kids := el.Children()
				for i, cnt := 0, kids.Num(); i < cnt; i++ {
					if el := kids.Index(i); !strings.EqualFold(el.TagName, "atom") {
						err = errutil.Append(err, errutil.New("input has unexpected child", el.TagName))
					} else if atomType := TypeName(el.GetAttribute("type").String()); len(atomType) == 0 {
						err = errutil.Append(err, errutil.New("atom has no type", el.TagName))
					} else if numInputs, e := m.addAtom(reg, atomType); e != nil {
						err = errutil.Append(err, e)
					} else {
						ret += numInputs
					}
				}
			}
		}
	}
	return
}
