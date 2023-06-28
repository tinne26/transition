package utils

import "image/color"

func RescaleAlphaRGBA(rgba color.RGBA, newAlpha uint8) color.RGBA {
	if rgba.A == newAlpha { return rgba }

	factor := float64(newAlpha)/float64(rgba.A)
	return color.RGBA{
		R: uint8(float64(rgba.R)*factor),
		G: uint8(float64(rgba.G)*factor),
		B: uint8(float64(rgba.B)*factor),
		A: newAlpha,
	}
}

func ToRGBAf32(clr color.Color) (r, g, b, a float32) {
	r16, g16, b16, a16 := clr.RGBA()
	return float32(r16)/65535.0, float32(g16)/65535.0, float32(b16)/65535.0, float32(a16)/65535.0
}

func RGBA8ToRGBAf32(clr color.RGBA) (r, g, b, a float32) {
	return float32(clr.R)/255.0, float32(clr.G)/255.0, float32(clr.B)/255.0, float32(clr.A)/255.0
}
