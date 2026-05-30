package types_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/testutil/helpers"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func TestEthereumExternalZeroAddressRejected(t *testing.T) {
	zeroAddress := common.Address{}.Hex()
	relayerAddress := helpers.GenAccAddress().String()

	require.Error(t, gravitytypes.ValidateExternalAddr(bsctypes.ModuleName, zeroAddress))

	bondedRelayer := &gravitytypes.MsgBondedRelayer{
		RelayerAddress:  relayerAddress,
		ExternalAddress: zeroAddress,
		DelegateAmount:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(1)),
		ChainName:       bsctypes.ModuleName,
	}
	require.Error(t, bondedRelayer.ValidateBasic())

	relayerSetConfirm := &gravitytypes.MsgRelayerSetConfirm{
		Nonce:           1,
		RelayerAddress:  relayerAddress,
		ExternalAddress: zeroAddress,
		Signature:       hex.EncodeToString([]byte{1}),
		ChainName:       bsctypes.ModuleName,
	}
	require.Error(t, relayerSetConfirm.ValidateBasic())

	bridgeValidator := &gravitytypes.BridgeValidator{
		Power:           gravitytypes.PowerBase,
		ExternalAddress: zeroAddress,
	}
	require.Error(t, bridgeValidator.ValidateBasic())
}
