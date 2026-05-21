package keeper

import (
	"context"
	"encoding/hex"
	wnfttypes "github.com/openmetaearth/me-hub/x/wnft/types"
	"slices"
	"strings"

	"cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"
	didtypes "github.com/openmetaearth/me-hub/x/did/types"
	"github.com/openmetaearth/me-hub/x/kyc/types"
	stktypes "github.com/openmetaearth/me-hub/x/wstaking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) Approve(goCtx context.Context, msg *types.MsgApprove) (*types.MsgApproveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgApproveResponse{}, didtypes.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgApproveResponse{}, sdkerrors.Wrap(didtypes.ErrInvalidIssuer, msg.Issuer)
	}

	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgApproveResponse{}, didtypes.ErrIssuerNotActive
	}

	// check holder did
	holderInfo, found := m.GetDidInfo(ctx, msg.Did)
	if found && holderInfo.Status == didtypes.DID_STATUS_ACTIVE {
		return &types.MsgApproveResponse{}, didtypes.ErrDidExists
	}

	// check region
	if _, found := m.stkKeeper.GetRegion(ctx, msg.RegionId); !found {
		return &types.MsgApproveResponse{}, stktypes.ErrRegionNotExist
	}

	// check holder address and pubkey
	address := sdk.MustAccAddressFromBech32(msg.Address)
	did, found := m.GetDID(ctx, address)
	if found && did != msg.Did {
		// notice: holder must have not DID
		return &types.MsgApproveResponse{}, didtypes.ErrHolderExists
	}

	// create DID
	m.SetDID(ctx, address, msg.Did)
	m.SetDidInfo(ctx, msg.Did, didtypes.DidInfo{
		Did:      msg.Did,
		Address:  msg.Address,
		Pubkey:   msg.Pubkey,
		RegionId: msg.RegionId,
		KycLevel: msg.Level,
		Status:   didtypes.DID_STATUS_ACTIVE,
	})

	// create KYC
	kyc := msg.GetKYC()
	m.SetKYC(ctx, msg.Did, kyc)

	// add region filter to KYC
	m.AddFilters(ctx, msg.Did, [][]byte{[]byte(msg.RegionId)}, kyc)

	// add reward to KYC holder and inviter
	if err := m.stkKeeper.KycReward(ctx, address, msg.RegionId, issuer); err != nil {
		return &types.MsgApproveResponse{}, errors.Wrap(err, "set reward failed")
	}

	if msg.Level >= didtypes.KYC_LEVEL_TWO {
		if err := m.stkKeeper.SendInviteReward(ctx, msg.Inviter, msg.Address, msg.RegionId); err != nil {
			return &types.MsgApproveResponse{}, sdkerrors.Wrap(types.ErrInviteReward, err.Error())
		}
	}

	// add account if not exists
	if !m.accountKeeper.HasAccount(ctx, address) {
		m.accountKeeper.SetAccount(ctx, m.accountKeeper.NewAccountWithAddress(ctx, address))
	}

	return &types.MsgApproveResponse{}, nil
}

func (m msgServer) Update(goCtx context.Context, msg *types.MsgUpdate) (*types.MsgUpdateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	m.Logger(ctx).Debug("call Update", "msg", msg)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgUpdateResponse{}, didtypes.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgUpdateResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgUpdateResponse{}, sdkerrors.Wrap(didtypes.ErrInvalidIssuer, msg.Issuer)
	}

	// check holder did
	holderInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found || holderInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgUpdateResponse{}, didtypes.ErrHolderNotActive
	}
	preKyc, found := m.GetKYC(ctx, msg.Did)
	if !found {
		return &types.MsgUpdateResponse{}, didtypes.ErrCredentialNotFound
	}
	perRegionId := string(preKyc.Data)
	perLevel := holderInfo.KycLevel

	if strings.EqualFold(perRegionId, stktypes.ExperienceRegionName) || strings.EqualFold(msg.RegionId, stktypes.ExperienceRegionName) {
		return nil, types.ErrTransferRegion.Wrap("from region or to region is experience region")
	}

	// check region
	if _, found := m.stkKeeper.GetRegion(ctx, msg.RegionId); !found {
		return &types.MsgUpdateResponse{}, stktypes.ErrRegionNotExist
	}

	// update KYC level
	//if msg.Level == didtypes.KYC_LEVEL_NONE {
	//	return &types.MsgUpdateResponse{}, errors.Wrap(didtypes.ErrParameter, "KYC level must be greater than 0")
	//}

	holderInfo.RegionId = msg.RegionId
	holderInfo.KycLevel = msg.Level
	m.SetDidInfo(ctx, msg.Did, holderInfo)
	address := sdk.MustAccAddressFromBech32(holderInfo.Address)
	// update KYC
	kyc := msg.GetKYC()
	m.SetKYC(ctx, msg.Did, kyc)

	// update KYC region filter
	m.DeleteFilters(ctx, msg.Did, [][]byte{[]byte(perRegionId)})
	m.AddFilters(ctx, msg.Did, [][]byte{[]byte(msg.RegionId)}, kyc)

	// change reward
	if err := m.TransferKycRegion(ctx, address.String(), msg.Issuer, perRegionId, msg.RegionId); err != nil {
		return &types.MsgUpdateResponse{}, sdkerrors.Wrap(types.ErrTransferRegion, err.Error())
	}

	if perLevel == didtypes.KYC_LEVEL_ONE && msg.Level >= didtypes.KYC_LEVEL_TWO {
		if err := m.stkKeeper.SendInviteReward(ctx, msg.Inviter, address.String(), msg.RegionId); err != nil {
			return &types.MsgUpdateResponse{}, sdkerrors.Wrap(types.ErrInviteReward, err.Error())
		}
	}

	// add event
	event := sdk.NewEvent(types.EventTypeUpdate,
		sdk.NewAttribute(types.AttributeKeyAddress, address.String()),
		sdk.NewAttribute(types.AttributeKeyRegionId, perRegionId),
		sdk.NewAttribute(types.AttributeKeyRegionIdChanged, msg.RegionId),
		sdk.NewAttribute(types.AttributeKeyLevel, perLevel.String()),
		sdk.NewAttribute(types.AttributeKeyLevelChanged, msg.Level.String()),
		sdk.NewAttribute(types.AttributeKeyInviter, msg.Inviter),
	)
	ctx.EventManager().EmitEvent(event)

	// event post-handler
	err := m.handlerReg.HandleEvent(ctx, types.EventTypeUpdate, event)
	if err != nil {
		return &types.MsgUpdateResponse{}, err
	}

	return &types.MsgUpdateResponse{}, nil
}

func (m msgServer) Remove(goCtx context.Context, msg *types.MsgRemove) (*types.MsgRemoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	m.Logger(ctx).Debug("call Remove", "msg", msg)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgRemoveResponse{}, didtypes.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgRemoveResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgRemoveResponse{}, didtypes.ErrIssuerNotActive
	}

	// disable did
	didInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found {
		return &types.MsgRemoveResponse{}, didtypes.ErrHolderNotFound
	}
	didInfo.RegionId = ""
	didInfo.KycLevel = 0
	didInfo.Status = didtypes.DID_STATUS_INACTIVE
	m.SetDidInfo(ctx, msg.Did, didInfo)

	// delete KYC
	kyc, found := m.GetKYC(ctx, msg.Did)
	if !found {
		return &types.MsgRemoveResponse{}, didtypes.ErrCredentialNotFound
	}
	m.DeleteKYC(ctx, msg.Did)

	// delete KYC region filter
	filters, _ := m.GetFilters(ctx, msg.Did)
	m.DeleteFilters(ctx, msg.Did, filters)

	// cancel reward
	address := sdk.MustAccAddressFromBech32(didInfo.Address)
	if err := m.DeleteApproveReward(ctx, address.String(), string(kyc.Data)); err != nil {
		return &types.MsgRemoveResponse{}, errors.Wrap(err, "delete reward failed")
	}

	return &types.MsgRemoveResponse{}, nil
}

func (m msgServer) CreateSBT(goCtx context.Context, msg *types.MsgCreateSBT) (*types.MsgCreateSBTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	m.Logger(ctx).Debug("call CreateSBT", "msg", msg)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrIssuerNotActive
	}

	// check holder did
	holderInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found || holderInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrHolderNotFound
	}
	if !m.HasKYC(ctx, msg.Did) {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrCredentialNotFound //
	}

	// mint SBT to KYC module address
	nftData, err := codectypes.NewAnyWithValue(&wnfttypes.Extension{Data: hex.EncodeToString(msg.Data)})
	if err != nil {
		return &types.MsgCreateSBTResponse{}, err
	}
	sbt := nft.NFT{
		ClassId: types.ModuleName,
		Id:      msg.Did,
		Uri:     msg.Uri,
		UriHash: msg.UriHash,
		Data:    nftData, // todo: check for encode
	}

	if err := m.SetSBT(ctx, sbt, sdk.MustAccAddressFromBech32(holderInfo.Address)); err != nil {
		return &types.MsgCreateSBTResponse{}, errors.Wrap(err, "mint SBT failed")
	}

	ctx.EventManager().EmitEvent(types.NewSbtEvent(types.EventTypeCreateSBT, msg.Did, msg.Uri, msg.UriHash, holderInfo.RegionId, holderInfo.KycLevel.String(), holderInfo.Address))
	return &types.MsgCreateSBTResponse{}, nil
}

func (m msgServer) UpdateSBT(goCtx context.Context, msg *types.MsgUpdateSBT) (*types.MsgUpdateSBTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	m.Logger(ctx).Debug("call UpdateSBT", "msg", msg)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgUpdateSBTResponse{}, didtypes.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgUpdateSBTResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgUpdateSBTResponse{}, didtypes.ErrIssuerNotActive
	}

	// check holder did
	holderInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found {
		return &types.MsgUpdateSBTResponse{}, didtypes.ErrHolderNotFound
	}
	if !m.HasKYC(ctx, msg.Did) {
		return &types.MsgUpdateSBTResponse{}, didtypes.ErrCredentialNotFound
	}

	// check SBT is existed
	sbt, found := m.GetSBT(ctx, msg.Did)
	if !found {
		return &types.MsgUpdateSBTResponse{}, types.ErrSbtNotFound
	}

	// update SBT
	nftData, err := codectypes.NewAnyWithValue(&wnfttypes.Extension{Data: hex.EncodeToString(msg.Data)})
	if err != nil {
		return nil, err
	}
	sbt.Uri = msg.Uri
	sbt.UriHash = msg.UriHash
	sbt.Data = nftData

	if err := m.nftKeeper.Update(ctx, sbt); err != nil {
		return &types.MsgUpdateSBTResponse{}, errors.Wrap(err, "update SBT failed")
	}

	ctx.EventManager().EmitEvent(types.NewSbtEvent(types.EventTypeUpdateSBT, msg.Did, msg.Uri, msg.UriHash, holderInfo.RegionId, holderInfo.KycLevel.String(), holderInfo.Address))
	return &types.MsgUpdateSBTResponse{}, nil
}

func (m msgServer) DeleteSBT(goCtx context.Context, msg *types.MsgDeleteSBT) (*types.MsgDeleteSBTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrIssuerNotActive
	}

	// check holder did
	holderInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrHolderNotFound
	}

	// remove SBT
	if err := m.RemoveSBT(ctx, msg.Did); err != nil {
		return &types.MsgDeleteSBTResponse{}, errors.Wrap(err, "burn SBT failed")
	}

	ctx.EventManager().EmitEvent(types.NewSbtEvent(types.EventTypeDeleteSBT, msg.Did, "", "", holderInfo.RegionId, holderInfo.KycLevel.String(), holderInfo.Address))
	return &types.MsgDeleteSBTResponse{}, nil
}
