import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgDeleteSBT } from "./types/metaearth/kyc/tx";
import { MsgApprove } from "./types/metaearth/kyc/tx";
import { MsgRemove } from "./types/metaearth/kyc/tx";
import { MsgUpdateSBT } from "./types/metaearth/kyc/tx";
import { MsgUpdate } from "./types/metaearth/kyc/tx";
import { MsgCreateSBT } from "./types/metaearth/kyc/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.kyc.MsgDeleteSBT", MsgDeleteSBT],
    ["/metaearth.kyc.MsgApprove", MsgApprove],
    ["/metaearth.kyc.MsgRemove", MsgRemove],
    ["/metaearth.kyc.MsgUpdateSBT", MsgUpdateSBT],
    ["/metaearth.kyc.MsgUpdate", MsgUpdate],
    ["/metaearth.kyc.MsgCreateSBT", MsgCreateSBT],
    
];

export { msgTypes }