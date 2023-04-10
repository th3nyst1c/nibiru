package app

import (
	"github.com/cosmos/cosmos-sdk/std"

	"github.com/NibiruChain/nibiru/app/params"
)

// MakeEncodingConfig creates a new EncodingConfig with all modules registered
// This function should be used only in tests or when creating a new app
// instance (NewApp*()). The App user shouldn't create new codecs but should
// instead use the app.AppCodec instead.
// [DEPRECATED]
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
