package audio

var stepTickBackoff uint8
func Update() error {
	// maaaassive haaaacks
	if stepTickBackoff > 0 {
		stepTickBackoff -= 1
	}

	return nil
}
