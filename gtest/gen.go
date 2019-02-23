package gtest

import "strconv"

type orderedGenerator struct {
	name string
	i    int
}

func (o *orderedGenerator) NewId() string {
	o.i++
	return o.name + strconv.Itoa(o.i)
}
