package decor

type Fn func(*Context) string

type Context interface {
	Parent() Context
	IsConnected() bool
	HasPrev() bool
	HasNext() bool
}
