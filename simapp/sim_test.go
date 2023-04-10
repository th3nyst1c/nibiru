package simapp

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	sdksimapp "cosmossdk.io/simapp"
	"github.com/NibiruChain/nibiru/x/common/testutil/sim"
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/NibiruChain/nibiru/app"
	"github.com/NibiruChain/nibiru/x/common/testutil/testapp"
)

func init() {
	sdksimapp.GetSimulatorFlags()
}

func TestFullAppSimulation(tb *testing.T) {
	config, db, dir, _, skip, err := sdksimapp.SetupSimulation("goleveldb-app-sim", "Simulation")
	if skip {
		tb.Skip("skipping application simulation")
	}
	require.NoError(tb, err, "simulation setup failed")

	defer func() {
		db.Close()
		err = os.RemoveAll(dir)
		if err != nil {
			tb.Fatal(err)
		}
	}()

	encoding := app.MakeTestEncodingConfig()
	app := testapp.NewNibiruTestApp(app.NewDefaultGenesisState(encoding.Marshaler))

	// Run randomized simulation:
	_, simParams, simErr := simulation.SimulateFromSeed(
		/* tb */ tb,
		/* w */ os.Stdout,
		/* app */ app.BaseApp,
		/* appStateFn */ AppStateFn(app.AppCodec(), app.SimulationManager()),
		/* randAccFn */ simulationtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
		/* ops */ sdksimapp.SimulationOperations(app, app.AppCodec(), config), // Run all registered operations
		/* blockedAddrs */ app.ModuleAccountAddrs(),
		/* config */ config,
		/* cdc */ app.AppCodec(),
	)

	// export state and simParams before the simulation error is checked
	if err = sdksimapp.CheckExportSimulation(app, config, simParams); err != nil {
		tb.Fatal(err)
	}

	if simErr != nil {
		tb.Fatal(simErr)
	}

	if config.Commit {
		sdksimapp.PrintStats(db)
	}
}

func TestAppStateDeterminism(t *testing.T) {
	if !sdksimapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	encoding := app.MakeTestEncodingConfig()

	config := sdksimapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = sim.SimAppChainID

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			db := dbm.NewMemDB()
			app := testapp.NewNibiruTestApp(app.NewDefaultGenesisState(encoding.Marshaler))

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				app.BaseApp,
				AppStateFn(app.AppCodec(), app.SimulationManager()),
				simtypes.RandomAccounts, // Replace with own random account function if using keys other than secp256k1
				sdksimapp.SimulationOperations(app, app.AppCodec(), config),
				app.ModuleAccountAddrs(),
				config,
				app.AppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				sdksimapp.PrintStats(db)
			}

			appHash := app.LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, string(appHashList[0]), string(appHashList[j]),
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
