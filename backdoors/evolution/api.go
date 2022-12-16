// Package evolution implements the "/sid" endpoint, expected by Evolution for testing.
// https://studio.evolutiongaming.com/api/evo-std-rest/docs/index.html#sid
package evolution

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

type sidRequest struct {
	Sid     string `json:"sid"`
	UserID  string `json:"userId" binding:"required"`
	UUID    string `json:"uuid"`
	Channel struct {
		Type string `json:"type"`
	} `json:"channel"`
}

type sidResponse struct {
	Status string `json:"status"`
	Sid    string `json:"sid,omitempty"`
	UUID   string `json:"uuid,omitempty"`
}

var Sid = func(eds datastore.ExtendedDatastore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get expected Evolution API Key
		res, err := eds.GetProviderApiKey("Evolution")
		if err != nil {
			log.Fatalf("evolution api key not configured in datastore: %v", err)
		}
		// Verify the API Key
		if token := c.Query("authToken"); token == "" || token != res.ApiKey {
			return c.Status(http.StatusUnauthorized).JSON(sidResponse{
				Status: "UNKNOWN_ERROR",
			})
		}

		var req sidRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(sidResponse{
				Status: "INVALID_PARAMETER",
			})
		}

		// We made to actual logic - invent a session and stick it in the datastore
		sid := createPlayerAndSession(c.UserContext(), eds, req.UserID)

		return c.Status(http.StatusOK).JSON(sidResponse{
			Sid:    sid,
			Status: "OK",
			UUID:   req.UUID,
		})
	}
}

func createPlayerAndSession(ctx context.Context, eds datastore.ExtendedDatastore, playerId string) string {
	sessionToken := rndStr()
	pla, err := eds.GetPlayer(ctx, playerId)
	if err != nil {
		eds.AddPlayer(datastore.Player{Id: rndInt(), PlayerIdentifier: playerId})
		pla, _ = eds.GetPlayer(ctx, playerId)
	}
	acc, err := eds.GetAccount(ctx, playerId, "EUR")
	if err != nil {
		eds.AddAccount(datastore.Account{Id: rndInt(), PlayerIdentifier: playerId, Currency: "EUR", CashAmount: 1000})
		acc, _ = eds.GetAccount(ctx, playerId, "EUR")
	}
	eds.AddSession(datastore.Session{
		Key:              sessionToken,
		PlayerIdentifier: pla.PlayerIdentifier,
		Provider:         "Evolution",
		Currency:         acc.Currency,
		Country:          acc.Country,
		Language:         acc.Language,
		Timestamp:        time.Now(),
		Timeout:          eds.GetSessionTimeout(),
	})
	return sessionToken
}

func rndStr() string {
	return utils.RandomString(10)
}

func rndInt() int {
	return utils.RandomInt()
}
