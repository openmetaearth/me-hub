import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdateDao } from "./types/metaearth/dao/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.dao.MsgUpdateDao", toCosmjsGeneratedType(MsgUpdateDao)],
    
];

export { msgTypes }