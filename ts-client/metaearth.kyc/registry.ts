import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgDeleteSBT } from "./types/metaearth/kyc/tx";
import { MsgApprove } from "./types/metaearth/kyc/tx";
import { MsgRemove } from "./types/metaearth/kyc/tx";
import { MsgUpdateSBT } from "./types/metaearth/kyc/tx";
import { MsgUpdate } from "./types/metaearth/kyc/tx";
import { MsgUpdateKycIssuer } from "./types/metaearth/kyc/tx";
import { MsgCreateSBT } from "./types/metaearth/kyc/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.kyc.MsgDeleteSBT", toCosmjsGeneratedType(MsgDeleteSBT)],
    ["/metaearth.kyc.MsgApprove", toCosmjsGeneratedType(MsgApprove)],
    ["/metaearth.kyc.MsgRemove", toCosmjsGeneratedType(MsgRemove)],
    ["/metaearth.kyc.MsgUpdateSBT", toCosmjsGeneratedType(MsgUpdateSBT)],
    ["/metaearth.kyc.MsgUpdate", toCosmjsGeneratedType(MsgUpdate)],
    ["/metaearth.kyc.MsgUpdateKycIssuer", toCosmjsGeneratedType(MsgUpdateKycIssuer)],
    ["/metaearth.kyc.MsgCreateSBT", toCosmjsGeneratedType(MsgCreateSBT)],
    
];

export { msgTypes }