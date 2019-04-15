package mutant

import "github.com/ionous/gblocks/block"

// Quark - a block in a mutation ui palette.
// supports iteration to the next quark in the palette.
type Quark interface {
	// type name without including owner mutation
	Name() string
	// type name scoped to the owner mutation
	BlockType() string
	// mui block display name
	Label() string
	// the quark (names) this can connect to
	LimitsOfNext() block.Limits
	// expand this quark into a bundle of workspace items (fields and inputs)
	// item names are prefixed with "scope" to make them unique.
	Atomize(scope string, atomizer Atomizer) (block.Args, error)
	// next quark in the mutation
	NextQuark() (Quark, bool)
}
