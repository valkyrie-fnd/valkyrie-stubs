// Package broken provides a stub server which can be used to test error handling
package broken

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

var hardFaults = []string{"timeout 5s", "connection close"}

func RunServer(addr string, regularRoutes func(fiber.Router)) *fiber.App {
	hostname, _ := os.Hostname()
	log.Info().Msgf("Starting broken stub server. Admin on http://%s%s/broken", hostname, addr)

	// Initialize standard Go html template engine
	engine := html.New("./broken/views", ".html")

	app := utils.HangingStart(addr, fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
		Views:                 engine,
	},
		func(app *fiber.App) {
			errorCases := make(chan scenario, 2)
			signal := make(chan any)

			app.Hooks().OnShutdown(func() error {
				close(errorCases)
				close(signal)
				return nil
			})

			app.Group("/broken").
				Get("/", func(c *fiber.Ctx) error {
					return c.Render("list", fiber.Map{
						"scenarios":  predefinedScenarios,
						"hardFaults": hardFaults,
					})
				}).
				Post("/", addError(errorCases, signal))

			// Setup middleware which injects errors
			app.Use(injectFault(errorCases, signal))

			// Add regular routes
			regularRoutes(app)
		},
	)

	return app
}

func addError(q chan<- scenario, signal chan any) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hard := c.FormValue("hard")
		scenarioName := c.FormValue("scenario")

		if scenarioName != "" {
			s := predefinedScenarios[scenarioName]

			if hard != "" {
				s.HardError = hard
			}

			log.Info().Msgf("Queuing error: %s", s)

			q <- s

			<-signal // wait for completion signal
		}

		return c.Redirect("/broken")
	}
}

func injectFault(q chan scenario, signal chan<- any) fiber.Handler {
	return func(c *fiber.Ctx) error {
		select {
		case pe := <-q:
			if !pe.match(c.Request()) {
				// This error is not for this path
				q <- pe
				return c.Next()
			}
			log.Info().Msgf("Injecting fault: %s", pe)
			signal <- nil
			if pe.HardError != "" {
				return breakage(c, pe.HardError)
			} else {
				return c.Status(500).JSON(pe.Response)
			}
		default:
			return c.Next()
		}
	}
}
