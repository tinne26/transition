package debug

import "strconv"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/ebitenutil"

func Draw(canvas *ebiten.Image) {
	// NOTE: use DebugPrintAt and x, y if multiple debug fragments are needed
	if debugPerformance {
		fps := strconv.FormatFloat(ebiten.ActualFPS(), 'f', 2, 64)
		ebitenutil.DebugPrint(canvas, fps + "FPS")
	}
}
