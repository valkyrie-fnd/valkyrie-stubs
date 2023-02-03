package genericpam

import (
	"github.com/gofiber/fiber/v2"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

func RunServer(ds datastore.DataStore, config Config) *fiber.App {
	app := utils.HangingStart(config.Address, fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
	}, func(app *fiber.App) {
		ConfigureLogging(config.LogConfig)
		SetUpRoutes(app, ds, config.PamApiKey, config.ProviderTokens)

	})

	return app
}

func Routes(
	ds datastore.DataStore,
	pamApiKey string,
	providerTokens map[string]string) func(router fiber.Router) {
	return func(router fiber.Router) {
		SetUpRoutes(router, ds, pamApiKey, providerTokens)
	}
}

func SetUpRoutes(router fiber.Router, ds datastore.DataStore, pamApiKey string, providerTokens map[string]string) {
	controller := NewController(ds)
	controller.registerMiddlewares(router, pamApiKey, providerTokens)

	RegisterHandlers(router, NewStrictHandler(controller, nil))
}
