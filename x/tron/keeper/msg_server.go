package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	gravitykeeper "github.com/st-chain/me-hub/x/gravity/keeper"
	gravitytypes "github.com/st-chain/me-hub/x/gravity/types"
	trontypes "github.com/st-chain/me-hub/x/tron/types"
)

var _ gravitytypes.MsgServer = msgServer{}

type msgServer struct {
	gravitykeeper.MsgServer
}

func NewMsgServerImpl(keeper Keeper) gravitytypes.MsgServer {
	return &msgServer{gravitykeeper.MsgServer{Keeper: keeper.Keeper}}
}

func (s msgServer) ConfirmBatch(c context.Context, msg *gravitytypes.MsgConfirmBatch) (*gravitytypes.MsgConfirmBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := s.GetOutgoingTxBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, errorsmod.Wrap(gravitytypes.ErrInvalid, "couldn't find batch")
	}

	checkpoint, err := trontypes.GetCheckpointConfirmBatch(batch, s.GetGravityID(ctx))
	if err != nil {
		return nil, errorsmod.Wrap(gravitytypes.ErrInvalid, "checkpoint generation")
	}

	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	err = s.confirmHandlerCommon(ctx, relayerAddress, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetBatchConfirm(ctx, msg.TokenContract, msg.Nonce, relayerAddress) != nil {
		return nil, errorsmod.Wrap(gravitytypes.ErrDuplicate, "signature")
	}
	s.SetBatchConfirm(ctx, relayerAddress, msg)
	return &gravitytypes.MsgConfirmBatchResponse{}, nil
}

// RelayerSetConfirm handles MsgRelayerSetConfirm
func (s msgServer) RelayerSetConfirm(c context.Context, msg *gravitytypes.MsgRelayerSetConfirm) (*gravitytypes.MsgRelayerSetConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	relayerSet := s.GetRelayerSet(ctx, msg.Nonce)
	if relayerSet == nil {
		return nil, errorsmod.Wrap(gravitytypes.ErrInvalid, "couldn't find relayerSet")
	}

	checkpoint, err := trontypes.GetCheckpointRelayerSet(relayerSet, s.GetGravityID(ctx))
	if err != nil {
		return nil, err
	}

	relayerAddress := sdk.MustAccAddressFromBech32(msg.RelayerAddress)
	err = s.confirmHandlerCommon(ctx, relayerAddress, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}

	// check if we already have this confirm
	if s.GetRelayerSetConfirm(ctx, msg.Nonce, relayerAddress) != nil {
		return nil, errorsmod.Wrap(gravitytypes.ErrDuplicate, "signature")
	}
	s.SetRelayerSetConfirm(ctx, relayerAddress, msg)
	return &gravitytypes.MsgRelayerSetConfirmResponse{}, nil
}

func (s msgServer) confirmHandlerCommon(ctx sdk.Context, relayerAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) error {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, "signature decoding failed")
	}

	relayer, found := s.GetRelayer(ctx, relayerAddr)
	if !found {
		return gravitytypes.ErrNotFoundRelayer
	}

	if relayer.ExternalAddress != signatureAddr {
		return errorsmod.Wrapf(gravitytypes.ErrExternalAddressNotMatch, "got %s, expected %s", signatureAddr, relayer.ExternalAddress)
	}

	if err = trontypes.ValidateTronSignature(checkpoint, sigBytes, relayer.ExternalAddress); err != nil {
		return errorsmod.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", relayer.ExternalAddress, hex.EncodeToString(checkpoint), signature))
	}
	return nil
}
