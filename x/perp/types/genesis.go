package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/common"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		VaultBalance:   nil,
		PerpEfBalance:  nil,
		FeePoolBalance: nil,
		PairMetadata: []*PairMetadata{
			{
				Pair: common.PairBTCStable.String(),
				CumulativePremiumFractions: []sdk.Dec{
					sdk.ZeroDec(),
				},
			},
		},
		Positions:            nil,
		PrepaidBadDebts:      nil,
		WhitelistedAddresses: nil,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
