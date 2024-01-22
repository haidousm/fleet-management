package main

import (
	"fmt"
	"math/rand"

	"github.com/haidousm/fleets/internal/maps"
)

type Robot struct {
	Name      string
	Location  maps.Location
	FloorMap  maps.Map
	Direction int
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
	for !r.FloorMap.IsLocationValid(r.getLocation(increment, r.Direction)) || r.FloorMap.IsColliding(r.getLocation(increment, r.Direction)) {
		r.Direction = rand.Intn(4)
	}
	r.Location = r.getLocation(increment, r.Direction)
}

func (r *Robot) UpdateMap(floorMap maps.Map) {
	r.FloorMap = floorMap
}

func NewRobot(id int) Robot {
	return Robot{
		Name:      fmt.Sprintf("Robot %d", id),
		Location:  maps.Location{X: 5, Y: 5},
		Direction: 1,
	}
}
