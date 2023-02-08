package transaction

import (
	"errors"
	"sort"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

var NegAmountError = errors.New("negative amount")
var InsufficientFundsError = errors.New("insufficient funds")
var CurrencyError = errors.New("mismatching currencies")
var TransactionTypeError = errors.New("bad trans type")
var GameError = errors.New("game not found")
var InputError = errors.New("bad input")
var PlayerNotFoundError = errors.New("player not found")
var CancelNonExistent = errors.New("can not cancel non existent transaction")
var CancelNonWithdraw = errors.New("can not cancel non withdraw transaction")
var CancelAlreadyCancelled = errors.New("can not cancel already cancelled")
var CancelAlreadySettled = errors.New("can not cancel already settled")
var DepositNotMatched = errors.New("can not deposit matching withdraw")
var DepositAlreadySettled = errors.New("can not deposit settled round")
var DepositAlreadyCancelled = errors.New("can not deposit cancelled round")
var AccountBlockedError = errors.New("account is blocked")
var SessionExpiredError = errors.New("session is expired")
var DuplicateTransactionError = errors.New("duplicate transaction")
var PromoOverdraftError = errors.New("promo overdraft")
var BonusOverdraftError = errors.New("bonus overdraft")

// Defines values for TransactionType.
const (
	CANCEL        string = "CANCEL"
	DEPOSIT       string = "DEPOSIT"
	PROMOCANCEL   string = "PROMOCANCEL"
	PROMODEPOSIT  string = "PROMODEPOSIT"
	PROMOWITHDRAW string = "PROMOWITHDRAW"
	WITHDRAW      string = "WITHDRAW"
)

func isValidForBlock(transType string) bool {
	switch transType {
	case WITHDRAW, PROMOWITHDRAW:
		return true
	case DEPOSIT, PROMODEPOSIT, CANCEL, PROMOCANCEL:
		return false
	default: // assume new transaction types are valid for blocking
		return true
	}
}

// validateAccount Validate the account is allowed to make transaction
func (ts *TransactionService) validateAccount(trans datastore.Transaction) error {
	_, err := ts.ds.GetPlayer(ts.ctx, trans.PlayerIdentifier)
	if err != nil {
		if errors.Is(err, datastore.EntryNotFoundError) {
			return PlayerNotFoundError
		} else {
			return err
		}
	}

	account, err := ts.ds.GetAccount(ts.ctx, trans.PlayerIdentifier, trans.Currency)
	if err != nil {
		if errors.Is(err, datastore.EntryNotFoundError) {
			// Our interpretation here is that the currency is not matching
			return CurrencyError
		} else {
			return err
		}
	}

	if account.IsBlocked && isValidForBlock(trans.TransactionType) {
		return AccountBlockedError
	}

	// Check that amount is not negative
	if trans.CashAmount < 0 || trans.BonusAmount < 0 || trans.PromoAmount < 0 {
		return NegAmountError
	}

	// Check that we are not over drafting
	if trans.TransactionType == WITHDRAW || trans.TransactionType == PROMOWITHDRAW {
		if account.CashAmount < trans.CashAmount {
			return InsufficientFundsError
		}
		if account.BonusAmount < trans.BonusAmount {
			return BonusOverdraftError
		}
		if account.PromoAmount < trans.PromoAmount {
			return PromoOverdraftError
		}
	}

	return nil
}

/*
validateTransaction Validates logic regarding transactions for CANCEL, PROMOCANCEL and DEPOSIT

Updates t amounts in case of cancellation to match cancelled transaction's amounts
*/
func (ts *TransactionService) validateTransaction(t *datastore.Transaction) error {
	// Handle cancellations
	if t.TransactionType == CANCEL || t.TransactionType == PROMOCANCEL {
		err := ts.cancellation(t)
		if err != nil {
			return err
		}
	}

	// Handle pure deposits (i.e. not promo deposits)
	if t.TransactionType == DEPOSIT {
		err := ts.deposit(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TransactionService) cancellation(cancelTrans *datastore.Transaction) error {
	precedingTrans, err := ts.getPreviousTransaction(cancelTrans, false)
	if err != nil {
		return err
	}
	if precedingTrans == nil {
		// Nothing to cancel - error
		return CancelNonExistent
	} else {
		switch precedingTrans.TransactionType {
		case WITHDRAW, PROMOWITHDRAW:
			// This is what we want
			break
		case CANCEL, PROMOCANCEL:
			// Another cancellation than our's has already done the job - error
			return CancelAlreadyCancelled

		case DEPOSIT, PROMODEPOSIT:
			// Withdrawal has already been settled - error
			return CancelAlreadySettled
		default:
			return TransactionTypeError
		}
	}

	// Check the game round so we are not trying to cancel a settled round
	gameRound, err := ts.ds.GetGameRound(ts.ctx, cancelTrans.PlayerIdentifier, *cancelTrans.ProviderRoundID)
	if err != nil {
		return err
	}

	if gameRound.EndTime != nil {
		return CancelAlreadySettled
	}

	// Finally, update the transaction amounts to corresponding amounts of the cancelled WITHDRAW/PROMOWITHDRAW
	cancelTrans.CashAmount = precedingTrans.CashAmount
	cancelTrans.BonusAmount = precedingTrans.BonusAmount
	cancelTrans.PromoAmount = precedingTrans.PromoAmount

	return nil
}

// deposit Method for validation of pure deposit validation (i.e. not promo deposits)
func (ts *TransactionService) deposit(depositTrans *datastore.Transaction) error {

	// First check the game round status. Game round needs to exist and has to be open
	gameRound, err := ts.ds.GetGameRound(ts.ctx, depositTrans.PlayerIdentifier, *depositTrans.ProviderRoundID)
	if err != nil {
		if errors.Is(err, datastore.EntryNotFoundError) {
			return DepositNotMatched
		} else {
			return err
		}
	}

	if gameRound.EndTime != nil {
		// Trouble
		return DepositAlreadySettled
	}

	// Extra check for fishy, preceding transactions
	precedingTrans, err := ts.getPreviousTransaction(depositTrans, false)
	if err != nil {
		return err
	}
	if precedingTrans != nil {
		// If we can find matching previous transactions - check the logic
		switch precedingTrans.TransactionType {
		case WITHDRAW:
			// This one is OK, break
			break
		case CANCEL:
			return DepositAlreadyCancelled
		default:
			return DepositNotMatched
		}
	}

	return nil
}

// transSortedByTime type has the necessary methods to implement the sort.Interface required
// for using the sort.Sort method.
type transSortedByTime []datastore.Transaction

func (t transSortedByTime) Len() int {
	return len(t)
}
func (t transSortedByTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t transSortedByTime) Less(i, j int) bool {
	return t[i].TransactionDateTime.Before(t[j].TransactionDateTime)
}

/*
getPreviousTransaction
In case exactly one matching transaction is found - it's returned.

In case many matching transactions exist
- the latest one is returned

Use matchType true to also match on the transactionType
*/
func (ts *TransactionService) getPreviousTransaction(trans *datastore.Transaction, matchType bool) (*datastore.Transaction, error) {
	var err error
	var trx transSortedByTime
	var twinTrans *datastore.Transaction
	trx, err = ts.ds.GetTransactionsByID(ts.ctx, trans.ProviderTransactionID, trans.ProviderName)
	if err == nil && trx != nil {
		for _, tx := range trx {
			if tx.PlayerIdentifier != trans.PlayerIdentifier ||
				tx.ProviderGameID != trans.ProviderGameID ||
				utils.OrZeroValue(tx.ProviderRoundID) != utils.OrZeroValue(trans.ProviderRoundID) {
				return nil, DuplicateTransactionError
			}
		}
	}
	if trans.ProviderBetRef != nil {
		trx, err = ts.ds.GetTransactionsByRef(ts.ctx, *trans.ProviderBetRef, trans.ProviderName)
	}
	sort.Sort(trx) // code below relies on transactions being sorted
	if err == nil && trx != nil {
		if matchType {
			for _, tx := range trx {
				if tx.TransactionType == trans.TransactionType {
					twinTrans = &tx
					break
				}
			}
		} else if len(trx) > 0 {
			twinTrans = &trx[len(trx)-1]
		}
	}

	return twinTrans, nil
}

// isValidTransactionType Checks if provided string is a valid transactionType
func isValidTransactionType(transType string) bool {
	switch transType {
	case PROMOCANCEL, PROMODEPOSIT, PROMOWITHDRAW:
		return true
	case CANCEL, DEPOSIT, WITHDRAW:
		return true
	default:
		return false
	}
}
