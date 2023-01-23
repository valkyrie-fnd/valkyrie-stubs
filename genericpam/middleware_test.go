package genericpam

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/memorydatastore"
)

func Test_CheckPlayerToken(t *testing.T) {
	tests := []struct {
		name     string
		session  string
		provider string
		status   int
	}{
		{
			"success",
			"key",
			"provider",
			http.StatusOK,
		},
		{
			"unauthorized invalid session",
			"invalid",
			"provider",
			http.StatusUnauthorized,
		},
		{
			"unauthorized invalid provider",
			"key",
			"invalid",
			http.StatusUnauthorized,
		},
		{
			"success provider recon token",
			"recon_token",
			"provider",
			http.StatusOK,
		},
	}

	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Sessions: []datastore.Session{{Key: "key", Provider: "provider"}},
	})
	app := fiber.New()
	app.Use(Controller{ds}.getCheckPlayerToken(map[string]string{"provider": "recon_token"})).
		All("", func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(http.StatusOK)
		})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?%s=%s", QueryProvider, test.provider), nil)
			req.Header.Set(XPlayerToken, test.session)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, test.status, resp.StatusCode)
		})
	}
}

func Test_CheckPamApiToken(t *testing.T) {
	tests := []struct {
		name      string
		pamAPIKey string
		status    int
	}{
		{
			"success",
			"pam_token",
			http.StatusOK,
		},
		{
			"unauthorized api key",
			"invalid",
			http.StatusUnauthorized,
		},
	}

	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{})
	app := fiber.New()
	app.Use(Controller{ds}.getCheckPamApiToken("pam_token")).
		All("", func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(http.StatusOK)
		})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", test.pamAPIKey))
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, test.status, resp.StatusCode)
		})
	}
}

func Test_checkPlayerId(t *testing.T) {
	tests := []struct {
		name      string
		playerId  string
		errorCode ErrorCode
	}{
		{
			"success",
			"id",
			"OK",
		},
		{
			"player not found",
			"invalid",
			PAMERRPLAYERNOTFOUND,
		},
	}

	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{
				Id:               1,
				PlayerIdentifier: "id",
			},
		},
	})
	app := fiber.New()
	app.Use("/players/:playerId/+", Controller{ds}.checkPlayerId).
		All("/players/:playerId/+", func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(http.StatusOK)
		})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s/foo", test.playerId), nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			responseBody, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(responseBody), test.errorCode)
		})
	}
}

func Test_checkProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		status   int
	}{
		{
			"success",
			"provider",
			http.StatusOK,
		},
		{
			"provider not found",
			"invalid",
			http.StatusBadRequest,
		},
	}

	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{Providers: []datastore.Provider{{
		ProviderId: 0,
		Provider:   "provider",
	}}})
	app := fiber.New()
	app.Use(Controller{ds}.checkProvider).
		All("", func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(http.StatusOK)
		})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?%s=%s", QueryProvider, test.provider), nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, test.status, resp.StatusCode)
		})
	}
}
