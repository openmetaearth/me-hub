import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdateState } from "./types/metaearth/rollapp/tx";
import { MsgSkipDelayRollapp } from "./types/metaearth/rollapp/tx";
import { MsgUpdateRollapp } from "./types/metaearth/rollapp/tx";
import { MsgCreateRollapp } from "./types/metaearth/rollapp/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.rollapp.MsgUpdateState", MsgUpdateState],
    ["/metaearth.rollapp.MsgSkipDelayRollapp", MsgSkipDelayRollapp],
    ["/metaearth.rollapp.MsgUpdateRollapp", MsgUpdateRollapp],
    ["/metaearth.rollapp.MsgCreateRollapp", MsgCreateRollapp],
    
];

export { msgTypes }