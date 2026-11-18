import { toCosmjsGeneratedType } from "../registry-codec";
import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgWithdrawFromGlobalDaoFeePool } from "./types/metaearth/wstaking/tx";
import { MsgNewRecord } from "./types/metaearth/wstaking/tx";
import { MsgIbcTransferFromRegionTreasure } from "./types/metaearth/wstaking/tx";
import { MsgNewRegion } from "./types/metaearth/wstaking/tx";
import { MsgUnstake } from "./types/metaearth/wstaking/tx";
import { MsgSetFixedDepositCfgStatus } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawFixedDeposit } from "./types/metaearth/wstaking/tx";
import { MsgRemoveRegion } from "./types/metaearth/wstaking/tx";
import { MsgRemoveFixedDepositCfg } from "./types/metaearth/wstaking/tx";
import { MsgStake } from "./types/metaearth/wstaking/tx";
import { MsgWithdrawDelegatorReward } from "./types/metaearth/wstaking/tx";
import { MsgNewFixedDepositCfg } from "./types/metaearth/wstaking/tx";
import { MsgSetFixedDepositCfgRate } from "./types/metaearth/wstaking/tx";
import { MsgDoFixedDeposit } from "./types/metaearth/wstaking/tx";
import { MsgReviewRecord } from "./types/metaearth/wstaking/tx";
import { MsgTransferRegion } from "./types/metaearth/wstaking/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/metaearth.wstaking.MsgWithdrawFromGlobalDaoFeePool", toCosmjsGeneratedType(MsgWithdrawFromGlobalDaoFeePool)],
    ["/metaearth.wstaking.MsgNewRecord", toCosmjsGeneratedType(MsgNewRecord)],
    ["/metaearth.wstaking.MsgIbcTransferFromRegionTreasure", toCosmjsGeneratedType(MsgIbcTransferFromRegionTreasure)],
    ["/metaearth.wstaking.MsgNewRegion", toCosmjsGeneratedType(MsgNewRegion)],
    ["/metaearth.wstaking.MsgUnstake", toCosmjsGeneratedType(MsgUnstake)],
    ["/metaearth.wstaking.MsgSetFixedDepositCfgStatus", toCosmjsGeneratedType(MsgSetFixedDepositCfgStatus)],
    ["/metaearth.wstaking.MsgWithdrawFixedDeposit", toCosmjsGeneratedType(MsgWithdrawFixedDeposit)],
    ["/metaearth.wstaking.MsgRemoveRegion", toCosmjsGeneratedType(MsgRemoveRegion)],
    ["/metaearth.wstaking.MsgRemoveFixedDepositCfg", toCosmjsGeneratedType(MsgRemoveFixedDepositCfg)],
    ["/metaearth.wstaking.MsgStake", toCosmjsGeneratedType(MsgStake)],
    ["/metaearth.wstaking.MsgWithdrawDelegatorReward", toCosmjsGeneratedType(MsgWithdrawDelegatorReward)],
    ["/metaearth.wstaking.MsgNewFixedDepositCfg", toCosmjsGeneratedType(MsgNewFixedDepositCfg)],
    ["/metaearth.wstaking.MsgSetFixedDepositCfgRate", toCosmjsGeneratedType(MsgSetFixedDepositCfgRate)],
    ["/metaearth.wstaking.MsgDoFixedDeposit", toCosmjsGeneratedType(MsgDoFixedDeposit)],
    ["/metaearth.wstaking.MsgReviewRecord", toCosmjsGeneratedType(MsgReviewRecord)],
    ["/metaearth.wstaking.MsgTransferRegion", toCosmjsGeneratedType(MsgTransferRegion)],
    
];

export { msgTypes }