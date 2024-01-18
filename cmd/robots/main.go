package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/haidousm/fleets/internal/mqtt"
)

func main() {
	num_robots := flag.Int("num_robots", 4, "number of robots")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	robots := make([]Robot, *num_robots)
	for i := 0; i < *num_robots; i++ {
		robots[i] = NewRobot(i)
	}

	simulateRobots(ctx, robots)
}

func simulateRobots(ctx context.Context, robots []Robot) {
	robotsLocationTopic := "robots/locations"
	mqttClient := mqtt.Client(ctx, robotsLocationTopic)
	for _, robot := range robots {
		go simulateRobot(ctx, robot, mqttClient, robotsLocationTopic)
	}
	<-mqttClient.Done()
}

func simulateRobot(ctx context.Context, robot Robot, mqttClient *autopaho.ConnectionManager, robotsLocationTopic string) {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			direction := rand.Intn(4) - 2
			robot.move(5, direction)

			jsonBytes, err := json.Marshal(robot)
			if err != nil {
				fmt.Printf("failed to marshal robot: %s\n", err)
			}
			_, err = mqttClient.Publish(context.Background(), &paho.Publish{
				Topic:   robotsLocationTopic,
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
}
