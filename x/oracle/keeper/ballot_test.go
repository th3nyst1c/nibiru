package keeper

import (
	"sort"
	"testing"

	"github.com/NibiruChain/nibiru/collections/keys"

	"github.com/NibiruChain/nibiru/x/common"

	"github.com/stretchr/testify/require"

	"github.com/NibiruChain/nibiru/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestOrganizeAggregate(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	btcBallot := types.ExchangeRateBallot{
		types.NewBallotVoteForTally(sdk.NewDec(17), common.PairBTCStable.String(), ValAddrs[0], power),
		types.NewBallotVoteForTally(sdk.NewDec(10), common.PairBTCStable.String(), ValAddrs[1], power),
		types.NewBallotVoteForTally(sdk.NewDec(6), common.PairBTCStable.String(), ValAddrs[2], power),
	}
	ethBallot := types.ExchangeRateBallot{
		types.NewBallotVoteForTally(sdk.NewDec(1000), common.PairETHStable.String(), ValAddrs[0], power),
		types.NewBallotVoteForTally(sdk.NewDec(1300), common.PairETHStable.String(), ValAddrs[1], power),
		types.NewBallotVoteForTally(sdk.NewDec(2000), common.PairETHStable.String(), ValAddrs[2], power),
	}

	for i := range btcBallot {
		input.OracleKeeper.Votes.Insert(
			input.Ctx,
			keys.String(ValAddrs[i].String()),
			types.NewAggregateExchangeRateVote(
				types.ExchangeRateTuples{
					{Pair: btcBallot[i].Pair, ExchangeRate: btcBallot[i].ExchangeRate},
					{Pair: ethBallot[i].Pair, ExchangeRate: ethBallot[i].ExchangeRate},
				},
				ValAddrs[i],
			),
		)
	}

	// organize votes by pair
	ballotMap := input.OracleKeeper.OrganizeBallotByPair(input.Ctx, map[string]types.ValidatorPerformance{
		ValAddrs[0].String(): {
			Power:      power,
			WinCount:   0,
			ValAddress: ValAddrs[0],
		},
		ValAddrs[1].String(): {
			Power:      power,
			WinCount:   0,
			ValAddress: ValAddrs[1],
		},
		ValAddrs[2].String(): {
			Power:      power,
			WinCount:   0,
			ValAddress: ValAddrs[2],
		},
	})

	// sort each ballot for comparison
	sort.Sort(btcBallot)
	sort.Sort(ethBallot)
	sort.Sort(ballotMap[common.PairBTCStable.String()])
	sort.Sort(ballotMap[common.PairETHStable.String()])

	require.Equal(t, btcBallot, ballotMap[common.PairBTCStable.String()])
	require.Equal(t, ethBallot, ballotMap[common.PairETHStable.String()])
}

func TestClearBallots(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	btcBallot := types.ExchangeRateBallot{
		types.NewBallotVoteForTally(sdk.NewDec(17), common.PairBTCStable.String(), ValAddrs[0], power),
		types.NewBallotVoteForTally(sdk.NewDec(10), common.PairBTCStable.String(), ValAddrs[1], power),
		types.NewBallotVoteForTally(sdk.NewDec(6), common.PairBTCStable.String(), ValAddrs[2], power),
	}
	ethBallot := types.ExchangeRateBallot{
		types.NewBallotVoteForTally(sdk.NewDec(1000), common.PairETHStable.String(), ValAddrs[0], power),
		types.NewBallotVoteForTally(sdk.NewDec(1300), common.PairETHStable.String(), ValAddrs[1], power),
		types.NewBallotVoteForTally(sdk.NewDec(2000), common.PairETHStable.String(), ValAddrs[2], power),
	}

	for i := range btcBallot {
		input.OracleKeeper.Prevotes.Insert(input.Ctx, keys.String(ValAddrs[i].String()), types.AggregateExchangeRatePrevote{
			Hash:        "",
			Voter:       ValAddrs[i].String(),
			SubmitBlock: uint64(input.Ctx.BlockHeight()),
		})

		input.OracleKeeper.Votes.Insert(input.Ctx, keys.String(ValAddrs[i].String()),
			types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
				{Pair: btcBallot[i].Pair, ExchangeRate: btcBallot[i].ExchangeRate},
				{Pair: ethBallot[i].Pair, ExchangeRate: ethBallot[i].ExchangeRate},
			}, ValAddrs[i]))
	}

	input.OracleKeeper.ClearBallots(input.Ctx, 5)

	prevoteCounter := len(input.OracleKeeper.Prevotes.Iterate(input.Ctx, keys.NewRange[keys.StringKey]()).Keys())
	voteCounter := len(input.OracleKeeper.Votes.Iterate(ctx, keys.NewRange[keys.StringKey]()).Keys())

	require.Equal(t, prevoteCounter, 3)
	require.Equal(t, voteCounter, 0)

	input.OracleKeeper.ClearBallots(input.Ctx.WithBlockHeight(input.Ctx.BlockHeight()+6), 5)

	prevoteCounter = len(input.OracleKeeper.Prevotes.Iterate(input.Ctx, keys.NewRange[keys.StringKey]()).Keys())
	require.Equal(t, prevoteCounter, 0)
}

func TestApplyWhitelist(t *testing.T) {
	input := CreateTestInput(t)

	whitelist := types.PairList{
		types.Pair{
			Name: "nibi:usd",
		},
		types.Pair{
			Name: "btc:usd",
		},
	}

	// prepare test by resetting the genesis pairs
	for _, p := range input.OracleKeeper.Pairs.Iterate(input.Ctx, keys.NewRange[keys.StringKey]()).Keys() {
		input.OracleKeeper.Pairs.Delete(input.Ctx, p)
	}
	for _, p := range whitelist {
		input.OracleKeeper.Pairs.Insert(input.Ctx, keys.String(p.Name))
	}

	voteTargets := map[string]struct{}{
		"nibi:usd": {},
		"btc:usd":  {},
	}
	// no updates case
	input.OracleKeeper.ApplyWhitelist(input.Ctx, whitelist, voteTargets)

	gotPairs := types.PairList{}

	for _, k := range input.OracleKeeper.Pairs.Iterate(input.Ctx, keys.NewRange[keys.StringKey]()).Keys() {
		gotPairs = append(gotPairs, types.Pair{Name: string(k)})
	}

	sort.Slice(whitelist, func(i, j int) bool {
		return whitelist[i].Name < whitelist[j].Name
	})
	require.Equal(t, whitelist, gotPairs)

	// len update (fast path)
	whitelist = append(whitelist, types.Pair{Name: "nibi:eth"})
	input.OracleKeeper.ApplyWhitelist(input.Ctx, whitelist, voteTargets)

	gotPairs = types.PairList{}

	for _, k := range input.OracleKeeper.Pairs.Iterate(input.Ctx, keys.NewRange[keys.StringKey]()).Keys() {
		gotPairs = append(gotPairs, types.Pair{Name: string(k)})
	}

	sort.Slice(whitelist, func(i, j int) bool {
		return whitelist[i].Name < whitelist[j].Name
	})
	require.Equal(t, whitelist, gotPairs)

	// diff update (slow path)
	voteTargets["nibi:eth"] = struct{}{}         // add previous pair
	whitelist[0] = types.Pair{Name: "nibi:usdt"} // update first pair
	input.OracleKeeper.ApplyWhitelist(input.Ctx, whitelist, voteTargets)

	gotPairs = types.PairList{}

	for _, k := range input.OracleKeeper.Pairs.Iterate(input.Ctx, keys.NewRange[keys.StringKey]()).Keys() {
		gotPairs = append(gotPairs, types.Pair{Name: string(k)})
	}

	sort.Slice(whitelist, func(i, j int) bool {
		return whitelist[i].Name < whitelist[j].Name
	})
	require.Equal(t, whitelist, gotPairs)
}
