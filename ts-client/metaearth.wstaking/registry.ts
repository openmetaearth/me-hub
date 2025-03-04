import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgRemoveFixedDepositCfg } from "./types/metaearth/wstaking/tx";
import { MsgDoFixedDeposit } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawFromRegion } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawDelegatorReward } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawFromGlobalDaoFeePool } from "./types/metaearth/wstaking/tx";
import { MsgNewRecord } from "./types/metaearth/wstaking/tx";
import { MsgRemoveRegion } from "./types/metaearth/wstaking/tx";
import { MsgReviewRecord } from "./types/metaearth/wstaking/tx";
import { MsgSetFixedDepositCfgStatus } from "./types/metaearth/wstaking/tx";
import { MsgIbcTransferFromRegionTreasure } from "./types/metaearth/wstaking/tx";
import { MsgTransferRegion } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawFixedDeposit } from "./types/metaearth/wstaking/tx";
import { MsgNewFixedDepositCfg } from "./types/metaearth/wstaking/tx";
import { MsgStake } from "./types/metaearth/wstaking/tx";
import { MsgUnstake } from "./types/metaearth/wstaking/tx";
import { MsgNewRegion } from "./types/metaearth/wstaking/tx";
import { MsgSetFixedDepositCfgRate } from "./types/metaearth/wstaking/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.wstaking.MsgRemoveFixedDepositCfg", MsgRemoveFixedDepositCfg],
    ["/metaearth.wstaking.MsgDoFixedDeposit", MsgDoFixedDeposit],
    ["/metaearth.wstaking.MsgWithdrawFromRegion", MsgWithdrawFromRegion],
    ["/metaearth.wstaking.MsgWithdrawDelegatorReward", MsgWithdrawDelegatorReward],
    ["/metaearth.wstaking.MsgWithdrawFromGlobalDaoFeePool", MsgWithdrawFromGlobalDaoFeePool],
    ["/metaearth.wstaking.MsgNewRecord", MsgNewRecord],
    ["/metaearth.wstaking.MsgRemoveRegion", MsgRemoveRegion],
    ["/metaearth.wstaking.MsgReviewRecord", MsgReviewRecord],
    ["/metaearth.wstaking.MsgSetFixedDepositCfgStatus", MsgSetFixedDepositCfgStatus],
    ["/metaearth.wstaking.MsgIbcTransferFromRegionTreasure", MsgIbcTransferFromRegionTreasure],
    ["/metaearth.wstaking.MsgTransferRegion", MsgTransferRegion],
    ["/metaearth.wstaking.MsgWithdrawFixedDeposit", MsgWithdrawFixedDeposit],
    ["/metaearth.wstaking.MsgNewFixedDepositCfg", MsgNewFixedDepositCfg],
    ["/metaearth.wstaking.MsgStake", MsgStake],
    ["/metaearth.wstaking.MsgUnstake", MsgUnstake],
    ["/metaearth.wstaking.MsgNewRegion", MsgNewRegion],
    ["/metaearth.wstaking.MsgSetFixedDepositCfgRate", MsgSetFixedDepositCfgRate],
    
];

export { msgTypes }