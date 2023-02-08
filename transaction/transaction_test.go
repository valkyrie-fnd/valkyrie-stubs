package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/memorydatastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

var ctx = context.Background()

func Test_fail_when_invalid_transactionType(t *testing.T) {
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{})
	tr := datastore.Transaction{TransactionType: "NoExistType"}
	ts := NewTransactionService(ctx, ds)
	res, err := ts.AddTransaction(tr)
	assert.Equal(t, 0, res)
	assert.ErrorAs(t, err, &TransactionTypeError)
}

func Test_fail_if_player_does_not_exist(t *testing.T) {
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{})
	tr := datastore.Transaction{TransactionType: "WITHDRAW"}
	ts := NewTransactionService(ctx, ds)
	res, err := ts.AddTransaction(tr)

	assert.Equal(t, 0, res)

	assert.ErrorAs(t, err, &PlayerNotFoundError)
}

func Test_fail_if_account_does_not_exist(t *testing.T) {
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
	})
	tr := datastore.Transaction{
		TransactionType:  "WITHDRAW",
		PlayerIdentifier: "101",
		Currency:         "USD",
	}
	ts := NewTransactionService(ctx, ds)
	res, err := ts.AddTransaction(tr)

	assert.Equal(t, 0, res)
	assert.ErrorAs(t, err, &CurrencyError)
}

func Test_fail_if_account_is_blocked(t *testing.T) {
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: true},
		},
	})
	tr := datastore.Transaction{
		TransactionType:  "WITHDRAW",
		PlayerIdentifier: "101",
		Currency:         "USD",
	}
	ts := NewTransactionService(ctx, ds)
	res, err := ts.AddTransaction(tr)

	assert.Equal(t, 0, res)
	assert.ErrorAs(t, err, &AccountBlockedError)
}

func Test_fail_if_transaction_has_negative_amounts(t *testing.T) {
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false},
		},
	})
	tr := datastore.Transaction{
		TransactionType:  "WITHDRAW",
		PlayerIdentifier: "101",
		Currency:         "USD",
		CashAmount:       -10,
	}
	ts := NewTransactionService(ctx, ds)
	res, err := ts.AddTransaction(tr)

	assert.Equal(t, 0, res)
	assert.ErrorAs(t, err, &NegAmountError)
}

func Test_fail_if_account_does_not_have_enough_funds(t *testing.T) {
	acc := datastore.Account{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false}
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			acc,
		},
	})
	tests := []struct {
		name        string
		transaction datastore.Transaction
		cashAmount  float64
		bonusAmount float64
		promoAmount float64
		err         error
	}{
		{"Not enough cash funds", datastore.Transaction{TransactionType: "WITHDRAW", PlayerIdentifier: "101", Currency: "USD", CashAmount: 10}, 5, 0, 0, InsufficientFundsError},
		{"Not enough bonus funds", datastore.Transaction{TransactionType: "WITHDRAW", PlayerIdentifier: "101", Currency: "USD", BonusAmount: 10}, 0, 5, 0, BonusOverdraftError},
		{"not enough promo funds", datastore.Transaction{TransactionType: "WITHDRAW", PlayerIdentifier: "101", Currency: "USD", PromoAmount: 10}, 0, 0, 5, PromoOverdraftError},
	}

	ts := NewTransactionService(ctx, ds)
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			acc.CashAmount = test.cashAmount
			acc.BonusAmount = test.bonusAmount
			acc.PromoAmount = test.promoAmount
			res, err := ts.AddTransaction(test.transaction)
			assert.Equal(tt, 0, res)
			assert.ErrorAs(tt, err, &test.err)
		})
	}
}

func Test_fail_if_transactionID_is_the_same_as_a_previous_transaction(t *testing.T) {
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false, CashAmount: 100},
		},
		Transactions: []datastore.Transaction{
			{ID: 1, ProviderTransactionID: "abc123", PlayerIdentifier: "some1else"},
		},
	})
	tr := datastore.Transaction{
		ProviderTransactionID: "abc123",
		TransactionType:       "WITHDRAW",
		PlayerIdentifier:      "101",
		Currency:              "USD",
		CashAmount:            10,
	}
	ts := NewTransactionService(ctx, ds)
	res, err := ts.AddTransaction(tr)

	assert.Equal(t, 0, res)
	assert.ErrorAs(t, err, &DuplicateTransactionError)
}

func Test_fail_if_cancelTransaction(t *testing.T) {
	acc := datastore.Account{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false}
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			acc,
		},
	})
	tests := []struct {
		name          string
		transactionID string
		setup         func(t datastore.Transaction)
		err           error
	}{
		{"No transaction To Cancel", "abc1", func(tr datastore.Transaction) {}, CancelNonExistent},
		{"Trying to cancel a promo cancel transaction", "abc2",
			func(tr datastore.Transaction) {
				_ = ds.AddTransaction(ctx, &datastore.Transaction{ProviderTransactionID: "abc123", PlayerIdentifier: "101", TransactionType: "PROMOCANCEL", ProviderGameID: "slot", ProviderRoundID: tr.ProviderRoundID, TransactionDateTime: time.Now()})
			}, CancelAlreadyCancelled},
		{"Trying to cancel settled transaction", "abc3", func(t datastore.Transaction) {
			_ = ds.AddTransaction(ctx, &datastore.Transaction{ProviderTransactionID: "abc123Settle", PlayerIdentifier: "101", TransactionType: "DEPOSIT", ProviderGameID: "slot", ProviderBetRef: t.ProviderBetRef, ProviderRoundID: t.ProviderRoundID, TransactionDateTime: time.Now()})
		}, CancelAlreadySettled},
		{"Trying to cancel transaction with ended gameround", "abc4", func(t datastore.Transaction) {
			n := time.Now()
			_ = ds.AddGameRound(ctx, datastore.GameRound{ProviderGameID: t.ProviderGameID, PlayerID: t.PlayerIdentifier, ProviderRoundID: *t.ProviderRoundID, StartTime: t.TransactionDateTime, EndTime: &n})
			_ = ds.AddTransaction(ctx, &datastore.Transaction{ProviderTransactionID: "abc123Settle", PlayerIdentifier: "101", TransactionType: "WITHDRAW", ProviderGameID: "slot", ProviderBetRef: t.ProviderBetRef, ProviderRoundID: t.ProviderRoundID, TransactionDateTime: time.Now()})
		}, CancelAlreadySettled},
	}
	roundID := "Round1"
	tr := datastore.Transaction{
		ProviderGameID:      "slot",
		ProviderRoundID:     utils.Ptr(roundID),
		TransactionType:     "CANCEL",
		PlayerIdentifier:    "101",
		Currency:            "USD",
		CashAmount:          10,
		TransactionDateTime: time.Now(),
	}
	ts := NewTransactionService(ctx, ds)
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			transaction := tr
			transaction.ProviderTransactionID = test.transactionID
			transaction.ProviderBetRef = utils.Ptr(test.transactionID + "betRef")
			test.setup(transaction)
			res, err := ts.AddTransaction(transaction)
			assert.Equal(tt, 0, res)
			assert.ErrorAs(tt, err, &test.err)
		})
	}
}

func Test_fail_if_depositTransaction(t *testing.T) {
	acc := datastore.Account{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false}
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			acc,
		},
	})
	tests := []struct {
		name          string
		transactionID string
		setup         func(t datastore.Transaction)
		err           error
	}{
		{"Fail when there is no game round", "abc1", func(t datastore.Transaction) {}, DepositNotMatched},
		{"Fail when game round has ended", "abc1", func(t datastore.Transaction) {
			n := time.Now()
			_ = ds.AddGameRound(ctx, datastore.GameRound{ProviderGameID: t.ProviderGameID, PlayerID: t.PlayerIdentifier, ProviderRoundID: *t.ProviderRoundID, StartTime: t.TransactionDateTime, EndTime: &n})
		}, DepositAlreadySettled},
		{"Fail when previous transaction is cancelled", "abc1", func(t datastore.Transaction) {
			_ = ds.AddGameRound(ctx, datastore.GameRound{ProviderGameID: t.ProviderGameID, PlayerID: t.PlayerIdentifier, ProviderRoundID: *t.ProviderRoundID, StartTime: t.TransactionDateTime})
			_ = ds.AddTransaction(ctx, &datastore.Transaction{ProviderTransactionID: "abc123Settle", PlayerIdentifier: "101", TransactionType: "CANCEL", ProviderGameID: "slot", ProviderBetRef: t.ProviderBetRef, ProviderRoundID: t.ProviderRoundID, TransactionDateTime: time.Now()})
		}, DepositAlreadyCancelled},
		{"Fail when previous transaction is not compatible", "abc1", func(t datastore.Transaction) {
			_ = ds.AddGameRound(ctx, datastore.GameRound{ProviderGameID: t.ProviderGameID, PlayerID: t.PlayerIdentifier, ProviderRoundID: *t.ProviderRoundID, StartTime: t.TransactionDateTime})
			_ = ds.AddTransaction(ctx, &datastore.Transaction{ProviderTransactionID: "abc123Settle", PlayerIdentifier: "101", TransactionType: "PROMOWITHDRAW", ProviderGameID: "slot", ProviderBetRef: t.ProviderBetRef, ProviderRoundID: t.ProviderRoundID, TransactionDateTime: time.Now()})
		}, DepositNotMatched},
	}
	roundID := "Round1"
	tr := datastore.Transaction{
		ProviderGameID:      "slot",
		ProviderRoundID:     utils.Ptr(roundID),
		TransactionType:     "DEPOSIT",
		PlayerIdentifier:    "101",
		Currency:            "USD",
		CashAmount:          10,
		TransactionDateTime: time.Now(),
	}
	ts := NewTransactionService(ctx, ds)
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			transaction := tr
			transaction.ProviderTransactionID = test.transactionID
			transaction.ProviderBetRef = utils.Ptr(test.transactionID + "betRef")
			test.setup(transaction)
			res, err := ts.AddTransaction(transaction)
			assert.Equal(tt, 0, res)
			assert.ErrorAs(tt, err, &test.err)
		})
	}
}

func Test_fail_if_game_does_not_exist(t *testing.T) {
	roundID := "Round1"
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false, CashAmount: 100},
		},
		GameRounds: []datastore.GameRound{
			{ProviderGameID: "slot", PlayerID: "101", ProviderRoundID: roundID, StartTime: time.Now()},
		},
	})
	tr := datastore.Transaction{
		ProviderTransactionID: "abc123",
		ProviderBetRef:        utils.Ptr("ABCBetRef"),
		ProviderGameID:        "slot",
		ProviderRoundID:       utils.Ptr(roundID),
		TransactionType:       "DEPOSIT",
		PlayerIdentifier:      "101",
		Currency:              "USD",
		CashAmount:            10,
		TransactionDateTime:   time.Now(),
	}

	ts := NewTransactionService(ctx, ds)
	res, err := ts.AddTransaction(tr)

	assert.Equal(t, 0, res)
	assert.ErrorAs(t, err, &GameError)
}

func Test_success_Account_balance_Updated(t *testing.T) {
	roundID := "Round1"
	startAmount := 100.0
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false, CashAmount: startAmount, BonusAmount: startAmount, PromoAmount: startAmount},
		},
		GameRounds: []datastore.GameRound{
			{ProviderGameID: "slot", PlayerID: "101", ProviderRoundID: roundID, StartTime: time.Now()},
		},
		Games: []datastore.Game{
			{ID: 1, ProviderGameID: "slot", ProviderName: "Valkyrie"},
		},
	})

	tr := datastore.Transaction{
		ProviderTransactionID: "abc123",
		ProviderBetRef:        utils.Ptr("ABCBetRef"),
		ProviderGameID:        "slot",
		ProviderRoundID:       utils.Ptr(roundID),
		TransactionType:       "WITHDRAW",
		PlayerIdentifier:      "101",
		Currency:              "USD",
		CashAmount:            10,
		BonusAmount:           15,
		PromoAmount:           5,
		TransactionDateTime:   time.Now(),
	}

	ts := NewTransactionService(ctx, ds)
	_, err := ts.AddTransaction(tr)
	assert.NoError(t, err)
	acc, _ := ds.GetAccount(ctx, "101", "USD")
	assert.Equal(t, startAmount-10, acc.CashAmount)
	assert.Equal(t, startAmount-15, acc.BonusAmount)
	assert.Equal(t, startAmount-5, acc.PromoAmount)
	tr = datastore.Transaction{
		ProviderTransactionID: "abc123",
		ProviderBetRef:        utils.Ptr("ABCBetRef"),
		ProviderGameID:        "slot",
		ProviderRoundID:       utils.Ptr(roundID),
		TransactionType:       "DEPOSIT",
		PlayerIdentifier:      "101",
		Currency:              "USD",
		CashAmount:            20,
		BonusAmount:           15,
		PromoAmount:           1,
		TransactionDateTime:   time.Now(),
		IsGameOver:            true,
	}
	_, err = ts.AddTransaction(tr)
	assert.NoError(t, err)
	acc, _ = ds.GetAccount(ctx, "101", "USD")
	assert.Equal(t, startAmount+10, acc.CashAmount)
	assert.Equal(t, startAmount, acc.BonusAmount)
	assert.Equal(t, startAmount-4, acc.PromoAmount)
}
func Test_success_GameRound_created(t *testing.T) {
	roundID := "Round1"
	startAmount := 100.0
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false, CashAmount: startAmount, BonusAmount: startAmount, PromoAmount: startAmount},
		},
		Games: []datastore.Game{
			{ID: 1, ProviderGameID: "slot", ProviderName: "Valkyrie"},
		},
	})
	tr := datastore.Transaction{
		ProviderTransactionID: "abc123",
		ProviderBetRef:        utils.Ptr("ABCBetRef"),
		ProviderGameID:        "slot",
		ProviderRoundID:       utils.Ptr(roundID),
		TransactionType:       "WITHDRAW",
		PlayerIdentifier:      "101",
		Currency:              "USD",
		CashAmount:            10,
		BonusAmount:           15,
		PromoAmount:           5,
		TransactionDateTime:   time.Now(),
	}

	ts := NewTransactionService(ctx, ds)
	_, err := ts.AddTransaction(tr)
	assert.NoError(t, err)
	gr, _ := ds.GetGameRound(ctx, "101", roundID)
	assert.Equal(t, gr.ProviderRoundID, roundID)

}
func Test_success_Gameround_ended(t *testing.T) {
	roundID := "Round1"
	startAmount := 100.0
	ds := memorydatastore.NewMapDataStore(&memorydatastore.Config{
		Players: []datastore.Player{
			{ID: 101, PlayerIdentifier: "101"},
		},
		Accounts: []datastore.Account{
			{ID: 1, PlayerIdentifier: "101", Currency: "USD", IsBlocked: false, CashAmount: startAmount, BonusAmount: startAmount, PromoAmount: startAmount},
		},
		Games: []datastore.Game{
			{ID: 1, ProviderGameID: "slot", ProviderName: "Valkyrie"},
		},
	})
	tr := datastore.Transaction{
		ProviderTransactionID: "abc123",
		ProviderBetRef:        utils.Ptr("ABCBetRef"),
		ProviderGameID:        "slot",
		ProviderRoundID:       utils.Ptr(roundID),
		TransactionType:       "WITHDRAW",
		PlayerIdentifier:      "101",
		Currency:              "USD",
		CashAmount:            10,
		BonusAmount:           15,
		PromoAmount:           5,
		TransactionDateTime:   time.Now(),
	}
	ts := NewTransactionService(ctx, ds)
	_, err := ts.AddTransaction(tr)
	assert.NoError(t, err)
	tr = datastore.Transaction{
		ProviderTransactionID: "abc123",
		ProviderBetRef:        utils.Ptr("ABCBetRef"),
		ProviderGameID:        "slot",
		ProviderRoundID:       utils.Ptr(roundID),
		TransactionType:       "DEPOSIT",
		PlayerIdentifier:      "101",
		Currency:              "USD",
		CashAmount:            10,
		BonusAmount:           15,
		PromoAmount:           5,
		TransactionDateTime:   time.Now(),
		IsGameOver:            true,
	}
	_, err = ts.AddTransaction(tr)
	assert.NoError(t, err)

	gr, _ := ds.GetGameRound(ctx, "101", roundID)
	assert.Equal(t, gr.ProviderRoundID, roundID)
	assert.NotNil(t, gr.EndTime)
}
