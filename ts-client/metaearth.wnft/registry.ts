import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgNewClass } from "./types/metaearth/wnft/tx";
import { MsgMintNFT } from "./types/metaearth/wnft/tx";
import { MsgSend } from "./types/metaearth/wnft/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.wnft.MsgNewClass", MsgNewClass],
    ["/metaearth.wnft.MsgMintNFT", MsgMintNFT],
    ["/metaearth.wnft.MsgSend", MsgSend],
    
];

export { msgTypes }