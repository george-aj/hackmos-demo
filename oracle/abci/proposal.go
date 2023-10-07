package abci

import (
	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/facundomedica/oracle/keeper"
	"github.com/facundomedica/oracle/mockprovider"
)

type Provider interface {
	GetDelegatorDelegations(query stakingTypes.QueryDelegatorDelegationsRequest) (stakingTypes.QueryDelegatorDelegationsResponse, error)
	GetDelegatorRewards(query distTypes.QueryDelegationTotalRewardsRequest) (distTypes.QueryDelegationTotalRewardsResponse, error)
	GetWasmQuery() (mockprovider.DataResponse, error)
}

type ProposalHandler struct {
	logger   log.Logger
	keeper   keeper.Keeper
	valStore baseapp.ValidatorStore
}

func NewProposalHandler(logger log.Logger, keeper keeper.Keeper, valStore baseapp.ValidatorStore) *ProposalHandler {
	return &ProposalHandler{
		logger:   logger,
		keeper:   keeper,
		valStore: valStore,
	}
}

func (h *ProposalHandler) PrepareProposal() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
		if err != nil {
			return nil, err
		}
		/*

			var proposalTxs [][]byte


			if req.Height >= ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight && req.Height > 2 {
				stakeWeightedPrices, err := h.computeStakeWeightedOraclePrices(ctx, req.LocalLastCommit)
				if err != nil {
					return nil, errors.New("failed to compute stake-weighted oracle prices")
				}

				injectedVoteExtTx := StakeWeightedPrices{
					StakeWeightedPrices: stakeWeightedPrices,
					ExtendedCommitInfo:  req.LocalLastCommit,
				}

				// NOTE: We use stdlib JSON encoding, but an application may choose to use
				// a performant mechanism. This is for demo purposes only.
				bz, err := json.Marshal(injectedVoteExtTx)
				if err != nil {
					h.logger.Error("failed to encode injected vote extension tx", "err", err)
					return nil, errors.New("failed to encode injected vote extension tx")
				}

				// Inject a "fake" tx into the proposal s.t. validators can decode, verify,
				// and store the canonical stake-weighted average prices.
				proposalTxs = append(proposalTxs, bz)
			}

			// proceed with normal block proposal construction, e.g. POB, normal txs, etc...
			proposalTxs = append(proposalTxs, req.Txs...)

		*/

		return &abci.ResponsePrepareProposal{
			Txs: req.Txs,
		}, nil
	}
}

func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		if len(req.Txs) == 0 {
			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
		}

		/*
			var injectedVoteExtTx StakeWeightedPrices
			if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
				h.logger.Error("failed to decode injected vote extension tx", "err", err)
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}

			err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), injectedVoteExtTx.ExtendedCommitInfo)
			if err != nil {
				return nil, err
			}

			// Verify the proposer's stake-weighted oracle prices by computing the same
			// calculation and comparing the results. We omit verification for brevity
			// and demo purposes.
			stakeWeightedPrices, err := h.computeStakeWeightedOraclePrices(ctx, injectedVoteExtTx.ExtendedCommitInfo)
			if err != nil {
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}
			if err := compareOraclePrices(injectedVoteExtTx.StakeWeightedPrices, stakeWeightedPrices); err != nil {
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}*/

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

func (h *ProposalHandler) PreBlocker(_ sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	res := &sdk.ResponsePreBlock{}
	if len(req.Txs) == 0 {
		return res, nil
	}

	/*

		var injectedVoteExtTx StakeWeightedPrices
		if err := json.Unmarshal(req.Txs[0], &injectedVoteExtTx); err != nil {
			h.logger.Error("failed to decode injected vote extension tx", "err", err)
			return nil, err
		}

		// set oracle prices using the passed in context, which will make these prices available in the current block
		if err := h.keeper.SetOraclePrices(ctx, injectedVoteExtTx.StakeWeightedPrices); err != nil {
			return nil, err
		}
	*/
	return res, nil
}
