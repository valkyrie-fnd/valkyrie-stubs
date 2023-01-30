package genericpam

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
)

func RunServer(ds datastore.DataStore, config Config) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
	})

	ConfigureLogging(config.LogConfig)
	SetUpRoutes(app, ds, config.PamApiKey, config.ProviderTokens)

	go func() {
		err := app.Listen(config.Address)
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
