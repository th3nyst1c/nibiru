package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// TODO(sahith): check key values with param declarations
var (
	KeyStopped                 = []byte("Stopped")
	KeyFeePoolFeeRatio         = []byte("FeePoolFeeRatio")
	KeyEcosystemFundFeeRatio   = []byte("EcosystemFundFeeRatio")
	KeyLiquidationFeeRatio     = []byte("LiquidationFeeRatio")
	KeyPartialLiquidationRatio = []byte("PartialLiquidationRatio")
	KeyFundingRateInterval     = []byte("FundingRateInterval")
	KeyTwapLookbackWindow      = []byte("TwapLookbackWindow")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			KeyStopped,
			&p.Stopped,
			validateStopped,
		),
		paramtypes.NewParamSetPair(
			KeyFeePoolFeeRatio,
			&p.FeePoolFeeRatio,
			validatePercentageRatio,
		),
		paramtypes.NewParamSetPair(
			KeyEcosystemFundFeeRatio,
			&p.EcosystemFundFeeRatio,
			validatePercentageRatio,
		),
		paramtypes.NewParamSetPair(
			KeyLiquidationFeeRatio,
			&p.LiquidationFeeRatio,
			validatePercentageRatio,
		),
		paramtypes.NewParamSetPair(
			KeyPartialLiquidationRatio,
			&p.PartialLiquidationRatio,
			validatePercentageRatio,
		),
		paramtypes.NewParamSetPair(
			KeyFundingRateInterval,
			&p.FundingRateInterval,
			validateFundingRateInterval,
		),
		paramtypes.NewParamSetPair(
			KeyTwapLookbackWindow,
			&p.TwapLookbackWindow,
			validateTwapLookbackWindow,
		),
	}
}

// NewParams creates a new Params instance
func NewParams(
	stopped bool,
	feePoolFeeRatio sdk.Dec,
	ecosystemFundFeeRatio sdk.Dec,
	liquidationFeeRatio sdk.Dec,
	partialLiquidationRatio sdk.Dec,
	fundingRateInterval string,
	twapLookbackWindow time.Duration,
) Params {
	return Params{
		Stopped:                 stopped,
		FeePoolFeeRatio:         feePoolFeeRatio,
		EcosystemFundFeeRatio:   ecosystemFundFeeRatio,
		LiquidationFeeRatio:     liquidationFeeRatio,
		PartialLiquidationRatio: partialLiquidationRatio,
		FundingRateInterval:     fundingRateInterval,
		TwapLookbackWindow:      twapLookbackWindow,
	}
}

// DefaultParams returns the default parameters for the x/perp module.
func DefaultParams() Params {
	return NewParams(
		/* stopped */ false,
		/* feePoolFeeRatio */ sdk.NewDecWithPrec(1, 3), // 10 bps
		/* ecosystemFundFeeRatio */ sdk.NewDecWithPrec(1, 3), // 10 bps
		/* liquidationFee */ sdk.NewDecWithPrec(25, 3), // 250 bps
		/* partialLiquidationRatio */ sdk.NewDecWithPrec(25, 2),
		/* epochIdentifier */ "30 min",
		/* twapLookbackWindow */ 15*time.Minute,
	)
}

// Validate validates the set of params
func (p *Params) Validate() error {
	err := validateStopped(p.Stopped)
	if err != nil {
		return err
	}

	err = validatePercentageRatio(p.LiquidationFeeRatio)
	if err != nil {
		return err
	}

	err = validatePercentageRatio(p.FeePoolFeeRatio)
	if err != nil {
		return err
	}

	err = validatePercentageRatio(p.PartialLiquidationRatio)
	if err != nil {
		return err
	}

	return validatePercentageRatio(p.EcosystemFundFeeRatio)
}

func validatePercentageRatio(i interface{}) error {
	ratio, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if ratio.GT(sdk.OneDec()) {
		return fmt.Errorf("ratio is above max value(1.00): %s", ratio.String())
	} else if ratio.IsNegative() {
		return fmt.Errorf("ratio is negative: %s", ratio.String())
	}

	return nil
}

func validateFundingRateInterval(i interface{}) error {
	_, err := getAsString(i)
	if err != nil {
		return err
	}
	return nil
}

func validateStopped(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func getAsString(i interface{}) (string, error) {
	value, ok := i.(string)
	if !ok {
		return "invalid", fmt.Errorf("invalid parameter type: %T", i)
	}
	return value, nil
}

func validateTwapLookbackWindow(i interface{}) error {
	val, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if val <= 0 {
		return fmt.Errorf("twap lookback window must be positive, current value is %s", val.String())
	}
	return nil
}
