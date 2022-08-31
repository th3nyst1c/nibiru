package oracle

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/NibiruChain/nibiru/simapp"
	"github.com/NibiruChain/nibiru/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	reflectionv2 "github.com/cosmos/cosmos-sdk/server/grpc/reflection/v2alpha1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
	"log"
	"math/big"
)

var (
	MaxSaltNumber = big.NewInt(100_000)
)

type PrevotesCache interface {
	SetPrevote(salt string, exchangeRatesStr, feeder string)
	GetPrevote() (salt, exchangeRatesStr, feeder string, ok bool)
}

type TxClient struct {
	feeder    sdk.AccAddress
	validator sdk.ValAddress
	prevotes  PrevotesCache

	authClient   authtypes.QueryClient
	chainID      string
	keyBase      keyring.Keyring
	txConfig     client.TxConfig
	newTxBuilder func() client.TxBuilder
	codec        codec.Codec
	ir           codectypes.InterfaceRegistry

	// for cleanup
	conn *grpc.ClientConn
}

type SymbolPrice struct {
	Symbol string
	Price  float64
}

func NewTxClient(grpcEndpoint string, validator sdk.ValAddress, feeder sdk.AccAddress, cache PrevotesCache, keyRing keyring.Keyring) (*TxClient, error) {
	// dial grpc
	conn, err := grpc.Dial(grpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	rc := reflectionv2.NewReflectionServiceClient(conn)
	chain, err := rc.GetChainDescriptor(context.Background(), &reflectionv2.GetChainDescriptorRequest{})
	if err != nil {
		return nil, err
	}

	// TODO(mercilex): use real app
	encConf := simapp.MakeTestEncodingConfig()

	// assert no errors in keybase

	return &TxClient{
		feeder:       feeder,
		validator:    validator,
		prevotes:     cache,
		authClient:   authtypes.NewQueryClient(conn),
		chainID:      chain.Chain.Id,
		keyBase:      keyRing,
		txConfig:     encConf.TxConfig,
		newTxBuilder: encConf.TxConfig.NewTxBuilder,
		codec:        encConf.Marshaler,
		ir:           encConf.InterfaceRegistry,
		conn:         conn,
	}, nil
}

func (c *TxClient) SendPrices(symbolPrices []SymbolPrice) error {
	// generate prevotes
	prevoteMsg := c.prevotesMsg(symbolPrices)
	voteMsg := c.voteMsg()

	for {
		log.Printf("sending prevote and vote:\n\t%s,\n\t %s", prevoteMsg, voteMsg)
		err := c.sendTx(prevoteMsg, voteMsg)
		if err != nil {
			log.Printf("failed sending tx: %s", err)

		}
	}
}

func (c *TxClient) prevotesMsg(prices []SymbolPrice) *types.MsgAggregateExchangeRatePrevote {
	tuple := make(types.ExchangeRateTuples, len(prices))
	for i, price := range prices {
		tuple[i] = types.ExchangeRateTuple{
			Pair:         price.Symbol,
			ExchangeRate: float64ToDec(price.Price),
		}
	}

	strTuple, err := tuple.ToString()
	if err != nil {
		panic(err)
	}

	nBig, err := rand.Int(rand.Reader, MaxSaltNumber)
	if err != nil {
		panic(err)
	}

	hash := types.GetAggregateVoteHash(nBig.String(), strTuple, c.validator)
	return types.NewMsgAggregateExchangeRatePrevote(hash, c.feeder, c.validator)
}

func (c *TxClient) voteMsg() *types.MsgAggregateExchangeRateVote {
	salt, exchangeRatesStr, feeder, ok := c.prevotes.GetPrevote()
	if !ok {
		return nil
	}

	// case where there's a feeder change there's nothing we can do
	// TODO(mercilex): we could support multi priv key feeders...
	if feeder != c.feeder.String() {
		return nil
	}

	return &types.MsgAggregateExchangeRateVote{
		Salt:          salt,
		ExchangeRates: exchangeRatesStr,
		Feeder:        feeder,
		Validator:     c.validator.String(),
	}
}

func (c *TxClient) sendTx(msgs ...sdk.Msg) error {
	// get key from keybase, can't fail
	keyInfo, err := c.keyBase.KeyByAddress(c.feeder)
	if err != nil {
		panic(err)
	}

	// set msgs, can't fail
	txBuilder := c.newTxBuilder()
	err = txBuilder.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}

	// get fee, can fail
	fee, gasLimit, err := c.getFee(msgs)
	if err != nil {
		return err
	}

	txBuilder.SetFeeAmount(fee)
	txBuilder.SetGasLimit(gasLimit)

	// get acc info, can fail
	accNum, sequence, err := c.getAccount()
	if err != nil {
		return err
	}

	txFactory := tx.Factory{}.
		WithChainID(c.chainID).
		WithKeybase(c.keyBase).
		WithTxConfig(c.txConfig).
		WithAccountNumber(accNum).
		WithSequence(sequence)

	// sign tx, can't fail
	err = tx.Sign(txFactory, keyInfo.GetName(), txBuilder, true)
	if err != nil {
		panic(err)
	}

	return nil
}

func (c *TxClient) getAccount() (uint64, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	accRaw, err := c.authClient.Account(ctx, &authtypes.QueryAccountRequest{Address: c.feeder.String()})
	if err != nil {
		panic(err)
	}

	var acc authtypes.AccountI
	err = c.ir.UnpackAny(accRaw.Account, &acc)
	if err != nil {
		panic(err)
	}

	return acc.GetAccountNumber(), acc.GetSequence(), nil
}

func (c *TxClient) getFee(_ []sdk.Msg) (sdk.Coins, uint64, error) {
	return sdk.NewCoins(sdk.NewInt64Coin("unibi", 1_000)), 1_000_000, nil
}

func (c *TxClient) Close() {
	_ = c.conn.Close()
}

func float64ToDec(price float64) sdk.Dec {
	// TODO(mercilex): precision for numbers with a lot of decimal digits
	return sdk.MustNewDecFromStr(fmt.Sprintf("%f", price))
}