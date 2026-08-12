import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdateDao } from "./types/metaearth/dao/tx";
import { MsgFreeGasAccount } from "./types/metaearth/dao/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.dao.MsgUpdateDao", MsgUpdateDao],
    ["/metaearth.dao.MsgFreeGasAccount", MsgFreeGasAccount],
    
];

export { msgTypes }