package cli

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	gravitytypes "github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"github.com/stretchr/testify/require"
)

func TestConfirmExternalMaterialUsesTronDomain(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	checkpoint := []byte("01234567890123456789012345678901")

	externalAddress, signature, err := confirmExternalMaterial(trontypes.ModuleName, privateKey, checkpoint)
	require.NoError(t, err)

	require.Equal(t, tronaddress.PubkeyToAddress(privateKey.PublicKey).String(), externalAddress)
	require.NoError(t, trontypes.ValidateTronSignature(checkpoint, append([]byte(nil), signature...), externalAddress))
}

func TestConfirmExternalMaterialUsesEthereumDomainForEVM(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	checkpoint := []byte("01234567890123456789012345678901")

	externalAddress, signature, err := confirmExternalMaterial("bsc", privateKey, checkpoint)
	require.NoError(t, err)

	require.Equal(t, crypto.PubkeyToAddress(privateKey.PublicKey).String(), externalAddress)
	require.NoError(t, gravitytypes.ValidateEthereumSignature(checkpoint, append([]byte(nil), signature...), externalAddress))
}
