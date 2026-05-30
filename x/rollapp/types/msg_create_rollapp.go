package types

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Compile-time interface check.
var _ sdk.Msg = &MsgCreateRollapp{}

// Constants for validation limits and security boundaries.
const (
	// MinSequencers defines the minimum allowed number of sequencers.
	MinSequencers uint64 = 1

	// MaxSequencersUpperBound defines the maximum allowed number of sequencers.
	MaxSequencersUpperBound uint64 = 1000

	// MaxRollappIDLength defines the maximum length of a rollapp ID after trimming.
	MaxRollappIDLength int = 128

	// MaxPermissionedAddresses defines a sensible upper limit to prevent gas exhaustion.
	MaxPermissionedAddresses int = 100

	// rollappIDPattern defines the allowed character set for rollapp IDs.
	// Only alphanumeric characters, hyphens, underscores, dots, and colons are permitted.
	rollappIDPattern = `^[a-zA-Z0-9_\-.:]+$`
)

// rollappIDRegexp validates the allowed characters in a rollapp ID.
var rollappIDRegexp = regexp.MustCompile(rollappIDPattern)

// MsgCreateRollapp defines a Cosmos SDK message to create a new rollapp.
//
// After successful validation via ValidateBasic(), the RollappId field is guaranteed to be
// whitespace-trimmed, canonical (valid chain ID per NewChainID), and free of leading/trailing
// whitespace. This prevents namespace squatting and EIP155 index conflicts.
type MsgCreateRollapp struct {
	// Creator is the bech32 address of the rollapp creator.
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`

	// RollappId is the unique identifier for the rollapp. Must be a valid chain ID.
	// ValidateBasic() normalizes this field by trimming whitespace and validating format.
	RollappId string `protobuf:"bytes,2,opt,name=rollappId,proto3" json:"rollappId,omitempty"`

	// MaxSequencers is the maximum number of sequencers allowed for this rollapp.
	// Must be within [MinSequencers, MaxSequencersUpperBound].
	MaxSequencers uint64 `protobuf:"varint,3,opt,name=maxSequencers,proto3" json:"maxSequencers,omitempty"`

	// PermissionedAddresses is a list of bech32 addresses allowed to operate as sequencers.
	// Must contain at least one address; all addresses must be unique and valid.
	PermissionedAddresses []string `protobuf:"bytes,4,rep,name=permissionedAddresses,proto3" json:"permissionedAddresses,omitempty"`
}

// NewMsgCreateRollapp creates a new MsgCreateRollapp instance.
// The caller should call ValidateBasic before using the message in any state transition.
func NewMsgCreateRollapp(creator string, rollappId string, maxSequencers uint64, permissionedAddresses []string) *MsgCreateRollapp {
	return &MsgCreateRollapp{
		Creator:               creator,
		RollappId:             rollappId,
		MaxSequencers:         maxSequencers,
		PermissionedAddresses: permissionedAddresses,
	}
}

// Route returns the message router key for routing to the rollapp module.
func (msg MsgCreateRollapp) Route() string {
	return RouterKey
}

// Type returns the message type for routing and identification.
func (msg MsgCreateRollapp) Type() string {
	return "create_rollapp"
}

// GetSigners returns the signer(s) required to authorize the message.
// Precondition: ValidateBasic has been called so Creator is a valid bech32 address.
// If the creator is invalid (should not happen after validation), returns an empty slice.
func (msg MsgCreateRollapp) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		// This should never occur if ValidateBasic was called; return empty as fallback.
		return []sdk.AccAddress{}
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs comprehensive validation and normalization of the message fields.
// It mutates the message to trim whitespace from RollappId and Creator, and enforces all
// security constraints. This method must be called before processing the message further.
//
// Returns an error if any field is invalid, wrapped with appropriate SDK error codes.
func (msg *MsgCreateRollapp) ValidateBasic() error {
	// 1. Normalize and validate RollappId
	if err := msg.normalizeAndValidateRollappID(); err != nil {
		return err
	}
	// 2. Normalize and validate Creator address
	if err := msg.normalizeAndValidateCreator(); err != nil {
		return err
	}
	// 3. Validate MaxSequencers
	if err := msg.validateMaxSequencers(); err != nil {
		return err
	}
	// 4. Validate PermissionedAddresses
	if err := msg.validatePermissionedAddresses(); err != nil {
		return err
	}
	return nil
}

// NormalizeRollappID is a public helper that trims whitespace, validates format,
// and checks against NewChainID. It can be used both by this message and the keeper.
func NormalizeRollappID(id string) (string, error) {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return "", sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"rollapp ID cannot be empty after trimming whitespace")
	}
	if len(trimmed) > MaxRollappIDLength {
		return "", sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest,
			"rollapp ID length exceeds maximum %d characters", MaxRollappIDLength)
	}
	if !rollappIDRegexp.MatchString(trimmed) {
		return "", sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"rollapp ID must only contain alphanumeric characters, hyphens, underscores, dots, or colons")
	}
	// Validate as a proper chain ID (including EIP155 if applicable).
	if _, err := NewChainID(trimmed); err != nil {
		return "", sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid rollapp chain ID '%s': %s", trimmed, err)
	}
	return trimmed, nil
}

// normalizeAndValidateRollappID uses NormalizeRollappID and updates the message field.
func (msg *MsgCreateRollapp) normalizeAndValidateRollappID() error {
	normalized, err := NormalizeRollappID(msg.RollappId)
	if err != nil {
		return err
	}
	msg.RollappId = normalized
	return nil
}

// normalizeAndValidateCreator trims whitespace and validates the creator bech32 address.
func (msg *MsgCreateRollapp) normalizeAndValidateCreator() error {
	msg.Creator = strings.TrimSpace(msg.Creator)
	if msg.Creator == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress,
			"creator address cannot be empty after trimming whitespace")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress,
			"invalid creator address: %s", err)
	}
	return nil
}

// validateMaxSequencers checks that MaxSequencers is within the allowed range.
func (msg *MsgCreateRollapp) validateMaxSequencers() error {
	if msg.MaxSequencers < MinSequencers {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest,
			"max sequencers must be at least %d, got %d", MinSequencers, msg.MaxSequencers)
	}
	if msg.MaxSequencers > MaxSequencersUpperBound {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest,
			"max sequencers cannot exceed %d, got %d", MaxSequencersUpperBound, msg.MaxSequencers)
	}
	return nil
}

// validatePermissionedAddresses checks that addresses are non-empty, valid bech32,
// contain no duplicates, and respect the upper bound.
func (msg *MsgCreateRollapp) validatePermissionedAddresses() error {
	if len(msg.PermissionedAddresses) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"at least one permissioned address is required")
	}
	if len(msg.PermissionedAddresses) > MaxPermissionedAddresses {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest,
			"too many permissioned addresses: %d (max %d)", len(msg.PermissionedAddresses), MaxPermissionedAddresses)
	}
	addrSet := make(map[string]struct{}, len(msg.PermissionedAddresses))
	for i, addr := range msg.PermissionedAddresses {
		trimmedAddr := strings.TrimSpace(addr)
		if trimmedAddr == "" {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress,
				"permissioned address at index %d is empty after trimming whitespace", i)
		}
		if _, err := sdk.AccAddressFromBech32(trimmedAddr); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress,
				"invalid permissioned address at index %d: %s", i, err)
		}
		if _, exists := addrSet[trimmedAddr]; exists {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest,
				"duplicate permissioned address at index %d: %s", i, trimmedAddr)
		}
		addrSet[trimmedAddr] = struct{}{}
	}
	return nil
}

// GetSignBytes returns the canonical bytes for signing the message.
func (msg MsgCreateRollapp) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetRollappId returns the normalized rollapp ID (after calling ValidateBasic).
func (msg *MsgCreateRollapp) GetRollappId() string {
	// If ValidateBasic has been called, this is already trimmed.
	// Otherwise we still trim for safety.
	return strings.TrimSpace(msg.RollappId)
}

// GetCreator returns the normalized creator address.
func (msg *MsgCreateRollapp) GetCreator() string {
	return strings.TrimSpace(msg.Creator)
}

// String returns a human-readable representation of the message.
func (msg MsgCreateRollapp) String() string {
	return fmt.Sprintf("MsgCreateRollapp{Creator: %s, RollappId: %s, MaxSequencers: %d, PermissionedAddresses: %v}",
		msg.Creator, msg.RollappId, msg.MaxSequencers, msg.PermissionedAddresses)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces.
// Required for protobuf compatibility.
func (msg *MsgCreateRollapp) UnpackInterfaces(_ context.Context) error {
	// No nested interfaces to unpack.
	return nil
}