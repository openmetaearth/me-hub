import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgSkipDelayRollapp } from "./types/metaearth/rollapp/tx";
import { MsgCreateRollapp } from "./types/metaearth/rollapp/tx";
import { MsgUpdateState } from "./types/metaearth/rollapp/tx";
import { MsgUpdateApp } from "./types/metaearth/rollapp/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.rollapp.MsgSkipDelayRollapp", toCosmjsGeneratedType(MsgSkipDelayRollapp)],
    ["/metaearth.rollapp.MsgCreateRollapp", toCosmjsGeneratedType(MsgCreateRollapp)],
    ["/metaearth.rollapp.MsgUpdateState", toCosmjsGeneratedType(MsgUpdateState)],
    ["/metaearth.rollapp.MsgUpdateApp", toCosmjsGeneratedType(MsgUpdateApp)],
];

export { msgTypes }
