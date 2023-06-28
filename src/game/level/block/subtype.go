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

// Precondition: rects can't be empty or malformed
// func (self Subtype) HasCollision(blockRect, collider RectU16, flags Flags) bool {
// 	intersection := blockRect.NonEmptyIntersect(collider)
// 	switch self {
// 	case SubtypeBlock:
// 		if intersection.Max.X <= blockRect.Min.X + 3 || intersection.Min.X >= blockRect.Max.X - 3 {
// 			if intersection.Max.Y <= 3 { return false }
// 			if intersection.Min.Y >= blockRect.Max.Y - 3 { return false }
// 		}
// 	case SubtypeLeftRect:
// 		if intersection.Min.X >= blockRect.Max.X - 3 {
// 			if intersection.Max.Y <= 3 { return false }
// 			if intersection.Min.Y >= blockRect.Max.Y - 3 { return false }
// 		}
// 	case SubtypeRightRect:
// 		if intersection.Max.X <= blockRect.Min.X + 3 {
// 			if intersection.Max.Y <= 3 { return false }
// 			if intersection.Min.Y >= blockRect.Max.Y - 3 { return false }
// 		}
// 	case SubtypeSpikes:
// 		if intersection.Max.X <= blockRect.Min.X - 6 { return false }
// 		if intersection.Min.X >= blockRect.Max.X - 6 { return false }
// 		if intersection.Max.Y <= blockRect.Min.Y - 6 { return false }
// 		if intersection.Min.Y >= blockRect.Max.Y - 6 { return false }
// 		intersection.Min.X += 6
// 		intersection.Max.X -= 6
// 		intersection.Min.Y += 6
// 		intersection.Max.Y -= 6
// 		if intersection.Max.X <= blockRect.Min.X + 2 || intersection.Min.X >= blockRect.Max.X - 2 {
// 			if intersection.Max.Y <= 2 { return false }
// 			if intersection.Min.Y >= blockRect.Max.Y - 2 { return false }
// 		}
// 	case SubtypeThinBlock:
// 		if flags & FlagInertiaDown == 0 { return false }
// 		if collider.Max.Y > blockRect.Max.Y { return false }
// 		if intersection.Max.X <= blockRect.Min.X + 1 || intersection.Min.X >= blockRect.Max.X - 1 {
// 			if intersection.Max.Y <= 1 { return false }
// 			if intersection.Min.Y >= blockRect.Max.Y - 1 { return false }
// 		}
// 	case SubtypePlantSpikyA:
// 		//if flags & FlagFalling == 0 { return false }
// 		//if collider.Max.Y > blockRect.Max.Y { return false }
// 		panic("unimplemented")
// 	case SubtypePlantSpikyB:
// 		//if flags & FlagFalling == 0 { return false }
// 		//if collider.Max.Y > blockRect.Max.Y { return false }
// 		panic("unimplemented")
// 	}
// 	return true
// }

// func (self Subtype) HitDamage(blockRect, collider RectU16, flags Flags, haveCollision bool) uint8 {
// 	switch self {
// 	case SubtypeSpikes:
// 		if haveCollision { return 255 }
// 	case SubtypePlantSpikyA:
// 		panic("unimplemented") // this one does actually need to recalculate collisions in sophisticated ways
// 	case SubtypePlantSpikyB:
// 		panic("unimplemented")
// 	}
// 	return 0 // no damage by default
// }
