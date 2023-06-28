package debug

import "os"

import "github.com/hajimehoshi/ebiten/v2"

var debugPerformance = false
var debugTrace = false

func DetectAndSetUp() {
	for _, arg := range os.Args {
		switch arg {
		case "--debug": // all other options together
			debugPerformance = true
			debugTrace = true
			ebiten.SetVsyncEnabled(false)
			ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		case "--trace": // show debug.Trace() and debug.Tracef() messages
			debugTrace = true
		case "--maxfps": // unlock and show fps
			debugPerformance = true
			ebiten.SetVsyncEnabled(false)
		case "--resizable", "--resize": // allow window resizing
			ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		// ...
		}
	}
}
