package keeper

import (
	"context"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/st-chain/me-hub/x/did/types"
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

func (k msgServer) CreateDid(goCtx context.Context, msg *types.MsgCreateDid) (*types.MsgCreateDidResponse, error) {

	// API inactive
	return &types.MsgCreateDidResponse{}, errors.Wrap(types.ErrApiInactive, "use the Approve method of the KYC module to create a DID")
}

func (k msgServer) UpdateDidStatus(goCtx context.Context, msg *types.MsgUpdateDidStatus) (*types.MsgUpdateDidStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return &types.MsgUpdateDidStatusResponse{}, types.ErrPermissionDenial
	}

	info, found := k.GetDidInfo(ctx, msg.Did)
	if !found {
		return &types.MsgUpdateDidStatusResponse{}, types.ErrDidNotFound
	}
	if info.Status == msg.Status {
		return &types.MsgUpdateDidStatusResponse{}, types.ErrSameDidStatus
	}

	info.Status = msg.Status
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

func (k msgServer) CreateService(goCtx context.Context, msg *types.MsgCreateService) (*types.MsgCreateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return &types.MsgCreateServiceResponse{}, types.ErrPermissionDenial
	}

	// check service
	_, found := k.GetService(ctx, msg.Sid)
	if found {
		return &types.MsgCreateServiceResponse{}, types.ErrServiceExists
	}

	// check issuer
	if info, found := k.GetDidInfo(ctx, msg.Issuer); !found || info.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgCreateServiceResponse{}, types.ErrIssuerNotActive
	}

	service := types.NewService(msg.Sid, msg.Name, msg.Description, types.SERVICE_STATUS_DEACTIVE, msg.Issuer)
	k.SetService(ctx, msg.Sid, service)

	return &types.MsgCreateServiceResponse{}, nil
}

func (k msgServer) UpdateServiceStatus(goCtx context.Context, msg *types.MsgUpdateServiceStatus) (*types.MsgUpdateServiceStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.daoKeeper.IsGlobalDao(ctx, msg.Creator) {
		return &types.MsgUpdateServiceStatusResponse{}, types.ErrPermissionDenial
	}

	svc, found := k.GetService(ctx, msg.Sid)
	if !found {
		return &types.MsgUpdateServiceStatusResponse{}, types.ErrServiceNotFound
	}
	if svc.Status == msg.Status {
		return &types.MsgUpdateServiceStatusResponse{}, types.ErrSameServiceStatus
	}

	svc.Status = msg.Status
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

func (k msgServer) CreateVC(goCtx context.Context, msg *types.MsgCreateVC) (*types.MsgCreateVCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	svc, found := k.GetService(ctx, msg.Sid)
	if !found {
		return &types.MsgCreateVCResponse{}, types.ErrServiceNotFound
	}

	// check issuer
	issuer := sdk.MustAccAddressFromBech32(msg.Issuer)
	if did, found := k.GetDID(ctx, issuer); !found || did != svc.Issuer {
		return &types.MsgCreateVCResponse{}, types.ErrIssuerNotFound
	}

	// check service status
	if svc.Status != types.SERVICE_STATUS_ACTIVE {
		return &types.MsgCreateVCResponse{}, types.ErrServiceNotActive
	}

	// check holder
	if doc, found := k.GetDidDocument(ctx, msg.Did); !found || doc.Info.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgCreateVCResponse{}, types.ErrHolderNotFound
	}

	vc := types.NewCredential(msg.Did, msg.Sid, msg.Hash, msg.Uri)
	k.SetCredential(ctx, msg.Did, msg.Sid, vc)

	return &types.MsgCreateVCResponse{}, nil
}

func (k msgServer) UpdateVC(goCtx context.Context, msg *types.MsgUpdateVC) (*types.MsgUpdateVCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check vc
	vc, found := k.GetCredential(ctx, msg.Did, msg.Hash)
	if !found || msg.Sid != vc.Sid {
		return &types.MsgUpdateVCResponse{}, types.ErrCredentialNotFound
	}

	// check svc
	svc, found := k.GetService(ctx, msg.Sid)
	if !found || svc.Status != types.SERVICE_STATUS_ACTIVE {
		return &types.MsgUpdateVCResponse{}, types.ErrServiceNotActive
	}

	// check issuer
	issuer := sdk.MustAccAddressFromBech32(msg.Issuer)
	if did, found := k.GetDID(ctx, issuer); !found || did != svc.Issuer {
		return &types.MsgUpdateVCResponse{}, types.ErrIssuerNotFound
	}

	// check holder
	if doc, found := k.GetDidDocument(ctx, msg.Did); !found || doc.Info.Status != types.DID_STATUS_ACTIVE {
		return &types.MsgUpdateVCResponse{}, types.ErrHolderNotActive
	}

	vc = types.NewCredential(msg.Did, msg.Sid, msg.Hash, msg.Uri)
	k.SetCredential(ctx, msg.Did, msg.Sid, vc)

	return &types.MsgUpdateVCResponse{}, nil
}

func (k msgServer) RemoveVC(goCtx context.Context, msg *types.MsgRemoveVC) (*types.MsgRemoveVCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	svc, found := k.GetService(ctx, msg.Sid)
	if !found {
		return &types.MsgRemoveVCResponse{}, types.ErrIssuerNotFound
	}

	// check issuer
	addr := sdk.MustAccAddressFromBech32(msg.Issuer)
	if did, found := k.GetDID(ctx, addr); !found || did != svc.Issuer {
		return &types.MsgRemoveVCResponse{}, types.ErrIssuerNotFound
	}

	k.DeleteCredential(ctx, msg.Did, msg.Sid)

	return &types.MsgRemoveVCResponse{}, nil
}
