package player

import "io/fs"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/tinne26/transition/src/utils"

var UICorruptionStages *ebiten.Image
var UIPowerFrame *ebiten.Image

func LoadUIGraphics(filesys fs.FS) error {
	var err error
	
	UICorruptionStages, err = utils.LoadFsEbiImage(filesys, "assets/graphics/ui/corruption_stages.png")
	if err != nil { return err }

	UIPowerFrame, err = utils.LoadFsEbiImage(filesys, "assets/graphics/ui/power_frame.png")
	if err != nil { return err }

	return nil
}
