package shaders

import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed background.kage
var backgroundSrc []byte

//go:embed masked_coloring.kage
var maskedColoringSrc []byte

//go:embed sword_challenge.kage
var swordChallengeSrc []byte

//go:embed title.kage
var titleSrc []byte

//go:embed projection.kage
var projectionSrc []byte

//go:embed parallax_projection.kage
var parallaxProjectionSrc []byte

var Background *ebiten.Shader
var MaskedColoring *ebiten.Shader
var SwordChallenge *ebiten.Shader
var Title *ebiten.Shader
var Projection *ebiten.Shader
var ParallaxProjection *ebiten.Shader

func LoadAll() error {
	var err error
	
	// background shader
	Background, err = ebiten.NewShader(backgroundSrc)
	if err != nil { return err }

	// masked coloring shader
	MaskedColoring, err = ebiten.NewShader(maskedColoringSrc)
	if err != nil { return err }

	// sword challenge shader
	SwordChallenge, err = ebiten.NewShader(swordChallengeSrc)
	if err != nil { return err }

	// sword challenge shader
	Title, err = ebiten.NewShader(titleSrc)
	if err != nil { return err }

	// projection shader
	Projection, err = ebiten.NewShader(projectionSrc)
	if err != nil { return err }
	
	// parallax projection shader
	ParallaxProjection, err = ebiten.NewShader(parallaxProjectionSrc)
	if err != nil { return err }

	return nil
}
