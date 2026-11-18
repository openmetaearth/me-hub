import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCreateSequencer } from "./types/metaearth/sequencer/tx";
import { MsgUnbond } from "./types/metaearth/sequencer/tx";
import { MsgReplaceProposerRequest } from "./types/metaearth/sequencer/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.sequencer.MsgCreateSequencer", toCosmjsGeneratedType(MsgCreateSequencer)],
    ["/metaearth.sequencer.MsgUnbond", toCosmjsGeneratedType(MsgUnbond)],
    ["/metaearth.sequencer.MsgReplaceProposerRequest", toCosmjsGeneratedType(MsgReplaceProposerRequest)],
];

export { msgTypes }
