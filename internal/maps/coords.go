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
