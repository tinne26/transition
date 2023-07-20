package audio

import "time"
import "io/fs"

import "github.com/hajimehoshi/ebiten/v2/audio"
import "github.com/tinne26/edau"

func Initialize(soundscape *Soundscape, filesys fs.FS) error {
	var err error
	var sfx *SfxPlayer

	// initialize context if necessary
	if audio.CurrentContext() == nil {
		audio.NewContext(44100)
	}
	
	// load sfxs
	sfx, err = loadWavMultiSFX(filesys, "assets/audio/sfx/step*.wav", '1', '3')
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.52)
	sfx.SetBackoff(time.Millisecond*85)
	SfxStep = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSoundEffect(filesys, "assets/audio/sfx/wings.wav")
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.92)
	SfxJump = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSoundEffect(filesys, "assets/audio/sfx/death.wav")
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.85)
	SfxDeath = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSoundEffect(filesys, "assets/audio/sfx/interact.wav")
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.54)
	SfxInteract = soundscape.RegisterSFX(sfx)

	sfx, err = loadWavSoundEffect(filesys, "assets/audio/sfx/reverse.wav")
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.6)
	SfxReverse = soundscape.RegisterSFX(sfx)
	
	sfx, err = loadWavSoundEffect(filesys, "assets/audio/sfx/sword_end.wav")
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.64)
	SfxSwordEnd = soundscape.RegisterSFX(sfx)
	
	sfx, err = loadOggMultiSFX(filesys, "assets/audio/sfx/sword_tap*.ogg", '1', '4')
	if err != nil { return err }
	sfx.SetVolumeCorrectorFactor(0.37)
	SfxSwordTap = soundscape.RegisterSFX(sfx)

	// load and set up bgms
	var loop1, loop2 *edau.Looper
	loop1, err = loadLooper(filesys, "assets/audio/bgm/background.ogg", 130, 5499197)
	if err != nil { return err }
	bgm := NewBgmFromLooper(loop1)
	bgm.SetVolumeCorrectorFactor(0.5)
	BgmBackground = soundscape.RegisterBGM(bgm)

	loop1, err = loadLooper(filesys, "assets/audio/bgm/challenge_base.ogg", 75938, 1942845)
	if err != nil { return err }
	loop2, err = loadLooper(filesys, "assets/audio/bgm/challenge_aux.ogg", 75938, 1942845)
	if err != nil { return err }
	auxFader := NewFader(NewProcessorInAdapterL16(loop2))
	bgm = NewBgmFromFader(NewFader(NewAdder(NewProcessorInAdapterL16(loop1), auxFader)))
	bgm.SetVolumeCorrectorFactor(0.5)
	BgmChallenge = soundscape.RegisterBGM(bgm)
	
	// put aux fader on soundscape's automation panel
	ResKeyChallengeFader = soundscape.AutomationPanel().StoreResource(auxFader)

	return nil
}
