// Package broken provides a stub server which can be used to test error handling
package broken

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/valkyrie-fnd/valkyrie-stubs/broken/views"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

var hardFaults = []string{"timeout 5s", "connection close"}

func RunServer(addr string, regularRoutes func(fiber.Router)) *fiber.App {
	hostname, _ := os.Hostname()
	log.Info().Msgf("Starting broken stub server. Admin on http://%s%s/broken/", hostname, addr)

	// Initialize standard Go html template engine
	engine := html.NewFileSystem(http.FS(views.Content), ".html")

	app := utils.HangingStart(addr, fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
		Views:                 engine,
	},
		func(app *fiber.App) {
			errorCases := make(chan scenario, 256)

			app.Hooks().OnShutdown(func() error {
				close(errorCases)
				return nil
			})

			bg := app.Group("/broken")
			bg.Get("/", func(c *fiber.Ctx) error {
				return c.Render("list", fiber.Map{
					"scenarios":  predefinedScenarios,
					"hardFaults": hardFaults,
				})
			})
			bg.Post("/", addError(errorCases))

			// Setup middleware which injects errors
			app.Use(injectFault(errorCases))

			// Add regular routes
			regularRoutes(app)
		},
	)

	return app
}

func addError(q chan<- scenario) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hard := c.FormValue("hard")
		scenarioName := c.FormValue("scenario")

		if scenarioName != "" {
			s := predefinedScenarios[scenarioName]

			if hard != "" {
				s.HardError = hard
			}

			select {
			case q <- s:
				log.Info().Msgf("Queued error: %s", s)
			default:
				log.Warn().Msg("Unable to queue error, channel full")
			}

		}

		return c.Redirect("/broken/")
	}
}

func injectFault(q chan scenario) fiber.Handler {
	return func(c *fiber.Ctx) error {
		select {
		case pe := <-q:
			if !pe.match(c.Request()) {
				// This error is not for this path
				q <- pe
				return c.Next()
			}
			log.Info().Msgf("Injecting fault: %s", pe)
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
