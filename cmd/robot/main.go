package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/paho"
)

type Location struct {
	X, Y int
}

type Robot struct {
	Name     string
	Location Location
}

const NumRobots = 10
const topic = "robots/locations"

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

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mqttClient := MQTTClient(ctx, topic)

	robots := make([]Robot, NumRobots)

	for i := 0; i < NumRobots; i++ {
		robots[i] = Robot{
			Name:     fmt.Sprintf("Robot %d", i),
			Location: Location{X: 0, Y: 0},
		}
	}

	for _, robot := range robots {
		go func(robot Robot) {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					direction := rand.Intn(4) - 2
					robot.move(1, direction)

					jsonBytes, err := json.Marshal(robot)
					if err != nil {
						fmt.Printf("failed to marshal robot: %s\n", err)
					}
					_, err = mqttClient.Publish(context.Background(), &paho.Publish{
						Topic:   topic,
						QoS:     0,
						Payload: jsonBytes,
					})
					if err != nil {
						fmt.Printf("failed to publish message: %s\n", err)
					}
				case <-ctx.Done():
					return
				}
			}
		}(robot)
	}
	<-mqttClient.Done()
}
