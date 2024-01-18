package mqtt

import (
	"context"
	"fmt"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
)

func Client(ctx context.Context) *autopaho.ConnectionManager {
	u, err := url.Parse("mqtt://localhost:1883")
	if err != nil {
		panic(err)
	}

	cliCfg := autopaho.ClientConfig{
		ServerUrls:     []*url.URL{u},
		KeepAlive:      20,
		OnConnectError: func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
	}

	c, err := autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		panic(err)
	}

	if err = c.AwaitConnection(ctx); err != nil {
		panic(err)
	}

	return c
}
