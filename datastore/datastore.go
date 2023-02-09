// Package datastore contains interface and struct to store wallet/player account data
package datastore

import (
	"context"
	"errors"
	"time"
)

type BalanceType string

const (
	Cash  BalanceType = "CASH"
	Bonus BalanceType = "BONUS"
	Promo BalanceType = "PROMO"
)

var EntryNotFoundError = errors.New("entry not found")

// ExtendedDatastore can be used in testing purposes to prepare datastore
type ExtendedDatastore interface {
	DataStore
	AddAccount(a Account)
	ClearSessionData()
	AddSession(s Session)
	AddPlayer(p Player)
	UpdateAccount(playerID string, a Account) error
	GetProviderAPIKey(provider string) (ProviderAPIKey, error)
	GetPamAPIToken() string
	GetSessionTimeout() int
	GetProviderTokens() map[string]string
}

// DataStore Data layer interface to fetch data from any data source
type DataStore interface {
	// GetSession gets session with sessionKey
	GetSession(ctx context.Context, sessionKey string) (*Session, error)
	// TouchSession keeps alive session with sessionKey
	TouchSession(ctx context.Context, sessionKey string) error
	// UpdateSession Update session with new token. Returns new session.
	UpdateSession(ctx context.Context, currentSessionToken, newSessionToken string) (*Session, error)
	// AddTransaction Adds transaction
	AddTransaction(ctx context.Context, t *Transaction) error
	// AddGameRound adds a gameround
	AddGameRound(ctx context.Context, gr GameRound) error
	// GetGame get the game with gameId for providerName
	GetGame(ctx context.Context, gameID, providerName string) (*Game, error)
	// GetAccount Gets the player account with specified currency
	GetAccount(ctx context.Context, player, currency string) (*Account, error)
	// GetPlayer Get player with id player
	GetPlayer(ctx context.Context, player string) (*Player, error)
	// GetAccountByToken returns account associated with specified session
	GetAccountByToken(ctx context.Context, sessionToken string) (*Account, error)
	// UpdateAccountBalance updates account balance of specified currency and balance type
	UpdateAccountBalance(ctx context.Context, player, currency string, balanceType BalanceType, amount float64) (float64, error)
	// EndGameRound Ends specified gameround
	EndGameRound(ctx context.Context, gr GameRound) error
	// GetTransactionsById Get transactions with providerTransactionId for specified provider
	GetTransactionsByID(ctx context.Context, providerTransactionID, providerName string) ([]Transaction, error)
	// GetTransactionsByRef Get transactions with providerBetRef for specified provider
	GetTransactionsByRef(ctx context.Context, providerBetRef, providerName string) ([]Transaction, error)
	// GetGameRound Gets gameround with roundId for player with playerId
	GetGameRound(ctx context.Context, playerID, roundID string) (*GameRound, error)
	// GetProvider gets provider with specified name
	GetProvider(ctx context.Context, provider string) (*Provider, error)
	// GetTransactionsByRoundId returns transactions for specified providerRoundId
	GetTransactionsByRoundID(ctx context.Context, providerRoundID string) ([]Transaction, error)
}

type Player struct {
	PlayerIdentifier string `yaml:"playerIdentifier"`
	ID               int
}

type Account struct {
	PlayerIdentifier string  `yaml:"playerIdentifier"`
	Currency         string  `yaml:"currency"`
	Country          string  `yaml:"country"`
	Language         string  `yaml:"language"`
	CashAmount       float64 `yaml:"cashAmount"`
	BonusAmount      float64 `yaml:"bonusAmount"`
	PromoAmount      float64 `yaml:"promoAmount"`
	IsBlocked        bool    `yaml:"isBlocked"`
	ID               int
}

type Transaction struct {
	TransactionDateTime   time.Time `yaml:"transactionDateTime" json:"transactionDateTime"`
	PlayerIdentifier      string    `yaml:"playerIdentifier" json:"playerIdentifier,omitempty"`
	Currency              string    `json:"currency,omitempty"`
	TransactionType       string    `yaml:"transactionType" json:"transactionType,omitempty"`
	ProviderTransactionID string    `yaml:"providerTransactionId" json:"providerTransactionId,omitempty"`
	ProviderBetRef        *string   `yaml:"providerBetRef" json:"providerBetRef,omitempty"`
	ProviderGameID        string    `yaml:"providerGameId" json:"providerGameId,omitempty"`
	ProviderRoundID       *string   `yaml:"providerRoundId" json:"providerRoundId,omitempty"`
	ProviderName          string    `yaml:"providerName" json:"providerName"`
	SessionKey            string    `yaml:"sessionKey" json:"sessionKey"`
	CashAmount            float64   `yaml:"cashAmount" json:"cashAmount,omitempty"`
	BonusAmount           float64   `yaml:"bonusAmount" json:"bonusAmount,omitempty"`
	PromoAmount           float64   `yaml:"promoAmount" json:"promoAmount,omitempty"`
	ID                    int       `json:"id,omitempty"`
	IsGameOver            bool      `yaml:"isGameOver" json:"isGameOver,omitempty"`
}

type Session struct {
	Timestamp        time.Time
	GameID           *string
	Key              string
	PlayerIdentifier string `yaml:"playerIdentifier"`
	Provider         string
	Currency         string
	Country          string
	Language         string
	// Timeout in seconds
	Timeout int
}

func (s *Session) IsExpired() bool {
	return s.Timestamp.Add(time.Duration(s.Timeout) * time.Second).Before(time.Now())
}

type Game struct {
	ProviderName   string `yaml:"providerName"`
	ProviderGameID string `yaml:"providerGameId"`
	ID             int
}

type GameRound struct {
	EndTime         *time.Time `yaml:"endTime"`
	StartTime       time.Time  `yaml:"startTime"`
	ProviderName    string     `yaml:"providerName"`
	ProviderGameID  string     `yaml:"providerGameId"`
	ProviderRoundID string     `yaml:"providerRoundId"`
	PlayerID        string     `yaml:"playerId"`
}

type Provider struct {
	Provider   string
	ProviderID int `yaml:"providerId"`
}

type ProviderAPIKey struct {
	APIKey   string `yaml:"apiKey"`
	Provider string
}
