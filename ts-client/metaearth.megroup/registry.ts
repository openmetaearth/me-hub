import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgDeleteGroup } from "./types/metaearth/megroup/tx";
import { MsgCreateGroup } from "./types/metaearth/megroup/tx";
import { MsgUpdateGroup } from "./types/metaearth/megroup/tx";
import { MsgJoinGroup } from "./types/metaearth/megroup/tx";
import { MsgLeaveGroupRequest } from "./types/metaearth/megroup/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.megroup.MsgDeleteGroup", toCosmjsGeneratedType(MsgDeleteGroup)],
    ["/metaearth.megroup.MsgCreateGroup", toCosmjsGeneratedType(MsgCreateGroup)],
    ["/metaearth.megroup.MsgUpdateGroup", toCosmjsGeneratedType(MsgUpdateGroup)],
    ["/metaearth.megroup.MsgJoinGroup", toCosmjsGeneratedType(MsgJoinGroup)],
    ["/metaearth.megroup.MsgLeaveGroupRequest", toCosmjsGeneratedType(MsgLeaveGroupRequest)],
    
];

export { msgTypes }