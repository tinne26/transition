package text

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

var BackColor  = color.RGBA{ 20,  20,  20, 255}
var FrontColor = color.RGBA{255, 255, 255, 255}

const LineHeight = 7
const LineInterspace = 2

// for long text and so on
func CenterRawDraw(canvas *ebiten.Image, lines []string, clr color.RGBA) {
	bounds := canvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	textHeight := len(lines)*LineHeight + (len(lines) - 1)*LineInterspace
	y := h/2 - textHeight/2
	for _, line := range lines {
		lineWidth := MeasureLineWidth(line)
		x := w/2 - lineWidth/2
		DrawLine(canvas, line, x, y, clr)
		y += LineHeight + LineInterspace
	}
}

func Draw(canvas *ebiten.Image, cx, cy int, msg *Message) {
	// helper function
	var fill = func(target *ebiten.Image, minX, minY, width, height int, clr color.RGBA) {
		target.SubImage(image.Rect(minX, minY, minX + width, minY + height)).(*ebiten.Image).Fill(clr)
	}
	
	// start calculating width and height
	height := 8*2 + 7
	if msg.HasTwoLines() { height += LineHeight + LineInterspace }
	width := 9*2

	// see which line is longest
	w1 := MeasureLineWidth(msg.FirstLine)
	w2 := MeasureLineWidth(msg.SecondLine)
	if w2 > w1 { w1 = w2 }
	width += w1

	// determine start point
	ox, oy := cx - w1/2, cy - height/2

	// draw main box
	fill(canvas, ox, oy, width, height, BackColor)

	// draw white border
	fill(canvas, ox + 1, oy + 1, width - 2, 1, FrontColor)
	fill(canvas, ox + 1, oy + height - 2, width - 2, 1, FrontColor)
	fill(canvas, ox + 1, oy + 2, 1, height - 4, FrontColor)
	fill(canvas, ox + width - 2, oy + 2, 1, height - 4, FrontColor)

	// draw first line
	DrawLine(canvas, msg.FirstLine, ox + 9, oy + 8, msg.Color)
	
	// draw first line
	if msg.HasTwoLines() {
		DrawLine(canvas, msg.SecondLine, ox + 9, oy + 8 + LineHeight + LineInterspace, msg.Color)
	}

	// apply skippable decoration
	if msg.IsSkippable {
		fill(canvas, ox + width - 14, oy + height, 11, 4, BackColor)
		fill(canvas, ox + width - 13, oy + height - 6, 9, 9, FrontColor)
		fill(canvas, ox + width - 12, oy + height - 5, 7, 7, BackColor)
		DrawLine(canvas, string(KeyMsgI), ox + width - 10, oy + height - 5, msg.Color)
	}

	// apply dialogue decoration
	if msg.IsDialogue {
		fill(canvas, ox + 3, oy - 2, 11, 2, BackColor)
		fill(canvas, ox + 4, oy - 1, 9, 5, FrontColor)
		fill(canvas, ox + 5, oy, 7, 3, BackColor)
		DrawLine(canvas, "...", ox + 6, oy - 4, msg.Color)
	}
}

func MeasureLineWidth(line string) int {
	var prevIsSpace bool
	width := 0
	for i, codePoint := range line {
		if codePoint == ' ' {
			width += 4
			prevIsSpace = true
		} else {
			if i != 0 && !prevIsSpace { width += 1 }
			prevIsSpace = false
			width += pkgBitmaps[codePoint].Bounds().Dx() // panic here => missing letter bitmap
		}
	}
	return width
}

func DrawLine(canvas *ebiten.Image, line string, ox, oy int, textColor color.RGBA) {
	var prevIsSpace bool
	x := 0
	opts := ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleWithColor(textColor)
	for i, codePoint := range line {
		if codePoint == ' ' {
			x += 4
			prevIsSpace = true
		} else {
			if i != 0 && !prevIsSpace { x += 1 }
			prevIsSpace = false
			img := pkgBitmaps[codePoint]
			opts.GeoM.Translate(float64(ox + x), float64(oy))
			canvas.DrawImage(img, &opts)
			opts.GeoM.Reset()
			x += img.Bounds().Dx()
		}
	}
}
