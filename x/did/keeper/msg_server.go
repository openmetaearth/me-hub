package keeper

import (
	"context"
	"slices"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/x/did/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {

	// API inactive
	return &types.MsgCreateDidResponse{}, errors.Wrap(types.ErrApiInactive, "use the Approve method of the KYC module to create a DID")
}

func (m msgServer) UpdateDidStatus(goCtx context.Context, msg *types.MsgUpdateDidStatus) (*types.MsgUpdateDidStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return &types.MsgUpdateDidStatusResponse{}, types.ErrPermissionDenial
	}

	info, found := m.GetDidInfo(ctx, msg.Did)
	if !found {
		return &types.MsgUpdateDidStatusResponse{}, types.ErrDidNotFound
	}
	if info.Status == msg.Status {
		return &types.MsgUpdateDidStatusResponse{}, types.ErrSameDidStatus
	}

	info.Status = msg.Status
	m.SetDidInfo(ctx, info.Did, info)

	ctx.EventManager().EmitEvent(types.NewDidEvent(types.EventTypeUpdateDidStatus, info.Did, info.Address, info.Status.String()))
	return &types.MsgUpdateDidStatusResponse{}, nil
}

//func (k msgServer) RemoveDid(goCtx context.Context, msg *types.MsgRemoveDid) (*types.MsgRemoveDidResponse, error) {
//	ctx := sdk.UnwrapSDKContext(goCtx)
//
//	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
//		return &types.MsgRemoveDidResponse{}, types.ErrPermissionDenial
//	}
//
//	if info, found := k.GetDidInfo(ctx, msg.Did); !found || info.Status != types.DID_STATUS_ACTIVE {
//		return &types.MsgRemoveDidResponse{}, types.ErrDidNotActive
//	}
//
//	k.DeleteDidDocument(ctx, msg.Did)
//
//	// API inactive
//	return &types.MsgRemoveDidResponse{}, nil
//}

func (m msgServer) CreateService(goCtx context.Context, msg *types.MsgCreateService) (*types.MsgCreateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return &types.MsgCreateServiceResponse{}, types.ErrPermissionDenial
	}

	// check service
	_, found := m.GetService(ctx, msg.Sid)
	if found {
		return &types.MsgCreateServiceResponse{}, types.ErrServiceExists
	}

	// check issuer
	for _, issuer := range msg.Issuers {
		if info, found := m.GetDidInfo(ctx, issuer); !found || info.Status != types.DID_STATUS_ACTIVE {
			return &types.MsgCreateServiceResponse{}, types.ErrIssuerNotActive
		}
	}

	svc := types.NewService(msg.Sid, msg.Name, msg.Description, types.SERVICE_STATUS_ACTIVE, msg.Issuers)
	m.SetService(ctx, msg.Sid, svc)

	ctx.EventManager().EmitEvent(types.NewServiceEvent(types.EventTypeCreateService, svc.Sid, svc.Name, svc.Status.String(), svc.Issuers))
	return &types.MsgCreateServiceResponse{}, nil
}

func (m msgServer) UpdateServiceStatus(goCtx context.Context, msg *types.MsgUpdateServiceStatus) (*types.MsgUpdateServiceStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return &types.MsgUpdateServiceStatusResponse{}, types.ErrPermissionDenial
	}

	svc, found := m.GetService(ctx, msg.Sid)
	if !found {
		return &types.MsgUpdateServiceStatusResponse{}, types.ErrServiceNotFound
	}
	if svc.Status == msg.Status {
		return &types.MsgUpdateServiceStatusResponse{}, types.ErrSameServiceStatus
	}

	svc.Status = msg.Status
	m.SetService(ctx, msg.Sid, svc)

	ctx.EventManager().EmitEvent(types.NewServiceEvent(types.EventTypeUpdateServiceStatus, svc.Sid, svc.Name, svc.Status.String(), svc.Issuers))
	return &types.MsgUpdateServiceStatusResponse{}, nil
}

//func (k msgServer) RemoveService(goCtx context.Context, msg *types.MsgRemoveService) (*types.MsgRemoveServiceResponse, error) {
//	ctx := sdk.UnwrapSDKContext(goCtx)
//
//	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
//		return &types.MsgRemoveServiceResponse{}, types.ErrPermissionDenial
//	}
//
//	if info, found := k.GetService(ctx, msg.Sid); !found || info.Status != types.SERVICE_STATUS_ACTIVE {
//		return &types.MsgRemoveServiceResponse{}, types.ErrServiceNotActive
//	}
//
//	k.DeleteService(ctx, msg.Sid)
//
//	return &types.MsgRemoveServiceResponse{}, nil
//}

func (m msgServer) CreateVC(goCtx context.Context, msg *types.MsgCreateVC) (*types.MsgCreateVCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx, msg.Sid)
	if !found || svc.Status != types.SERVICE_STATUS_ACTIVE {
		return &types.MsgCreateVCResponse{}, types.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgCreateVCResponse{}, types.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgCreateVCResponse{}, types.ErrIssuerNotActive
	}

	// check holder did
	if holderInfo, found := m.GetDidInfo(ctx, msg.Did); !found || holderInfo.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgCreateVCResponse{}, types.ErrHolderNotFound
	}

	// check vc
	if m.HasCredential(ctx, msg.Did, msg.Sid) {
		return &types.MsgCreateVCResponse{}, types.ErrCredentialExists
	}

	// create VC
	vc := msg.GetCredential()
	m.SetCredential(ctx, msg.Did, msg.Sid, vc)

	// add filters to VC
	m.AddFilters(ctx, msg.Did, msg.Sid, msg.Filters, vc)

	ctx.EventManager().EmitEvent(types.NewVcEvent(types.EventTypeCreateVC, vc.Sid, vc.Did, vc.Hash, vc.Uri))
	return &types.MsgCreateVCResponse{}, nil
}

func (m msgServer) UpdateVC(goCtx context.Context, msg *types.MsgUpdateVC) (*types.MsgUpdateVCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check vc
	if found := m.HasCredential(ctx, msg.Did, msg.Sid); !found {
		return &types.MsgUpdateVCResponse{}, types.ErrCredentialNotFound
	}

	// check credential service
	svc, found := m.GetService(ctx, msg.Sid)
	if !found || svc.Status != types.SERVICE_STATUS_ACTIVE {
		return &types.MsgUpdateVCResponse{}, types.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgUpdateVCResponse{}, types.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgUpdateVCResponse{}, types.ErrIssuerNotActive
	}

	// check holder
	if info, found := m.GetDidInfo(ctx, msg.Did); !found || info.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgUpdateVCResponse{}, types.ErrHolderNotActive
	}

	// update VC
	vc := msg.GetCredential()
	m.SetCredential(ctx, msg.Did, msg.Sid, vc)

	// update filters to VC
	m.AddFilters(ctx, msg.Did, msg.Sid, msg.Filters, vc)

	ctx.EventManager().EmitEvent(types.NewVcEvent(types.EventTypeUpdateVC, vc.Sid, vc.Did, vc.Hash, vc.Uri))
	return &types.MsgUpdateVCResponse{}, nil
}

func (m msgServer) RemoveVC(goCtx context.Context, msg *types.MsgRemoveVC) (*types.MsgRemoveVCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check credential service
	svc, found := m.GetService(ctx, msg.Sid)
	if !found || svc.Status != types.SERVICE_STATUS_ACTIVE {
		return &types.MsgRemoveVCResponse{}, types.ErrServiceNotActive
	}

	// check issuer did
	issuer, found := m.GetDID(ctx, sdk.MustAccAddressFromBech32(msg.Issuer))
	if !found || !slices.Contains(svc.Issuers, issuer) {
		return &types.MsgRemoveVCResponse{}, types.ErrIssuerNotFound
	}
	issuerInfo, found := m.GetDidInfo(ctx, issuer)
	if !found || issuerInfo.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgRemoveVCResponse{}, types.ErrIssuerNotActive
	}

	m.DeleteCredential(ctx, msg.Did, msg.Sid)

	filters, _ := m.GetFilters(ctx, msg.Did, msg.Sid)
	m.DeleteFilters(ctx, msg.Did, msg.Sid, filters)

	ctx.EventManager().EmitEvent(types.NewVcEvent(types.EventTypeRemoveVC, msg.Sid, msg.Did, "", ""))
	return &types.MsgRemoveVCResponse{}, nil
}
