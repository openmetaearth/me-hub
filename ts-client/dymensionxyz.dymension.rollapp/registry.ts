import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdateState } from "./types/dymensionxyz/dymension/rollapp/tx";
import { MsgCreateRollapp } from "./types/dymensionxyz/dymension/rollapp/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/dymensionxyz.dymension.rollapp.MsgUpdateState", MsgUpdateState],
    ["/dymensionxyz.dymension.rollapp.MsgCreateRollapp", MsgCreateRollapp],
    
];

export { msgTypes }