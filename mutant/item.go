package mutant

import (
	r "reflect"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/enum"
)

// interface for creating, listening to atom expansion
type Atomizer interface {
	// return choices for the passed (pre-registered) enum.
	GetPairs(string) []enum.Pair
	// lookup the limits of the passed block type
	GetTerms(string) (block.Limits, error)
	// lookup the limits of the passed block type
	GetStatements(string) (block.Limits, error)

	// a work in progress: interfaces arent registered to blockly
	// so they cant be found by type name; we need typename tho for loading from xml.
	GetTermsByType(r.Type) block.Limits
	GetStatementsByType(r.Type) block.Limits
}
