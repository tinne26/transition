package shaders

var AnimRespawn *Animation

// Called from loadAnimationEffects()
func loadAnimRespawn() {
	// TODO: maybe I should write a couple variations so it doesn't look always the same
	//       when you die. It's not terrible as it is, though.
	AnimRespawn = NewAnim()
	AnimRespawn.SetStartValue(RadialBlur, "RadialBlurAmount", 4.0)
	_ = AnimRespawn.AddPt(RadialBlur, "RadialBlurAmount", 0, 100, 0.0, InterpSine)

	AnimRespawn.SetStartValue(ChromaticAberration, "BlurSofteningFactor", 1.12)
	AnimRespawn.SetStartValue(ChromaticAberration, "HorzAberrationMinRange", 0.0)
	AnimRespawn.SetStartValue(ChromaticAberration, "HorzAberrationMaxRange", 0.0)
	t0 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", 0, 8, 0.0, InterpLinear)
	t1 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t0, 36, 0.02, InterpLinear)
	t2 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t1, 6, 0.1, InterpExpo)
	t3 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t2, 6, 0.32, InterpExpo)
	t4 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t3, 14, 0.04, InterpExpo)
	t5 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t4, 8, 0.2, InterpSine)
	t6 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t5, 8, 0.0, InterpSine)
	t7 := AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t6, 8, 0.08, InterpExpo)
	_   = AnimRespawn.AddPt(ChromaticAberration, "HorzAberrationMaxRange", t7, 8, 0.0, InterpSine)

	AnimRespawn.SetStartValue(VignetteBW, "MinRadius", 0.2)
	AnimRespawn.SetStartValue(VignetteBW, "MaxRadius", 1.1)
	AnimRespawn.SetStartValue(VignetteBW, "MixLevel", 1.0)
	t1vig := AnimRespawn.AddPt(VignetteBW, "MixLevel", 0, 80, 1.0, InterpLinear)
	_   = AnimRespawn.AddPt(VignetteBW, "MixLevel", t1vig, 160, 0.0, InterpLinear)
}
