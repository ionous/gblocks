package mutant

import "github.com/ionous/gblocks/block"

func MutatorName(blockType string) string {
	return block.Scope("mutates", blockType)
}

func ContainerName(blockType string) string {
	return block.Scope("mui", blockType)
}

// a single input's mutation
// ex. tin.Mutable
type InMutation interface {
	// name of the original mutation type ( in blockly friendly format )
	Name() string
	// blocks which can directly attach to the mutation's input
	Limits() block.Limits
	// returns an optional fixed-in-place first quark
	FirstBlock() (Quark, bool)
	// iterator over all quarks (mui bocks) which can be used for this mutaiton
	Quarks(paletteOnly bool) (Quark, bool)
}
