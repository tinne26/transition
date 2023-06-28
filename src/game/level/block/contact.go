package block

import "strconv"

type ContactType uint8
const (
	ContactNone ContactType = iota
	ContactGround
	ContactSlipIntoFall
	ContactTightFront1
	ContactTightFront2
	ContactTightBack1
	ContactTightBack2
	ContactWallStick
	
	ContactClonk
	ContactSideBlock
	
	ContactStepUp
	ContactStepDown

	ContactDeath
	//ContactHurt
)

func (self ContactType) String() string {
	switch self {
	case ContactNone: return "ContactNone"
	case ContactGround: return "ContactGround"
	case ContactSlipIntoFall: return "ContactSlipIntoFall"
	case ContactTightFront1: return "ContactTightFront1"
	case ContactTightFront2: return "ContactTightFront2"
	case ContactTightBack1: return "ContactTightBack1"
	case ContactTightBack2: return "ContactTightBack2"
	case ContactWallStick: return "ContactWallStick"
	case ContactClonk: return "ContactClonk"
	case ContactSideBlock: return "ContactSideBlock"
	case ContactStepUp: return "ContactStepUp"
	case ContactStepDown: return "ContactStepDown"
	case ContactDeath: return "ContactDeath"
	default:
		return "ContactType#" + strconv.Itoa(int(self))
	}
}
