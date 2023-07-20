package game

import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/level/lvlkey"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/text"

func (self *Game) HandleTriggerResponse(response any) {
	switch typedResponse := response.(type) {
	case *text.Message:
		self.textMessage = typedResponse
	case hint.Hint:
		self.activeHint = &typedResponse
	case lvlkey.EntryKey:
		// changing savepoint key
		key := typedResponse
		lvl, _ := level.GetEntryPoint(key)
		lvl.DisableSavepoints()
		lvl.EnableSavepoint(key)
		self.player.UnblockInteractionAfter(16)
	case trigger.Transfer:
		self.fadeInTicksLeft = FadeTicks
		lvl, pt := level.GetEntryPoint(typedResponse.Key)
		self.transferPlayer(lvl, pt)
	case []string:
		self.longText = typedResponse
	case float64: // treat as opacity level, forcing fade outs and stuff
		self.forcefulFadeOutLevel = typedResponse
	case *sword.Challenge:
		challenge := typedResponse
		self.player.SetBlockedForInteraction()
		self.swordChallenge = challenge
		self.camera.SetStaticTarget(float64(challenge.X), float64(challenge.Y))
		self.camera.RequireMustMatch()
	default:
		panic(response)
	}
}
