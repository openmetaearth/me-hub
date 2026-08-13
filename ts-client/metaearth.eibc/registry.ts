import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgFulfillOrder } from "./types/metaearth/eibc/tx";
import { MsgUpdateDemandOrder } from "./types/metaearth/eibc/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.eibc.MsgFulfillOrder", MsgFulfillOrder],
    ["/metaearth.eibc.MsgUpdateDemandOrder", MsgUpdateDemandOrder],
];

export { msgTypes }
