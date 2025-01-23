import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgLeaveGroupRequest } from "./types/metaearth/megroup/tx";
import { MsgCreateGroup } from "./types/metaearth/megroup/tx";
import { MsgJoinGroup } from "./types/metaearth/megroup/tx";
import { MsgUpdateGroup } from "./types/metaearth/megroup/tx";
import { MsgDeleteGroup } from "./types/metaearth/megroup/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.megroup.MsgLeaveGroupRequest", MsgLeaveGroupRequest],
    ["/metaearth.megroup.MsgCreateGroup", MsgCreateGroup],
    ["/metaearth.megroup.MsgJoinGroup", MsgJoinGroup],
    ["/metaearth.megroup.MsgUpdateGroup", MsgUpdateGroup],
    ["/metaearth.megroup.MsgDeleteGroup", MsgDeleteGroup],
    
];

export { msgTypes }