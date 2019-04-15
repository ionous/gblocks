package block

// iterate over all blocks stacked in this input
// returns true if the callback function doesnt exit early
func VisitStack(in Input, cb func(b Shape) (keepGoing bool)) (exhausted bool) {
	earlyOut := false
	// get the input's connection information
	if c := in.Connection(); c != nil {
		// for every block connected to the input...
		for b := c.TargetBlock(); b != nil; {
			if !cb(b) {
				earlyOut = true
				break
			}

			// move to the next
			if c := b.NextConnection(); c != nil {
				b = c.TargetBlock()
			} else {
				break
			}
		}
	}
	return !earlyOut
}
