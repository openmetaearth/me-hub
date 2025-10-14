package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/st-chain/me-hub/x/blacklist/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateBlacklist allows authorized addresses to update the blacklist by adding and removing addresses
func (k msgServer) UpdateBlacklist(goCtx context.Context, msg *types.MsgUpdateBlacklist) (*types.MsgUpdateBlacklistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get current blacklist
	blacklist, found := k.GetBlacklist(ctx)
	if !found {
		// Initialize empty blacklist if not found
		blacklist = types.Blacklist{
			Addresses: []string{},
		}
	}

	// Create map for efficient removal
	toRemoveMap := make(map[string]struct{})
	for _, addr := range msg.AddressesToRemove {
		toRemoveMap[addr] = struct{}{}
	}

	// Remove addresses from blacklist
	var newAddrList []string
	for _, addr := range blacklist.Addresses {
		if _, found := toRemoveMap[addr]; !found {
			newAddrList = append(newAddrList, addr)
		}
	}

	// Create map to avoid duplicates
	existingAddrMap := make(map[string]struct{})
	for _, addr := range newAddrList {
		existingAddrMap[addr] = struct{}{}
	}

	// Add new addresses to blacklist
	for _, addrToAdd := range msg.AddressesToAdd {
		if _, found := existingAddrMap[addrToAdd]; !found {
			newAddrList = append(newAddrList, addrToAdd)
			existingAddrMap[addrToAdd] = struct{}{} // Handle duplicates in AddressesToAdd
		}
	}

	// Update blacklist
	blacklist.Addresses = newAddrList

	// Store updated blacklist
	if err := k.SetBlacklist(ctx, blacklist); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateBlacklist,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyAddressesAdded, fmt.Sprintf("%v", msg.AddressesToAdd)),
			sdk.NewAttribute(types.AttributeKeyAddressesRemoved, fmt.Sprintf("%v", msg.AddressesToRemove)),
		),
	)

	return &types.MsgUpdateBlacklistResponse{}, nil
}
