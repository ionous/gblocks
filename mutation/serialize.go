package mutation

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
	"strings"
)

// Serialize the workspace block's mutations, returning XML with the number and type of user-defined inputs.
// note: this is used by blockly during block creation, updates to generate block.xml.
// this function cannot return nil or blockly will throw an error.
func (mb *mutableBlock) mutationToDom() (ret *dom.Element, err error) {
	out := dom.NewElement("mutation")
	for i, cnt := 0, mb.block.NumInputs(); i < cnt; i++ {
		if mi, ok := mb.mutationByIndex(i); ok {
			// if there are atoms, create a node for the data.
			if cnt := len(mi.atoms); cnt > 0 {
				atoms := dom.NewElement("atoms", dom.Attrs{"name": mi.name.String()})
				// add the specific atoms to the outptu
				for i := 0; i < cnt; i++ {
					atom := mi.atoms[i]
					atoms.AppendChild(dom.NewElement("atom", dom.Attrs{"type": atom.typeName.String()}))
				}
				out.AppendChild(atoms)
			}
		}
	}
	ret = out
	return
}

// Deserialize mutations by expanding XML into atoms.
// The atoms generate inputs for the block, which are filled by blockly during the next stage of deserializing.
// Returns the total number of inputs added
func (mb *mutableBlock) domToMutation(dom *dom.Element) (ret int, err error) {
	// we are "reloading" the mutations; remove all dynamic inputs
	mb.removeAtoms()
	// rebuild from the dom
	mutations := dom.Children()
	for i, cnt := 0, mutations.Num(); i < cnt; i++ {
		if el := mutations.Index(i); !strings.EqualFold(el.TagName, "atoms") {
			err = errutil.Append(err, errutil.New("mutation has unexpected child", el.TagName))
		} else {
			itemName := block.Item(el.GetAttribute("name").String())
			if mi, inputIndex := mb.mutationByName(itemName); inputIndex < 0 {
				err = errutil.New("input is not a mutation", itemName)
			} else {
				atoms := el.Children()
				for i, cnt := 0, atoms.Num(); i < cnt; i++ {
					if el := atoms.Index(i); !strings.EqualFold(el.TagName, "atom") {
						err = errutil.Append(err, errutil.New("input has unexpected child", el.TagName))
					} else if atomType := block.Type(el.GetAttribute("type").String()); len(atomType) == 0 {
						err = errutil.Append(err, errutil.New("atom has no type", el.TagName))
					} else if numInputs, e := mi.addAtom(atomType); e != nil {
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
