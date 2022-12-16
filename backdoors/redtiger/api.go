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
	UserId   *string `json:"userId,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

type SessionResponse struct {
	Success bool   `json:"success"`
	Result  Result `json:"result"`
}

type Result struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

var Session = func(db datastore.ExtendedDatastore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SessionRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		var userId = utils.RandomString(10) // default generate userId
		if req.UserId != nil {
			userId = *req.UserId
		}
		var currency = "EUR" // default currency
		if req.Currency != nil {
			currency = *req.Currency
		}

		token := createPlayerAndSession(c.UserContext(), db, userId, currency)

		return c.Status(http.StatusOK).JSON(SessionResponse{
			Success: true,
			Result: Result{
				Token:  token,
				UserId: userId,
			},
		})
	}
}

func createPlayerAndSession(ctx context.Context, ds datastore.ExtendedDatastore, playerId, currency string) string {
	sessionToken := utils.RandomString(redtigerTokenSize)
	pla, err := ds.GetPlayer(ctx, playerId)
	if err != nil {
		ds.AddPlayer(datastore.Player{Id: utils.RandomInt(), PlayerIdentifier: playerId})
		pla, _ = ds.GetPlayer(ctx, playerId)
	}
	acc, err := ds.GetAccount(ctx, playerId, currency)
	if err != nil {
		ds.AddAccount(datastore.Account{Id: utils.RandomInt(), PlayerIdentifier: playerId, Currency: currency, CashAmount: 1000, PromoAmount: 1000})
		acc, _ = ds.GetAccount(ctx, playerId, currency)
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
