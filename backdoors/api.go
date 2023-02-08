package backdoors

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

type SessionRequest struct {
	IsBlocked   *bool    `json:"isBlocked,omitempty"`
	CashAmount  *float64 `json:"cashAmount,omitempty"`
	BonusAmount *float64 `json:"bonusAmount,omitempty"`
	PromoAmount *float64 `json:"promoAmount,omitempty"`
	UserID      *string  `json:"userId,omitempty"`
	GameID      *string  `json:"gameId"`
	Currency    *string  `json:"currency,omitempty"`
	Provider    string   `json:"provider"`
}

type SessionResponse struct {
	Result  Result `json:"result"`
	Success bool   `json:"success"`
}

type Result struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

var Session = func(ds datastore.ExtendedDatastore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SessionRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		if req.UserID == nil {
			req.UserID = utils.Ptr(utils.RandomString(10)) // default generate userID
		}
		if req.Currency == nil {
			req.Currency = utils.Ptr("EUR") // default currency
		}
		if req.GameID == nil {
			req.GameID = utils.Ptr(utils.RandomString(10)) // default just generate gameID
		}

		_, err := ds.GetPlayer(c.UserContext(), *req.UserID)
		if err != nil {
			ds.AddPlayer(datastore.Player{ID: utils.RandomInt(), PlayerIdentifier: *req.UserID})
		}
		acc := upsertAccount(c.UserContext(), ds, req)
		token := createSession(ds, req.Provider, *req.GameID, acc)

		return c.Status(http.StatusOK).JSON(SessionResponse{
			Success: true,
			Result: Result{
				Token:  token,
				UserID: *req.UserID,
			},
		})
	}
}

// upsertAccount creates or updates account
func upsertAccount(ctx context.Context, ds datastore.ExtendedDatastore, req SessionRequest) datastore.Account {
	account, err := ds.GetAccount(ctx, *req.UserID, *req.Currency)
	if err != nil {
		newAccount := datastore.Account{
			ID:               utils.RandomInt(),
			PlayerIdentifier: *req.UserID,
			Currency:         *req.Currency,
			CashAmount:       utils.OrDefault(req.CashAmount, 1000.0),
			BonusAmount:      utils.OrDefault(req.BonusAmount, 0.0),
			PromoAmount:      utils.OrDefault(req.PromoAmount, 0.0),
			IsBlocked:        utils.OrDefault(req.IsBlocked, false),
		}
		ds.AddAccount(newAccount)
		return newAccount
	} else {
		account.CashAmount = utils.OrDefault(req.CashAmount, account.CashAmount)
		account.BonusAmount = utils.OrDefault(req.BonusAmount, account.BonusAmount)
		account.PromoAmount = utils.OrDefault(req.PromoAmount, account.PromoAmount)
		account.IsBlocked = utils.OrDefault(req.IsBlocked, account.IsBlocked)

		_ = ds.UpdateAccount(account.PlayerIdentifier, *account)
		return *account
	}
}

func createSession(ds datastore.ExtendedDatastore, provider, gameID string, acc datastore.Account) string {
	sessionToken := utils.RandomString(32)
	ds.AddSession(datastore.Session{
		Key:              sessionToken,
		PlayerIdentifier: acc.PlayerIdentifier,
		Provider:         provider,
		Currency:         acc.Currency,
		Country:          acc.Country,
		Language:         acc.Language,
		Timestamp:        time.Now(),
		Timeout:          ds.GetSessionTimeout(),
		GameID:           utils.Ptr(gameID),
	})
	return sessionToken
}
