package test

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

type CheckNext struct {
	NextStatement *CheckNext
}

// // see also: blockly/tests/jsunit/block_test.js setupStackBlocks.
type StackBlock struct {
	NextStatement interface{}
}

// see also: blockly/tests/jsunit/block_test.js setUpRowBlocks.
type RowBlock struct {
	Input *RowBlock
}

// Output - implement a generic output
func (b *RowBlock) Output() interface{} {
	return b
}

type FieldBlock struct {
	Number float32
}

// Output - implement a generic output
func (b *FieldBlock) Output() interface{} {
	return b
}

type InputBlock struct {
	Value     *RowBlock      `input:"value"`
	Statement *InputBlock    `input:"statement"`
	Mutation  *BlockMutation `input:"mutation"`
}

type MutableBlock struct {
	Input  *MutableBlock
	Mutant *BlockMutation `input:"mutation"`
	Field  string
}

type BlockMutation struct {
	ExtraField    int
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

// an atom that contains an input with an interface
type AtomWithInterface struct {
	Input PinInterface
}

type PinInterface interface{ Output() PinInterface }

func (*AtomWithInterface) NextAtom() NextAtom { return nil }

// can plug into a PinInterface Input
type InterfacingTerm struct{}

func (t *InterfacingTerm) Output() PinInterface { return t }
