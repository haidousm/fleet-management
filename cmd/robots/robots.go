package main

import (
	"fmt"

	"github.com/haidousm/fleets/internal/maps"
)

type Robot struct {
	Name     string
	Location maps.Location
}

func (r *Robot) move(increment int, direction int) {
	switch direction {
	case 1:
		r.Location.X += increment
	case -1:
		r.Location.X -= increment
	case 2:
		r.Location.Y += increment
	case -2:
		r.Location.Y -= increment
	}
}

func NewRobot(id int) Robot {
	return Robot{
		Name:     fmt.Sprintf("Robot %d", id),
		Location: maps.Location{X: 0, Y: 0},
	}
}
