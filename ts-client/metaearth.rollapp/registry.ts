import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgSkipDelayRollapp } from "./types/metaearth/rollapp/tx";
import { MsgCreateRollapp } from "./types/metaearth/rollapp/tx";
import { MsgUpdateState } from "./types/metaearth/rollapp/tx";
import { MsgUpdateRollapp } from "./types/metaearth/rollapp/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.rollapp.MsgSkipDelayRollapp", MsgSkipDelayRollapp],
    ["/metaearth.rollapp.MsgCreateRollapp", MsgCreateRollapp],
    ["/metaearth.rollapp.MsgUpdateState", MsgUpdateState],
    ["/metaearth.rollapp.MsgUpdateRollapp", MsgUpdateRollapp],
];

export { msgTypes }
