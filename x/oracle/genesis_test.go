package oracle_test

import (
	"testing"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/NibiruChain/nibiru/collections/keys"

	"github.com/stretchr/testify/require"

	"github.com/NibiruChain/nibiru/x/oracle"
	"github.com/NibiruChain/nibiru/x/oracle/keeper"
	"github.com/NibiruChain/nibiru/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestExportInitGenesis(t *testing.T) {
	input := keeper.CreateTestInput(t)

	input.OracleKeeper.FeederDelegations.Insert(input.Ctx, keys.String(keeper.ValAddrs[0].String()), gogotypes.BytesValue{Value: keeper.Addrs[0].Bytes()})
	input.OracleKeeper.ExchangeRates.Insert(input.Ctx, "pair1:pair2", sdk.DecProto{Dec: sdk.NewDec(123)})
	input.OracleKeeper.Prevotes.Insert(input.Ctx, keys.String(keeper.ValAddrs[0].String()), types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{123}, keeper.ValAddrs[0], uint64(2)))
	input.OracleKeeper.Votes.Insert(input.Ctx, keys.String(keeper.ValAddrs[0].String()), types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Pair: "foo", ExchangeRate: sdk.NewDec(123)}}, keeper.ValAddrs[0]))
	input.OracleKeeper.Pairs.Insert(input.Ctx, "pair1:pair1")
	input.OracleKeeper.Pairs.Insert(input.Ctx, "pair2:pair2")
	input.OracleKeeper.MissCounters.Insert(input.Ctx, keys.String(keeper.ValAddrs[0].String()), gogotypes.UInt64Value{Value: 10})
	input.OracleKeeper.SetPairReward(input.Ctx, &types.PairReward{
		Pair: "pair1:pair2",
	})
	genesis := oracle.ExportGenesis(input.Ctx, input.OracleKeeper)

	newInput := keeper.CreateTestInput(t)
	oracle.InitGenesis(newInput.Ctx, newInput.OracleKeeper, genesis)
	newGenesis := oracle.ExportGenesis(newInput.Ctx, newInput.OracleKeeper)

	require.Equal(t, genesis, newGenesis)
}

func TestInitGenesis(t *testing.T) {
	input := keeper.CreateTestInput(t)
	genesis := types.DefaultGenesisState()
	require.NotPanics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.FeederDelegations = []types.FeederDelegation{{
		FeederAddress:    keeper.Addrs[0].String(),
		ValidatorAddress: "invalid",
	}}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.FeederDelegations = []types.FeederDelegation{{
		FeederAddress:    "invalid",
		ValidatorAddress: keeper.ValAddrs[0].String(),
	}}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.FeederDelegations = []types.FeederDelegation{{
		FeederAddress:    keeper.Addrs[0].String(),
		ValidatorAddress: keeper.ValAddrs[0].String(),
	}}

	genesis.MissCounters = []types.MissCounter{
		{
			ValidatorAddress: "invalid",
			MissCounter:      10,
		},
	}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.MissCounters = []types.MissCounter{
		{
			ValidatorAddress: keeper.ValAddrs[0].String(),
			MissCounter:      10,
		},
	}

	genesis.AggregateExchangeRatePrevotes = []types.AggregateExchangeRatePrevote{
		{
			Hash:        "hash",
			Voter:       "invalid",
			SubmitBlock: 100,
		},
	}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.AggregateExchangeRatePrevotes = []types.AggregateExchangeRatePrevote{
		{
			Hash:        "hash",
			Voter:       keeper.ValAddrs[0].String(),
			SubmitBlock: 100,
		},
	}

	genesis.AggregateExchangeRateVotes = []types.AggregateExchangeRateVote{
		{
			ExchangeRateTuples: []types.ExchangeRateTuple{
				{
					Pair:         "nibi:usd",
					ExchangeRate: sdk.NewDec(10),
				},
			},
			Voter: "invalid",
		},
	}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.AggregateExchangeRateVotes = []types.AggregateExchangeRateVote{
		{
			ExchangeRateTuples: []types.ExchangeRateTuple{
				{
					Pair:         "nibi:usd",
					ExchangeRate: sdk.NewDec(10),
				},
			},
			Voter: keeper.ValAddrs[0].String(),
		},
	}

	require.NotPanics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})
}
