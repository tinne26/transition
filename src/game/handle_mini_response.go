package game

import "github.com/tinne26/transition/src/game/player/miniscene"
import "github.com/tinne26/transition/src/game/player/motion"
import "github.com/tinne26/transition/src/game/player/comm"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/level/lvlkey"
import "github.com/tinne26/transition/src/game/flash"
import "github.com/tinne26/transition/src/shaders"

// (handle mini*scene* response)
func (self *Game) HandleMiniResponse(response any) {
	if response == nil { return }
	switch typedResponse := response.(type) {
	case miniscene.OverFlags:
		self.mini = nil
		flags := typedResponse
		if flags.HasToRestorePlayerOnCam() {
			self.camera.SetTarget(self.player)
		}
		if flags.HasToUnblockPlayer() {
			self.player.UnblockInteractionAfter(flags.GetUnblockTicks())
		}
	case motion.Pair:
		self.player.SetMotionPair(typedResponse, self.ctx)
	case comm.Action:
		self.player.ReceiveAction(typedResponse, self.ctx)
	case lvlkey.EntryKey:
		key := typedResponse
		self.ctx.State.LastSaveEntryKey = typedResponse
		lvl, _ := level.GetEntryPoint(key)
		lvl.DisableSavepoints()
		lvl.EnableSavepoint(key)
	case *shaders.Animation:
		self.gfxAnim = typedResponse
	case *flash.Flash:
		self.flash = typedResponse
	default:
		panic(response)
	}
}
