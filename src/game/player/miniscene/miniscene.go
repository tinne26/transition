package miniscene

import "github.com/tinne26/transition/src/project"
import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/player/comm"

type Scene interface {
	Update(*context.Context, *camera.Camera, comm.Status) (any, error)
	CurrentText() *text.Message
	BackDraw(*project.Projector)
	
	// ...other draw types may be required, e.g.:
	// Draw(*project.Projector)
}
