import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgFulfillOrder } from "./types/metaearth/eibc/tx";
import { MsgUpdateDemandOrder } from "./types/metaearth/eibc/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.eibc.MsgFulfillOrder", toCosmjsGeneratedType(MsgFulfillOrder)],
    ["/metaearth.eibc.MsgUpdateDemandOrder", toCosmjsGeneratedType(MsgUpdateDemandOrder)],
];

export { msgTypes }
