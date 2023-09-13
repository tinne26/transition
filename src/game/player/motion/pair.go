package motion

type Pair struct {
	State State
	Animation *Animation
}

func NewPair(state State, anim *Animation) Pair {
	return Pair{
		State: state,
		Animation: anim,
	}
}
