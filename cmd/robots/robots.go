package main

import (
	"fmt"
	"math/rand"

	"github.com/haidousm/fleets/internal/maps"
)

type Robot struct {
	Name     string
	Location maps.Location
	FloorMap maps.Map
}

func (r *Robot) getLocation(increment int, direction int) maps.Location {
	new_location := r.Location
	switch direction {
	case 0:
		new_location.X += increment
	case 1:
		new_location.Y += increment
	case 2:
		new_location.X -= increment
	case 3:
		new_location.Y -= increment
	}
	return new_location
}

func (r *Robot) move(increment int) {
	direction := 1
	for !r.FloorMap.IsLocationValid(r.getLocation(increment, direction)) {
		direction = rand.Intn(4)
	}
	fmt.Printf("%d\n", direction)
	r.Location = r.getLocation(increment, direction)
}

func (r *Robot) UpdateMap(floorMap maps.Map) {
	r.FloorMap = floorMap
}

func NewRobot(id int) Robot {
	return Robot{
		Name:     fmt.Sprintf("Robot %d", id),
		Location: maps.Location{X: 5, Y: 5},
	}
}
