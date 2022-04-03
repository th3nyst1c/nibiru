package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/MatrixDao/matrix/x/stablecoin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	ParamSubspace paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	priceKeeper   types.PriceKeeper
}

// NewKeeper Creates a new x/stablecoin Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	paramSubspace paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	priceKeeper types.PriceKeeper,
) Keeper {

	// Ensure that the module account is set.
	if moduleAcc := accountKeeper.GetModuleAddress(types.ModuleName); moduleAcc == nil {
		panic("The stablecoin module account has not been set")
	}

	// Set param.types.'KeyTable' if it has not already been set
	if !paramSubspace.HasKeyTable() {
		paramSubspace = paramSubspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		ParamSubspace: paramSubspace,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		priceKeeper:   priceKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
