package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

// InputMutation --
// b/c we can't easily extend the existing blockly input types
// we hang optional mutation data off of dummy inputs.
// each mutation generally represent some sort of array --
// each element in the array generates some number of inputs.
// alt: completely refactor blockly to compose blocks out of "atoms";
//      merging the concept of arg#/message#, and mutations.
type InputMutation struct {
	*js.Object
	name        string     `js:"name"`        // name of the mutation, ex. TestMutation
	totalInputs int        `js:"totalInputs"` // total number of inputs generated by this mutation
	atoms       *js.Object `js:"atoms"`
}

type Atom struct {
	*js.Object
	connections *js.Object `js:"connections"` // []*Connection
}

// Reset - delete all tracking info about inputs drawn by this mutation.
// Returns the previous total inputs.
func (m *InputMutation) Reset() (ret int) {
	ret, m.totalInputs = m.totalInputs, 0
	m.atoms = js.MakeWrapper([]Atom{})
	return
}

// TotalIputs - total number of inputs used by all of the atoms generated by this mutation.
func (m *InputMutation) TotalInputs() int {
	return m.totalInputs
}

// AddAtom - some number of contiguous inputs (already added to the parent block).
func (m *InputMutation) AddAtom(numInputs int) {
	connections := js.MakeWrapper(make([]*Connection, numInputs))
	atom := js.MakeWrapper(&Atom{connections: connections})
	m.atoms.SetIndex(m.atoms.Length(), atom)
	m.totalInputs += numInputs
}

// Atoms - number of sub-blocks used by this mutation.
func (m *InputMutation) Atoms() (ret int) {
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

// Connections - the number of contiguous inputs used by this atom.
func (a *Atom) Connections() int {
	return a.connections.Length()
}

func (a *Atom) SaveConnection(i int, in *Input) {
	var target *Connection
	//blockInput.Connection.TargetConnection
	if in != nil {
		if c := in.Connection(); c != nil {
			target = c.TargetConnection()
		}
	}
	a.connections.SetIndex(i, target)
}

func (a *Atom) Connection(i int) (ret *Connection) {
	if obj := a.connections.Index(i); obj != nil && obj.Bool() {
		ret = &Connection{Object: obj}
	}
	return
}
