package utils

import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

var vertices []ebiten.Vertex
var mask1x1 *ebiten.Image
var stdTriOpts ebiten.DrawTrianglesOptions
var linearTriOpts ebiten.DrawTrianglesOptions

func init() {
	mask3x3 := ebiten.NewImage(3, 3)
	mask3x3.Fill(color.White)
	mask1x1 = mask3x3.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
	vertices = make([]ebiten.Vertex, 4)
	for i := 0; i < 4; i++ {
		vertices[i].SrcX = 1.0
		vertices[i].SrcY = 1.0
	}
	linearTriOpts.Filter = ebiten.FilterLinear
}

var rawWhiteMaskBuffer = make([]byte, 0, 1024)
func RawAlphaMaskToWhiteMask(width int, mask []byte) *ebiten.Image {
	for i := 0; i < len(mask); i++ {
		if mask[i] > 0 {
			rawWhiteMaskBuffer = append(rawWhiteMaskBuffer, 255, 255, 255, 255)
		} else {
			rawWhiteMaskBuffer = append(rawWhiteMaskBuffer, 0, 0, 0, 0)
		}
	}
	
	img := ebiten.NewImage(width, len(mask)/width)
	img.WritePixels(rawWhiteMaskBuffer)
	rawWhiteMaskBuffer = rawWhiteMaskBuffer[:0]
	return img
}

func DrawRectF32(target *ebiten.Image, minX, minY, maxX, maxY float32, r, g, b, a float32) {
	for i := 0; i < 4; i++ {
		vertices[i].ColorR = r
		vertices[i].ColorG = g
		vertices[i].ColorB = b
		vertices[i].ColorA = a
	}

	vertices[0].DstX = minX
	vertices[0].DstY = minY
	vertices[1].DstX = maxX
	vertices[1].DstY = minY
	vertices[2].DstX = maxX
	vertices[2].DstY = maxY
	vertices[3].DstX = minX
	vertices[3].DstY = maxY
	target.DrawTriangles(vertices[0 : 4], []uint16{0, 1, 2, 2, 3, 0}, mask1x1, &linearTriOpts)
}

// Similar to Ebitengine's Image.Fill(), but doesn't override the content but draw
// on top instead. Used when we want to support transparency on fills.
func FillOver(target *ebiten.Image, fillColor color.Color) {
	bounds := target.Bounds()
	if bounds.Empty() { return }

	r, g, b, a := fillColor.RGBA()
	if a == 0 { return }
	fr, fg, fb, fa := float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535
	for i := 0; i < 4; i++ {
		vertices[i].ColorR = fr
		vertices[i].ColorG = fg
		vertices[i].ColorB = fb
		vertices[i].ColorA = fa
	}

	minX, minY := float32(bounds.Min.X), float32(bounds.Min.Y)
	maxX, maxY := float32(bounds.Max.X), float32(bounds.Max.Y)
	vertices[0].DstX = minX
	vertices[0].DstY = minY
	vertices[1].DstX = maxX
	vertices[1].DstY = minY
	vertices[2].DstX = maxX
	vertices[2].DstY = maxY
	vertices[3].DstX = minX
	vertices[3].DstY = maxY

	target.DrawTriangles(vertices[0 : 4], []uint16{0, 1, 2, 2, 3, 0}, mask1x1, &stdTriOpts)
}

func FillOverWithBlend(target *ebiten.Image, fillColor color.Color, blendMode ebiten.Blend) {
	stdTriOpts.Blend = blendMode
	FillOver(target, fillColor)
	stdTriOpts.Blend = ebiten.Blend{}
}

func FillOverF32(target *ebiten.Image, r, g, b, a float32) {
	bounds := target.Bounds()
	if bounds.Empty() { return }

	for i := 0; i < 4; i++ {
		vertices[i].ColorR = r
		vertices[i].ColorG = g
		vertices[i].ColorB = b
		vertices[i].ColorA = a
	}

	minX, minY := float32(bounds.Min.X), float32(bounds.Min.Y)
	maxX, maxY := float32(bounds.Max.X), float32(bounds.Max.Y)
	vertices[0].DstX = minX
	vertices[0].DstY = minY
	vertices[1].DstX = maxX
	vertices[1].DstY = minY
	vertices[2].DstX = maxX
	vertices[2].DstY = maxY
	vertices[3].DstX = minX
	vertices[3].DstY = maxY

	target.DrawTriangles(vertices[0 : 4], []uint16{0, 1, 2, 2, 3, 0}, mask1x1, &stdTriOpts)
}
