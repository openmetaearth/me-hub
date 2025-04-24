import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgMintNFT } from "./types/metaearth/wnft/tx";
import { MsgSend } from "./types/metaearth/wnft/tx";
import { MsgNewClass } from "./types/metaearth/wnft/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.wnft.MsgMintNFT", MsgMintNFT],
    ["/metaearth.wnft.MsgSend", MsgSend],
    ["/metaearth.wnft.MsgNewClass", MsgNewClass],
    
];

export { msgTypes }