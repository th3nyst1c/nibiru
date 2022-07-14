package types

import (
	"github.com/NibiruChain/nibiru/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Vpools: []*Pool{
			{
				Pair:                  common.PairBTCStable,
				BaseAssetReserve:      sdk.NewDec(500),        // 500 btc
				QuoteAssetReserve:     sdk.NewDec(10_000_000), // 10 million unusd
				TradeLimitRatio:       sdk.OneDec(),
				FluctuationLimitRatio: sdk.OneDec(),
				MaxOracleSpreadRatio:  sdk.OneDec(),
			},
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
