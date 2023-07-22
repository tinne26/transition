package context

import "io/fs"

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/state"

type Context struct {
	State *state.State
	Input *input.Input
	Audio *audio.Soundscape
}

func NewContext(filesys fs.FS) (*Context, error) {
	soundscape := audio.NewSoundscape()
	err := audio.Initialize(soundscape, filesys)
	if err != nil { return nil, err }

	return &Context{
		State: state.New(),
		Input: input.NewInput(input.StdKeyboardMappingWASD, input.StdGamepadMapping),
		Audio: soundscape,
	}, nil
}

func (self *Context) Update() error {
	var err error

	err = self.Audio.Update()
	if err != nil { return err }
	err = self.Input.Update()
	if err != nil { return err }

	return nil
}
