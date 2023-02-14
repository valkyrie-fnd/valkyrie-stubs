// Package provider provides a stubbed implementation of provider endpoints.
package provider

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/valkyrie-fnd/valkyrie-stubs/provider/caleta"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

type providerEndpointsStubs struct {
	caletaSignatureVerifier SignatureVerifier
	cannedTransactions      func() (caleta.TransactionResponse, error)
}

type Option func(*providerEndpointsStubs)

func Create(ctx context.Context, addr string, options ...Option) {
	p := &providerEndpointsStubs{
		caletaSignatureVerifier: func(signature string, payload []byte) error { return nil },
	}

	for _, option := range options {
		option(p)
	}

	app := utils.HangingStart(addr, fiber.Config{
		DisableStartupMessage: true,
		Immutable:             true,
	}, func(app *fiber.App) {

		app.Get("/evo/game/launch", evoGameLaunch())
		app.Get("/caleta/api/game/url", caletaGameLaunch())
		app.Get("/caleta/api/transactions/round", caletaGetTransactions(p))
	})

	go func() {
		<-ctx.Done()
		_ = app.Shutdown()
	}()
}

func WithCaletaSignatureVerifier(v SignatureVerifier) Option {
	return func(p *providerEndpointsStubs) {
		p.caletaSignatureVerifier = v
	}
}

func WithCannedTransactions(v func() (caleta.TransactionResponse, error)) Option {
	return func(p *providerEndpointsStubs) {
		p.cannedTransactions = v
	}
}

func evoGameLaunch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	}
}

type SignatureVerifier func(signature string, payload []byte) error

func caletaGameLaunch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	}
}

func caletaGetTransactions(p *providerEndpointsStubs) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t, err := p.cannedTransactions()
		if err != nil {
			return err
		}
		return c.JSON(t)
	}
}
