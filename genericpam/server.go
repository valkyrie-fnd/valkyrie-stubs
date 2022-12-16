package genericpam

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
)

func RunServer(ds datastore.DataStore, addr, pamApiKey string, providerTokens map[string]string) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
	})

	SetUpRoutes(app, ds, pamApiKey, providerTokens)

	go func() {
		err := app.Listen(addr)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start")
		}
	}()

	return app
}

func SetUpRoutes(router fiber.Router, ds datastore.DataStore, pamApiKey string, providerTokens map[string]string) {
	controller := NewController(ds)
	controller.registerMiddlewares(router, pamApiKey, providerTokens)

	RegisterHandlers(router, NewStrictHandler(controller, nil))
}
