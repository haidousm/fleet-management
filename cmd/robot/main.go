package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/paho"
)

type Location struct {
	x, y int
}

type Robot struct {
	name     string
	location Location
}

const NumRobots = 10
const topic = "robots/locations"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mqttClient := MQTTClient(ctx, topic)

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
					_, err := mqttClient.Publish(context.Background(), &paho.Publish{
						Topic: topic,
						QoS:   0,
						Payload: []byte(fmt.Sprintf("%s is at (%d, %d)",
							robot.name, robot.location.x, robot.location.y)),
					})
					if err != nil {
						fmt.Printf("failed to publish message: %s\n", err)
					}
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
