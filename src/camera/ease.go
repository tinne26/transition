package camera

import "math"

// Speeds are given in pixels per tick.
func easeMove(dist, maxDist, maxSpeed, baseSpeed float64) float64 {
	// don't move if already at goal
	if dist == 0 { return 0 }
	
	absDist := dist
	if absDist < 0 { absDist = -absDist }
	if absDist > maxDist { absDist = maxDist }
	
	var move float64
	if absDist >= maxDist {
		move = absDist
	} else {
		move = baseSpeed + easeFunc(absDist/maxDist)*(maxSpeed - baseSpeed)
	}
	
	if dist < 0 { return -move }
	return move
}

func easeFunc(t float64) float64 {
	return easeCubic(t) // easeQuad also seems fair to me, easeSine seems visibly worse
}

func easeCubic(t float64) float64 {
	t *= 2.0
	if t < 1.0 {
		return 0.5*t*t*t
	} else {
		t -= 2.0
		return 0.5*(t*t*t + 2.0)
	}
}

func easeQuad(t float64) float64 {
	if t < 0.5 {
		return 2.0*t*t
	} else {
		t = 2.0*t - 1.0
		return -0.5*(t*(t - 2) - 1)
	}
}

func easeSine(t float64) float64 {
	return -0.5*(math.Cos(math.Pi*t) - 1.0)
}
