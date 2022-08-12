#!/bin/sh
#set -e

# Console log text colour
red=$(tput setaf 9)
green=$(tput setaf 10)
blue=$(tput setaf 12)
reset=$(tput sgr0)

echo_info() {
  echo "${blue}"
  echo "$1"
  echo "${reset}"
}

echo_error() {
  echo "${red}"
  echo "$1"
  echo "${reset}"
}

echo_success() {
  echo "${green}"
  echo "$1"
  echo "${reset}"
}

echo_info "Building from source..."
if make install; then
  echo_success "Successfully built binary"
else
  echo_error "Could not build binary. Failed to make install."
  exit 1
fi

# Set localnet settings
BINARY=nibid
CHAIN_ID=nibiru-localnet-0
RPC_PORT=26657
GRPC_PORT=9090
VALIDATOR_MNEMONIC="guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"
SHRIMP_MNEMONIC="item segment elevator fork swim tone search hope enough asthma apology pact embody extra trash educate deposit raccoon giant gift essay able female develop"
WHALE_MNEMONIC="throw oblige vague twist clutch grunt physical sell conduct blossom owner delay suspect square kidney joy define book boss outside reason silk success you"
LIQUIDATOR_MNEMONIC="oxygen tattoo pond upgrade barely sudden wheat correct dumb roast glance conduct scene profit female health speak hire north grab allow provide depend away"
ORACLE_MNEMONIC="abandon wave reason april rival valid saddle cargo aspect toe tomato stomach zero quick side potato artwork mixture over basket sort churn palace cherry"
GENESIS_COINS=10000000000000000000unibi,10000000000000000000unusd

# Stop nibid if it is already running
if pgrep -x "$BINARY" >/dev/null; then
  echo_error "Terminating $BINARY..."
  killall nibid
fi

# Remove previous data
echo_info "Removing previous chain data from $HOME/.nibid..."
rm -rf $HOME/.nibid

# Add directory for chain, exit if error
if ! mkdir -p $HOME/.nibid 2>/dev/null; then
  echo_error "Failed to create chain folder. Aborting..."
  exit 1
fi

# Initialize nibid with "localnet" chain id
echo_info "Initializing $CHAIN_ID..."
if $BINARY init nibiru-localnet-0 --chain-id $CHAIN_ID --overwrite; then
  echo_success "Successfully initialized $CHAIN_ID"
else
  echo_error "Failed to initialize $CHAIN_ID"
fi

# Configure keyring-backend to "test"
echo_info "Configuring keyring-backend..."
if $BINARY config keyring-backend test; then
  echo_success "Successfully configured keyring-backend"
else
  echo_error "Failed to configure keyring-backend"
fi

# Configure chain-id
echo_info "Configuring chain-id..."
if $BINARY config chain-id $CHAIN_ID; then
  echo_success "Successfully configured chain-id"
else
  echo_error "Failed to configure chain-id"
fi

# Configure broadcast mode
echo_info "Configuring broadcast mode..."
if $BINARY config broadcast-mode block; then
  echo_success "Successfully configured broadcast-mode"
else
  echo_error "Failed to configure broadcast mode"
fi

# Configure output mode
echo_info "Configuring output mode..."
if $BINARY config output json; then
  echo_success "Successfully configured output mode"
else
  echo_error "Failed to configure output mode"
fi

# Enable API Server
echo_info "Enabling API server"
if sed -i '' '/\[api\]/,+3 s/enable = false/enable = true/' $HOME/.nibid/config/app.toml; then
  echo_success "Successfully enabled API server"
else
  echo_error "Failed to enable API server"
fi

# Enable Swagger Docs
echo_info "Enabling Swagger Docs"
if sed -i '' 's/swagger = false/swagger = true/' $HOME/.nibid/config/app.toml; then
  echo_success "Successfully enabled Swagger Docs"
else
  echo_error "Failed to enable Swagger Docs"
fi

# Enable CORS for grpc http server
echo_info "Enabling CORS"
if sed -i '' 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' $HOME/.nibid/config/app.toml; then
  echo_success "Successfully enabled CORS"
else
  echo_error "Failed to enable CORS"
fi

# Enable CORS for tendermint grpc server
echo_info "Enabling CORS"
if sed -i '' 's/cors_allowed_origins = \[]/cors_allowed_origins = ["*"]/' $HOME/.nibid/config/config.toml; then
  echo_success "Successfully enabled CORS"
else
  echo_error "Failed to enable CORS"
fi

echo_info "Adding genesis accounts..."

# validator
echo "$VALIDATOR_MNEMONIC" | $BINARY keys add validator --recover
$BINARY add-genesis-account $($BINARY keys show validator -a) $GENESIS_COINS

# shrimp
echo "$SHRIMP_MNEMONIC" | $BINARY keys add shrimp --recover
$BINARY add-genesis-account $($BINARY keys show shrimp -a) $GENESIS_COINS

# whale
echo "$WHALE_MNEMONIC" | $BINARY keys add whale --recover
$BINARY add-genesis-account $($BINARY keys show whale -a) $GENESIS_COINS

# liquidator
echo "$LIQUIDATOR_MNEMONIC" | $BINARY keys add liquidator --recover
$BINARY add-genesis-account $($BINARY keys show liquidator -a) $GENESIS_COINS

# oracle
echo "$ORACLE_MNEMONIC" | $BINARY keys add oracle --recover --keyring-backend test
$BINARY add-genesis-account $($BINARY keys show oracle -a --keyring-backend test) $GENESIS_COINS

echo_success "Genesis accounts added"

echo_info "Adding gentx validator..."
if $BINARY gentx validator 900000000unibi --chain-id $CHAIN_ID; then
  echo_success "Successfully added gentx"
else
  echo_error "Failed to add gentx"
fi

echo_info "Collecting gentx..."
if $BINARY collect-gentxs; then
  echo_success "Successfully collected genesis txs into genesis.json"
else
  echo_error "Failed to collect genesis txs"
fi

echo_info "Configuring genesis params"
# x/vpool
# BTC:NUSD
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].pair = {token0:"ubtc",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].base_asset_reserve = "10000000000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json      # 50k BTC
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].quote_asset_reserve = "200000000000000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json # 1 billion NUSD
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].trade_limit_ratio = "0.1"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].fluctuation_limit_ratio = "0.1"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].max_oracle_spread_ratio = "0.1"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].maintenance_margin_ratio = "0.0625"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[0].max_leverage = "15"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
# ETH:NUSD
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].pair = {token0:"ueth",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].base_asset_reserve = "666666000000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json      # 666k ETH
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].quote_asset_reserve = "1000000000000000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json # 1 billion NUSD
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].trade_limit_ratio = "0.1"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].fluctuation_limit_ratio = "0.1"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].max_oracle_spread_ratio = "0.1"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].maintenance_margin_ratio = "0.0625"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.vpool.vpools[1].max_leverage = "12"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
# x/perp
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.params.stopped = false' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.params.fee_pool_fee_ratio = "0.001"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.params.ecosystem_fund_fee_ratio = "0.001"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.params.liquidation_fee_ratio = "0.025"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.params.partial_liquidation_ratio = "0.25"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.params.funding_rate_interval = "30 min"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.params.twap_lookback_window = "900s"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.pair_metadata[0].pair = {token0:"ubtc",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.pair_metadata[0].cumulative_premium_fractions = ["0"]' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.pair_metadata[1].pair = {token0:"ueth",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.perp.pair_metadata[1].cumulative_premium_fractions = ["0"]' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
# x/pricefeed
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.params.pairs[0] = {token0:"unibi",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.params.pairs[1] = {token0:"uusdc",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.params.pairs[2] = {token0:"ubtc",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.params.pairs[3] = {token0:"ueth",token1:"unusd"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.params.twap_lookback_window = "900s"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.genesis_oracles = ["nibi1zaavvzxez0elundtn32qnk9lkm8kmcsz44g7xl", "nibi1pzd5e402eld9kcc3h78tmfrm5rpzlzk6hnxkvu"]' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.posted_prices[0].pair_id = "ubtc:unusd"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.posted_prices[0].oracle = "nibi1pzd5e402eld9kcc3h78tmfrm5rpzlzk6hnxkvu"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.posted_prices[0].price = "20000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.pricefeed.posted_prices[0].expiry = "2123-08-08T14:48:13.591905Z"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
# x/stablecoin
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.coll_ratio = "1000000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.fee_ratio = "2000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.ef_fee_ratio = "500000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.bonus_rate_recoll = "2000"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.distr_epoch_identifier = "15 min"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.adjustment_step = "2500"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.price_upper_bound = "999900"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.price_upper_bound = "1000100"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.stablecoin.params.is_collateral_ratio_valid = false' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
# x/dex
cat $HOME/.nibid/config/genesis.json | jq '.app_state.dex.params.starting_pool_number = "1"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.dex.params.pool_creation_fee[0] = {denom:"unibi",amount:"1000000000"}' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json
cat $HOME/.nibid/config/genesis.json | jq '.app_state.dex.params.whitelisted_asset = ["uusdc","unibi","unusd"]' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json

# gov module
# Update voting period params
cat $HOME/.nibid/config/genesis.json | jq '.app_state.gov.voting_params.voting_period = "10s"' >$HOME/.nibid/config/tmp_genesis.json && mv $HOME/.nibid/config/tmp_genesis.json $HOME/.nibid/config/genesis.json

# Start the network
echo_info "Starting $CHAIN_ID in $HOME/.nibid..."
$BINARY start
