package block

import "strconv"

// Subtypes are a hacky way to implement collisions, damage, more sophisticate
// collision or damage boxes or whatever as a simple big switch instead of
// a general system that can work with everything nicely.
type Subtype uint8

const (
	SubtypeNone Subtype = iota
	
	SubtypeBlock // has a 3 pixel corner
	SubtypeThinBlock // can jump into it from below, has a single pixel corner
	SubtypeThinStep // can step up or down from any side
	SubtypeThinStepOnLeft // can step down or up on the left side
	SubtypeThinStepOnRight // can step down or up on the right side
	SubtypeSpikes // the corners don't have spikes
	SubtypePlantSpikyA // varies based on reversal being into effect or not
	SubtypePlantSpikyB
	SubtypeDarkFloor
)

func (self Subtype) String() string {
	switch self {
	case SubtypeNone: return "SubtypeNone"
	case SubtypeBlock: return "SubtypeBlock"
	case SubtypeThinBlock: return "SubtypeThinBlock"
	case SubtypeThinStep: return "SubtypeThinStep"
	case SubtypeThinStepOnLeft: return "SubtypeThinStepOnLeft"
	case SubtypeThinStepOnRight: return "SubtypeThinStepOnRight"
	case SubtypeSpikes: return "SubtypeSpikes"
	case SubtypePlantSpikyA: return "SubtypePlantSpikyA"
	case SubtypePlantSpikyB: return "SubtypePlantSpikyB"
	case SubtypeDarkFloor: return "SubtypeDarkFloor"
	default:
		return "Subtype#" + strconv.Itoa(int(self))
	}
}
