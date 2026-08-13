import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCreateSequencer } from "./types/metaearth/sequencer/tx";
import { MsgUnbond } from "./types/metaearth/sequencer/tx";
import { MsgReplaceProposerRequest } from "./types/metaearth/sequencer/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.sequencer.MsgCreateSequencer", MsgCreateSequencer],
    ["/metaearth.sequencer.MsgUnbond", MsgUnbond],
    ["/metaearth.sequencer.MsgReplaceProposerRequest", MsgReplaceProposerRequest],
];

export { msgTypes }
