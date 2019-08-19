package block

type Project interface {
	IsBlockRegistered(blockType string) bool
	RegisterBlock(blockType string, desc Dict) error
	RegisterMutator(name string, mutator Mutator) error
	GenerateUniqueName() string
}

type Workspace interface {
	NewBlock(blockType string) (Shape, error)
	NewBlockWithId(blockId, blockType string) (Shape, error)
	OnDelete(OnDelete)
}

type OnDelete interface{ OnDelete(blockId string) }

// input access
type Inputs interface {
	NumInputs() int
	Input(int) Input
	InputByName(string) (Input, int)
	RemoveInput(string)
	SetInput(int, Input)
	Interpolate(msg string, args []Dict)
}

//go:generate stringer -type=Flag
type Flag int

const (
	Deletable Flag = iota
	Movable
	Shadow
	InsertionMaker
	Editable
	Collapsed
	Enabled      // blockly uses isEnabled, setDisabled
	InputsInline // getInputsInline, setInputsInline
)

type Flags interface {
	GetFlag(Flag) bool
	SetFlag(Flag, bool)
}

// an instance of a block
type Shape interface {
	BlockId() string
	BlockType() string

	// indicates if the block has been disposed ( deleted by the user. )
	HasWorkspace() bool
	BlockWorkspace() Workspace
	InitSvg() // wish this happened on demand
	Dispose()

	// note: blockly puts all "fields" ( editable and non-editable items ) inside of inputs.
	Inputs
	Flags

	// for statement blocks, the following block in a stack of blocks
	PreviousConnection() Connection
	NextConnection() Connection
}

type Input interface {
	InputName() string
	InputType() string
	SetInvisible()
	Connection() Connection
}

type Connection interface {
	SourceBlock() Shape
	TargetBlock() Shape
	TargetConnection() Connection
	IsConnected() bool
	Connect(Connection)
	Disconnect()
}

type Mutator interface {
	// list of block names to appear in the mui palette
	Quarks() []string
	// callback after a block instance containing a mutator has been created
	PostMixin(main Shape) error
	//
	MutationToDom(main Shape) (dom string, err error)
	DomToMutation(main Shape, dom string) error
	Decompose(main Shape, popup Workspace) (mui Shape, err error)
	Compose(main, mui Shape) error
	SaveConnections(main, mui Shape) error
}
