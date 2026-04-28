package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/gravity/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
)

// claimHandlerCommon is an internal function that provides common code for processing claims once they are
// translated from the message to the Ethereum claim interface
func (s MsgServer) claimHandlerCommon(ctx sdk.Context, msg types.ExternalClaim) (err error) {
	bridgerAddr := msg.GetClaimer()
	if err := s.checkIsRelayer(ctx, bridgerAddr); err != nil {
		return err
	}

	// Add the claim to the store
	if _, err := s.Attest(ctx, bridgerAddr, msg); err != nil {
		return err
	}
	return nil
}

func (s MsgServer) confirmHandlerCommon(ctx sdk.Context, relayerAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) error {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, "signature decoding failed")
	}

	relayer, found := s.GetRelayer(ctx, relayerAddr)
	if !found {
		return types.ErrNotFoundRelayer
	}

	if relayer.ExternalAddress != signatureAddr {
		return errorsmod.Wrapf(types.ErrExternalAddressNotMatch, "got %s, expected %s", signatureAddr, relayer.ExternalAddress)
	}

	if s.moduleName == trontypes.ModuleName {
		if err = trontypes.ValidateTronSignature(checkpoint, sigBytes, relayer.ExternalAddress); err != nil {
			return errorsmod.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s and sig %s", relayer.ExternalAddress, hex.EncodeToString(checkpoint), signature))
		}
	} else {
		if err = types.ValidateEthereumSignature(checkpoint, sigBytes, relayer.ExternalAddress); err != nil {
			return errorsmod.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s and sig %s", relayer.ExternalAddress, hex.EncodeToString(checkpoint), signature))
		}
	}
	return nil
}
