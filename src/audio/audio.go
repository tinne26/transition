package audio

import "io/fs"
import "io/ioutil"
import "math/rand"

import "github.com/hajimehoshi/ebiten/v2/audio"
import "github.com/hajimehoshi/ebiten/v2/audio/wav"
import "github.com/hajimehoshi/ebiten/v2/audio/vorbis"

import "github.com/tinne26/edau"

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

var ctx = audio.NewContext(44100)

func PlayBGM(filesys fs.FS) error {
	const LoopStartSample = 200
	const LoopEndSample   = 5510869

	file, err := filesys.Open("assets/audio/bgm/background.ogg")
	if err != nil { return err }
	stream, err := vorbis.DecodeWithSampleRate(44100, file)
	if err != nil { return err }
	looper := edau.NewLooper(stream, LoopStartSample*4, LoopEndSample*4)
	player, err := audio.NewPlayer(ctx, looper)
	player.SetVolume(0.4)
	if err != nil { return err }
	player.Play()
	return nil
}

func LoadSFX(filesys fs.FS) error {
	// load sfx
	var err error
	SfxJump, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/wings.wav")
	if err != nil { return err }
	SfxJump.SetVolume(0.9)

	SfxStep1, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/step1.wav")
	if err != nil { return err }
	SfxStep1.SetVolume(0.5)
	
	SfxStep2, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/step2.wav")
	if err != nil { return err }
	SfxStep2.SetVolume(0.5)

	SfxStep3, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/step3.wav")
	if err != nil { return err }
	SfxStep3.SetVolume(0.5)

	SfxDeath, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/death.wav")
	if err != nil { return err }
	SfxDeath.SetVolume(0.74)

	SfxInteract, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/interact.wav")
	if err != nil { return err }
	SfxInteract.SetVolume(0.42)

	SfxReverse, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/reverse.wav")
	if err != nil { return err }
	SfxReverse.SetVolume(0.5)
	
	SfxSword, err = loadWavSoundEffect(ctx, filesys, "assets/audio/sfx/sword_end.wav")
	if err != nil { return err }
	SfxSword.SetVolume(0.44)

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
