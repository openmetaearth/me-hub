import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdateDidStatus } from "./types/metaearth/did/tx";
import { MsgUpdateVC } from "./types/metaearth/did/tx";
import { MsgUpdateServiceStatus } from "./types/metaearth/did/tx";
import { MsgCreateService } from "./types/metaearth/did/tx";
import { MsgCreateVC } from "./types/metaearth/did/tx";
import { MsgRemoveService } from "./types/metaearth/did/tx";
import { MsgRemoveVC } from "./types/metaearth/did/tx";
import { MsgCreateDid } from "./types/metaearth/did/tx";
import { MsgRemoveDid } from "./types/metaearth/did/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.did.MsgUpdateDidStatus", MsgUpdateDidStatus],
    ["/metaearth.did.MsgUpdateVC", MsgUpdateVC],
    ["/metaearth.did.MsgUpdateServiceStatus", MsgUpdateServiceStatus],
    ["/metaearth.did.MsgCreateService", MsgCreateService],
    ["/metaearth.did.MsgCreateVC", MsgCreateVC],
    ["/metaearth.did.MsgRemoveService", MsgRemoveService],
    ["/metaearth.did.MsgRemoveVC", MsgRemoveVC],
    ["/metaearth.did.MsgCreateDid", MsgCreateDid],
    ["/metaearth.did.MsgRemoveDid", MsgRemoveDid],
    
];

export { msgTypes }