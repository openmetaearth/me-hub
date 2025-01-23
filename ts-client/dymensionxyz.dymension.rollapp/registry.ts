import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCreateRollapp } from "./types/dymensionxyz/dymension/rollapp/tx";
import { MsgUpdateState } from "./types/dymensionxyz/dymension/rollapp/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/dymensionxyz.dymension.rollapp.MsgCreateRollapp", MsgCreateRollapp],
    ["/dymensionxyz.dymension.rollapp.MsgUpdateState", MsgUpdateState],
    
];

export { msgTypes }