package tin

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/pascal"
)

type Model int

//go:generate stringer -type=Model
const (
	UnknownModel Model = iota
	TopBlock
	MidBlock
	TermBlock
)

// ptr should be pointer to a struct or an interface
func (c Model) PtrInfo(ptr interface{}) (ret *TypeInfo, err error) {
	return c.TypeInfo(r.TypeOf(ptr))
}

func (c Model) TypeInfo(ptrType r.Type) (ret *TypeInfo, err error) {
	if ptrType.Kind() != r.Ptr || ptrType.Elem().Kind() != r.Struct {
		err = errutil.New("expected pointer to struct")
	} else {
		name := pascal.ToUnderscore(ptrType.Elem().Name())
		ret = &TypeInfo{name, c, ptrType}
	}
	return
}
