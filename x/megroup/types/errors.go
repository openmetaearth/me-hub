package types

// DONTCOVER

import (
	"encoding/json"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/megroup module sentinel errors
var (
	ErrCheckGlobalDao        = sdkerrors.Register(ModuleName, 1100, "check global dao error")
	ErrCreate                = sdkerrors.Register(ModuleName, 1101, "create group err")
	ErrMeidNotExists         = sdkerrors.Register(ModuleName, 1102, "meid not exists")
	ErrDeleteGroup           = sdkerrors.Register(ModuleName, 1103, "delete group error")
	ErrPermissionDenied      = sdkerrors.Register(ModuleName, 1104, "permission denied")
	ErrGroupNotExist         = sdkerrors.Register(ModuleName, 1105, "group not exist")
	ErrProcData              = sdkerrors.Register(ModuleName, 1106, "process data error")
	ErrGroupNumberExist      = sdkerrors.Register(ModuleName, 1107, "group's member has benn exist")
	ErrGroupCreateRepeated   = sdkerrors.Register(ModuleName, 1108, "create group repeat")
	ErrGroupMemberRepeated   = sdkerrors.Register(ModuleName, 1109, "join group repeat")
	ErrRegionNotExist        = sdkerrors.Register(ModuleName, 1110, "region not exist")
	ErrNotSupport            = sdkerrors.Register(ModuleName, 1111, "function not support")
	ErrExecute               = sdkerrors.Register(ModuleName, 1112, "execute error")
	ErrGroupMemberNotExist   = sdkerrors.Register(ModuleName, 1113, "group member not exist")
	ErrGroupExceededInRegion = sdkerrors.Register(ModuleName, 1115, "group has been exceeded in region")
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
