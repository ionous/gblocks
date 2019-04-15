package mock

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
)

func CreateBlock(id string, desc block.Dict) (ret *MockBlock) {
	typeName := "unknown"
	if n, ok := desc[option.Type]; ok {
		typeName = n.(string)
	}
	b := &MockBlock{Id: id, Type: typeName, Flags: make(Flags)}
	if args, ok := desc[option.Args(0)]; ok {
		args := args.([]block.Dict)
		for _, arg := range args {
			itemName := arg[option.Name].(string)
			itemType := arg[option.Type].(string) // value, dummy, etc.
			var c *MockConnection
			switch itemType {
			case block.ValueInput, block.StatementInput:
				c = &MockConnection{Name: itemName, Source: b}
			}

			in := &MockInput{
				Name: itemName,
				Type: itemType,
				Next: c,
			}
			b.Inputs = append(b.Inputs, in)
		}
	}
	if _, ok := desc[option.Next]; ok {
		b.Next = &MockConnection{Name: option.Next, Source: b}
	}
	if _, ok := desc[option.Prev]; ok {
		b.Prev = &MockConnection{Name: option.Prev, Source: b}
	}
	ret = b
	return
}
