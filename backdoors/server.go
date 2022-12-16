// Package backdoors contains functionality for tests to inject test data, for example to create a session.
package backdoors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valkyrie-fnd/valkyrie-stubs/backdoors/evolution"
	"github.com/valkyrie-fnd/valkyrie-stubs/backdoors/redtiger"
	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
)

func BackdoorServer(eds datastore.ExtendedDatastore, addr string) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
	})

	app.Post("/backdoors/evolution/sid", evolution.Sid(eds))
	app.Post("/backdoors/redtiger/session", redtiger.Session(eds))
	app.Post("/backdoors/datastore/session/reset", func(ctx *fiber.Ctx) error {
		eds.ClearSessionData()
		return nil
	})
	app.Post("/backdoors/session", Session(eds))

	go func() {
		err := app.Listen(addr)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start")
		}
	}()

	return app
}
