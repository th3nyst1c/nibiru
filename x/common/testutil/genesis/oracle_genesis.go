package genesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/app"
	"github.com/NibiruChain/nibiru/x/common/asset"
	"github.com/NibiruChain/nibiru/x/common/denoms"
	oracletypes "github.com/NibiruChain/nibiru/x/oracle/types"
)

func AddOracleGenesis(gen app.GenesisState) app.GenesisState {
	oracleGenesis := oracletypes.DefaultGenesisState()
	oracleGenesis.ExchangeRates = []oracletypes.ExchangeRateTuple{
		{Pair: asset.Registry.Pair(denoms.ETH, denoms.NUSD), ExchangeRate: sdk.NewDec(1_000)},
		{Pair: asset.Registry.Pair(denoms.NIBI, denoms.NUSD), ExchangeRate: sdk.NewDec(10)},
	}
	oracleGenesis.Params.VotePeriod = 1_000

	gen[oracletypes.ModuleName] = app.MakeEncodingConfig().Marshaler.
		MustMarshalJSON(oracleGenesis)
	return gen
}
