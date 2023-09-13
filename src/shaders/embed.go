package shaders

import _ "embed"

import "github.com/hajimehoshi/ebiten/v2"

//go:embed identity.kage
var identitySrc []byte

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

//go:embed vignette_bw.kage
var vignetteBwSrc []byte

//go:embed chromatic_aberration.kage
var chromaticAberrationSrc []byte

//go:embed radial_blur.kage
var radialBlurSrc []byte

//go:embed reversal_disk.kage
var reversalDiskSrc []byte

//go:embed transparency_highlighter.kage
var transparencyHighlighterSrc []byte

var Identity *ebiten.Shader
var Background *ebiten.Shader
var MaskedColoring *ebiten.Shader
var SwordChallenge *ebiten.Shader
var Title *ebiten.Shader
var Projection *ebiten.Shader
var ParallaxProjection *ebiten.Shader
var VignetteBW *ebiten.Shader
var ChromaticAberration *ebiten.Shader
var RadialBlur *ebiten.Shader
var ReversalDisk *ebiten.Shader
var TransparencyHighlighter *ebiten.Shader

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

	// vignette black and white
	VignetteBW, err = ebiten.NewShader(vignetteBwSrc)
	if err != nil { return err }

	// chromatic aberration
	ChromaticAberration, err = ebiten.NewShader(chromaticAberrationSrc)
	if err != nil { return err }

	// poor man's radial blur
	RadialBlur, err = ebiten.NewShader(radialBlurSrc)
	if err != nil { return err }

	// poor man's radial blur
	ReversalDisk, err = ebiten.NewShader(reversalDiskSrc)
	if err != nil { return err }

	// identity
	Identity, err = ebiten.NewShader(identitySrc)
	if err != nil { return err }

	// transparency highlighter
	TransparencyHighlighter, err = ebiten.NewShader(transparencyHighlighterSrc)
	if err != nil { return err }

	// load animation effects
	loadAnimationEffects()

	return nil
}
