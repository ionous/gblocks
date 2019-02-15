package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/named"
	"strconv"
	"strings"
)

// InputMutation --
// b/c we can't easily extend the existing blockly input types
// we hang optional mutation data off of dummy inputs.
type InputMutation struct {
	*js.Object
	input        *js.Object `js:"input"` // *Blockly.Input containing this mutation
	MutationName named.Type `js:"name"`  // name of the mutation, ex. TestMutation
	atoms        *js.Object `js:"atoms"` // []
	TotalInputs  int        `js:"totalInputs"`
}

type Atom struct {
	*js.Object
	Type      named.Type `js:"type"`
	NumInputs int        `js:"totalInputs"`
}

func NewInputMutation(in *Input, name named.Type) *InputMutation {
	m := &InputMutation{Object: new(js.Object)}
	m.input = in.Object
	m.MutationName = name
	m.ResetAtoms()
	return m
}

// Input - return the mutation that owns this.
func (m *InputMutation) Input() *Input {
	return &Input{Object: m.input}
}

// Reset - delete all tracking info about inputs drawn by this mutation.
// Returns the previous total inputs.
func (m *InputMutation) ResetAtoms() (ret int) {
	ret, m.TotalInputs = m.TotalInputs, 0
	m.atoms = js.Global.Get("Array").New()
	return
}

// Path - return a unique name for the atom: "MUTANT/i/"
// the members of the Mutation struct are MUTANT/0/memberName
// the first atom connected via the struct's NextStatement becomes MUTANT/1/memberName, and so on.
func (m *InputMutation) Path(depth int) string {
	itemName := m.Input().Name
	return strings.Join([]string{itemName.String(), strconv.Itoa(depth), ""}, "/")
}

// NumAtoms - number of sub-blocks used by this mutation.
// there is a one-to-one correspondence between atoms in a workspace block, and blocks in a mutation ui.
func (m *InputMutation) NumAtoms() (ret int) {
	if m.atoms != nil && m.atoms.Bool() {
		ret = m.atoms.Length()
	}
	return
}

// Atom - return a single element of the mutation.
func (m *InputMutation) Atom(i int) (ret *Atom) {
	if obj := m.atoms.Index(i); obj != nil && obj.Bool() {
		ret = &Atom{Object: obj}
	}
	return
}

// AddAtom - some number of contiguous inputs (already added to the parent block).
func (m *InputMutation) addAtom(reg *Registry, atomType named.Type) (ret int, err error) {
	// find the atom type in order to generate inputs
	if rtype, exists := reg.types[atomType]; !exists {
		err = errutil.New("atom not registered", atomType)
	} else {
		in := m.Input()
		b := in.Block()
		// the atom inputs will be placed directly after this input
		if _, mindex := b.InputByName(in.Name); mindex < 0 {
			err = errutil.New("input missing from owner block", in.Name)
		} else {
			atomIndex := m.NumAtoms()
			// expansion of atom into blockly inputs, etc.
			if args, e := reg.buildArgs(rtype, m.Path(atomIndex+1)); e != nil {
				err = e
			} else {
				// generate new inputs from the atom
				oldLen := b.NumInputs()
				b.interpolate(args.message(), args.list)
				newLen := b.NumInputs()
				numInputs := newLen - oldLen
				// reorder inputs so that the atom's inputs follow the mutation's inputs.
				if numInputs > 0 {
					// record the desired order
					scratch := make([]*Input, 0, newLen)
					// there are three sections
					// 1. up-to-and-including the mutation input
					// 2. the atom's added inputs ( which were appened to the input list )
					// 3. the inputs originally following the mutation input
					end := mindex + m.TotalInputs + 1
					for _, rng := range [][]int{{0, end}, {oldLen, newLen}, {end, oldLen}} {
						for i, last := rng[0], rng[1]; i < last; i++ {
							scratch = append(scratch, b.Input(i))
						}
					}
					// rewrite the input order
					for i, in := range scratch {
						b.setInput(i, in)
					}
				}
				// record the atom
				a := &Atom{Object: new(js.Object)}
				a.Type = atomType
				a.NumInputs = numInputs
				m.atoms.SetIndex(atomIndex, a)
				m.TotalInputs += numInputs
				ret = numInputs
			}
		}
	}
	return
}
