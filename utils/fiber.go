package utils

import (
	"errors"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Routing = func(app *fiber.App)

// waitUntil starts Fiber app and waits until it is ready to accept connections.
func HangingStart(adr string, cfg fiber.Config, rt Routing) *fiber.App {

	app := fiber.New(cfg)

	// register routes
	rt(app)

	// register startup listener
	var wg sync.WaitGroup
	wg.Add(1)
	app.Hooks().OnListen(func(_ fiber.ListenData) error {
		wg.Done()
		return nil
	})

	go func() {
		err := app.Listen(adr)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start")
		}
	}()

	if err := waitUntil(&wg, 3*time.Second); err != nil {
		log.Fatal().Err(err).Send()
	}

	return app
}

func waitUntil(wg *sync.WaitGroup, timeout time.Duration) error {
	rdy := make(chan struct{})

	go func() {
		defer close(rdy)
		wg.Wait()
	}()

	select {
	case <-rdy:
		return nil
	case <-time.After(timeout):
		return errors.New("timeout waiting for fiber app to start")
	}
}
