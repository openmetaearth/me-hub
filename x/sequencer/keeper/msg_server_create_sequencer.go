package keeper

import (
	"context"
	"slices"
	"strconv"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/sequencer/types"

	errorsmod "cosmossdk.io/errors"
)

// CreateSequencer defines a method for creating a new sequencer
func (k msgServer) CreateSequencer(goCtx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.DymintPubKey == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidPubKey, "sequencer pubkey can not be empty")
	}

	// check to see if the sequencer has been registered before
	if _, found := k.GetSequencer(ctx, msg.Creator); found {
		return nil, types.ErrSequencerExists
	}

	// check to see if the rollapp has been registered before
	rollapp, found := k.rollappKeeper.GetRollapp(ctx, msg.RollappId)
	if !found {
		return nil, types.ErrUnknownRollappID
	}
	if rollapp.Frozen {
		return nil, types.ErrRollappJailed
	}

	// check if there are permissionedAddresses.
	// if the list is not empty, it means that only permissioned sequencers can be added
	permissionedAddresses := rollapp.PermissionedAddresses
	if 0 < len(permissionedAddresses) && !slices.Contains(permissionedAddresses, msg.Creator) {
		return nil, types.ErrSequencerNotPermissioned
	}

	// check to see if the sequencer has enough balance and deduct the bond
	seqAcc, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	bondedSequencers := k.GetSequencersByRollappByStatus(ctx, msg.RollappId, types.Bonded)
	unbondingSequencers := k.GetSequencersByRollappByStatus(ctx, msg.RollappId, types.Unbonding)
	if err := k.validateUniqueDymintPubKey(msg.DymintPubKey, bondedSequencers, unbondingSequencers); err != nil {
		return nil, err
	}

	bond := sdk.Coins{}
	minBond := k.GetParams(ctx).MinBond
	if !minBond.IsNil() && !minBond.IsZero() {
		if msg.Bond.Denom != minBond.Denom {
			return nil, errorsmod.Wrapf(
				types.ErrInvalidCoinDenom, "got %s, expected %s", msg.Bond.Denom, minBond.Denom,
			)
		}

		if msg.Bond.Amount.LT(minBond.Amount) {
			return nil, errorsmod.Wrapf(
				types.ErrInsufficientBond, "got %s, expected %s", msg.Bond.Amount, minBond,
			)
		}

		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, seqAcc, types.ModuleName, sdk.NewCoins(msg.Bond))
		if err != nil {
			return nil, err
		}
		bond = sdk.NewCoins(msg.Bond)
	}

	sequencer := types.Sequencer{
		SequencerAddress: msg.Creator,
		DymintPubKey:     msg.DymintPubKey,
		RollappId:        msg.RollappId,
		Description:      msg.Description,
		Status:           types.Bonded,
		Proposer:         false,
		Tokens:           bond,
	}

	// check to see if we reached the maximum number of sequencers for this rollapp
	currentNumOfSequencers := len(bondedSequencers) + len(unbondingSequencers)
	if rollapp.MaxSequencers > 0 && uint64(currentNumOfSequencers) >= rollapp.MaxSequencers {
		return nil, types.ErrMaxSequencersLimit
	}
	// this is the first sequencer, make it a PROPOSER
	if len(bondedSequencers) == 0 {
		sequencer.Proposer = true
	}

	k.SetSequencer(ctx, sequencer)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateSequencer,
			sdk.NewAttribute(types.AttributeKeySequencer, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyProposer, strconv.FormatBool(sequencer.Proposer)),
		),
	)

	return &types.MsgCreateSequencerResponse{}, nil
}

func (k msgServer) validateUniqueDymintPubKey(dymintPubKey *codectypes.Any, sequencerGroups ...[]types.Sequencer) error {
	var pubKey cryptotypes.PubKey
	if err := k.cdc.UnpackAny(dymintPubKey, &pubKey); err != nil {
		return errorsmod.Wrapf(types.ErrInvalidPubKey, "invalid sequencer pubkey(%s)", err)
	}
	if pubKey == nil {
		return errorsmod.Wrapf(types.ErrInvalidPubKey, "sequencer pubkey can not be empty")
	}

	for _, sequencers := range sequencerGroups {
		for _, sequencer := range sequencers {
			if sequencer.DymintPubKey == nil {
				continue
			}

			var existingPubKey cryptotypes.PubKey
			if err := k.cdc.UnpackAny(sequencer.DymintPubKey, &existingPubKey); err != nil {
				return errorsmod.Wrapf(
					types.ErrInvalidPubKey,
					"registered sequencer %s has invalid pubkey(%s)",
					sequencer.SequencerAddress,
					err,
				)
			}

			if pubKey.Equals(existingPubKey) {
				return errorsmod.Wrapf(
					types.ErrInvalidPubKey,
					"dymint pubkey already registered by sequencer %s",
					sequencer.SequencerAddress,
				)
			}
		}
	}

	return nil
}
