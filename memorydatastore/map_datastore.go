package memorydatastore

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

type gameRounds struct {
	Mux sync.RWMutex
	Map map[string]*datastore.GameRound
}

type sessions struct {
	Mux sync.RWMutex
	Map map[string]*datastore.Session
}

type games struct {
	Mux sync.RWMutex
	Map map[string]*datastore.Game
}

type providers struct {
	Mux sync.RWMutex
	Map map[string]*datastore.Provider
}

type providerApiKeys struct {
	Mux sync.RWMutex
	Map map[string]*datastore.ProviderApiKey
}

type accountKey struct {
	Player   string
	Currency string
}

type accounts struct {
	Mux sync.RWMutex
	Map map[accountKey]*datastore.Account
}

type players struct {
	Mux sync.RWMutex
	Map map[string]*datastore.Player
}

type transactions struct {
	Mux sync.RWMutex
	Map map[int]*datastore.Transaction
}

type mapStorage struct {
	Players          *players
	Sessions         *sessions
	Accounts         *accounts
	Transactions     *transactions
	Games            *games
	GameRounds       *gameRounds
	Providers        *providers
	ProviderApiKeys  *providerApiKeys
	ProviderSessions *sessions
	PamApiToken      string
	SessionTimeout   *int
}
type MapDataStore struct {
	mapStorage
}

func NewMapDataStore(config *Config) *MapDataStore {
	dataStore := &MapDataStore{mapStorage: mapStorage{
		Sessions: &sessions{
			Map: map[string]*datastore.Session{},
		},
		Players: &players{
			Map: map[string]*datastore.Player{},
		},
		Accounts: &accounts{
			Map: map[accountKey]*datastore.Account{},
		},
		Transactions: &transactions{
			Map: map[int]*datastore.Transaction{},
		},
		Games: &games{
			Map: map[string]*datastore.Game{},
		},
		GameRounds: &gameRounds{
			Map: map[string]*datastore.GameRound{},
		},
		Providers: &providers{
			Map: map[string]*datastore.Provider{},
		},
		ProviderApiKeys: &providerApiKeys{
			Map: map[string]*datastore.ProviderApiKey{},
		},
		ProviderSessions: &sessions{
			Map: map[string]*datastore.Session{},
		},
		SessionTimeout: config.SessionTimeout,
	},
	}

	dataStore.configure(config)

	return dataStore
}
func (ds *MapDataStore) AddPamApiToken(t string) {
	ds.PamApiToken = t
}

func (ds *MapDataStore) GetSessionTimeout() int {
	if ds.SessionTimeout != nil {
		return *ds.SessionTimeout
	}
	return 5 * 60
}
func (ds *MapDataStore) addProvider(p datastore.Provider) {
	ds.Providers.Mux.Lock()
	defer ds.Providers.Mux.Unlock()

	ds.Providers.Map[p.Provider] = &p
}

func (ds *MapDataStore) addProviderApiKey(p datastore.ProviderApiKey) {
	ds.ProviderApiKeys.Mux.Lock()
	defer ds.ProviderApiKeys.Mux.Unlock()

	ds.ProviderApiKeys.Map[p.Provider] = &p
}

func (ds *MapDataStore) addGame(g datastore.Game) {
	ds.Games.Mux.Lock()
	defer ds.Games.Mux.Unlock()

	ds.Games.Map[g.ProviderGameId] = &g
}

func (ds *MapDataStore) configure(config *Config) {
	ds.AddPamApiToken(config.PamApiToken)
	for _, p := range config.Providers {
		ds.addProvider(p)
	}
	for _, p := range config.ProviderApiKeys {
		ds.addProviderApiKey(p)
	}
	for _, s := range config.ProviderSessions {
		ds.AddProviderSession(s)
	}
	for _, g := range config.Games {
		ds.addGame(g)
	}
	for _, gr := range config.GameRounds {
		_ = ds.AddGameRound(context.Background(), gr)
	}
	for _, p := range config.Players {
		ds.AddPlayer(p)
	}
	for _, a := range config.Accounts {
		ds.AddAccount(a)
	}
	for _, s := range config.Sessions {
		ds.AddSession(s)
	}
	for _, t := range config.Transactions {
		_ = ds.AddTransaction(context.Background(), &t)
	}
}

func (ds *MapDataStore) AddPlayer(p datastore.Player) {
	ds.Players.Mux.Lock()
	defer ds.Players.Mux.Unlock()
	ds.Players.Map[p.PlayerIdentifier] = &p
}

func (ds *MapDataStore) AddAccount(a datastore.Account) {
	ds.Accounts.Mux.Lock()
	defer ds.Accounts.Mux.Unlock()
	accKey := accountKey{
		Player:   a.PlayerIdentifier,
		Currency: a.Currency,
	}
	ds.Accounts.Map[accKey] = &a
}

func (ds *MapDataStore) AddSession(s datastore.Session) {
	ds.Sessions.Mux.Lock()
	defer ds.Sessions.Mux.Unlock()

	ds.Sessions.Map[s.Key] = &s
}

func (ds *MapDataStore) TouchSession(_ context.Context, sessionKey string) error {
	ds.Sessions.Mux.Lock()
	defer ds.Sessions.Mux.Unlock()
	if s, found := ds.Sessions.Map[sessionKey]; found {
		s.Timestamp = time.Now()
		return nil
	} else {
		return datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) UpdateSession(_ context.Context, currKey, newKey string) (*datastore.Session, error) {
	ds.Sessions.Mux.Lock()
	defer ds.Sessions.Mux.Unlock()

	if s, found := ds.Sessions.Map[currKey]; found {
		newSession := &datastore.Session{
			Key:              newKey,
			PlayerIdentifier: s.PlayerIdentifier,
			Provider:         s.Provider,
			Currency:         s.Currency,
			Country:          s.Country,
			Language:         s.Language,
			Timestamp:        time.Now(),
			Timeout:          s.Timeout,
			GameId:           s.GameId,
		}
		s.Timeout = 0
		ds.Sessions.Map[newKey] = newSession

		return newSession, nil
	} else {
		return nil, datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) AddTransaction(_ context.Context, t *datastore.Transaction) error {
	if t.Id == 0 {
		t.Id = utils.RandomInt()
	}

	ds.Transactions.Mux.Lock()
	defer ds.Transactions.Mux.Unlock()

	ds.Transactions.Map[t.Id] = t
	return nil
}

func (ds *MapDataStore) AddGameRound(_ context.Context, gr datastore.GameRound) error {
	ds.GameRounds.Mux.Lock()
	defer ds.GameRounds.Mux.Unlock()

	ds.GameRounds.Map[gr.ProviderRoundId] = &gr
	return nil
}

func (ds *MapDataStore) GetGame(_ context.Context, gameId, providerName string) (*datastore.Game, error) {
	ds.Games.Mux.RLock()
	defer ds.Games.Mux.RUnlock()

	if v, found := ds.Games.Map[gameId]; found {
		return v, nil
	} else {
		return nil, datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) GetPlayer(_ context.Context, playerIdentifier string) (*datastore.Player, error) {
	ds.Players.Mux.RLock()
	defer ds.Players.Mux.RUnlock()
	if v, found := ds.Players.Map[playerIdentifier]; found {
		return v, nil
	} else {
		return nil, datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) GetAccount(_ context.Context, player, currency string) (*datastore.Account, error) {
	ds.Accounts.Mux.RLock()
	defer ds.Accounts.Mux.RUnlock()
	accKey := accountKey{
		Player:   player,
		Currency: currency,
	}
	if v, found := ds.Accounts.Map[accKey]; found {
		return v, nil
	} else {
		return nil, datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) GetAccountByToken(ctx context.Context, sessionToken string) (*datastore.Account, error) {
	var player, currency string
	if v, err := ds.GetSession(ctx, sessionToken); err == nil {
		player = v.PlayerIdentifier
		currency = v.Currency
	} else {
		return nil, datastore.EntryNotFoundError
	}

	if v, err := ds.GetAccount(ctx, player, currency); err == nil {
		return v, nil
	} else {
		return nil, datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) UpdateAccountBalance(_ context.Context, player, currency string, balanceType datastore.BalanceType, amount float64) (float64, error) {
	ds.Accounts.Mux.Lock()
	defer ds.Accounts.Mux.Unlock()
	accKey := accountKey{
		Player:   player,
		Currency: currency,
	}
	if v, found := ds.Accounts.Map[accKey]; found {
		switch balanceType {
		case datastore.Cash:
			v.CashAmount += amount
			return v.CashAmount, nil
		case datastore.Bonus:
			v.BonusAmount += amount
			return v.BonusAmount, nil
		case datastore.Promo:
			v.PromoAmount += amount
			return v.PromoAmount, nil
		default:
			return 0, fmt.Errorf("unknown balance type '%s'", balanceType)
		}
	} else {
		return 0, errors.New("failed to update account balance")
	}
}

func (ds *MapDataStore) UpdateAccount(playerId string, acc datastore.Account) error {
	ds.Accounts.Mux.Lock()
	defer ds.Accounts.Mux.Unlock()

	accKey := accountKey{
		Player:   playerId,
		Currency: acc.Currency,
	}

	if v, found := ds.Accounts.Map[accKey]; found {
		v.CashAmount = acc.CashAmount
		v.BonusAmount = acc.BonusAmount
		v.Currency = acc.Currency
		v.Language = acc.Language
		v.Country = acc.Country
		v.IsBlocked = acc.IsBlocked

		// omitting v.Player and v.Id since changing these would break datastore relations
		return nil
	} else {
		return errors.New("failed to update account")
	}
}

func (ds *MapDataStore) EndGameRound(_ context.Context, gr datastore.GameRound) error {
	ds.GameRounds.Mux.Lock()
	defer ds.GameRounds.Mux.Unlock()

	if v, found := ds.GameRounds.Map[gr.ProviderRoundId]; found {
		t := time.Now()
		v.EndTime = &t
	} else {
		return errors.New("failed to update game round")
	}
	return nil
}

func (ds *MapDataStore) GetSession(_ context.Context, key string) (*datastore.Session, error) {
	ds.Sessions.Mux.RLock()
	defer ds.Sessions.Mux.RUnlock()

	if v, found := ds.Sessions.Map[key]; found {
		return v, nil
	} else {
		return nil, datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) GetProviderApiKey(provider string) (datastore.ProviderApiKey, error) {
	ds.ProviderApiKeys.Mux.RLock()
	defer ds.ProviderApiKeys.Mux.RUnlock()

	if v, found := ds.ProviderApiKeys.Map[provider]; found {
		return *v, nil
	} else {
		return datastore.ProviderApiKey{}, datastore.EntryNotFoundError
	}
}

// Note that transaction not found is not regarded as an error
func (ds *MapDataStore) GetTransactionsById(_ context.Context, providerTransactionId, providerName string) ([]datastore.Transaction, error) {
	ds.Transactions.Mux.RLock()
	defer ds.Transactions.Mux.RUnlock()
	var trx []datastore.Transaction

	for _, t := range ds.Transactions.Map {
		if t.ProviderTransactionId == providerTransactionId {
			trx = append(trx, *t)
		}
	}
	return trx, nil
}

// Note that transaction not found is not regarded as an error
func (ds *MapDataStore) GetTransactionsByRef(_ context.Context, providerBetRef, providerName string) ([]datastore.Transaction, error) {
	ds.Transactions.Mux.RLock()
	defer ds.Transactions.Mux.RUnlock()

	var trx []datastore.Transaction
	for _, t := range ds.Transactions.Map {
		if t.ProviderBetRef != nil && *t.ProviderBetRef == providerBetRef {
			trx = append(trx, *t)
		}
	}
	return trx, nil
}

func (ds *MapDataStore) GetTransactionsByRoundId(_ context.Context, providerRoundId string) ([]datastore.Transaction, error) {
	ds.Transactions.Mux.RLock()
	defer ds.Transactions.Mux.RUnlock()

	var trx []datastore.Transaction
	for _, t := range ds.Transactions.Map {
		if t.ProviderRoundId != nil && *t.ProviderRoundId == providerRoundId {
			trx = append(trx, *t)
		}
	}
	return trx, nil
}

func (ds *MapDataStore) GetGameRound(_ context.Context, playerId, roundId string) (*datastore.GameRound, error) {
	ds.GameRounds.Mux.RLock()
	defer ds.GameRounds.Mux.RUnlock()

	if v, found := ds.GameRounds.Map[roundId]; found {
		return v, nil
	} else {
		return nil, datastore.EntryNotFoundError

	}
}

func (ds *MapDataStore) GetProvider(_ context.Context, provider string) (*datastore.Provider, error) {
	ds.Providers.Mux.RLock()
	defer ds.Providers.Mux.RUnlock()

	if v, found := ds.Providers.Map[provider]; found {
		return v, nil
	} else {
		return nil, datastore.EntryNotFoundError
	}
}

func (ds *MapDataStore) GetPamApiToken() string {
	return ds.PamApiToken
}

func (ds *MapDataStore) ClearSessionData() {
	ds.Sessions.Mux.Lock()
	defer ds.Sessions.Mux.Unlock()
	ds.mapStorage.Sessions.Map = map[string]*datastore.Session{}
}

func (ds *MapDataStore) AddProviderSession(s datastore.Session) {
	ds.ProviderSessions.Mux.Lock()
	defer ds.ProviderSessions.Mux.Unlock()

	ds.ProviderSessions.Map[s.Provider] = &s
}

func (ds *MapDataStore) GetProviderTokens() map[string]string {
	ds.ProviderSessions.Mux.RLock()
	defer ds.ProviderSessions.Mux.RUnlock()
	res := make(map[string]string)

	for k, v := range ds.ProviderSessions.Map {
		res[k] = v.Key
	}
	return res
}
