package genericpam

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/valkyrie-fnd/valkyrie-stubs/transaction"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
)

type Amt float64

type Controller struct {
	dataStore datastore.DataStore
}

func NewController(ds datastore.DataStore) *Controller {
	return &Controller{dataStore: ds}
}

func (c Controller) GetSession(ctx context.Context, request GetSessionRequestObject) (GetSessionResponseObject, error) {
	session, err := c.dataStore.GetSession(ctx, request.Params.XPlayerToken)
	if err != nil {
		return GetSession200JSONResponse{
			Error: &PamError{
				Code:    PAMERRSESSIONNOTFOUND,
				Message: "session not found",
			},
			Status: ERROR,
		}, nil
	}

	account, err := c.dataStore.GetAccount(ctx, session.PlayerIdentifier, session.Currency)
	if err != nil {
		return GetSession200JSONResponse{
			Error: &PamError{
				Code:    PAMERRACCNOTFOUND,
				Message: "account not found",
			},
			Status: ERROR,
		}, nil
	}

	return GetSession200JSONResponse{
		Session: &Session{
			Country:  account.Country,
			Currency: account.Currency,
			Language: account.Language,
			PlayerId: account.PlayerIdentifier,
			Token:    request.Params.XPlayerToken,
			GameId:   session.GameID,
		},
		Status: OK,
	}, nil
}

func (c Controller) RefreshSession(ctx context.Context, request RefreshSessionRequestObject) (RefreshSessionResponseObject, error) {
	_, err := checkSessionValidity(ctx, c.dataStore, request.Params.XPlayerToken)
	if err != nil {
		return RefreshSession200JSONResponse{
			Error:  toPamError(err),
			Status: ERROR,
		}, nil
	}

	session, err := c.dataStore.UpdateSession(ctx, request.Params.XPlayerToken, utils.RandomString(32))
	if err != nil {
		return RefreshSession200JSONResponse{
			Error: &PamError{
				Code:    PAMERRSESSIONNOTFOUND,
				Message: "session update failed",
			},
			Status: ERROR,
		}, nil
	}

	account, err := c.dataStore.GetAccount(ctx, session.PlayerIdentifier, session.Currency)
	if err != nil {
		return RefreshSession200JSONResponse{
			Error: &PamError{
				Code:    PAMERRACCNOTFOUND,
				Message: "account not found",
			},
			Status: ERROR,
		}, nil
	}

	return RefreshSession200JSONResponse{
		Session: &Session{
			Country:  account.Country,
			Currency: account.Currency,
			Language: account.Language,
			PlayerId: account.PlayerIdentifier,
			Token:    session.Key,
		},
		Status: OK,
	}, nil
}

func (c Controller) GetBalance(ctx context.Context, request GetBalanceRequestObject) (GetBalanceResponseObject, error) {
	session, err := checkSessionValidity(ctx, c.dataStore, request.Params.XPlayerToken)
	if err != nil {
		return GetBalance200JSONResponse{
			Error:  toPamError(err),
			Status: ERROR,
		}, nil
	}

	account, err := c.dataStore.GetAccount(ctx, session.PlayerIdentifier, session.Currency)

	if err != nil {
		return GetBalance200JSONResponse{
			Error: &PamError{
				Code:    PAMERRACCNOTFOUND,
				Message: "failed to fetch account",
			},
			Status: ERROR,
		}, nil
	}

	return GetBalance200JSONResponse{
		Balance: &Balance{
			BonusAmount: Amount(account.BonusAmount),
			CashAmount:  Amount(account.CashAmount),
			PromoAmount: Amount(account.PromoAmount),
		},
		Status: OK,
	}, nil
}

type GetGameRound404JSONResponse GameRoundResponse

func (response GetGameRound404JSONResponse) VisitGetGameRoundResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(http.StatusNotFound)

	return ctx.JSON(&response)
}

func (c Controller) GetGameRound(ctx context.Context, request GetGameRoundRequestObject) (GetGameRoundResponseObject, error) {
	gameRound, err := c.dataStore.GetGameRound(ctx, request.PlayerId, request.ProviderGameRoundId)

	if err != nil || gameRound == nil {
		return GetGameRound404JSONResponse{
			Error: &PamError{
				Code:    PAMERRROUNDNOTFOUND,
				Message: "failed to fetch game round",
			},
			Status: ERROR,
		}, nil
	}
	return GetGameRound200JSONResponse{
		Gameround: &GameRound{
			ProviderGameId:  gameRound.ProviderGameID,
			ProviderRoundId: gameRound.ProviderRoundID,
			StartTime:       gameRound.StartTime,
			EndTime:         gameRound.EndTime,
		},
		Status: OK,
	}, nil
}

type GetTransactions404JSONResponse GetTransactionsResponse

func (response GetTransactions404JSONResponse) VisitGetTransactionsResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(http.StatusNotFound)

	return ctx.JSON(&response)
}

func (c Controller) GetTransactions(ctx context.Context, request GetTransactionsRequestObject) (GetTransactionsResponseObject, error) {
	var err error
	var trx []datastore.Transaction
	// Prioritize ref. If not there - go for ID
	if request.Params.ProviderBetRef != nil {
		trx, err = c.dataStore.GetTransactionsByRef(ctx, *request.Params.ProviderBetRef, request.Params.Provider)
		if len(trx) == 1 && trx[0].ID == 0 {
			trx = nil
		}
	} else {
		trx, err = c.dataStore.GetTransactionsByID(ctx, *request.Params.ProviderTransactionId, request.Params.Provider)
		if len(trx) == 1 && trx[0].ID == 0 {
			trx = nil
		}
	}

	if err != nil || trx != nil && len(trx) == 0 {
		return GetTransactions404JSONResponse{
			Error: &PamError{
				Code:    PAMERRTRANSNOTFOUND,
				Message: "transaction not found",
			},
			Status: ERROR,
		}, nil
	}

	transactions := make([]Transaction, 0)
	for _, trans := range trx {
		transactions = append(transactions, mapTransaction(request.Params.Provider, trans))
	}
	return GetTransactions200JSONResponse{
		Transactions: &transactions,
		Status:       OK,
	}, nil
}

func (c Controller) AddTransaction(ctx context.Context, request AddTransactionRequestObject) (AddTransactionResponseObject, error) {
	if string(request.Body.TransactionType) == "WITHDRAW" || string(request.Body.TransactionType) == "PROMOWITHDRAW" {
		_, err := checkSessionValidity(ctx, c.dataStore, request.Params.XPlayerToken)
		if err != nil {
			return AddTransaction200JSONResponse{
				Error:  toPamError(err),
				Status: ERROR,
			}, nil
		}
	}
	ts := transaction.NewTransactionService(ctx, c.dataStore)

	transactionID, err := ts.AddTransaction(datastore.Transaction{
		PlayerIdentifier:      request.PlayerId,
		CashAmount:            float64(request.Body.CashAmount),
		BonusAmount:           float64(request.Body.BonusAmount),
		PromoAmount:           float64(request.Body.PromoAmount),
		Currency:              request.Body.Currency,
		TransactionType:       string(request.Body.TransactionType),
		ProviderTransactionID: request.Body.ProviderTransactionId,
		ProviderBetRef:        request.Body.ProviderBetRef,
		ProviderGameID:        utils.OrZeroValue(request.Body.ProviderGameId),
		ProviderRoundID:       request.Body.ProviderRoundId,
		ProviderName:          request.Params.Provider,
		IsGameOver:            utils.OrZeroValue(request.Body.IsGameOver),
		TransactionDateTime:   request.Body.TransactionDateTime,
	})

	balance, balanceErr := c.handleGetBalanceByPlayer(ctx, request.PlayerId, request.Body.Currency)
	if balanceErr != nil {
		return AddTransaction200JSONResponse{
			Error:  toPamError(balanceErr),
			Status: ERROR,
		}, nil
	}

	if err != nil {
		return AddTransaction200JSONResponse{
			TransactionResult: &TransactionResult{
				Balance:       balance,
				TransactionId: nil,
			},
			Error:  c.handleTransactionError(err),
			Status: ERROR,
		}, nil
	}

	respTransID := strconv.Itoa(transactionID)
	return AddTransaction200JSONResponse{
		TransactionResult: &TransactionResult{
			Balance:       balance,
			TransactionId: &respTransID,
		},
		Status: OK,
	}, nil
}

func mapTransaction(provider string, transaction datastore.Transaction) Transaction {
	isGameOver := false
	transactionType := TransactionType(transaction.TransactionType)
	if transactionType == "DEPOSIT" {
		isGameOver = true
	}

	return Transaction{
		Currency:              transaction.Currency,
		Provider:              provider,
		TransactionType:       transactionType,
		CashAmount:            Amount(transaction.CashAmount),
		BonusAmount:           Amount(transaction.BonusAmount),
		PromoAmount:           Amount(transaction.PromoAmount),
		TransactionDateTime:   time.Now(),
		ProviderTransactionId: transaction.ProviderTransactionID,
		ProviderBetRef:        transaction.ProviderBetRef,
		ProviderGameId:        &transaction.ProviderGameID,
		ProviderRoundId:       transaction.ProviderRoundID,
		IsGameOver:            &isGameOver,
	}
}

func checkSessionValidity(ctx context.Context, ds datastore.DataStore, sessionToken string) (*datastore.Session, error) {
	session, err := ds.GetSession(ctx, sessionToken)
	if err != nil {
		return nil, PamError{
			Code:    PAMERRSESSIONNOTFOUND,
			Message: "session not found",
		}
	}
	if session.IsExpired() {
		return nil, PamError{
			Code:    PAMERRSESSIONEXPIRED,
			Message: "session expired",
		}
	}
	err = ds.TouchSession(ctx, sessionToken)
	if err != nil {
		return nil, PamError{
			Code:    PAMERRSESSIONNOTFOUND,
			Message: "session not found",
		}
	}
	return session, nil
}
func (c Controller) handleTransactionError(err error) *PamError {
	// Switch on common err to carve out correct error to report back
	var code ErrorCode
	switch {
	case errors.Is(err, transaction.NegAmountError):
		code = PAMERRNEGATIVESTAKE
	case errors.Is(err, transaction.CurrencyError):
		code = PAMERRTRANSCURRENCY
	case errors.Is(err, transaction.InputError), errors.Is(err, transaction.TransactionTypeError):
		code = PAMERRUNDEFINED
	case errors.Is(err, transaction.GameError):
		code = PAMERRGAMENOTFOUND
	case errors.Is(err, transaction.InsufficientFundsError):
		code = PAMERRCASHOVERDRAFT
	case errors.Is(err, transaction.CancelNonExistent):
		code = PAMERRCANCELNOTFOUND
	case errors.Is(err, transaction.CancelAlreadyCancelled):
		code = PAMERRTRANSALREADYCANCELLED
	case errors.Is(err, transaction.CancelAlreadySettled):
		code = PAMERRTRANSALREADYSETTLED
	case errors.Is(err, transaction.CancelNonWithdraw):
		code = PAMERRCANCELNONWITHDRAW
	case errors.Is(err, transaction.DepositAlreadyCancelled):
		code = PAMERRTRANSALREADYCANCELLED
	case errors.Is(err, transaction.DepositAlreadySettled):
		code = PAMERRTRANSALREADYSETTLED
	case errors.Is(err, transaction.DepositNotMatched):
		code = PAMERRTRANSNOTFOUND
	case errors.Is(err, transaction.AccountBlockedError):
		code = PAMERRBETNOTALLOWED
	case errors.Is(err, transaction.DuplicateTransactionError):
		code = PAMERRDUPLICATETRANS
	case errors.Is(err, transaction.BonusOverdraftError):
		code = PAMERRBONUSOVERDRAFT
	case errors.Is(err, transaction.PromoOverdraftError):
		code = PAMERRPROMOOVERDRAFT
	default:
		code = PAMERRUNDEFINED
	}

	return &PamError{
		Code:    code,
		Message: err.Error(),
	}
}

func (c Controller) handleGetBalanceByPlayer(ctx context.Context, playerIdentifier, currency string) (*Balance, error) {
	account, err := c.dataStore.GetAccount(ctx, playerIdentifier, currency)

	if err != nil {
		// Check that the player actually exists, If account is then not found we interpret it as currency account is missing
		if _, err = c.dataStore.GetPlayer(ctx, playerIdentifier); err != nil {
			return nil, PamError{
				Code:    PAMERRACCNOTFOUND,
				Message: "failed to fetch account",
			}
		}

		return nil, PamError{
			Code:    PAMERRTRANSCURRENCY,
			Message: "faulty currency",
		}
	}

	return &Balance{
		BonusAmount: Amount(account.BonusAmount),
		CashAmount:  Amount(account.CashAmount),
		PromoAmount: Amount(account.PromoAmount),
	}, nil
}
