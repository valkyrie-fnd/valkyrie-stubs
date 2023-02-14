// Package caleta contains the implementation of the caleta provider stub.
package caleta

import (
	"time"
)

type TransactionResponse struct {
	RoundTransactions *[]RoundTransaction `json:"transactions,omitempty"`
	Message           string              `json:"message"`
	RoundID           string              `json:"round_id"`
	Code              int                 `json:"code"`
}

type RoundTransaction struct {
	CreatedTime  time.Time `json:"created_time"`
	ClosedTime   time.Time `json:"closed_time"`
	TxnUUID      string    `json:"txn_uuid"`
	Payload      Payload   `json:"payload"`
	ID           int       `json:"id"`
	RoundID      int       `json:"round_id"`
	TxnType      int       `json:"txn_type"`
	Status       int       `json:"status"`
	CacheEntryID int       `json:"cache_entry_id"`
	Amount       int       `json:"amount"`
}

type Payload struct {
	Bet                      string `json:"bet"`
	Round                    string `json:"round"`
	Token                    string `json:"token"`
	Currency                 string `json:"currency"`
	GameCode                 string `json:"game_code"`
	RequestUUID              string `json:"request_uuid"`
	SupplierUser             string `json:"supplier_user"`
	TransactionUUID          string `json:"transaction_uuid"`
	ReferenceTransactionUUID string `json:"reference_transaction_uuid"`
	GameID                   string `json:"game_id"`
	JackpotContribution      int    `json:"jackpot_contribution"`
	Amount                   int    `json:"amount"`
	RoundClosed              bool   `json:"round_closed"`
	IsFree                   bool   `json:"is_free"`
}
