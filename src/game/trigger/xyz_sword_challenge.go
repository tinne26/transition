package trigger

import "image/color"

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/text"

var white color.RGBA = color.RGBA{255, 255, 255, 255}
var apology = []*text.Message{
	text.NewSkippableMsg2(
		"SORRY ABOUT IT, THIS IS WHERE YOU WOULD ACTUALLY GET",
		"YOUR MAIN POWER AND THE GAME WOULD START PROPERLY, BUT...", white),
	text.NewSkippableMsg1("LIFE HAPPENED", white),
	text.NewSkippableMsg2(
		"FEEL FREE TO ROAM AROUND IF YOU WANT AND ENJOY THE SCENERY,",
		"BUT THERE'S NOTHING MORE TO SEE IN TERMS OF CONTENT", white),
}

var _ Trigger = (*TrigSwordChallenge)(nil)

type TrigSwordChallenge struct {
	done bool
	area u16.Rect
	hint any // hint.Hint
	challenge any // *sword.Challenge
	apologyIndex int
}

func NewSwordChallenge(area u16.Rect, hint any, challenge any) Trigger {
	return &TrigSwordChallenge{
		area: area,
		hint: hint,
		challenge: challenge,
	}
}

func (self *TrigSwordChallenge) Update(playerRect u16.Rect, state *State) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	
	if !self.done {
		if input.Trigger(input.ActionInteract) {
			self.done = true
			return self.challenge, nil
		}
		return self.hint, nil
	}

	if state.SwordChallengesSolved > 0 {
		if input.Trigger(input.ActionInteract) {
			self.apologyIndex += 1
		}
		if self.apologyIndex < len(apology) {
			return apology[self.apologyIndex], nil
		}
	}

	return nil, nil
}

func (self *TrigSwordChallenge) OnLevelEnter(state *State) {}
func (self *TrigSwordChallenge) OnLevelExit(state *State) {}
func (self *TrigSwordChallenge) OnDeath(state *State) {}
