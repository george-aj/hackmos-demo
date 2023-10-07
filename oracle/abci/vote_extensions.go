package abci

import (
	"encoding/json"
	"fmt"
	distTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/facundomedica/oracle/mockprovider"
	"time"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/facundomedica/oracle/keeper"
	"golang.org/x/sync/errgroup"
)

type VoteExtHandler struct {
	logger       log.Logger
	currentBlock int64         // current block height
	timeout      time.Duration // timeout for executing a query
	querier      mockprovider.Provider

	Keeper keeper.Keeper
}

func NewVoteExtHandler(
	logger log.Logger,
	timeout time.Duration,
	keeper keeper.Keeper,
	querier Provider,
) *VoteExtHandler {
	return &VoteExtHandler{
		logger:  logger,
		timeout: timeout,
		Keeper:  keeper,
		querier: querier,
	}
}

// QueryVoteExtension defines the canonical vote extension structure.
type QueryVoteExtension struct {
	Height      int64
	QueryResult interface{}
}

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		h.currentBlock = req.Height
		h.logger.Info("Executing query", "height", req.Height)

		g := new(errgroup.Group)

		// How an application determines which providers to use and for which pairs
		// can be done in a variety of ways. For demo purposes, we presume they are
		// locally configured. However, providers can be governed by governance.
		queryProvider := h.querier
		var (
			delResp stakingTypes.QueryDelegatorDelegationsResponse
			rewResp distTypes.QueryDelegationTotalRewardsResponse
			scSmart mockprovider.DataResponse
			err     error
		)

		// Launch a goroutine to fetch ticker prices from this oracle provider.
		// Recall, vote extensions are not required to be deterministic.
		g.Go(func() error {
			doneCh := make(chan bool, 1)
			errCh := make(chan error, 1)

			go func() {
				delResp, err = queryProvider.GetDelegatorDelegations(keeper.QueryDelegations)
				if err != nil {
					h.logger.Error("QUERY RESPONSE FAILURE: ", "err", err)
					errCh <- err
				}

				rewResp, err = queryProvider.GetDelegatorRewards(keeper.QueryDelegatorRewards)
				if err != nil {
					h.logger.Error("QUERY RESPONSE FAILURE: ", "err", err)
					errCh <- err
				}

				scSmart, err = queryProvider.GetWasmQuery()
				if err != nil {
					h.logger.Error("QUERY RESPONSE FAILURE: ", "err", err)
					errCh <- err
				}

				doneCh <- true
			}()

			select {
			case <-doneCh:
				break

			case err := <-errCh:
				return err

			case <-time.After(h.timeout):
				return fmt.Errorf("query timed out")
			}

			return nil
		})

		if err := g.Wait(); err != nil {
			h.logger.Error("failed to get all queries prices")
		}

		// produce a canonical vote extension
		voteExt := QueryVoteExtension{
			Height:      req.Height,
			QueryResult: delResp.String(),
		}

		h.logger.Info("QUERY RESPONSE SUCCESS: ", "delegations", delResp)
		h.logger.Info("QUERY RESPONSE SUCCESS: ", "rewards", rewResp)
		h.logger.Info("QUERY RESPONSE SUCCESS: ", "smart query", scSmart)

		// NOTE: We use stdlib JSON encoding, but an application may choose to use
		// a performant mechanism. This is for demo purposes only.
		bz, err := json.Marshal(voteExt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
		}

		return &abci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
