package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/testutil/mock"
	"github.com/NibiruChain/nibiru/x/vpool/types"
)

func TestSwapQuoteForBase(t *testing.T) {
	tests := []struct {
		name                 string
		pair                 common.TokenPair
		direction            types.Direction
		quoteAmount          sdk.Dec
		baseLimit            sdk.Dec
		expectedQuoteReserve sdk.Dec
		expectedBaseReserve  sdk.Dec
		expectedBaseAmount   sdk.Dec
		expectedErr          error
	}{
		{
			name:                 "quote amount == 0",
			pair:                 NUSDPair,
			direction:            types.Direction_ADD_TO_POOL,
			quoteAmount:          sdk.NewDec(0),
			baseLimit:            sdk.NewDec(10),
			expectedQuoteReserve: sdk.NewDec(10_000_000),
			expectedBaseReserve:  sdk.NewDec(5_000_000),
			expectedBaseAmount:   sdk.ZeroDec(),
		},
		{
			name:                 "normal swap add",
			pair:                 NUSDPair,
			direction:            types.Direction_ADD_TO_POOL,
			quoteAmount:          sdk.NewDec(100_000),
			baseLimit:            sdk.NewDec(49504),
			expectedQuoteReserve: sdk.NewDec(10_100_000),
			expectedBaseReserve:  sdk.MustNewDecFromStr("4950495.049504950495049505"),
			expectedBaseAmount:   sdk.MustNewDecFromStr("49504.950495049504950495"),
		},
		{
			name:                 "normal swap remove",
			pair:                 NUSDPair,
			direction:            types.Direction_REMOVE_FROM_POOL,
			quoteAmount:          sdk.NewDec(100_000),
			baseLimit:            sdk.NewDec(50506),
			expectedQuoteReserve: sdk.NewDec(9_900_000),
			expectedBaseReserve:  sdk.MustNewDecFromStr("5050505.050505050505050505"),
			expectedBaseAmount:   sdk.MustNewDecFromStr("50505.050505050505050505"),
		},
		{
			name:        "pair not supported",
			pair:        "xxx:yyy",
			direction:   types.Direction_ADD_TO_POOL,
			quoteAmount: sdk.NewDec(10),
			baseLimit:   sdk.NewDec(10),
			expectedErr: types.ErrPairNotSupported,
		},
		{
			name:        "base amount less than base limit in Long",
			pair:        NUSDPair,
			direction:   types.Direction_ADD_TO_POOL,
			quoteAmount: sdk.NewDec(500_000),
			baseLimit:   sdk.NewDec(454_500),
			expectedErr: types.ErrAssetOverUserLimit,
		},
		{
			name:        "base amount more than base limit in Short",
			pair:        NUSDPair,
			direction:   types.Direction_REMOVE_FROM_POOL,
			quoteAmount: sdk.NewDec(1_000_000),
			baseLimit:   sdk.NewDec(454_500),
			expectedErr: types.ErrAssetOverUserLimit,
		},
		{
			name:        "quote input bigger than reserve ratio",
			pair:        NUSDPair,
			direction:   types.Direction_REMOVE_FROM_POOL,
			quoteAmount: sdk.NewDec(10_000_000),
			baseLimit:   sdk.NewDec(10),
			expectedErr: types.ErrOverTradingLimit,
		},
		{
			name:        "over fluctuation limit fails",
			pair:        NUSDPair,
			direction:   types.Direction_ADD_TO_POOL,
			quoteAmount: sdk.NewDec(1_000_000),
			baseLimit:   sdk.NewDec(454_544),
			expectedErr: types.ErrOverFluctuationLimit,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vpoolKeeper, ctx := VpoolKeeper(t,
				mock.NewMockPriceKeeper(gomock.NewController(t)),
			)

			vpoolKeeper.CreatePool(
				ctx,
				NUSDPair,
				sdk.MustNewDecFromStr("0.9"), // 0.9 ratio
				sdk.NewDec(10_000_000),       // 10 tokens
				sdk.NewDec(5_000_000),        // 5 tokens
				sdk.MustNewDecFromStr("0.1"), // 0.1 ratio
				sdk.MustNewDecFromStr("0.1"),
			)

			res, err := vpoolKeeper.SwapQuoteForBase(
				ctx,
				tc.pair,
				tc.direction,
				tc.quoteAmount,
				tc.baseLimit,
			)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedBaseAmount, res)

				pool, err := vpoolKeeper.getPool(ctx, NUSDPair)
				require.NoError(t, err)
				require.Equal(t, tc.expectedQuoteReserve, pool.QuoteAssetReserve)
				require.Equal(t, tc.expectedBaseReserve, pool.BaseAssetReserve)

				snapshot, _, err := vpoolKeeper.getLatestReserveSnapshot(ctx, NUSDPair)
				require.NoError(t, err)
				require.EqualValues(t, tc.expectedQuoteReserve, snapshot.QuoteAssetReserve)
				require.EqualValues(t, tc.expectedBaseReserve, snapshot.BaseAssetReserve)
			}
		})
	}
}

func TestSwapBaseForQuote(t *testing.T) {
	tests := []struct {
		name                     string
		initialQuoteReserve      sdk.Dec
		initialBaseReserve       sdk.Dec
		direction                types.Direction
		baseAssetAmount          sdk.Dec
		quoteAssetLimit          sdk.Dec
		expectedQuoteReserve     sdk.Dec
		expectedBaseReserve      sdk.Dec
		expectedQuoteAssetAmount sdk.Dec
		expectedErr              error
	}{
		{
			name:                     "zero base asset swap",
			initialQuoteReserve:      sdk.NewDec(10_000_000),
			initialBaseReserve:       sdk.NewDec(5_000_000),
			direction:                types.Direction_ADD_TO_POOL,
			baseAssetAmount:          sdk.ZeroDec(),
			quoteAssetLimit:          sdk.ZeroDec(),
			expectedQuoteReserve:     sdk.NewDec(10_000_000),
			expectedBaseReserve:      sdk.NewDec(5_000_000),
			expectedQuoteAssetAmount: sdk.ZeroDec(),
		},
		{
			name:                     "add base asset swap",
			initialQuoteReserve:      sdk.NewDec(10_000_000),
			initialBaseReserve:       sdk.NewDec(5_000_000),
			direction:                types.Direction_ADD_TO_POOL,
			baseAssetAmount:          sdk.NewDec(1_000_000),
			quoteAssetLimit:          sdk.NewDec(1_666_666),
			expectedQuoteReserve:     sdk.MustNewDecFromStr("8333333.333333333333333333"),
			expectedBaseReserve:      sdk.NewDec(6_000_000),
			expectedQuoteAssetAmount: sdk.MustNewDecFromStr("1666666.666666666666666667"),
		},
		{
			name:                     "remove base asset",
			initialQuoteReserve:      sdk.NewDec(10_000_000),
			initialBaseReserve:       sdk.NewDec(5_000_000),
			direction:                types.Direction_REMOVE_FROM_POOL,
			baseAssetAmount:          sdk.NewDec(1_000_000),
			quoteAssetLimit:          sdk.NewDec(2_500_001),
			expectedQuoteReserve:     sdk.NewDec(12_500_000),
			expectedBaseReserve:      sdk.NewDec(4_000_000),
			expectedQuoteAssetAmount: sdk.NewDec(2_500_000),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			vpoolKeeper, ctx := VpoolKeeper(t,
				mock.NewMockPriceKeeper(gomock.NewController(t)),
			)

			vpoolKeeper.CreatePool(
				ctx,
				NUSDPair,
				sdk.OneDec(),
				tc.initialQuoteReserve,
				tc.initialBaseReserve,
				sdk.OneDec(),
				sdk.OneDec(),
			)

			quoteAssetAmount, err := vpoolKeeper.SwapBaseForQuote(
				ctx,
				NUSDPair,
				tc.direction,
				tc.baseAssetAmount,
				tc.quoteAssetLimit,
			)

			if tc.expectedErr != nil {
				require.Error(t, err)
			} else {
				pool, err := vpoolKeeper.getPool(ctx, NUSDPair)
				require.NoError(t, err)

				require.EqualValuesf(t, tc.expectedQuoteAssetAmount, quoteAssetAmount,
					"expected %s; got %s", tc.expectedQuoteAssetAmount.String(), quoteAssetAmount.String())
				require.NoError(t, err)
				require.Equal(t, tc.expectedQuoteReserve, pool.QuoteAssetReserve)
				require.Equal(t, tc.expectedBaseReserve, pool.BaseAssetReserve)

				snapshot, _, err := vpoolKeeper.getLatestReserveSnapshot(ctx, NUSDPair)
				require.NoError(t, err)
				require.EqualValues(t, tc.expectedQuoteReserve, snapshot.QuoteAssetReserve)
				require.EqualValues(t, tc.expectedBaseReserve, snapshot.BaseAssetReserve)
			}
		})
	}
}