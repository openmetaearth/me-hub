import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgJoinGroup } from "./types/metaearth/megroup/tx";
import { MsgDeleteGroup } from "./types/metaearth/megroup/tx";
import { MsgCreateGroup } from "./types/metaearth/megroup/tx";
import { MsgUpdateGroup } from "./types/metaearth/megroup/tx";
import { MsgLeaveGroupRequest } from "./types/metaearth/megroup/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.megroup.MsgJoinGroup", MsgJoinGroup],
    ["/metaearth.megroup.MsgDeleteGroup", MsgDeleteGroup],
    ["/metaearth.megroup.MsgCreateGroup", MsgCreateGroup],
    ["/metaearth.megroup.MsgUpdateGroup", MsgUpdateGroup],
    ["/metaearth.megroup.MsgLeaveGroupRequest", MsgLeaveGroupRequest],
    
];

export { msgTypes }