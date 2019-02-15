package decor

type Fn func(Context) string

// ContextType - reflects the hierarchy of go language descriptions used for blockly.
type ContextType int

const (
	BlockContext    ContextType = iota
	MutationContext             // parent is the block, future: siblings are other mutations
	AtomContext                 // parent is the mutation, sibilings are other atoms
	ItemContext                 // parent can be any other type
)

// fix:
type Context interface {
	Parent() Context
	ContextType() ContextType
	String() string
	IsConnected() bool
	HasPrev() bool
	HasNext() bool
}
