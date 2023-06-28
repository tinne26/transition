package camera

type StaticTarget struct {
	X float64
	Y float64
}

func (self StaticTarget) GetCameraTargetPos() (float64, float64) {
	return self.X, self.Y
}
