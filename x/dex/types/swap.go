package types

import (
	"fmt"

	"github.com/NibiruChain/nibiru/x/dex/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
Calculates the amount of tokenOut given tokenIn, deducting the swap fee.
Solved using the SolveConstantProductInvariant AMM curve.
Only supports single asset swaps.

args:
  - tokenIn: the amount of tokens to swap
  - tokenOutDenom: the target token denom

ret:
  - tokenOut: the tokens received from the swap
  - err: error if any
*/
func (pool Pool) CalcOutAmtGivenIn(tokenIn sdk.Coin, tokenOutDenom string) (
	tokenOut sdk.Coin, err error,
) {
	_, poolAssetIn, err := pool.getPoolAssetAndIndex(tokenIn.Denom)
	if err != nil {
		return tokenOut, err
	}

	_, poolAssetOut, err := pool.getPoolAssetAndIndex(tokenOutDenom)
	if err != nil {
		return tokenOut, err
	}

	tokenAmountInAfterFee := tokenIn.Amount.ToDec().Mul(sdk.OneDec().Sub(pool.PoolParams.SwapFee))
	poolTokenInBalance := poolAssetIn.Token.Amount.ToDec()
	poolTokenInBalancePostSwap := poolTokenInBalance.Add(tokenAmountInAfterFee)

	// deduct swapfee on the in asset
	// delta balanceOut is positive(tokens inside the pool decreases)
	tokenAmountOut := math.SolveConstantProductInvariant(
		/*xPrior=*/ poolTokenInBalance,
		/*xAfter=*/ poolTokenInBalancePostSwap,
		/*xWeight=*/ poolAssetIn.Weight.ToDec(),
		/*yPrior=*/ poolAssetOut.Token.Amount.ToDec(),
		/*yWeight=*/ poolAssetOut.Weight.ToDec(),
	).TruncateInt()
	return sdk.NewCoin(tokenOutDenom, tokenAmountOut), nil
}

/*
Calculates the amount of tokenIn required to obtain tokenOut coins from a swap,
accounting for additional fees.
Only supports single asset swaps.
This function is the inverse of CalcOutAmtGivenIn.

args:
  - tokenOut: the amount of tokens to swap
  - tokenInDenom: the target token denom

ret:
  - tokenIn: the tokens received from the swap
  - err: error if any
*/
func (pool Pool) CalcInAmtGivenOut(tokenOut sdk.Coin, tokenInDenom string) (
	tokenIn sdk.Coin, err error,
) {
	_, poolAssetOut, err := pool.getPoolAssetAndIndex(tokenOut.Denom)
	if err != nil {
		return tokenIn, err
	}

	_, poolAssetIn, err := pool.getPoolAssetAndIndex(tokenInDenom)
	if err != nil {
		return tokenIn, err
	}

	// assuming the user wishes to withdraw 'tokenOut', the balance of 'tokenOut' post swap will be lower
	poolTokenOutBalance := poolAssetOut.Token.Amount.ToDec()
	poolTokenOutBalancePostSwap := poolTokenOutBalance.Sub(tokenOut.Amount.ToDec())
	// (x_0)(y_0) = (x_0 + in)(y_0 - out)
	tokenAmountIn := math.SolveConstantProductInvariant(
		/*xPrior=*/ poolTokenOutBalance,
		/*xAfter=*/ poolTokenOutBalancePostSwap,
		/*xWeight=*/ poolAssetOut.Weight.ToDec(),
		/*yPrior=*/ poolAssetIn.Token.Amount.ToDec(),
		/*yWeight=*/ poolAssetIn.Weight.ToDec(),
	).Neg()

	// We deduct a swap fee on the input asset. The swap happens by following the invariant curve on the input * (1 - swap fee)
	// and then the swap fee is added to the pool.
	// Thus in order to give X amount out, we solve the invariant for the invariant input. However invariant input = (1 - swapfee) * trade input.
	// Therefore we divide by (1 - swapfee) here
	tokenAmountInBeforeFee := tokenAmountIn.Quo(sdk.OneDec().Sub(pool.PoolParams.SwapFee)).Ceil().TruncateInt()
	return sdk.NewCoin(tokenInDenom, tokenAmountInBeforeFee), nil
}

/*
Applies a swap to the pool by adding tokenIn and removing tokenOut from pool asset balances.

args:
  - tokenIn: the amount of token to deposit
  - tokenOut: the amount of token to withdraw

ret:
  - err: error if any
*/
func (pool *Pool) ApplySwap(tokenIn sdk.Coin, tokenOut sdk.Coin) (err error) {
	if tokenIn.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("tokenIn (%s) cannot be zero", tokenIn.Denom)
	}
	if tokenOut.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("tokenOut (%s) cannot be zero", tokenOut.Denom)
	}

	_, poolAssetIn, err := pool.getPoolAssetAndIndex(tokenIn.Denom)
	if err != nil {
		return err
	}

	_, poolAssetOut, err := pool.getPoolAssetAndIndex(tokenOut.Denom)
	if err != nil {
		return err
	}

	poolAssetIn.Token.Amount = poolAssetIn.Token.Amount.Add(tokenIn.Amount)
	poolAssetOut.Token.Amount = poolAssetOut.Token.Amount.Sub(tokenOut.Amount)

	return pool.updatePoolAssetBalances(sdk.NewCoins(
		poolAssetIn.Token,
		poolAssetOut.Token,
	))
}