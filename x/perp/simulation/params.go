package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/NibiruChain/nibiru/x/perp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const maxLookbackWindowMinutes = 7 * 24 * 60

func RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyStopped),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%v", randomGenesisStoppedValue(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyFeePoolFeeRatio),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genFeePoolFeeRatio(r))
			}),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyEcosystemFundFeeRatio),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genEcosystemFundFeeRatio(r))
			}),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyLiquidationFeeRatio),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genLiquidationFeeRatio(r))
			}),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyPartialLiquidationRatio),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genPartialLiquidationRatio(r))
			}),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyTwapLookbackWindow),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genTwapLookbackWindow(r))
			}),
	}
}

func randomGenesisStoppedValue(r *rand.Rand) bool {
	return r.Int63n(101) >= 90
}

func genFeePoolFeeRatio(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(30)), 2))
}

func genEcosystemFundFeeRatio(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(30)), 2))
}

func genLiquidationFeeRatio(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(30)), 2))
}

func genPartialLiquidationRatio(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 3).Add(sdk.NewDecWithPrec(int64(r.Intn(30)), 2))
}

func genTwapLookbackWindow(r *rand.Rand) time.Duration {
	return time.Duration(r.Intn(maxLookbackWindowMinutes)) * time.Minute
}
