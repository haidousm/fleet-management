package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/haidousm/fleets/internal/maps"
	"github.com/haidousm/fleets/internal/mqtt"
)

func main() {
	num_robots := flag.Int("num_robots", 4, "number of robots")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	robots := []*Robot{}
	for i := 0; i < *num_robots; i++ {
		r := NewRobot(i)
		robots = append(robots, &r)
	}

	mqttClient := mqtt.Client(ctx)
	updateRobotsMap(ctx, mqttClient, robots)
	simulateRobots(ctx, mqttClient, robots)
	<-mqttClient.Done()
}

func updateRobotsMap(ctx context.Context, client *autopaho.ConnectionManager, robots []*Robot) {
	floorMapTopic := "maps/floor"
	client.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{
				Topic: floorMapTopic,
				QoS:   0,
			},
		},
	})
	client.AddOnPublishReceived(func(pr autopaho.PublishReceived) (bool, error) {
		var floorMap maps.Map

		bytes := pr.Packet.Payload
		if err := json.Unmarshal(bytes, &floorMap); err != nil {
			fmt.Printf("failed to unmarshal floor map: %s\n", err)
			return false, nil
		}
		for _, robot := range robots {
			robot.UpdateMap(floorMap)
		}
		return true, nil
	})
}

func simulateRobots(ctx context.Context, client *autopaho.ConnectionManager, robots []*Robot) {
	robotsLocationTopic := "robots/locations"
	robotsLocationPub := mqtt.Client(ctx)
	for _, robot := range robots {
		go simulateRobot(ctx, client, robot, robotsLocationTopic)
	}
	<-robotsLocationPub.Done()
}

func simulateRobot(ctx context.Context, client *autopaho.ConnectionManager, robot *Robot, robotsLocationTopic string) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			robot.move(1)
			jsonBytes, err := json.Marshal(robot)
			if err != nil {
				fmt.Printf("failed to marshal robot: %s\n", err)
			}
			_, err = client.Publish(context.Background(), &paho.Publish{
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
