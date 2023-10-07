package mockprovider_test

import (
	"fmt"
	"github.com/facundomedica/oracle/keeper"
	"github.com/facundomedica/oracle/mockprovider"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuery(t *testing.T) {
	resp1, err1 := mockprovider.MockProvider{}.GetDelegatorDelegations(keeper.QueryDelegations)
	if err1 != nil {
		require.NoError(t, err1)
	}

	fmt.Println(resp1)

	resp2, err2 := mockprovider.MockProvider{}.GetDelegatorRewards(keeper.QueryDelegatorRewards)
	if err2 != nil {
		require.NoError(t, err2)
	}

	fmt.Println(resp2)

	resp3, err3 := mockprovider.MockProvider{}.GetWasmQuery()
	if err3 != nil {
		require.NoError(t, err3)
	}

	fmt.Println(resp3)
}
