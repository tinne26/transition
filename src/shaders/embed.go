package shaders

import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed background.kage
var backgroundSrc []byte

//go:embed masked_coloring.kage
var maskedColoringSrc []byte

//go:embed sword_challenge.kage
var swordChallengeSrc []byte

var Background *ebiten.Shader
var MaskedColoring *ebiten.Shader
var SwordChallenge *ebiten.Shader

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

	return nil
}
