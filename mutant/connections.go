package mutant

import "github.com/ionous/gblocks/block"

// FIX: storing fields
type Store struct {
	connections []block.Connection
}

func (s *Store) SaveConnection(in block.Input) {
	if c := in.Connection(); c != nil {
		var tgt block.Connection
		if c.IsConnected() {
			tgt = c.TargetConnection()
		}
		s.connections = append(s.connections, tgt)
	}
}

// FIX: this isnt going to be right for fields, etc.
func (s *Store) Restore(b block.Shape, start, numInputs int) {
	for i := 0; i < numInputs && i < len(s.connections); i++ {
		wsinput := b.Input(start + i)
		c := s.connections[i]
		reconnect(b, wsinput, c)
	}
}

// ported from blockly
func reconnect(b block.Shape, in block.Input, connectionChild block.Connection) (okay bool) {
	if isConnectionValid(connectionChild) {
		if connectionParent := in.Connection(); connectionParent != nil {
			currentParent := connectionChild.TargetBlock()
			if currentParent == nil || currentParent == b {
				if connectionParent.TargetConnection() != connectionChild {
					if connectionParent.IsConnected() {
						connectionParent.Disconnect()
					}
					connectionParent.Connect(connectionChild)
					okay = true
				}
			}
		}
	}
	return
}

func isConnectionValid(c block.Connection) (okay bool) {
	if c != nil {
		connectionOwner := c.SourceBlock()
		if connectionOwner != nil && connectionOwner.HasWorkspace() {
			okay = true
		}
	}
	return
}
