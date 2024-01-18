package maps

type Location struct {
	X, Y int
}

type Size struct {
	Width, Height int
}

type Line struct {
	Start, End Location
}

type Map struct {
	Lines []Line
	Size  Size
}

func (m *Map) IsLocationValid(location Location) bool {
	if location.X < 0 || location.Y < 0 {
		return false
	}
	if location.X >= m.Size.Width-10 || location.Y >= m.Size.Height-10 {
		return false
	}
	return true
}

func (m *Map) IsColliding(Location Location) bool {
	for _, line := range m.Lines {
		if line.Start == Location || line.End == Location {
			return true
		}
	}
	return false
}
