package block

import "fmt"

const hw, hh = 11, 43
func (self Subtype) GetContactType(hx, hy, bx, by, bw, bh uint16, flags Flags) ContactType {
	// NOTICE: this is only called if there's actual contact, otherwise
	//         this is never invoked, so we can skip redundant checks
	switch self {
	case SubtypeBlock:
		if hx > bx + bw { panic("unexpected") } // return ContactNone
		if hx + hw < bx { panic("unexpected") } // return ContactNone

		// side cases (symmetrical)
		if hx == bx + bw { // right side
			if flags.IsRightOriented() { return ContactNone }
			if hy + 17 >= by && hy + 39 <= by + bh { return ContactWallStick }
			return ContactSideBlock
		} else if hx + hw == bx { // left side
			if flags.IsLeftOriented() { return ContactNone }
			if hy + 17 >= by && hy + 39 <= by + bh { return ContactWallStick }
			return ContactSideBlock
		}
		
		// head to the block
		if hy == by + bh { return ContactClonk }
		
		// check top ground positionings
		if hy + hh > by { return ContactSlipIntoFall }
		
		leftNormalLimit  := bx
		rightNormalLimit := bx + bw - hw
		if hx >= leftNormalLimit && hx <= rightNormalLimit {
			return ContactGround
		}

		// check tight positionings
		rightSide := (hx + hw/2 >= bx + bw/2)
		if rightSide {
			if flags.IsRightOriented() { // right side forwards edge
				if hx >= rightNormalLimit + 5 { return ContactSlipIntoFall }
				if hx == rightNormalLimit + 4 { return ContactTightFront2 }
				if hx == rightNormalLimit + 3 { return ContactTightFront2 }
				if hx == rightNormalLimit + 2 { return ContactTightFront1 }
				if hx == rightNormalLimit + 1 { return ContactTightFront1 }
				panic("missing case")
			} else { // right side backwards edge
				if hx >= rightNormalLimit + 5 { return ContactSlipIntoFall }
				if hx == rightNormalLimit + 4 { return ContactTightBack2 }
				if hx == rightNormalLimit + 3 { return ContactTightBack2 }
				if hx == rightNormalLimit + 2 { return ContactTightBack1 }
				if hx == rightNormalLimit + 1 { return ContactTightBack1 }
				panic("missing case")
			}
		} else {
			// ---- symmetrical to the above ----
			if flags.IsLeftOriented() { // left side forwards edge
				if hx <= leftNormalLimit - 5 { return ContactSlipIntoFall }
				if hx == leftNormalLimit - 4 { return ContactTightFront2 }
				if hx == leftNormalLimit - 3 { return ContactTightFront2 }
				if hx == leftNormalLimit - 2 { return ContactTightFront1 }
				if hx == leftNormalLimit - 1 { return ContactTightFront1 }
				panic("missing case")
			} else { // left side backwards edge
				if hx <= leftNormalLimit - 5 { return ContactSlipIntoFall }
				if hx == leftNormalLimit - 4 { return ContactTightBack2 }
				if hx == leftNormalLimit - 3 { return ContactTightBack2 }
				if hx == leftNormalLimit - 2 { return ContactTightBack1 }
				if hx == leftNormalLimit - 1 { return ContactTightBack1 }
				panic("missing case")
			}
		}

		fmt.Printf("hx, hy, bx, by, bw, bh = %d, %d, %d, %d, %d, %d\n", hx, hy, bx, by, bw, bh)
		panic("unexpected end")
	case SubtypeThinStep:
		// side slipping
		if !flags.HasUpInertia() && hy + hh <= by + 4 {
			if hx + hw < bx + 7 { return ContactSlipIntoFall }
			if hx + 7 > bx + bw { return ContactSlipIntoFall }
		}

		// step up cases
		if hy + hh == by + 7 {
			if flags & FlagDownPressed == 0 {
				if flags.HasLeftInertia()  && hx + 4 == bx + bw { return ContactStepUp }
				if flags.HasRightInertia() && hx + hw == bx + 4 { return ContactStepUp }
			}
		}

		// ground or step down cases
		if hy + hh == by {
			if flags.HasLeftInertia()  && hx + 4 == bx { return ContactStepDown }
			if flags.HasRightInertia() && hx + hw == bx + bw + 4 { return ContactStepDown }
			if !flags.HasUpInertia() { return ContactGround }
		}

		return ContactNone
	case SubtypeThinBlock:
		if hy + hh == by {
			leftNormalLimit  := bx
			rightNormalLimit := bx + bw - hw
			if hx >= leftNormalLimit && hx <= rightNormalLimit {
				return ContactGround
			}

			// check tight positionings
			rightSide := (hx + hw/2 >= bx + bw/2)
			if rightSide {
				if flags.IsRightOriented() { // right side forwards edge
					if hx >= rightNormalLimit + 6 { return ContactSlipIntoFall }
					if hx == rightNormalLimit + 5 { return ContactTightFront2 }
					if hx == rightNormalLimit + 4 { return ContactTightFront2 }
					if hx == rightNormalLimit + 3 { return ContactTightFront1 }
					if hx == rightNormalLimit + 2 { return ContactTightFront1 }
				} else { // right side backwards edge
					if hx >= rightNormalLimit + 6 { return ContactSlipIntoFall }
					if hx == rightNormalLimit + 5 { return ContactTightBack2 }
					if hx == rightNormalLimit + 4 { return ContactTightBack2 }
					if hx == rightNormalLimit + 3 { return ContactTightBack1 }
					if hx == rightNormalLimit + 2 { return ContactTightBack1 }
				}
			} else {
				// ---- symmetrical to the above ----
				if flags.IsLeftOriented() { // left side forwards edge
					if hx <= leftNormalLimit - 6 { return ContactSlipIntoFall }
					if hx == leftNormalLimit - 5 { return ContactTightFront2 }
					if hx == leftNormalLimit - 4 { return ContactTightFront2 }
					if hx == leftNormalLimit - 3 { return ContactTightFront1 }
					if hx == leftNormalLimit - 2 { return ContactTightFront1 }
				} else { // left side backwards edge
					if hx <= leftNormalLimit - 6 { return ContactSlipIntoFall }
					if hx == leftNormalLimit - 5 { return ContactTightBack2 }
					if hx == leftNormalLimit - 4 { return ContactTightBack2 }
					if hx == leftNormalLimit - 3 { return ContactTightBack1 }
					if hx == leftNormalLimit - 2 { return ContactTightBack1 }
				}
			}

			return ContactGround
		}

		return ContactNone
	case SubtypeThinStepOnLeft: // MIX FROM SubtypeThinStep AND SubtypeThinBlock
		// side slipping
		if !flags.HasUpInertia() && hy + hh <= by + 4 {
			if hx + hw < bx + 7 { return ContactSlipIntoFall }
		}

		// step up cases
		if hy + hh == by + 7 {
			if flags & FlagDownPressed == 0 {
				if flags.HasRightInertia() && hx + hw == bx + 4 { return ContactStepUp }
			}
		}

		// step down cases, tight position or ground cases
		if hy + hh == by {
			if flags.HasLeftInertia() && hx + 4 == bx { return ContactStepDown }

			leftNormalLimit  := bx
			rightNormalLimit := bx + bw - hw
			if hx >= leftNormalLimit && hx <= rightNormalLimit {
				return ContactGround
			}

			// check tight positionings
			rightSide := (hx + hw/2 >= bx + bw/2)
			if rightSide {
				if flags.IsRightOriented() { // right side forwards edge
					if hx >= rightNormalLimit + 6 { return ContactSlipIntoFall }
					if hx == rightNormalLimit + 5 { return ContactTightFront2 }
					if hx == rightNormalLimit + 4 { return ContactTightFront2 }
					if hx == rightNormalLimit + 3 { return ContactTightFront1 }
					if hx == rightNormalLimit + 2 { return ContactTightFront1 }
				} else { // right side backwards edge
					if hx >= rightNormalLimit + 6 { return ContactSlipIntoFall }
					if hx == rightNormalLimit + 5 { return ContactTightBack2 }
					if hx == rightNormalLimit + 4 { return ContactTightBack2 }
					if hx == rightNormalLimit + 3 { return ContactTightBack1 }
					if hx == rightNormalLimit + 2 { return ContactTightBack1 }
				}
			} else {
				// right case, already handled by preceding step behavior
			}
			return ContactGround
		}

		return ContactNone
	case SubtypeThinStepOnRight: // MIRROR OF SubtypeThinStepOnLeft
		// side slipping
		if !flags.HasUpInertia() && hy + hh <= by + 4 {
			if hx + 7 > bx + bw { return ContactSlipIntoFall }
		}

		// step up cases
		if hy + hh == by + 7 {
			if flags & FlagDownPressed == 0 {
				if flags.HasLeftInertia()  && hx + 4 == bx + bw { return ContactStepUp }
			}
		}

		// step down cases, tight position or ground cases
		if hy + hh == by {
			if flags.HasRightInertia() && hx + hw == bx + bw + 4 { return ContactStepDown }
			
			leftNormalLimit  := bx
			rightNormalLimit := bx + bw - hw
			if hx >= leftNormalLimit && hx <= rightNormalLimit {
				return ContactGround
			}

			// check tight positionings
			rightSide := (hx + hw/2 >= bx + bw/2)
			if rightSide {
				// left case, already handled by preceding step behavior
			} else {
				if flags.IsLeftOriented() { // left side forwards edge
					if hx <= leftNormalLimit - 6 { return ContactSlipIntoFall }
					if hx == leftNormalLimit - 5 { return ContactTightFront2 }
					if hx == leftNormalLimit - 4 { return ContactTightFront2 }
					if hx == leftNormalLimit - 3 { return ContactTightFront1 }
					if hx == leftNormalLimit - 2 { return ContactTightFront1 }
				} else { // left side backwards edge
					if hx <= leftNormalLimit - 6 { return ContactSlipIntoFall }
					if hx == leftNormalLimit - 5 { return ContactTightBack2 }
					if hx == leftNormalLimit - 4 { return ContactTightBack2 }
					if hx == leftNormalLimit - 3 { return ContactTightBack1 }
					if hx == leftNormalLimit - 2 { return ContactTightBack1 }
				}
			}

			return ContactGround
		}

		return ContactNone

	default:
		panic("unimplemented subtype contact type for " + self.String())
	}
}
