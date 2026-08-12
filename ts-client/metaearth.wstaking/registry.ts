import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgWithdrawFixedDeposit } from "./types/metaearth/wstaking/tx";
import { MsgRevokeRegionWithdraw } from "./types/metaearth/wstaking/tx";
import { MsgDoFixedDeposit } from "./types/metaearth/wstaking/tx";
import { MsgNewRegion } from "./types/metaearth/wstaking/tx";
import { MsgSetFixedDepositCfgStatus } from "./types/metaearth/wstaking/tx";
import { MsgNewRecord } from "./types/metaearth/wstaking/tx";
import { MsgRemoveFixedDepositCfg } from "./types/metaearth/wstaking/tx";
import { MsgGrantRegionWithdraw } from "./types/metaearth/wstaking/tx";
import { MsgSetFixedDepositCfgRate } from "./types/metaearth/wstaking/tx";
import { MsgIbcTransferFromRegionTreasure } from "./types/metaearth/wstaking/tx";
import { MsgRemoveRegion } from "./types/metaearth/wstaking/tx";
import { MsgReviewRecord } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawDelegatorReward } from "./types/metaearth/wstaking/tx";
import { MsgStake } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawFromGlobalDaoFeePool } from "./types/metaearth/wstaking/tx";
import { MsgNewFixedDepositCfg } from "./types/metaearth/wstaking/tx";
import { MsgUpdateValidator } from "./types/metaearth/wstaking/tx";
import { MsgUnstake } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawFromRegion } from "./types/metaearth/wstaking/tx";
import { MsgSendToModule } from "./types/metaearth/wstaking/tx";
import { MsgTransferRegion } from "./types/metaearth/wstaking/tx";
import { MsgReplaceConsensusPubKeyRequest } from "./types/metaearth/wstaking/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.wstaking.MsgWithdrawFixedDeposit", MsgWithdrawFixedDeposit],
    ["/metaearth.wstaking.MsgRevokeRegionWithdraw", MsgRevokeRegionWithdraw],
    ["/metaearth.wstaking.MsgDoFixedDeposit", MsgDoFixedDeposit],
    ["/metaearth.wstaking.MsgNewRegion", MsgNewRegion],
    ["/metaearth.wstaking.MsgSetFixedDepositCfgStatus", MsgSetFixedDepositCfgStatus],
    ["/metaearth.wstaking.MsgNewRecord", MsgNewRecord],
    ["/metaearth.wstaking.MsgRemoveFixedDepositCfg", MsgRemoveFixedDepositCfg],
    ["/metaearth.wstaking.MsgGrantRegionWithdraw", MsgGrantRegionWithdraw],
    ["/metaearth.wstaking.MsgSetFixedDepositCfgRate", MsgSetFixedDepositCfgRate],
    ["/metaearth.wstaking.MsgIbcTransferFromRegionTreasure", MsgIbcTransferFromRegionTreasure],
    ["/metaearth.wstaking.MsgRemoveRegion", MsgRemoveRegion],
    ["/metaearth.wstaking.MsgReviewRecord", MsgReviewRecord],
    ["/metaearth.wstaking.MsgWithdrawDelegatorReward", MsgWithdrawDelegatorReward],
    ["/metaearth.wstaking.MsgStake", MsgStake],
    ["/metaearth.wstaking.MsgWithdrawFromGlobalDaoFeePool", MsgWithdrawFromGlobalDaoFeePool],
    ["/metaearth.wstaking.MsgNewFixedDepositCfg", MsgNewFixedDepositCfg],
    ["/metaearth.wstaking.MsgUpdateValidator", MsgUpdateValidator],
    ["/metaearth.wstaking.MsgUnstake", MsgUnstake],
    ["/metaearth.wstaking.MsgWithdrawFromRegion", MsgWithdrawFromRegion],
    ["/metaearth.wstaking.MsgSendToModule", MsgSendToModule],
    ["/metaearth.wstaking.MsgTransferRegion", MsgTransferRegion],
    ["/metaearth.wstaking.MsgReplaceConsensusPubKeyRequest", MsgReplaceConsensusPubKeyRequest],
    
];

export { msgTypes }