import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType, Registry } from "@cosmjs/proto-signing";
import {
  MsgSubmitProposal,
  MsgVote,
  MsgSubmitEquityUpdate,
  MsgSubmitFreezeEquity,
  MsgSubmitUnfreezeEquity,
  MsgApprovePendingAction,
  MsgRejectPendingAction,
  MsgFreezeRegion,
  MsgUnfreezeRegion,
  MsgFreezeEquity,
  MsgUnfreezeEquity,
} from "./types/metaearth/regiongov/tx";

export const msgTypes: Array<[string, GeneratedType]> = [
  ["/metaearth.regiongov.MsgSubmitProposal", toCosmjsGeneratedType(MsgSubmitProposal)],
  ["/metaearth.regiongov.MsgVote", toCosmjsGeneratedType(MsgVote)],
  ["/metaearth.regiongov.MsgSubmitEquityUpdate", toCosmjsGeneratedType(MsgSubmitEquityUpdate)],
  ["/metaearth.regiongov.MsgSubmitFreezeEquity", toCosmjsGeneratedType(MsgSubmitFreezeEquity)],
  ["/metaearth.regiongov.MsgSubmitUnfreezeEquity", toCosmjsGeneratedType(MsgSubmitUnfreezeEquity)],
  ["/metaearth.regiongov.MsgApprovePendingAction", toCosmjsGeneratedType(MsgApprovePendingAction)],
  ["/metaearth.regiongov.MsgRejectPendingAction", toCosmjsGeneratedType(MsgRejectPendingAction)],
  ["/metaearth.regiongov.MsgFreezeRegion", toCosmjsGeneratedType(MsgFreezeRegion)],
  ["/metaearth.regiongov.MsgUnfreezeRegion", toCosmjsGeneratedType(MsgUnfreezeRegion)],
  ["/metaearth.regiongov.MsgFreezeEquity", toCosmjsGeneratedType(MsgFreezeEquity)],
  ["/metaearth.regiongov.MsgUnfreezeEquity", toCosmjsGeneratedType(MsgUnfreezeEquity)],
];

export const registry = new Registry(msgTypes);
