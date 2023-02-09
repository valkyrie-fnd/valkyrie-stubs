// Package redtiger implements the "/session" endpoint, for the purpose of testing.
// spec: https://dev.redtigergaming.com/#!api/backend/session
package redtiger

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

const redtigerTokenSize = 32

type SessionRequest struct {
	UserID   *string `json:"userId,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

type SessionResponse struct {
	Result  Result `json:"result"`
	Success bool   `json:"success"`
}

type Result struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

var Session = func(db datastore.ExtendedDatastore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SessionRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		var userID = utils.RandomString(10) // default generate userId
		if req.UserID != nil {
			userID = *req.UserID
		}
		var currency = "EUR" // default currency
		if req.Currency != nil {
			currency = *req.Currency
		}

		token := createPlayerAndSession(c.UserContext(), db, userID, currency)

		return c.Status(http.StatusOK).JSON(SessionResponse{
			Success: true,
			Result: Result{
				Token:  token,
				UserID: userID,
			},
		})
	}
}

func createPlayerAndSession(ctx context.Context, ds datastore.ExtendedDatastore, playerID, currency string) string {
	sessionToken := utils.RandomString(redtigerTokenSize)
	pla, err := ds.GetPlayer(ctx, playerID)
	if err != nil {
		ds.AddPlayer(datastore.Player{ID: utils.RandomInt(), PlayerIdentifier: playerID})
		pla, _ = ds.GetPlayer(ctx, playerID)
	}
	acc, err := ds.GetAccount(ctx, playerID, currency)
	if err != nil {
		ds.AddAccount(datastore.Account{ID: utils.RandomInt(), PlayerIdentifier: playerID, Currency: currency, CashAmount: 1000, PromoAmount: 1000})
		acc, _ = ds.GetAccount(ctx, playerID, currency)
	}
	ds.AddSession(datastore.Session{
		Key:              sessionToken,
		PlayerIdentifier: pla.PlayerIdentifier,
		Provider:         "Red Tiger",
		Currency:         acc.Currency,
		Country:          acc.Country,
		Language:         acc.Language,
		Timestamp:        time.Now(),
		Timeout:          ds.GetSessionTimeout(),
	})
	return sessionToken
}
