import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdate } from "./types/metaearth/kyc/tx";
import { MsgRemove } from "./types/metaearth/kyc/tx";
import { MsgCreateSBT } from "./types/metaearth/kyc/tx";
import { MsgDeleteSBT } from "./types/metaearth/kyc/tx";
import { MsgUpdateSBT } from "./types/metaearth/kyc/tx";
import { MsgApprove } from "./types/metaearth/kyc/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.kyc.MsgUpdate", MsgUpdate],
    ["/metaearth.kyc.MsgRemove", MsgRemove],
    ["/metaearth.kyc.MsgCreateSBT", MsgCreateSBT],
    ["/metaearth.kyc.MsgDeleteSBT", MsgDeleteSBT],
    ["/metaearth.kyc.MsgUpdateSBT", MsgUpdateSBT],
    ["/metaearth.kyc.MsgApprove", MsgApprove],
    
];

export { msgTypes }