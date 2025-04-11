import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgDeleteGroup } from "./types/metaearth/megroup/tx";
import { MsgCreateGroup } from "./types/metaearth/megroup/tx";
import { MsgUpdateGroup } from "./types/metaearth/megroup/tx";
import { MsgJoinGroup } from "./types/metaearth/megroup/tx";
import { MsgLeaveGroupRequest } from "./types/metaearth/megroup/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.megroup.MsgDeleteGroup", MsgDeleteGroup],
    ["/metaearth.megroup.MsgCreateGroup", MsgCreateGroup],
    ["/metaearth.megroup.MsgUpdateGroup", MsgUpdateGroup],
    ["/metaearth.megroup.MsgJoinGroup", MsgJoinGroup],
    ["/metaearth.megroup.MsgLeaveGroupRequest", MsgLeaveGroupRequest],
    
];

export { msgTypes }