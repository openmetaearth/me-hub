package keeper

import (
	"bytes"
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/openmetaearth/me-hub/x/did/types"
)

func (m msgServer) BondZkCoreID(goCtx context.Context, req *types.MsgBondZkCoreIDRequest) (*types.MsgBondZkCoreIDResponse, error) {
	if req == nil {
		return nil, errors.Wrap(types.ErrParameter, "request is nil")
	}

	if nil == m.evmKeeper {
		return nil, errors.Wrap(types.ErrZkEVMCall, "evm keeper is not set")
	}

	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	creator, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, errors.Wrap(types.ErrParameter, "invalid creator address")
	}

	// check if the creator is a valid account and verify public key is ethsecp256k1
	account := m.accountKeeper.GetAccount(ctx, creator)
	if account == nil {
		return nil, errors.Wrap(types.ErrPermissionDenial, "account must be created and have a public key set")
	}
	pk := account.GetPubKey()
	if pk == nil {
		return nil, errors.Wrap(types.ErrPermissionDenial, " account public key is not set")
	}
	if _, ok := pk.(*ethsecp256k1.PubKey); !ok {
		return nil, errors.Wrap(types.ErrPermissionDenial, " account public key is not ethsecp256k1")
	}
	//
	did, found := m.GetDID(ctx, creator)
	if !found {
		did, found = m.GetSubAccountDidMap(ctx, req.Creator)
		if !found {
			return nil, types.ErrDidNotFound
		}
	}

	didInfo, found := m.GetDidInfo(ctx, did)
	if !found {
		return nil, errors.Wrapf(types.ErrDidNotFound, "did-info not found,did = %s",did)
	}
	if didInfo.Status != types.DID_STATUS_ACTIVE {
		return nil, types.ErrDidNotActive
	}
	if didInfo.KycLevel < types.KYC_LEVEL_TWO {
		return nil, errors.Wrap(types.ErrUnauthorized, "kyc level is not enough")
	}

	if _, found := m.GetZkCoreIDInfo(ctx, did); found {
		return nil, types.ErrDidAlreadyBound
	}
	if boundDID, found := m.GetDIDByZkCoreID(ctx, req.CoreInfo.CoreId); found {
		return nil, errors.Wrapf(types.ErrZkCoreIDAlreadyBound, "bound to DID %s", boundDID)
	}
	if req.CoreInfo.ExpiresAt <= ctx.BlockTime().Unix() {
		return nil, types.ErrZkBindingExpired
	}

	coreID, err := types.ZkCoreIDToInt(req.CoreInfo.CoreId)
	if err != nil {
		return nil, errors.Wrap(types.ErrInvalidZkCoreID, err.Error())
	}
	caller := common.BytesToAddress(creator.Bytes())
	if err := m.checkZkCoreIDExists(ctx, caller, coreID); err != nil {
		return nil, err
	}
	proofCoreID, proofChallenge, err := m.verifyZkAuthProof(ctx, caller, req.CoreInfo.AuthProof)
	if err != nil {
		return nil, err
	}
	if proofCoreID.Cmp(coreID) != 0 {
		return nil, types.ErrZkCoreIDMismatch
	}
	expectedChallenge, err := types.CalculateZkCoreIDBindingChallenge(
		ctx.ChainID(), req.Creator, did, req.CoreInfo.CoreId,
		req.CoreInfo.Nonce, req.CoreInfo.ExpiresAt,
	)
	if err != nil {
		return nil, errors.Wrap(types.ErrParameter, err.Error())
	}
	if !bytes.Equal(proofChallenge.Bytes(), expectedChallenge.Bytes()) {
		return nil, types.ErrZkChallengeMismatch
	}

	m.SetZkCoreIDBinding(ctx, did, *req.CoreInfo)
	ctx.EventManager().EmitEvent(types.NewBondZkCoreIDEvent(did, req.Creator, req.CoreInfo.CoreId,ctx.BlockHeight()))
	return &types.MsgBondZkCoreIDResponse{}, nil
}

func (m msgServer) SetZkContractAddress(goCtx context.Context, req *types.MsgSetZkContractAddressRequest) (*types.MsgSetZkContractAddressResponse, error) {
	if req == nil {
		return nil, errors.Wrap(types.ErrParameter, "request is required")
	}
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if m.daoKeeper == nil || !m.daoKeeper.IsGlobalDao(ctx, req.Creator) {
		return nil, types.ErrPermissionDenial
	}

	contractInfo := *req.ContractInfo
	contractInfo.ContractAddress = common.HexToAddress(contractInfo.ContractAddress).Hex()
	m.SetZkContractAddressInfo(ctx, contractInfo)
	ctx.EventManager().EmitEvent(types.NewSetZkContractAddressEvent(contractInfo,ctx.BlockHeight()))
	return &types.MsgSetZkContractAddressResponse{}, nil
}
