package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/NibiruChain/nibiru/x/perp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	FeePoolFeeRatio         = "feePoolFeeRatio"
	EcosystemFundFeeRatio   = "ecosystemFundFeeRatio"
	LiquidationFeeRatio     = "liquidationFeeRatio"
	PartialLiquidationRatio = "partialLiquidationRatio"
	FundingRateInterval     = "funding_rate_interval"
	TwapLookbackWindowKey   = "twap_lookback_window"
)

// RandomizedGenState generates a random genesis for pricefeed.
func RandomizedGenState(simState *module.SimulationState) {
	var stopped bool
	simState.AppParams.GetOrGenerate(
		simState.Cdc, string(types.KeyStopped), &stopped, simState.Rand,
		func(r *rand.Rand) { stopped = randomGenesisStoppedValue(r) },
	)
	var feePoolFeeRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, FeePoolFeeRatio, &feePoolFeeRatio, simState.Rand,
		func(r *rand.Rand) { feePoolFeeRatio = genFeePoolFeeRatio(r) },
	)
	var ecosystemFundFeeRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, EcosystemFundFeeRatio, &ecosystemFundFeeRatio, simState.Rand,
		func(r *rand.Rand) { ecosystemFundFeeRatio = genEcosystemFundFeeRatio(r) },
	)
	var liquidationFeeRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, LiquidationFeeRatio, &liquidationFeeRatio, simState.Rand,
		func(r *rand.Rand) { liquidationFeeRatio = genLiquidationFeeRatio(r) },
	)
	var partialLiquidationRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, PartialLiquidationRatio, &partialLiquidationRatio, simState.Rand,
		func(r *rand.Rand) { partialLiquidationRatio = genLiquidationFeeRatio(r) },
	)
	var twapLookbackWindow time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, TwapLookbackWindowKey, &twapLookbackWindow, simState.Rand,
		func(r *rand.Rand) { twapLookbackWindow = genTwapLookbackWindow(r) },
	)
	perpsGenesis := types.GenesisState{
		Params: types.Params{
			Stopped:                 stopped,
			FeePoolFeeRatio:         feePoolFeeRatio,
			EcosystemFundFeeRatio:   ecosystemFundFeeRatio,
			LiquidationFeeRatio:     liquidationFeeRatio,
			PartialLiquidationRatio: partialLiquidationRatio,
			//TODO: check why funding rate interval is string over duration
			FundingRateInterval: "",
			TwapLookbackWindow:  twapLookbackWindow,
		},
		VaultBalance:         sdk.Coins{},
		PerpEfBalance:        sdk.Coins{},
		FeePoolBalance:       sdk.Coins{},
		PairMetadata:         nil,
		Positions:            nil,
		PrepaidBadDebts:      nil,
		WhitelistedAddresses: nil,
	}

	bz, err := json.MarshalIndent(&perpsGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated perp parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&perpsGenesis)
}
