package shaders

var AnimSetRespawn *Animation

// Called from loadAnimationEffects()
func loadAnimSetRespawn() {
	anim := NewAnim()
	anim.SetStartValue(RadialBlur, "RadialBlurAmount", 0.0)
	var t0 uint16
	t1 := anim.AddPt(RadialBlur, "RadialBlurAmount", t0, 8, 1.4, InterpSine)
	_   = anim.AddPt(RadialBlur, "RadialBlurAmount", t1, 10, 0.0, InterpSine)

	anim.SetStartValue(ChromaticAberration, "BlurSofteningFactor", 1.12)
	anim.SetStartValue(ChromaticAberration, "HorzAberrationMinRange", 0.0)
	anim.SetStartValue(ChromaticAberration, "HorzAberrationMaxRange", 0.0)
	t1 = anim.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t0, 8, 0.4, InterpExpo)
	_  = anim.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t1, 8, 0.00, InterpExpo)

	AnimSetRespawn = anim
}
