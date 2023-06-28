//go:build nope

package level

import "math/rand"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/shaders"

var rectShaderOpts = ebiten.DrawRectShaderOptions{
	Uniforms: make(map[string]any, 1),
}

func FxParallaxNoise(target, source *ebiten.Image) {
	rectShaderOpts.Uniforms["Randomness"] = rand.Float32()
	rectShaderOpts.Images[0] = source
	bounds := target.Bounds()
	target.DrawRectShader(bounds.Dx(), bounds.Dy(), shaders.Background, &rectShaderOpts)
}

func FxParallaxNone(target, source *ebiten.Image) {
	target.DrawImage(source, nil)
}

