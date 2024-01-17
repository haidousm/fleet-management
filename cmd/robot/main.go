package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Location struct {
	x, y int
}

type Robot struct {
	name     string
	location Location
}

const NumRobots = 10

func (r *Robot) move(increment int, direction int) {
	switch direction {
	case 1:
		r.location.x += increment
	case -1:
		r.location.x -= increment
	case 2:
		r.location.y += increment
	case -2:
		r.location.y -= increment
	}
}

func main() {
	robots := make([]Robot, NumRobots)

	for i := 0; i < NumRobots; i++ {
		robots[i] = Robot{
			name:     fmt.Sprintf("Robot %d", i),
			location: Location{x: 0, y: 0},
		}
	}

	quit := make(chan struct{})
	for _, robot := range robots {
		go func(robot Robot) {
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ticker.C:
					direction := rand.Intn(4) - 2
					robot.move(1, direction)
					fmt.Println(robot)
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}(robot)
	}

	time.Sleep(10 * time.Second)
	close(quit)
}
