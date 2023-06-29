package audio

import "io/fs"
import "io/ioutil"
import "math/rand"

import "github.com/hajimehoshi/ebiten/v2/audio"
import "github.com/hajimehoshi/ebiten/v2/audio/wav"

var SfxStep1 *SfxPlayer
var SfxStep2 *SfxPlayer
var SfxStep3 *SfxPlayer
var SfxJump *SfxPlayer
var SfxDeath *SfxPlayer
var SfxInteract *SfxPlayer
var SfxReverse *SfxPlayer
var SfxSword *SfxPlayer

func PlayJump() { SfxJump.Play() }
func PlayDeath() { SfxDeath.Play() }

func PlayStep() {
	r := rand.Float64()
	if r < 0.33 {
		SfxStep1.Play()
	} else if r < 0.66 {
		SfxStep2.Play()
	} else {
		SfxStep3.Play()
	}
}

func PlayReverse() { SfxReverse.Play() }
func PlayInteract() { SfxInteract.Play() }

func PlaySwordEnd() { SfxSword.Play() }
func PlaySwordTap() {
	PlayStep()
}

func LoadSFX(filesys fs.FS) error {
	ctx := audio.NewContext(44100)

	// load sfx
	var err error
	SfxJump, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/wings.wav")
	if err != nil { return err }
	SfxJump.SetVolume(0.3)

	SfxStep1, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/step1.wav")
	if err != nil { return err }
	SfxStep1.SetVolume(0.1)
	
	SfxStep2, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/step2.wav")
	if err != nil { return err }
	SfxStep2.SetVolume(0.1)

	SfxStep3, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/step3.wav")
	if err != nil { return err }
	SfxStep3.SetVolume(0.1)

	SfxDeath, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/death.wav")
	if err != nil { return err }
	SfxDeath.SetVolume(0.34)

	SfxInteract, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/interact.wav")
	if err != nil { return err }
	SfxInteract.SetVolume(0.25)

	SfxReverse, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/reverse.wav")
	if err != nil { return err }
	SfxReverse.SetVolume(0.3)
	
	SfxSword, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/sword_end.wav")
	if err != nil { return err }
	SfxSword.SetVolume(0.2)

	return nil
}

func loadWavSoundEffect(ctx *audio.Context, filesys fs.FS, filename string) (*SfxPlayer, error) {
	file, err := filesys.Open(filename)
	if err != nil { return nil, err }
	stream, err := wav.DecodeWithSampleRate(44100, file)
	if err != nil { return nil, err }
	audioBytes, err := ioutil.ReadAll(stream)
	if err != nil { return nil, err }
	return NewSfxPlayer(ctx, audioBytes), nil
}
