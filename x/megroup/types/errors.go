package types

// DONTCOVER

import (
	"encoding/json"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/megroup module sentinel errors
var (
	//ErrSample           = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrCheckGlobalAdmin = sdkerrors.Register(ModuleName, 1100, "check global admin error")
	ErrCreate           = sdkerrors.Register(ModuleName, 1101, "create group err.")
	ErrMeidNotExists    = sdkerrors.Register(ModuleName, 1102, "meid not exists")
	ErrDeleteGroup      = sdkerrors.Register(ModuleName, 1103, "delete group error")
	ErrPermissionDenied = sdkerrors.Register(ModuleName, 1104, "permission denied")
	ErrGroupNotExist    = sdkerrors.Register(ModuleName, 1105, "group not exist")
)

func EmitNewGroupError(err string, msg *MsgCreateGroup) []byte {
	log := map[string]interface{}{
		"type":      "group",
		"creator":   msg.Creator,
		"GroupInfo": msg.GroupInfo,
		"err":       err,
	}
	logBytes, _ := json.Marshal(log)

	return logBytes
}
