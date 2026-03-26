package types

// DONTCOVER

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
)

// x/megroup module sentinel errors
var (
	//ErrSample           = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrCheckGlobalAdmin      = errorsmod.Register(ModuleName, 1100, "check global admin error")
	ErrCreate                = errorsmod.Register(ModuleName, 1101, "create group err.")
	ErrMeidNotExists         = errorsmod.Register(ModuleName, 1102, "meid not exists")
	ErrDeleteGroup           = errorsmod.Register(ModuleName, 1103, "delete group error")
	ErrPermissionDenied      = errorsmod.Register(ModuleName, 1104, "permission denied")
	ErrGroupNotExist         = errorsmod.Register(ModuleName, 1105, "group not exist")
	ErrProcData              = errorsmod.Register(ModuleName, 1106, "process data error")
	ErrGroupNumberExist      = errorsmod.Register(ModuleName, 1107, "group's member has benn exist")
	ErrGroupCreateRepeated   = errorsmod.Register(ModuleName, 1108, "create group repeat")
	ErrGroupMemberRepeated   = errorsmod.Register(ModuleName, 1109, "join group repeat")
	ErrRegionNotExist        = errorsmod.Register(ModuleName, 1110, "region not exist")
	ErrNotSupport            = errorsmod.Register(ModuleName, 1111, "function not support")
	ErrExcute                = errorsmod.Register(ModuleName, 1112, "excute error.")
	ErrGroupMemberNotExist   = errorsmod.Register(ModuleName, 1113, "group member not exist")
	ErrGroupExceededInRegion = errorsmod.Register(ModuleName, 1115, "group has been exceeded in region")
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
