package keeper

import (
	"context"
	"cosmossdk.io/errors"
	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	"github.com/st-chain/me-hub/x/kyc/types"
	stktypes "github.com/st-chain/me-hub/x/wstaking/types"
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
		return &types.MsgApproveResponse{}, didtypes.ErrServiceNotFound
	}

	// check issuer
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || issuer != svc.Issuer {
		return &types.MsgApproveResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerDoc, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerDoc.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgApproveResponse{}, didtypes.ErrIssuerNotActive
	}

	// check DID
	_, found = m.GetDidInfo(ctx, msg.Did)
	if found {
		return &types.MsgApproveResponse{}, didtypes.ErrDidExists
	}

	// check region
	if _, found := m.stkKeeper.GetRegion(ctx, msg.RegionId); !found {
		return &types.MsgApproveResponse{}, stktypes.ErrRegionNotExist
	}

	// check holder address and pubkey
	address := sdk.MustAccAddressFromBech32(msg.Address)
	pubkey, err := m.PubKeyFromString(msg.Pubkey)
	if err != nil {
		return &types.MsgApproveResponse{}, err
	}
	if !address.Equals(sdk.AccAddress(pubkey.Address())) {
		return &types.MsgApproveResponse{}, didtypes.ErrHolderNotFound
	}
	if m.HasDID(ctx, address) {
		// notice: holder must have not DID
		return &types.MsgApproveResponse{}, didtypes.ErrHolderExists
	}

	// create DID
	m.SetDID(ctx, address, msg.Did)
	m.SetDidInfo(ctx, msg.Did, didtypes.DidInfo{
		Did:    msg.Did,
		Pubkey: msg.Pubkey,
		Status: didtypes.DID_STATUS_ACTIVE,
	})

	// create KYC
	kyc := didtypes.NewCredential(msg.Did, types.ModuleName, msg.Hash, msg.Uri)
	m.SetKYC(ctx, msg.Did, kyc)

	// add region filter to KYC
	m.AddFilters(ctx, msg.Did, [][]byte{[]byte(msg.RegionId)}, kyc)

	// add reward to KYC holder and inviter
	if err := m.SetApproveReward(ctx, msg.Address, msg.Inviter, msg.Issuer, msg.RegionId); err != nil {
		return &types.MsgApproveResponse{}, err
	}

	return &types.MsgApproveResponse{}, nil
}

func (m msgServer) Update(goCtx context.Context, msg *types.MsgUpdate) (*types.MsgUpdateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgUpdateResponse{}, didtypes.ErrServiceNotFound
	}

	// check issuer
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || issuer != svc.Issuer {
		return &types.MsgUpdateResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerDoc, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerDoc.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgUpdateResponse{}, didtypes.ErrIssuerNotActive
	}

	// check DID
	didInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found {
		return &types.MsgUpdateResponse{}, didtypes.ErrHolderNotFound
	}

	// check region
	if _, found := m.stkKeeper.GetRegion(ctx, msg.RegionId); !found {
		return &types.MsgUpdateResponse{}, stktypes.ErrRegionNotExist
	}

	// update KYC
	kyc := didtypes.NewCredential(msg.Did, types.ModuleName, msg.Hash, msg.Uri)
	m.SetKYC(ctx, msg.Did, kyc)

	// update KYC region filter
	m.AddFilters(ctx, msg.Did, [][]byte{[]byte(msg.RegionId)}, kyc)

	// change reward
	address := m.MustAccAddressFromPubkeyString(didInfo.Pubkey).String()
	if err := m.DeleteApproveReward(ctx, address); err != nil {
		return &types.MsgUpdateResponse{}, err
	}
	if err := m.SetApproveReward(ctx, address, "", msg.Issuer, msg.RegionId); err != nil {
		return &types.MsgUpdateResponse{}, err
	}

	return &types.MsgUpdateResponse{}, nil
}

func (m msgServer) Remove(goCtx context.Context, msg *types.MsgRemove) (*types.MsgRemoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgRemoveResponse{}, didtypes.ErrServiceNotFound
	}

	// check issuer
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || issuer != svc.Issuer {
		return &types.MsgRemoveResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerDoc, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerDoc.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgRemoveResponse{}, didtypes.ErrIssuerNotActive
	}

	// disable did
	didInfo, found := m.GetDidInfo(ctx, msg.Did)
	if !found {
		return &types.MsgRemoveResponse{}, didtypes.ErrHolderNotFound
	}
	didInfo.Status = didtypes.DID_STATUS_DEACTIVE
	m.SetDidInfo(ctx, msg.Did, didInfo)

	// delete KYC
	if !m.HasKYC(ctx, msg.Did) {
		return &types.MsgRemoveResponse{}, didtypes.ErrCredentialNotFound
	}
	m.DeleteKYC(ctx, msg.Did)

	// delete KYC region filter
	filters, _ := m.GetFilters(ctx, msg.Did)
	m.DeleteFilters(ctx, msg.Did, filters)

	// cancel reward
	address := m.MustAccAddressFromPubkeyString(didInfo.Pubkey).String()
	if err := m.DeleteApproveReward(ctx, address); err != nil {
		return &types.MsgRemoveResponse{}, err
	}

	return &types.MsgRemoveResponse{}, nil
}

func (m msgServer) CreateSBT(goCtx context.Context, msg *types.MsgCreateSBT) (*types.MsgCreateSBTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrServiceNotFound
	}

	// check issuer
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || issuer != svc.Issuer {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerDoc, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerDoc.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrIssuerNotActive
	}

	// check holder
	if !m.HasDidInfo(ctx, msg.Did) {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrHolderNotFound
	}
	if !m.HasKYC(ctx, msg.Did) {
		return &types.MsgCreateSBTResponse{}, didtypes.ErrCredentialNotFound //
	}

	// mint SBT
	sbt := nft.NFT{
		ClassId: types.ModuleName,
		Id:      msg.Did,
		Uri:     msg.Uri,
		UriHash: msg.UriHash,
		Data:    types2.UnsafePackAny(msg.Data), // todo: check for encode
	}
	if err := m.SetSBT(ctx, sbt, sdk.MustAccAddressFromBech32(msg.Issuer)); err != nil {
		return &types.MsgCreateSBTResponse{}, err
	}

	return &types.MsgCreateSBTResponse{}, nil
}

func (m msgServer) DeleteSBT(goCtx context.Context, msg *types.MsgDeleteSBT) (*types.MsgDeleteSBTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx)
	if !found || svc.Status != didtypes.SERVICE_STATUS_ACTIVE {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrServiceNotFound
	}

	// check issuer
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || issuer != svc.Issuer {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrIssuerNotFound
	}
	issuerDoc, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerDoc.Status != didtypes.DID_STATUS_ACTIVE {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrIssuerNotActive
	}

	// check holder
	if !m.HasDidInfo(ctx, msg.Did) {
		return &types.MsgDeleteSBTResponse{}, didtypes.ErrHolderNotFound
	}

	// remove SBT
	if err := m.RemoveSBT(ctx, msg.Did); err != nil {
		return &types.MsgDeleteSBTResponse{}, errors.Wrap(err, "remove SBT failed")
	}

	return &types.MsgDeleteSBTResponse{}, nil
}
