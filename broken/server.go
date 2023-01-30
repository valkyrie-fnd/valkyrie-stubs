package broken

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/valkyrie-fnd/valkyrie-stubs/genericpam"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

var softErrors = []genericpam.PamError{
	{Code: genericpam.PAMERRACCNOTFOUND, Message: "Account not found"},
	{Code: genericpam.PAMERRTIMEOUT, Message: "Timeout"},
}

var hardErrors = []string{"timeout 5s", "connection close"}

func RunServer(addr string, regularRoutes func(fiber.Router)) *fiber.App {
	hostname, _ := os.Hostname()
	log.Info().Msgf("Starting broken stub server. Admin on http://%s%s/broken", hostname, addr)

	// Initialize standard Go html template engine
	engine := html.New("./broken/views", ".html")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true, // since we store values in-memory after handlers have returned
		Views:                 engine,
	})

	pamErrors := make(chan genericpam.PamError)
	breakingErrors := make(chan string)
	app.Hooks().OnShutdown(func() error {
		close(pamErrors)
		close(breakingErrors)
		return nil
	})

	app.Group("/broken").Get("/", func(c *fiber.Ctx) error {
		return c.Render("list", fiber.Map{
			"softErrors": softErrors,
			"hardErrors": hardErrors,
		})
	}).
		Get("error", softErrorRoute(pamErrors)).
		Get("hard", hardErrorRoute(breakingErrors))

	// Setup middleware which injects errors
	app.Use(injectPAMError(pamErrors))
	app.Use(injectBreakage(breakingErrors))

	// Add regular routes
	regularRoutes(app)

	go func() {
		err := app.Listen(addr)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start")
		}
	}()

	return app
}

func softErrorRoute(pamErrors chan<- genericpam.PamError) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code != "" {
			for _, e := range softErrors {
				if string(e.Code) == code {
					log.Info().Msgf("Queuing PAM error: %s", e)
					pamErrors <- e
					return c.Redirect("/broken")
				}
			}
		} else {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.Redirect("/broken")
	}
}

func hardErrorRoute(breakingErrors chan<- string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code != "" {
			for _, e := range hardErrors {
				if string(e) == code {
					log.Info().Msgf("Queuing hard error: %s", e)
					breakingErrors <- e
					return c.Redirect("/broken")
				}
			}
		} else {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.Redirect("/broken")
	}
}

func injectPAMError(q chan genericpam.PamError) fiber.Handler {
	return func(c *fiber.Ctx) error {
		select {
		case err := <-q:
			log.Info().Str("error", err.Message).Msg("Injecting error")
			return c.Status(500).JSON(err)
		default:
			return c.Next()
		}
	}
}

func injectBreakage(q chan string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		select {
		case err := <-q:
			return breakage(c, err)
		default:
			return c.Next()
		}
	}
}
