package main

import "embed"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/debug"
import "github.com/tinne26/transition/src/shaders"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/game"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/audio"

// Windows compilation:
// > go build -o game.exe -trimpath -ldflags "-w -s -H windowsgui" -tags "ebitenginesinglethread" main.go

//go:embed assets/*
var filesys embed.FS

func main() {
	// debug setup
	debug.DetectAndSetUp()
	
	// ebitengine basic config
	err := utils.OnWindowsPreferOpenGL()
	if err != nil { debug.Fatal(err) }
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowTitle("tinne/transition")
	ebiten.SetScreenClearedEveryFrame(false)
	if utils.OsArgReceived("--windowed") {
		utils.SetMaxMultRawWindowSize(640, 360, 128)
	} else {
		ebiten.SetWindowSize(640, 360)
		ebiten.SetFullscreen(true)
	}

	// load files. no loading screen, you can do it
	// better at home if you want
	err = shaders.LoadAll()
	if err != nil { debug.Fatal(err) }
	err = level.CreateAll(filesys)
	if err != nil { debug.Fatal(err) }
	err = hint.LoadHintGraphics(filesys)
	if err != nil { debug.Fatal(err) }
	err = audio.LoadSFX(filesys)
	if err != nil { debug.Fatal(err) }

	// create game and run it
	gg, err := game.New(filesys)
	if err != nil { debug.Fatal(err) }
	err = ebiten.RunGame(gg)
	if err != nil { debug.Fatal(err) }
}
