import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgMintNFT } from "./types/metaearth/wnft/tx";
import { MsgSend } from "./types/metaearth/wnft/tx";
import { MsgNewClass } from "./types/metaearth/wnft/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.wnft.MsgMintNFT", toCosmjsGeneratedType(MsgMintNFT)],
    ["/metaearth.wnft.MsgSend", toCosmjsGeneratedType(MsgSend)],
    ["/metaearth.wnft.MsgNewClass", toCosmjsGeneratedType(MsgNewClass)],
    
];

export { msgTypes }