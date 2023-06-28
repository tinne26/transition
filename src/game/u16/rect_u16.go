package u16

import "image"

type Rect struct {
	Min, Max Point
}

// Won't normalize even if coords are swapped. Intentional.
func NewRect(minX, minY, maxX, maxY uint16) Rect {
	return Rect{ Min: Point{ minX, minY }, Max: Point{ maxX, maxY }}
}

func FromImageRect(imgRect image.Rectangle) Rect {
	// unsafe, no range checks
	return NewRect(
		uint16(imgRect.Min.X), uint16(imgRect.Min.Y),
		uint16(imgRect.Max.X), uint16(imgRect.Max.Y),
	)
}

func (self Rect) PadEachFace(pad uint16) Rect {
	return NewRect(self.Min.X - pad, self.Min.Y - pad, self.Max.X + pad, self.Max.Y + pad)
}

func (self Rect) PadHorz(pad uint16) Rect {
	return NewRect(self.Min.X - pad, self.Min.Y, self.Max.X + pad, self.Max.Y)
}

func (self Rect) ToImageRect() image.Rectangle {
	return image.Rect(int(self.Min.X), int(self.Min.Y), int(self.Max.X), int(self.Max.Y))
}

func (self Rect) Intersect(other Rect) Rect {
	if self.Min.X < other.Min.X { self.Min.X = other.Min.X }
	if self.Max.X > other.Max.X { self.Max.X = other.Max.X }
	if self.Min.X >= self.Max.X { return Rect{} }
	
	if self.Min.Y < other.Min.Y { self.Min.Y = other.Min.Y }
	if self.Max.Y > other.Max.Y { self.Max.Y = other.Max.Y }
	if self.Min.Y >= self.Max.Y { return Rect{} }
	return self
}

func (self Rect) NonEmptyIntersect(other Rect) Rect {
	if self.Min.X < other.Min.X { self.Min.X = other.Min.X }
	if self.Max.X > other.Max.X { self.Max.X = other.Max.X }
	if self.Min.Y < other.Min.Y { self.Min.Y = other.Min.Y }
	if self.Max.Y > other.Max.Y { self.Max.Y = other.Max.Y }
	return self
}

func (self Rect) Empty() bool {
	return self.Min.X >= self.Max.X || self.Min.Y >= self.Max.Y
}

func (self Rect) NonEmptyOverlap(other Rect) bool {
	return self.Min.X < other.Max.X && other.Min.X < self.Max.X && self.Min.Y < other.Max.Y && other.Min.Y < self.Max.Y
}

func (self Rect) Overlap(other Rect) bool {
	return self.NonEmptyOverlap(other) && !self.Empty() && !other.Empty()
}
