package main

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	ns, err := server.NewServer(&server.Options{
		Host: "localhost",
		Port: 4222,
	})
	if err != nil {
		slog.Error("failed to initialize NATS Server: " + err.Error())
	}
	ns.Start()
	defer ns.Shutdown()
	time.Sleep(time.Second)
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		slog.Error("failed to connect to NATS Server: " + err.Error())
	}
	defer nc.Close()
	greeter, err := micro.AddService(nc, micro.Config{
		Name:    "Greeter",
		Version: "1.0.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "hello",
			Handler: micro.HandlerFunc(func(req micro.Request) {
				slog.Info("Received hello request: " + string(req.Data()))
				req.Respond([]byte("Hello, " + string(req.Data()) + "!"))
			}),
		},
	})
	if err != nil {
		slog.Error("failed to add Greeter service: " + err.Error())
	}
	defer greeter.Stop()
	slog.Info("Created service: " + greeter.Info().Name + " (" + greeter.Info().ID + ")")
	for {
		time.Sleep(time.Second)
	}
}
