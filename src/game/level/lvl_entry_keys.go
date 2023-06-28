package level

import "github.com/tinne26/transition/src/game/u16"

type EntryKey uint8
const (
	EntryStartSaveLeft EntryKey = iota
	EntryStartSaveRight
	EntryStartTransRight
	EntrySwordTransLeft
	EntrySwordTransRight
	EntrySwordSaveCenter
	EntryBasicsTransLeft
	EntryBasicsTransRight
	EntryGhostsTransLeft
	EntryGhostsTransRight
	EntryGhostsTransGate
	EntryGhostsSave
	EntrySpikesLeft
	EntrySpikesRight
	EntryGateTransGhosts

	//EntryPlantsLeft
	//EntryPlantsRight
	
	entryKeyEndSentinel
)

const numEntryKeys = entryKeyEndSentinel
var allEntries [numEntryKeys]u16.Point
var allEntryLevels [numEntryKeys]*Level

func GetEntryPoint(key EntryKey) (*Level, u16.Point) {
	lvl, pt := allEntryLevels[key], allEntries[key]
	if pt.X == 0 && pt.Y == 0 { panic(key) }
	return lvl, pt
}

func SetEntryPoint(key EntryKey, level *Level, x, y uint16) {
	allEntries[key] = u16.Point{X: x, Y: y}
	allEntryLevels[key] = level
}
