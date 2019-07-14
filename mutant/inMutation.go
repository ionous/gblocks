package mutant

import "github.com/ionous/gblocks/block"

func MutatorName(blockType string) string {
	return block.Scope("mutates", blockType)
}

func ContainerName(blockType string) string {
	return block.Scope("mui", blockType)
}

// the description of a single input's mutation
type Mutation interface {
	// name of the original mutation type ( in blockly friendly format )
	Name() string
	// quarks which can directly attach to the mutation's input
	Limits() block.Limits
	// returns an optional fixed-in-place first quark
	FirstBlock() (Quark, bool)
	// all quarks (mui bock types) available for this mutation
	Quarks(paletteOnly bool) (Quark, bool)
}
