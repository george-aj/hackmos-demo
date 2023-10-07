package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"

	distTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/facundomedica/oracle"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Schema  collections.Schema
	Params  collections.Item[oracle.Params]
	Counter collections.Map[string, uint64]
	Prices  collections.Map[string, []byte]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, oracle.ParamsKey, "params", codec.CollValue[oracle.Params](cdc)),
		Counter:      collections.NewMap(sb, oracle.CounterKey, "counter", collections.StringKey, collections.Uint64Value),
		Prices:       collections.NewMap(sb, oracle.PricesKey, "prices", collections.StringKey, collections.BytesValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) SetOraclePrices(ctx context.Context, prices map[string]math.LegacyDec) error {
	for b, q := range prices {
		bz, err := q.Marshal()
		if err != nil {
			return err
		}

		err = k.Prices.Set(ctx, b, bz)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) GetOraclePrices(ctx context.Context) (map[string]math.LegacyDec, error) {
	prices := make(map[string]math.LegacyDec)
	err := k.Prices.Walk(ctx, nil, func(key string, value []byte) (bool, error) {
		var q math.LegacyDec
		if err := q.Unmarshal(value); err != nil {
			return true, err
		}
		prices[key] = q
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return prices, nil
}

var QueryDelegations = stakingTypes.QueryDelegatorDelegationsRequest{
	DelegatorAddr: "juno1cvncg08uc6nxuem2w2wkkd74kd4n6gzfqexjla",
	Pagination:    nil,
}

var QueryDelegatorRewards = distTypes.QueryDelegationTotalRewardsRequest{
	DelegatorAddress: "juno1cvncg08uc6nxuem2w2wkkd74kd4n6gzfqexjla",
}

type (
	CurrencyPair struct {
		Base  string
		Quote string
	}

	TickerPrice struct {
		Price  math.LegacyDec // last trade price
		Volume math.LegacyDec // 24h volume
	}

	// AggregatedProviderPrices defines a type alias for a map of
	// provider -> asset -> TickerPrice (e.g. Binance -> ATOM/USD -> 11.98)
	AggregatedProviderPrices map[string]map[string]TickerPrice
)

func (cp CurrencyPair) String() string {
	return cp.Base + cp.Quote
}