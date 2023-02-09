package memorydatastore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valkyrie-fnd/valkyrie-stubs/datastore"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

func TestNewMapDataStore(t *testing.T) {
	config := &Config{
		Transactions: []datastore.Transaction{
			{ID: 1, ProviderRoundID: utils.Ptr("1")},
			{ID: 2, ProviderRoundID: utils.Ptr("1")},
		},
	}

	got := NewMapDataStore(config)

	ts, err := got.GetTransactionsByRoundID(context.TODO(), "1")
	assert.NoError(t, err)

	assert.Contains(t, ts, config.Transactions[0])
	assert.Contains(t, ts, config.Transactions[1])
}
