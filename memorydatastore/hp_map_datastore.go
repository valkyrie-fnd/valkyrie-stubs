package memorydatastore

import (
	"context"

	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
)

// High-performance version of the map datastore.
// Preferably used when doing targeted benchmarking as it could introduce functionality limitations.
// Embeds the default MapDataStore which provides the option to do custom implementations if required.
type HPMapDataStore struct {
	MapDataStore
}

func NewHPMapDataStore(config *Config) *HPMapDataStore {
	return &HPMapDataStore{*NewMapDataStore(config)}
}

func (ds *HPMapDataStore) AddTransaction(_ context.Context, t *datastore.Transaction) error {
	// Hard coded id will keep the tx map very slim
	t.Id = 123

	ds.Transactions.Mux.Lock()
	defer ds.Transactions.Mux.Unlock()

	ds.Transactions.Map[t.Id] = t
	return nil
}
