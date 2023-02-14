package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valkyrie-fnd/valkyrie-stubs/provider/caleta"
	"github.com/valkyrie-fnd/valkyrie-stubs/utils"
)

func TestCaletaTransactions(t *testing.T) {
	_, addr, err := utils.GetFreePort()
	require.NoError(t, err)

	expectedTransactions := func() (caleta.TransactionResponse, error) {
		return caleta.TransactionResponse{
			RoundTransactions: &[]caleta.RoundTransaction{
				{
					ID: 1,
				},
			},
		}, nil
	}

	Create(context.TODO(), addr, WithCannedTransactions(expectedTransactions))

	a := fiber.Get(fmt.Sprintf("http://%s/caleta/api/transactions/round", addr))
	res := caleta.TransactionResponse{}
	status, body, errs := a.Struct(&res)

	//TODO: use when with Go 1.20
	// require.NoError(t, errors.Join(errs...))
	require.Len(t, errs, 0)
	require.Equal(t, 200, status, string(body))
	assert.Len(t, *res.RoundTransactions, 1)
}
