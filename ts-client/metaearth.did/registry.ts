import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCreateVC } from "./types/metaearth/did/tx";
import { MsgUpdateDidStatus } from "./types/metaearth/did/tx";
import { MsgCreateService } from "./types/metaearth/did/tx";
import { MsgUpdateServiceStatus } from "./types/metaearth/did/tx";
import { MsgRemoveService } from "./types/metaearth/did/tx";
import { MsgUpdateVC } from "./types/metaearth/did/tx";
import { MsgCreateDid } from "./types/metaearth/did/tx";
import { MsgRemoveDid } from "./types/metaearth/did/tx";
import { MsgRemoveVC } from "./types/metaearth/did/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.did.MsgCreateVC", MsgCreateVC],
    ["/metaearth.did.MsgUpdateDidStatus", MsgUpdateDidStatus],
    ["/metaearth.did.MsgCreateService", MsgCreateService],
    ["/metaearth.did.MsgUpdateServiceStatus", MsgUpdateServiceStatus],
    ["/metaearth.did.MsgRemoveService", MsgRemoveService],
    ["/metaearth.did.MsgUpdateVC", MsgUpdateVC],
    ["/metaearth.did.MsgCreateDid", MsgCreateDid],
    ["/metaearth.did.MsgRemoveDid", MsgRemoveDid],
    ["/metaearth.did.MsgRemoveVC", MsgRemoveVC],
    
];

export { msgTypes }