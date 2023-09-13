package game

import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/game/player/miniscene"
import "github.com/tinne26/transition/src/text"

func (self *Game) HandleTriggerResponse(response any) {
	switch typedResponse := response.(type) {
	case *text.Message:
		self.textMessage = typedResponse
	case hint.Hint:
		self.activeHint = &typedResponse
	case trigger.Transfer:
		lvl, pt := level.GetEntryPoint(typedResponse.Key)
		self.transferPlayer(lvl, pt)
	case []string:
		self.longText = typedResponse
	case float64: // treat as opacity level, forcing fade outs and stuff
		if self.fader.IsFading() {
			self.fader.SetBlacknessIfBelow(typedResponse)
		} else {
			self.fader.SetBlackness(typedResponse)
		}
	case *sword.Challenge:
		challenge := typedResponse
		self.player.SetBlockedForInteraction()
		self.swordChallenge = challenge
		self.camera.SetStaticTarget(float64(challenge.X), float64(challenge.Y))
		self.camera.RequireMustMatch()
	case miniscene.Scene:
		self.mini = typedResponse
		self.player.SetBlockedForInteraction()
	case trigger.SwitchPoint:
		self.mini = miniscene.NewResetSwitchScene(typedResponse.X, typedResponse.Y, typedResponse.Key)
		self.player.SetBlockedForInteraction()
	default:
		panic(response)
	}
}
