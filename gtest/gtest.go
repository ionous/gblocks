package gtest

type Enum int

const (
	DefaultChoice Enum = iota
	AlternativeChoice
)

func (i Enum) String() (ret string) {
	switch i {
	case DefaultChoice:
		ret = "DefaultChoice"
	case AlternativeChoice:
		ret = "AlternativeChoice"
	}
	return
}

type EnumStatement struct {
	Enum
}

// // see also: blockly/tests/jsunit/block_test.js setupStackBlocks.
type StackBlock struct {
	PreviousStatement,
	NextStatement interface{}
}

// see also: blockly/tests/jsunit/block_test.js setUpRowBlocks.
type RowBlock struct {
	// implement an interface that generates an output
	Input interface {
		Output() interface{}
	}
}

// Output - implement a generic output
func (b *RowBlock) Output() interface{} {
	return b
}

type FieldBlock struct {
	Number float32
}

type MutableBlock struct {
	Input  *MutableBlock
	Mutant TestMutation
	Field  string
}

type TestMutation struct {
	NextStatement NextAtom
}

// Output - implement a generic output
func (n *MutableBlock) Output() *MutableBlock {
	return n
}

// NextAtom
type NextAtom interface {
	NextAtom() NextAtom
}

type AtomTest struct {
	AtomInput     *MutableBlock
	NextStatement NextAtom
}

func (a *AtomTest) NextAtom() NextAtom { return a.NextStatement }

type AtomAltTest struct {
	AtomField     string
	NextStatement NextAtom
}

func (a *AtomAltTest) NextAtom() NextAtom { return a.NextStatement }
