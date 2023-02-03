// Package backdoors contains functionality for tests to inject test data, for example to create a session.
package backdoors

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valkyrie-fnd/valkyrie-stubs/backdoors/evolution"
	"github.com/valkyrie-fnd/valkyrie-stubs/backdoors/redtiger"
	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

// BackdoorServer creates a new fiber app with the backdoor endpoints. Returns the
// app and the address it is listening on.
func BackdoorServer(eds datastore.ExtendedDatastore, addr string) (*fiber.App, string) {
	app := utils.HangingStart(addr, fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
	},
		func(app *fiber.App) {
			app.Post("/backdoors/evolution/sid", evolution.Sid(eds))
			app.Post("/backdoors/redtiger/session", redtiger.Session(eds))
			app.Post("/backdoors/datastore/session/reset", func(ctx *fiber.Ctx) error {
				eds.ClearSessionData()
				return nil
			})
			app.Post("/backdoors/session", Session(eds))
		},
	)

	return app, fmt.Sprintf("http://%s/backdoors/", addr)
}
