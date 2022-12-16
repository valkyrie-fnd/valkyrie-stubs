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
	UserId      *string  `json:"userId,omitempty"`
	Currency    *string  `json:"currency,omitempty"`
	CashAmount  *float64 `json:"cashAmount,omitempty"`
	BonusAmount *float64 `json:"bonusAmount,omitempty"`
	PromoAmount *float64 `json:"promoAmount,omitempty"`
	IsBlocked   *bool    `json:"isBlocked,omitempty"`
	Provider    string   `json:"provider"`
	GameId      *string  `json:"gameId"`
}

type SessionResponse struct {
	Success bool   `json:"success"`
	Result  Result `json:"result"`
}

type Result struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

var Session = func(ds datastore.ExtendedDatastore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SessionRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		if req.UserId == nil {
			req.UserId = utils.Ptr(utils.RandomString(10)) // default generate userid
		}
		if req.Currency == nil {
			req.Currency = utils.Ptr("EUR") // default currency
		}
		if req.GameId == nil {
			req.GameId = utils.Ptr(utils.RandomString(10)) // default just generate gameid
		}

		_, err := ds.GetPlayer(c.UserContext(), *req.UserId)
		if err != nil {
			ds.AddPlayer(datastore.Player{Id: utils.RandomInt(), PlayerIdentifier: *req.UserId})
		}
		acc := upsertAccount(c.UserContext(), ds, req)
		token := createSession(ds, req.Provider, *req.GameId, acc)

		return c.Status(http.StatusOK).JSON(SessionResponse{
			Success: true,
			Result: Result{
				Token:  token,
				UserId: *req.UserId,
			},
		})
	}
}

// upsertAccount creates or updates account
func upsertAccount(ctx context.Context, ds datastore.ExtendedDatastore, req SessionRequest) datastore.Account {
	currAccount, err := ds.GetAccount(ctx, *req.UserId, *req.Currency)
	if err != nil {
		newAccount := datastore.Account{
			Id:               utils.RandomInt(),
			PlayerIdentifier: *req.UserId,
			Currency:         *req.Currency,
			CashAmount:       utils.OrDefault(req.CashAmount, 1000.0),
			BonusAmount:      utils.OrDefault(req.BonusAmount, 0.0),
			PromoAmount:      utils.OrDefault(req.PromoAmount, 0.0),
			IsBlocked:        utils.OrDefault(req.IsBlocked, false),
		}
		ds.AddAccount(newAccount)
		return newAccount
	} else {
		currAccount.CashAmount = utils.OrDefault(req.CashAmount, currAccount.CashAmount)
		currAccount.BonusAmount = utils.OrDefault(req.BonusAmount, currAccount.BonusAmount)
		currAccount.PromoAmount = utils.OrDefault(req.PromoAmount, currAccount.PromoAmount)
		currAccount.IsBlocked = utils.OrDefault(req.IsBlocked, currAccount.IsBlocked)

		_ = ds.UpdateAccount(currAccount.PlayerIdentifier, *currAccount)
		return *currAccount
	}
}

func createSession(ds datastore.ExtendedDatastore, provider, gameId string, acc datastore.Account) string {
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
		GameId:           utils.Ptr(gameId),
	})
	return sessionToken
}
