package main

import "image"
import "embed"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/debug"
import "github.com/tinne26/transition/src/shaders"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/game"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/hint"

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
		// notice: setting a proper size before fullscreening
		//         is critical in case we later leave fullscreen
		w, h := utils.FindMaxMultRawWindowSize(640, 360, 128)
		ebiten.SetWindowSize(w, h)
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

	// set window icon
	ico16, err := utils.LoadFsImage(filesys, "assets/ico/16x16.png")
	if err != nil { debug.Fatal(err) }
	ico32, err := utils.LoadFsImage(filesys, "assets/ico/32x32.png")
	if err != nil { debug.Fatal(err) }
	ico48, err := utils.LoadFsImage(filesys, "assets/ico/48x48.png")
	if err != nil { debug.Fatal(err) }
	ebiten.SetWindowIcon([]image.Image{ico16, ico32, ico48})

	// create game and run it
	gg, err := game.New(filesys)
	if err != nil { debug.Fatal(err) }
	err = ebiten.RunGame(gg)
	if err != nil { debug.Fatal(err) }
}
