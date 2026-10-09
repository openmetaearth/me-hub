package app

import (
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"
	cryptocodec "github.com/evmos/ethermint/crypto/codec"
	"github.com/evmos/ethermint/ethereum/eip712"
	ethermint "github.com/evmos/ethermint/types"

	"github.com/openmetaearth/me-hub/app/params"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
)

// makeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func makeEncodingConfig() params.EncodingConfig {
	amino := codec.NewLegacyAmino()
	addrCodec := address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	signingOpts := signing.Options{
		AddressCodec:          addrCodec,
		ValidatorAddressCodec: address.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
	}
	registerWstakingCustomGetSigners(&signingOpts, addrCodec)

	interfaceRegistry, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles:     proto.HybridResolver,
		SigningOptions: signingOpts,
	})
	if err != nil {
		panic(err)
	}
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(cdc, tx.DefaultSignModes)

	return params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             cdc,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := makeEncodingConfig()

	RegisterLegacyAminoCodec(encodingConfig.Amino)
	// In SDK 0.50, sdk.Msg is an alias of proto.Message. Registering that
	// interface on amino makes every proto concrete an implementer and panics
	// on 4-byte name collisions (e.g. feegrant.MsgGrantAllowance). Register
	// module amino types without the Msg interface.
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	RegisterInterfaces(encodingConfig.InterfaceRegistry)

	gravitytypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// EIP-712 still uses the deprecated legacytx.StdSignBytes helper, which
	// panics unless this codec is set. Keep it in sync with the app amino.
	legacytx.RegressionTestingAminoCodec = encodingConfig.Amino
	eip712.SetEncodingConfig(encodingConfig.Amino, encodingConfig.InterfaceRegistry)

	return encodingConfig
}

// RegisterLegacyAminoCodec registers Interfaces from types, crypto, and SDK std.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cryptocodec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
}

// RegisterInterfaces registers Interfaces from types, crypto, and SDK std.
func RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	std.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	ethermint.RegisterInterfaces(interfaceRegistry)
}
