// Package transaction validate and add transaction to datastore
package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
)

type TransactionService struct {
	ds  datastore.DataStore
	ctx context.Context
}

func NewTransactionService(ctx context.Context, ds datastore.DataStore) *TransactionService {
	return &TransactionService{ds, ctx}
}

/*
Add adds transaction to provided datastore and updates balance.

Returns validation error in case add fails
*/
func (ts *TransactionService) AddTransaction(t datastore.Transaction) (int, error) {
	// Check transaction type
	if !isValidTransactionType(t.TransactionType) {
		return 0, TransactionTypeError
	}

	if err := ts.validateAccount(t); err != nil {
		return 0, err
	}

	if trans, err := ts.getPreviousTransaction(&t, true); err != nil {
		return 0, err
	} else if trans != nil {
		// If exactly the same transaction already exists (trans-ID and trans-type), simply return
		return trans.ID, nil
	}

	// Validate general transaction logic

	if err := ts.validateTransaction(&t); err != nil {
		return 0, err
	}

	// Check if game exist
	if t.TransactionType != PROMODEPOSIT {
		if _, err := ts.ds.GetGame(ts.ctx, t.ProviderGameID, t.ProviderName); err != nil {
			return 0, fmt.Errorf("%w - %s", GameError, err.Error())
		}
	}

	if err := ts.ds.AddTransaction(ts.ctx, &t); err != nil {
		return 0, err
	}

	if err := ts.updateAccountBalance(t); err != nil {
		return 0, err
	}

	if t.IsGameOver {
		n := time.Now()
		if err := ts.ds.EndGameRound(ts.ctx, datastore.GameRound{
			ProviderName: t.ProviderName, ProviderGameID: t.ProviderGameID, ProviderRoundID: *t.ProviderRoundID, EndTime: &n,
		}); err != nil {
			return 0, err
		}
	} else if t.ProviderRoundID != nil {
		if err := ts.ds.AddGameRound(ts.ctx, datastore.GameRound{
			ProviderName: t.ProviderName, PlayerID: t.PlayerIdentifier, ProviderGameID: t.ProviderGameID, ProviderRoundID: *t.ProviderRoundID, StartTime: time.Now(), EndTime: nil,
		}); err != nil {
			return 0, err
		}
	}
	return t.ID, nil
}

func (ts *TransactionService) updateAccountBalance(trans datastore.Transaction) error {

	if err := ts.updateAccountBalanceType(trans, trans.CashAmount, datastore.Cash); err != nil {
		return err
	}

	if err := ts.updateAccountBalanceType(trans, trans.BonusAmount, datastore.Bonus); err != nil {
		return err
	}

	if err := ts.updateAccountBalanceType(trans, trans.PromoAmount, datastore.Promo); err != nil {
		return err
	}

	return nil
}

func (ts *TransactionService) updateAccountBalanceType(trans datastore.Transaction, amount float64, balanceType datastore.BalanceType) error {
	if amount != 0.0 {
		transType := trans.TransactionType
		if transType == WITHDRAW || transType == PROMOWITHDRAW {
			amount = -1 * amount
		}
		if _, err := ts.ds.UpdateAccountBalance(ts.ctx, trans.PlayerIdentifier, trans.Currency, balanceType, amount); err != nil {
			return err
		}
	}
	return nil
}
