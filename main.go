package main

import (
	"log/slog"
	"strings"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

var (
	greeters = map[string]func(req micro.Request){
		"en": func(req micro.Request) {
			req.Respond([]byte("Hello, " + string(req.Data()) + "!"))
		},
		"es": func(req micro.Request) {
			req.Respond([]byte("¡Hola, " + string(req.Data()) + "!"))
		},
		"fr": func(req micro.Request) {
			req.Respond([]byte("Bonjour, " + string(req.Data()) + "!"))
		},
		"de": func(req micro.Request) {
			req.Respond([]byte("Hallo, " + string(req.Data()) + "!"))
		},
		"it": func(req micro.Request) {
			req.Respond([]byte("Ciao, " + string(req.Data()) + "!"))
		},
		"zh": func(req micro.Request) {
			req.Respond([]byte("你好，" + string(req.Data()) + "!"))
		},
		"ja": func(req micro.Request) {
			req.Respond([]byte("こんにちは、" + string(req.Data()) + "!"))
		},
		"ko": func(req micro.Request) {
			req.Respond([]byte("안녕하세요, " + string(req.Data()) + "!"))
		},
		"ar": func(req micro.Request) {
			req.Respond([]byte("مرحبًا، " + string(req.Data()) + "!"))
		},
		"hi": func(req micro.Request) {
			req.Respond([]byte("नमस्ते, " + string(req.Data()) + "!"))
		},
		"philly": func(req micro.Request) {
			req.Respond([]byte("Yo, " + string(req.Data()) + "! What's good, cuz?"))
		},
	}

	logReqs = func(next micro.HandlerFunc) micro.HandlerFunc {
		return func(req micro.Request) {
			slog.Info("Received " + strings.TrimPrefix(req.Subject(), "greet.") + " request. Data: " + string(req.Data()))
			next(req)
		}
	}
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
		Name:    "greet",
		Version: "0.1.0",
	})
	if err != nil {
		slog.Error("failed to add Greeter service: " + err.Error())
	}
	greet := greeter.AddGroup("greet")
	for lang, handler := range greeters {
		greet.AddEndpoint(lang, logReqs(handler))
	}
	defer greeter.Stop()
	slog.Info("Created service: " + greeter.Info().Name + " (" + greeter.Info().ID + ")")
	for {
		time.Sleep(time.Second)
	}
}
