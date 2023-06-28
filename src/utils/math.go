package utils

type Integer interface {
	~uint8 | ~uint16 | ~uint32  | ~uint64 |
	~int8 | ~int16 | ~int32 | ~int64 |
	~int | ~uint
}

type Float interface { ~float32 | ~float64 }

type Numeric interface { Integer | Float }

func Max[Num Numeric](a, b Num) Num {
	if a >= b { return a }
	return b
}

func Min[Num Numeric](a, b Num) Num {
	if a <= b { return a }
	return b
}

func Abs[Num Numeric](a Num) Num {
	if a >= 0 { return a }
	return -a
}

func MinMax[Num Numeric](a, b Num) (Num, Num) {
	if a <= b { return a, b }
	return b, a
}

func FastFloor[F Float](x F) F {
	return F(int(x))
}

func FastCeil[F Float](x F) F {
	backForth := F(int(x))
	if backForth == x { return backForth }
	return backForth + 1
}

func FastCeilInt[F Float](x F) int {
	return -int(-x)
}

func Compare[Num Numeric](a, b Num) int {
	if a == b { return 0 }
	if a > b { return 1 }
	return -1
}

func NextTowards[Int Integer](from, towards Int) Int {
	if from == towards { return from }
	if from < towards { return from + 1 }
	return from - 1
}
