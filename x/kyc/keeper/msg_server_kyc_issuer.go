package keeper

import (
	"context"
	"encoding/json"
	"slices"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
)

func (m msgServer) UpdateKycIssuer(goCtx context.Context, msg *types.MsgUpdateKycIssuer) (*types.MsgUpdateKycIssuerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !m.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return nil, types.ErrUnauthorized
	}
	service, found := m.GetService(ctx)
	if !found {
		return nil, didtypes.ErrServiceNotFound
	}

	// Work on a copy and persist only after every address has been checked.
	issuers := slices.Clone(service.Issuers)
	for _, address := range msg.Addresses {
		did, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(address))
		if !found {
			return nil, errorsmod.Wrap(didtypes.ErrDidNotFound, address)
		}
		info, found := m.GetDidInfo(ctx, did)
		if !found {
			return nil, errorsmod.Wrap(didtypes.ErrHolderNotFound, address)
		}
		if info.Address != address || info.Did != did {
			return nil, errorsmod.Wrap(didtypes.ErrInvalidIssuer, address)
		}

		if msg.OptionType == types.KycIssuerRemove {
			issuers = slices.DeleteFunc(issuers, func(issuer string) bool { return issuer == did })
			continue
		}

		if info.Status != didtypes.DID_STATUS_ACTIVE {
			return nil, errorsmod.Wrap(didtypes.ErrDidNotActive, address)
		}

		if !slices.Contains(issuers, did) {
			issuers = append(issuers, did)
		}
	}
	service.Issuers = issuers
	m.SetService(ctx, service)

	addresses, err := json.Marshal(msg.Addresses)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeUpdateKycIssuer,
		sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		sdk.NewAttribute(types.AttributeKeyAddresses, string(addresses)),
		sdk.NewAttribute(types.AttributeKeyOptionType, strconv.FormatInt(int64(msg.OptionType), 10)),
	))
	return &types.MsgUpdateKycIssuerResponse{}, nil
}
