package gblocks

import "github.com/ionous/gblocks/named"

type Constraints struct {
	constraints   []named.Type
	hasConnection bool
}

func (c *Constraints) GetConstraints() ([]named.Type, bool) {
	return c.constraints, c.hasConnection
}

// unlike in blockly proper, a single empty typename here means "any connection", rather than "no connection".
func (c *Constraints) AddConstraint(typeName named.Type) {
	if len(typeName) > 0 {
		c.constraints = append(c.constraints, typeName)
	}
	c.hasConnection = true
}

// does the next of this and the prev of that overlap?
// func (c *Constraints) ConnectsTo(c *Constraints) bool {
// }
