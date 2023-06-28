package level

import "io/fs"

import "github.com/tinne26/transition/src/game/level/block"

var allLevels []*Level

type Key uint8
var (
	LvlStart Key // medna's walk
	LvlSword Key // first sword
	LvlPlants Key // second sword, thorn's patch
	LvlGhosts Key // connects to gate too, Land of The Yahnon
	LvlSpikes Key // third sword
	LvlGate Key // The White Gate
)

func Get(key Key) *Level {
	return allLevels[key]
}

func CreateAll(filesys fs.FS) error {
	// create blocks first
	err := block.CreateAll(filesys)
	if err != nil { return err }

	// --- define level entries ---
	var lvl *Level
	
	// start level
	lvl = CreateStartLevel()
	LvlStart = Key(len(allLevels))
	allLevels = append(allLevels, lvl)

	// sword level
	lvl = CreateSwordLevel()
	LvlSword = Key(len(allLevels))
	allLevels = append(allLevels, lvl)

	// ...

	return nil
}
