package app_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/NibiruChain/nibiru/app"
	pricefeedtypes "github.com/NibiruChain/nibiru/x/pricefeed/types"
	"github.com/NibiruChain/nibiru/x/testutil"
	"github.com/NibiruChain/nibiru/x/testutil/sample"

	sdk "github.com/cosmos/cosmos-sdk/types"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"

	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	localhosttypes "github.com/cosmos/ibc-go/v3/modules/light-clients/09-localhost/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	ibcmock "github.com/cosmos/ibc-go/v3/testing/mock"

	"github.com/stretchr/testify/suite"
)

/* SetupTestingApp returns the TestingApp and default genesis state used to
   initialize the testing app. */
func SetupNibiruTestingApp() (
	testingApp ibctesting.TestingApp,
	defaultGenesis map[string]json.RawMessage,
) {
	// create testing app
	nibiruApp, ctx := testutil.NewNibiruApp(true)
	token0, token1 := "uatom", "unibi"
	oracle := sample.AccAddress()
	nibiruApp.PriceKeeper.SetParams(ctx, pricefeedtypes.Params{
		Pairs: []pricefeedtypes.Pair{
			{Token0: token0, Token1: token1,
				Oracles: []sdk.AccAddress{oracle}, Active: true},
		},
	})
	nibiruApp.PriceKeeper.SetPrice(
		ctx, oracle, token0, token1, sdk.OneDec(),
		ctx.BlockTime().Add(time.Hour),
	)
	nibiruApp.PriceKeeper.SetCurrentPrices(ctx, token0, token1)

	// Create genesis state
	encCdc := app.MakeTestEncodingConfig()
	genesisState := app.NewDefaultGenesisState(encCdc.Marshaler)

	return nibiruApp, genesisState
}

// init changes the value of 'DefaultTestingAppInit' to use custom initialization.
func init() {
	ibctesting.DefaultTestingAppInit = SetupNibiruTestingApp
}

// IBCTestSuite is a testing suite to test keeper functions.
type IBCTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	path *ibctesting.Path // chainA <---> chainB
}

// TestIBCTestSuite runs all the tests within this package.
func TestIBCTestSuite(t *testing.T) {
	suite.Run(t, new(IBCTestSuite))
}

/* NewIBCTestingTransferPath returns a "path" for testing.
   A path contains two endpoints, 'EndpointA' and 'EndpointB' that correspond
   to the order of the chains passed into the ibctesting.NewPath function.
   A path is a pointer, and its values will be filled in as necessary during
   the setup portion of testing.
*/
func NewIBCTestingTransferPath(
	chainA, chainB *ibctesting.TestChain,
) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointA.ChannelConfig.Version = "ics20-1"
	path.EndpointB.ChannelConfig.Version = "ics20-1"
	return path
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *IBCTestSuite) SetupTest() {
	// initializes 2 test chains
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	// clientID, connectionID, channelID empty
	suite.path = NewIBCTestingTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)

	suite.coordinator.SetupClients(suite.path)
	suite.Require().Equal("07-tendermint-0", suite.path.EndpointA.ClientID)

	suite.coordinator.SetupConnections(suite.path)
	suite.Require().Equal("connection-0", suite.path.EndpointA.ConnectionID)

	suite.coordinator.ChanOpenInitOnBothChains(suite.path)
	suite.Require().Equal("channel-0", suite.path.EndpointA.ChannelID)
	// clientID, connectionID, channelID filled

	/* NOTE: Investigate the difference between individual Setup calls and
	   suite.coordinator.Setup(suite.path)
	*/
}

func (suite IBCTestSuite) TestInitialization() {
	suite.SetupTest()

	var err error = suite.coordinator.ConnOpenInitOnBothChains(suite.path)
	suite.Require().NoError(err)
}

func (suite IBCTestSuite) TestClient_BeginBlocker() {

	// set localhost client
	setLocalHostClient := func() {
		revision := ibcclienttypes.ParseChainID(suite.chainA.GetContext().ChainID())
		localHostClient := localhosttypes.NewClientState(
			suite.chainA.GetContext().ChainID(),
			ibcclienttypes.NewHeight(revision, uint64(suite.chainA.GetContext().BlockHeight())),
		)
		suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(
			suite.chainA.GetContext(), ibcexported.Localhost, localHostClient)
	}
	setLocalHostClient()

	prevHeight := ibcclienttypes.GetSelfHeight(suite.chainA.GetContext())

	localHostClient := suite.chainA.GetClientState(ibcexported.Localhost)
	suite.Require().Equal(prevHeight, localHostClient.GetLatestHeight())

	for i := 0; i < 10; i++ {
		// increment height
		suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

		suite.Require().NotPanics(func() {
			ibcclient.BeginBlocker(
				suite.chainA.GetContext(), suite.chainA.App.GetIBCKeeper().ClientKeeper)
		}, "BeginBlocker shouldn't panic")

		localHostClient = suite.chainA.GetClientState(ibcexported.Localhost)
		suite.Require().Equal(prevHeight.Increment(), localHostClient.GetLatestHeight())
		prevHeight = localHostClient.GetLatestHeight().(ibcclienttypes.Height)
	}

}

func NewPacket(
	path *ibctesting.Path,
	sender string, receiver string,
	coin sdk.Coin,
	timeoutHeight ibcclienttypes.Height,
) channeltypes.Packet {

	transfer := transfertypes.NewFungibleTokenPacketData(
		coin.Denom, coin.Amount.String(), sender, receiver)
	bz := transfertypes.ModuleCdc.MustMarshalJSON(&transfer)
	packet := channeltypes.Packet{
		Data:               bz,
		Sequence:           1,
		SourcePort:         path.EndpointA.ChannelConfig.PortID,
		SourceChannel:      path.EndpointA.ChannelID,
		DestinationPort:    path.EndpointB.ChannelConfig.PortID,
		DestinationChannel: path.EndpointB.ChannelID,
		TimeoutHeight:      timeoutHeight,
		TimeoutTimestamp:   0}
	return packet
}

func (suite IBCTestSuite) TestSentPacket() {
	suite.SetupTest()

	timeoutHeight := ibcclienttypes.NewHeight(1000, 1000)
	ack := ibcmock.MockAcknowledgement

	// create packet 1
	sender := sample.AccAddress().String()
	receiver := sample.AccAddress().String()
	coin := sdk.NewInt64Coin("unibi", 1000)
	path := suite.path
	packet1 := NewPacket(path, sender, receiver, coin, timeoutHeight)

	fmt.Println(ack, packet1)

	// ---------------------------- Below is work in progress
	// --------------------------------------------------------

	// // send on endpointA
	// suite.coordinator.CreateConnections(path) // TODO fix: raises error

	// var err error
	// fmt.Println("client: ", path.EndpointB.ClientID)
	// fmt.Println("connection: ", path.EndpointB.ConnectionID)
	// fmt.Println("channel: ", path.EndpointB.ChannelID)
	// fmt.Println("packet destination channel: ", packet1.DestinationChannel)

	// err = path.EndpointA.SendPacket(packet1)
	// suite.Require().NoError(err)

	// // receive on endpointB
	// err = path.EndpointB.RecvPacket(packet1)
	// suite.Require().NoError(err)

	// // acknowledge the receipt of the packet
	// err = path.EndpointA.AcknowledgePacket(packet1, ack.Acknowledgement())
	// suite.Require().NoError(err)

	// err = simapp.FundModuleAccount(
	// 	/* bankKeeper */ suite.chainA.App.(*app.NibiruApp).BankKeeper,
	// 	/* ctx */ suite.chainB.GetContext(),
	// 	/* recipientModule */ stabletypes.ModuleName,
	// 	/* coins */ sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	// )
	// suite.Require().NoError(err)

	// // we can also relay
	// packet2 := channeltypes.NewPacket()

	// path.EndpointA.SendPacket(packet2)

	// path.Relay(packet2, expectedAck)

	// // if needed we can update our clients
	// path.EndpointB.UpdateClient()
}