package mock

import (
	r "reflect"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/enum"
)

// implements mutant.Atomizer
type MockDatabase struct {
	Atoms map[string][]MockAtom
}

type MockAtom struct {
	Name string // ex. NUM
	Type string // ex. "field_number", "input_value", etc.
}

func (db *MockDatabase) GetPairs(string) []enum.Pair {
	return nil
}

func (db *MockDatabase) GetTerms(targetType string) (block.Limits, error) {
	return block.MakeUnlimited(), nil
}

func (db *MockDatabase) GetStatements(targetType string) (block.Limits, error) {
	return block.MakeUnlimited(), nil
}

func (db *MockDatabase) GetTermsByType(r.Type) block.Limits {
	return block.MakeUnlimited()
}

func (db *MockDatabase) GetStatementsByType(r.Type) block.Limits {
	return block.MakeUnlimited()
}
