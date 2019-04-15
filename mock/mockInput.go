package mock

import (
	"fmt"

	"github.com/ionous/gblocks/block"
)

type MockInput struct {
	Name string
	Type string // statement, value, dummy,
	Next *MockConnection
}

func (in *MockInput) InputName() string {
	return in.Name
}

func (in *MockInput) InputType() string {
	return in.Type
}

func (in *MockInput) SetInvisible() {}

func (in *MockInput) Connection() (ret block.Connection) {
	// ugh. avoid returning a non-nil interface to a nil value
	// https://golang.org/doc/faq#nil_error
	if in.Next != nil {
		ret = in.Next
	}
	return
}

// log for debugging
func (in *MockInput) String() string {
	return fmt.Sprintf("%s:%s", in.Name, in.Type)
}
