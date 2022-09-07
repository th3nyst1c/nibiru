package oracle

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/NibiruChain/nibiru/feeder"
	"math/big"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txservice "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc"

	"github.com/NibiruChain/nibiru/simapp"
	oracletypes "github.com/NibiruChain/nibiru/x/oracle/types"
)

var (
	// MaxSaltNumber is the maximum salt number we can use for randomness.
	// NOTE: max length of the salt is 4.
	// TODO(mercilex): if we used digits + alphanumerics it's more randomized
	MaxSaltNumber = big.NewInt(9999) // NOTE(mercilex): max salt length is 4
)

var (
	_ feeder.TxClient = (*TxClient)(nil)
)

// TODO(mercilex): maybe prevote cache does not make any sense to exist
// considering that in case of oracle => stop/start then what's going
// to happen most likely is that the voting period will be over already
type PrevotesCache interface {
	SetPrevote(salt string, exchangeRatesStr, feeder string)
	GetPrevote() (salt, exchangeRatesStr, feeder string, ok bool)
}

type TxClient struct {
	feeder    sdk.AccAddress
	validator sdk.ValAddress
	prevotes  PrevotesCache

	authClient   authtypes.QueryClient
	oracleClient oracletypes.QueryClient
	txClient     txservice.ServiceClient
	chainID      string
	keyBase      keyring.Keyring
	txConfig     client.TxConfig
	newTxBuilder func() client.TxBuilder
	codec        codec.Codec
	ir           codectypes.InterfaceRegistry

	// for cleanup
	conn *grpc.ClientConn
}

func NewTxClient(grpcEndpoint string, validator sdk.ValAddress, feeder sdk.AccAddress, cache PrevotesCache, keyRing keyring.Keyring, chainID string) (*TxClient, error) {
	// dial grpc
	conn, err := grpc.Dial(grpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// TODO(mercilex): use real app
	encConf := simapp.MakeTestEncodingConfig()

	// assert no errors in keybase
	_, err = keyRing.KeyByAddress(feeder)
	if err != nil {
		return nil, err
	}

	return &TxClient{
		feeder:       feeder,
		validator:    validator,
		prevotes:     cache,
		authClient:   authtypes.NewQueryClient(conn),
		oracleClient: oracletypes.NewQueryClient(conn),
		txClient:     txservice.NewServiceClient(conn),
		chainID:      chainID,
		keyBase:      keyRing,
		txConfig:     encConf.TxConfig,
		newTxBuilder: encConf.TxConfig.NewTxBuilder,
		codec:        encConf.Marshaler,
		ir:           encConf.InterfaceRegistry,
		conn:         conn,
	}, nil
}

func (c *TxClient) SendPrices(ctx context.Context, symbolPrices []feeder.SymbolPrice) {
	// preparing the prevotes gets the new prices ready for the current voting period.
	prevoteMsg, salt, votesStr := c.preparePrevote(symbolPrices)

	// prepare to unveil votes of the last voting period.
	var voteMsg *oracletypes.MsgAggregateExchangeRateVote
	for {
		if ctxDone(ctx) {
			log.Error().Msg("context cancelled, exiting")
			break
		}

		votes, err := c.voteMsg(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to prepare votes")
			continue
		}

		voteMsg = votes
		break
	}

	var msgs []sdk.Msg
	if voteMsg != nil {
		// if order is inverted then current voting period prevotes will over-write the older ones.
		msgs = []sdk.Msg{voteMsg, prevoteMsg}
	} else {
		msgs = []sdk.Msg{prevoteMsg}
	}

	log.Info().Interface("prevote", prevoteMsg).Interface("vote", voteMsg).Msg("sending votes and prevotes")

	// attempt to send tx
	for {
		if ctxDone(ctx) {

		}

		err := c.sendTx(ctx, msgs...)
		if err != nil {
			log.Error().Err(err).Msg("failed to send tx")
			continue
		}
		break
	}

	c.prevotes.SetPrevote(salt, votesStr, c.feeder.String())
}

func (c *TxClient) preparePrevote(prices []feeder.SymbolPrice) (msg *oracletypes.MsgAggregateExchangeRatePrevote, salt, votesStr string) {
	tuple := make(oracletypes.ExchangeRateTuples, len(prices))
	for i, price := range prices {
		tuple[i] = oracletypes.ExchangeRateTuple{
			Pair:         price.Symbol,
			ExchangeRate: float64ToDec(price.Price),
		}
	}

	votesStr, err := tuple.ToString()
	if err != nil {
		panic(err)
	}
	nBig, err := rand.Int(rand.Reader, MaxSaltNumber)
	if err != nil {
		panic(err)
	}
	salt = nBig.String()

	hash := oracletypes.GetAggregateVoteHash(salt, votesStr, c.validator)

	return oracletypes.NewMsgAggregateExchangeRatePrevote(hash, c.feeder, c.validator), salt, votesStr
}

func (c *TxClient) voteMsg(ctx context.Context) (*oracletypes.MsgAggregateExchangeRateVote, error) {
	// there might be cases where due to downtimes the prevote
	// has expired. So we check if a prevote exists in the chain, if it does not
	// then we simply return.
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	resp, err := c.oracleClient.AggregatePrevote(ctx, &oracletypes.QueryAggregatePrevoteRequest{
		ValidatorAddr: c.validator.String(),
	})
	if err != nil {
		// TODO(mercilex): a better way?
		if strings.Contains(err.Error(), oracletypes.ErrNoAggregatePrevote.Error()) {
			log.Warn().Msg("no aggregate prevote found for this voting period")
			return nil, nil
		} else {
			return nil, err
		}
	}

	salt, exchangeRatesStr, feeder, ok := c.prevotes.GetPrevote()
	if !ok {
		return nil, nil
	}

	// assert equality between feeder's prevote and chain's prevote
	if localHash := oracletypes.GetAggregateVoteHash(salt, exchangeRatesStr, c.validator).String(); localHash != resp.AggregatePrevote.Hash {
		log.Warn().Str("chain hash", resp.AggregatePrevote.Hash).Str("local hash", localHash).Msg("chain and local prevote do not match")
		return nil, nil
	}

	// case where there's a feeder change there's nothing we can do
	// TODO(mercilex): we could support multi priv key feeders...
	if feeder != c.feeder.String() {
		return nil, nil
	}

	return &oracletypes.MsgAggregateExchangeRateVote{
		Salt:          salt,
		ExchangeRates: exchangeRatesStr,
		Feeder:        feeder,
		Validator:     c.validator.String(),
	}, nil
}

func (c *TxClient) sendTx(ctx context.Context, msgs ...sdk.Msg) error {
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

	txBytes, err := c.txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()
	resp, err := c.txClient.BroadcastTx(ctx, &txservice.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    txservice.BroadcastMode_BROADCAST_MODE_BLOCK,
	})
	if err != nil {
		return err
	}
	if resp.TxResponse.Code != tmtypes.CodeTypeOK {
		return fmt.Errorf("tx failed: %s", resp.TxResponse.RawLog)
	}

	return nil
}

func (c *TxClient) getAccount() (uint64, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	accRaw, err := c.authClient.Account(ctx, &authtypes.QueryAccountRequest{Address: c.feeder.String()})
	// TODO(mercilex): need to investigate account not found error.
	if err != nil {
		return 0, 0, err
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

func ctxDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}