package keeper

import (
	"bytes"
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	//"github.com/celestiaorg/celestia-app/pkg/appconsts"
	//celestiaBlob "github.com/celestiaorg/celestia-node/blob"
	tenderminttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/x/rollapp/types"
	rollupTypes "github.com/dymensionxyz/dymension/v3/x/rollup/types"
	"plugin"
	"strconv"
)

func (k Keeper) GetLastSubmitBlockInfo(goCtx context.Context, req *types.MsgLastSubmitBlkRequest) (*types.MsgLastSubmitBlkResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//ctx.Logger().Info("receviced SubmitBlockDAInfo", "msg", req.String())
	if !k.RollappsEnabled(ctx) {
		return nil, types.ErrRollappsDisabled
	}
	isFound := k.IsRollappExist(ctx, req.RollappId)
	if !isFound {
		return nil, types.ErrUnknownRollappID
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupBlockKeyPrefix(req.RollappId))
	lastCommitKey := store.Get([]byte(types.KeyLastRollupCommit))
	if lastCommitKey != nil {
		if startBlkHeight, number, errP := types.ParserRollupKey(string(lastCommitKey)); errP != nil {
			return nil, errorsmod.Wrapf(types.ErrParserData, errP.Error())
		} else {
			return &types.MsgLastSubmitBlkResponse{
				LastBatch: &types.MsgSubmitBlockBatch{
					StartHeight: startBlkHeight,
					NumBlocks:   number,
				}}, nil

		}
	} else {
		return &types.MsgLastSubmitBlkResponse{LastBatch: nil}, nil
	}
}

func (k Keeper) GetSubmitterBlockStatics(goCtx context.Context, req *types.MsgSubmitBlockStaticsRequest) (*types.MsgSubmitBlockStaticsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.RollappsEnabled(ctx) {
		return nil, types.ErrRollappsDisabled
	}
	isFound := k.IsRollappExist(ctx, req.RollappId)
	if !isFound {
		return nil, types.ErrUnknownRollappID
	}
	if req.StartHeight < 1 {
		return nil, errorsmod.Wrapf(types.ErrInputParams, "StartHeight < 1")
	}
	if req.EndHeight < req.StartHeight {
		return nil, errorsmod.Wrapf(types.ErrInputParams, "req.EndHeight < req.StartHeight")
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupBlockWithSubmitterKeyPrefix(req.RollappId))
	//recStore.Set(types.ConvertBlockHeightToKey(req.StartHeight), types.ConvertToRecordSubmitVal(req.Creator, req.NumBlocks))
	iterator := sdk.KVStorePrefixIterator(store, types.ConvertBlockHeightToKey(req.StartHeight))
	defer iterator.Close() // nolint: errcheck
	endBlkHeight := types.ConvertBlockHeightToKey(req.EndHeight)
	bIsFindEndIndex := false
	mapSubmitterBlkInfo := make(map[string][]*types.MsgSubmitBlockBatch)
	for ; iterator.Valid(); iterator.Next() {
		if bytes.Compare(iterator.Key(), endBlkHeight) >= 0 {
			bIsFindEndIndex = true
		}
		blkHeight, err := strconv.ParseUint(string(iterator.Key()), 10, 64)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrParserData,
				fmt.Sprintf("ParseUint error. key = %s,err = %s", string(iterator.Key()), err.Error()))
		}
		if submitter, number, errP := types.ParserRecordSubmitVal(string(iterator.Value())); errP != nil {
			return nil, errorsmod.Wrapf(types.ErrParserData,
				fmt.Sprintf("ParserRecordSubmitVal error. val = %s,err = %s", string(iterator.Value()), errP.Error()))
		} else {
			var sliSubmitBlock []*types.MsgSubmitBlockBatch
			if submitBlockList, ok := mapSubmitterBlkInfo[submitter]; ok {
				sliSubmitBlock = submitBlockList
			}
			sliSubmitBlock = append(sliSubmitBlock, &types.MsgSubmitBlockBatch{
				StartHeight: blkHeight,
				NumBlocks:   number,
			})
			mapSubmitterBlkInfo[submitter] = sliSubmitBlock
		}

		if bIsFindEndIndex {
			break
		}
	}
	//遍历map
	var submitBlkStaticsList []*types.MsgSubmitterStaticsInfo
	for key, val := range mapSubmitterBlkInfo {
		submitBlkStaticsList = append(submitBlkStaticsList, &types.MsgSubmitterStaticsInfo{
			Submitter:  key,
			SubmitList: val,
		})
	}
	return &types.MsgSubmitBlockStaticsResponse{
		SubmitterNumber:   uint32(len(submitBlkStaticsList)),
		SubmitStaticsList: submitBlkStaticsList,
	}, nil
}

func (k msgServer) SubmitBlockDAInfo(goCtx context.Context, req *types.MsgBlkDAInfo) (*types.MsgBlkDAResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.Logger().Info("receviced SubmitBlockDAInfo", "msg", req.String())

	if !k.RollappsEnabled(ctx) {
		return nil, types.ErrRollappsDisabled
	}

	// load rollapp object for stateful validations
	rollapp, isFound := k.GetRollapp(ctx, req.RollappId)
	if !isFound {
		return nil, types.ErrUnknownRollappID
	}

	// check rollapp version
	if rollapp.Version != req.Version {
		return nil, errorsmod.Wrapf(types.ErrVersionMismatch, "rollappId(%s) current version is %d, but got %d", req.RollappId, rollapp.Version, req.Version)
	}

	if len(req.Blocks.LightBlocks) != int(req.NumBlocks) {
		return nil, errorsmod.Wrapf(types.ErrVersionMismatch, "rollappId(%s)  LightBlocks's number(%d) != NumBlocks(%d)",
			req.RollappId, len(req.Blocks.LightBlocks), req.NumBlocks)
	}
	//TODO:明确这是干嘛
	err := k.hooks.BeforeUpdateState(ctx, req.Creator, req.RollappId)
	if err != nil {
		return nil, err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupBlockKeyPrefix(req.RollappId))
	lastCommitKey := store.Get([]byte(types.KeyLastRollupCommit))
	lastBlkHeight := uint64(0)
	if lastCommitKey != nil {
		if startBlkHeight, number, errP := types.ParserRollupKey(string(lastCommitKey)); errP != nil {
			return nil, errorsmod.Wrapf(types.ErrParserData, errP.Error())
		} else {
			if startBlkHeight < 1 || number < 1 {
				return nil, errorsmod.Wrapf(types.ErrParserData,
					fmt.Sprintf("data from lastCommitKey.startBlkHeight = %d, number = %d", startBlkHeight, number))
			}
			lastBlkHeight = startBlkHeight + uint64(number) - 1
		}

	}

	//if lastBlkHeight > 1 {
	if req.StartHeight != lastBlkHeight+1 {
		return nil, types.ErrWrongBlockHeight
	} else {
		if resStatus, err := k.verifyRollBlkIsAllowSubmit(goCtx, req, &rollapp); err != nil {
			if resStatus != types.SUBMIT_BLOCK_NORMAL_ERR {
				if err = k.punishSubmitter(goCtx, resStatus, rollapp.Creator, rollapp.RollappId); err != nil {
					return nil, err
				}
			}
			return nil, err
		}
		resp := &types.MsgBlkDAResponse{}

		if err = k.commitRollupBlockDAInfo(ctx, req); err != nil {
			return nil, err
		}
		//
		recStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupBlockWithSubmitterKeyPrefix(req.RollappId))
		recStore.Set(types.ConvertBlockHeightToKey(req.StartHeight), types.ConvertToRecordSubmitVal(req.Creator, req.NumBlocks))
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventSubmitBlockDA,
				sdk.NewAttribute("moduleName", types.ModuleName),
				sdk.NewAttribute("submitter", req.Creator),
				sdk.NewAttribute("startHeight", strconv.FormatUint(req.StartHeight, 10)),
				sdk.NewAttribute("number", strconv.FormatUint(uint64(req.NumBlocks), 10)),
			),
		)

		return resp, nil

	}

}

// 对submitter进行惩罚
func (k msgServer) punishSubmitter(ctx context.Context, status int, address, rollappID string) error {

	rate := uint32(0)
	//目前讨论的结果就是，3,6，10
	if status == types.SUBMIT_BLOCK_DA_VALIDATE_ERR {
		rate = 30
	} else if status == types.SUBMIT_BLOCK_DA_VERIFY_ERR {
		rate = 60
	} else if status == types.SUBMIT_BLOCK_DA_VERIFY_FAILED {
		rate = 100
	} else {
		return nil
	}
	if rate > 0 {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if err := k.rollupKeeper.Punishment(sdkCtx, address, rollappID, rate, 0); err != nil {
			return err
		} else {
			//如果有对提交者进行了惩罚，则需要对其进行资质重估
			errP := k.rollupKeeper.RevaluateSequencer(sdkCtx, address, rollappID)
			if errP != nil {
				return errP
			}
		}
	}
	return nil
}

func (k msgServer) commitRollupBlockDAInfo(ctx sdk.Context, msgBlkInfo *types.MsgBlkDAInfo) error {
	data, err := k.cdc.Marshal(msgBlkInfo)
	if err != nil {
		return errorsmod.Wrapf(types.ErrParserData,
			fmt.Sprintf("Marshal MsgBlkDAInfo error.err = %s", err.Error()))
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupBlockKeyPrefix(msgBlkInfo.RollappId))
	blockKey := types.GetRollupBlockKey(msgBlkInfo.StartHeight, msgBlkInfo.NumBlocks)
	store.Set(blockKey, data)
	store.Set([]byte(types.KeyLastRollupCommit), blockKey)

	return nil
	/*目前方案有修改，暂时不走缓存方式了，L2也提交也变成batch submitter
	//从append寻找还有没有后续连续的区块
	appendStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetRollupAppendBlockKeyPrefix(msgBlkInfo.RollappId))
	iterator := sdk.KVStorePrefixIterator(appendStore, []byte{})
	defer iterator.Close() // nolint: errcheck
	var deleteKey [][]byte
	endBlkHeight := msgBlkInfo.StartHeight + uint64(msgBlkInfo.NumBlocks) - 1
	for ; iterator.Valid(); iterator.Next() {
		if startHeight, number, err := types.ParserRollupKey(string(iterator.Key())); err != nil {
			return errorsmod.Wrapf(types.ErrParserData,
				fmt.Sprintf("ParserRollupKey error.err = %s", err.Error()))
		} else {
			if endBlkHeight <= msgBlkInfo.StartHeight {
				deleteKey = append(deleteKey, iterator.Key())
			} else if startHeight == endBlkHeight+1 { //如果是继续连续的

			}

		}

	}
	*/

}

// 校验rollup提交的区块是否允许被写入
func (k msgServer) verifyRollBlkIsAllowSubmit(ctx context.Context, submitBlock *types.MsgBlkDAInfo, pRollapp *types.Rollapp) (int, error) {
	//校验区块的提交者的是否质押数量足够，这样就可以将区块的提交者和竞选的Sequencer剥离开，只要满足最小质押金额，都可以提交。
	//这样也可以方便L2层设置一个专门提交区块的batch submitter
	resp, err := k.Keeper.rollupKeeper.QueryStake(ctx, &rollupTypes.QueryStakeRequest{
		RollappId: pRollapp.RollappId,
		Address:   pRollapp.Creator})
	if err != nil {
		return types.SUBMIT_BLOCK_NORMAL_ERR, errorsmod.Wrapf(types.ErrLogic,
			fmt.Sprintf("QueryStake error, err = %s,addr = %s", err.Error(), pRollapp.Creator))
	}
	//
	//这样操作的原因是因为解质押一个周期内申请，下周期竞选完成就可以解质押。如果在竞选前一天进行解质押，此时该提交者提交的数据
	//通过了基本校验，但是后续被提交了对应的欺诈证明挑战。而由于目前的欺诈证明方式需要几天的时间，那么此时就完全有可能剩下的金额
	//已经很少了，即使扣除也只能扣一点点。
	stakeAmount := resp.StakeInfo.StakeAmount - resp.StakeInfo.ApplyUnStakeAmount
	if stakeAmount < k.Keeper.rollupKeeper.GetParams(sdk.UnwrapSDKContext(ctx)).MinStakeAmount*rollupTypes.MecPrecision {
		return types.SUBMIT_BLOCK_NORMAL_ERR, errorsmod.Wrapf(types.ErrNotEnoughStake,
			fmt.Sprintf("StakeAmount = %d,ApplyUnStakeAmount = %d", resp.StakeInfo.StakeAmount, resp.StakeInfo.ApplyUnStakeAmount))
	}

	//校验所有的light block
	for _, val := range submitBlock.Blocks.LightBlocks {
		pLightBlock, err := tenderminttypes.LightBlockFromProto(val)
		if err != nil {
			return types.SUBMIT_BLOCK_NORMAL_ERR, errorsmod.Wrapf(types.ErrParserData,
				fmt.Sprintf("LightBlockFromProto, err = %s", err.Error()))
		}
		if err = k.verifyRollupBlkConsensus(ctx, submitBlock.RollappId, pLightBlock, pRollapp.PermissionedAddresses); err != nil {
			return types.SUBMIT_BLOCK_NORMAL_ERR, err
		}
	}
	//校验DA 的commitProof
	p, err := plugin.Open("celestiaPlugin.so")
	if err != nil {
		return types.SUBMIT_BLOCK_NORMAL_ERR, errorsmod.Wrapf(types.ErrLoadPlugin,
			fmt.Sprintf(" err = %s", err.Error()))
	}
	pVal, err := p.Lookup("VerifyDACommitmentProof")
	if err != nil {
		return types.SUBMIT_BLOCK_NORMAL_ERR, errorsmod.Wrapf(types.ErrLoadPlugin,
			fmt.Sprintf(" Lookup function error.err = %s", err.Error()))
	}
	daVerify, ok := pVal.(func([]byte, []byte) (int, error))
	if !ok {
		return types.SUBMIT_BLOCK_NORMAL_ERR, errorsmod.Wrapf(types.ErrLoadPlugin,
			" Lookup function typeAssert error.")
	}
	verifyRes, err := daVerify(submitBlock.CommitmentProof, submitBlock.DaRoot)
	if err != nil {
		return types.SUBMIT_BLOCK_SUCCESS, nil
	}
	return verifyRes, errorsmod.Wrapf(types.ErrCommitVerify,
		fmt.Sprintf("err = %s", err.Error()))

}

// 1、校验区块Header信息是否经过足够的签名教研,2、签名者是否有2f在Rollup的Election sequencer中
func (k msgServer) verifyRollupBlkConsensus(ctx context.Context, rollAppId string, lightBlock *tenderminttypes.LightBlock, allowAddress []string) error {
	err := verifyRollBlkInfo(rollAppId, lightBlock)
	if err != nil {
		return err
	}
	//校验该区块的验证者是否为投票选举出来的

	resp, err := k.rollupKeeper.QueryElectionResult(ctx, &rollupTypes.QueryElectionRequest{rollAppId})
	if err != nil {
		return errorsmod.Wrapf(types.ErrParserData,
			fmt.Sprintf("QueryElectionResult error.err = %s", err.Error()))
	}

	bIsInitSeq := false
	voteNumber := uint32(0)
	sequencerLen := uint32(0)
	if 0 == resp.ElectionTime || nil == resp.NodeStatusList { //表明没有进行过选举，此时应该采用最初创建RollApp中的Sequencer
		bIsInitSeq = true
		voteNumber, sequencerLen, err = calcElectSequencerVoteForBlock(lightBlock, nil, allowAddress)
	} else {
		voteNumber, sequencerLen, err = calcElectSequencerVoteForBlock(lightBlock, resp.NodeStatusList, nil)
	}
	if err != nil {
		return err
	}

	f := sequencerLen / 3
	if voteNumber >= 2*f { //如果>=2f个，则认为是允许的
		return nil
	} else {
		if !bIsInitSeq { //当前的竞选不是初始化的，才会有可能有前一次的选举
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			param := k.rollupKeeper.GetParams(sdkCtx)
			tmp := sdkCtx.BlockTime().Unix() - int64(resp.ElectionTime)
			if tmp <= int64(param.ElectionInterimTime) { //当前时间间隔<=选举过渡期，则允许使用之前的选举Sequencer进行判断
				preElect, err := k.rollupKeeper.GetPreviousElectionResult(ctx, rollAppId)
				if err != nil {
					return err
				}
				if preElect.NodeStatusList != nil && len(preElect.NodeStatusList) > 0 { //如果有存在上一次的数据
					voteNumber, sequencerLen, err = calcElectSequencerVoteForBlock(lightBlock, preElect.NodeStatusList, nil)
				} else { //然后没有找到之前的选举结果，说明刚刚那次选举是第一次选举，则用allowAddress来代替之前的选举
					voteNumber, sequencerLen, err = calcElectSequencerVoteForBlock(lightBlock, nil, allowAddress)
				}
				if err != nil {
					return fmt.Errorf("%s in second time", err.Error())
				}
				f = sequencerLen / 3
				if voteNumber >= 2*f { //如果>=2f个，则认为是允许的
					return nil
				} else {
					return errorsmod.Wrapf(types.ErrNotEnoughSequencerSign,
						fmt.Sprintf("in second time.sequencer vote's number = %d,total = %d,blkHeight = %d",
							voteNumber, sequencerLen, lightBlock.Height))
				}

			}
		}
		return errorsmod.Wrapf(types.ErrNotEnoughSequencerSign,
			fmt.Sprintf("sequencer vote's number = %d,total = %d,blkHeight = %d", voteNumber, sequencerLen, lightBlock.Height))
	}

}

// 计算为block投票的验证者也同时存在于选举后的Sequencer中的数量
func calcElectSequencerVoteForBlock(pBlock *tenderminttypes.LightBlock, electNodeList []*rollupTypes.ElectionNodeStatus, initAllowAddress []string) (uint32, uint32, error) {

	mapSequencer := make(map[string]struct{})
	if electNodeList != nil && len(electNodeList) > 0 {
		for _, val := range electNodeList {
			if val.Status != rollupTypes.NodeSequencer {
				continue
			}
			mapSequencer[val.Address] = struct{}{}
		}
	} else {
		for _, val := range initAllowAddress {
			mapSequencer[val] = struct{}{}
		}
	}
	if len(mapSequencer) < 1 {
		return 0, 0, errorsmod.Wrapf(types.ErrLogic, "electNodeList and initAllowAddress are both empty in calcElectSequencerVoteForBlock")
	}

	voteNumber := uint32(0)
	for _, cmtSig := range pBlock.Commit.Signatures {
		if cmtSig.ForBlock() {
			if _, ok := mapSequencer[cmtSig.ValidatorAddress.String()]; ok {
				voteNumber++
			}
		}
	}

	sequencerLen := len(mapSequencer)
	return voteNumber, uint32(sequencerLen), nil

}

func verifyRollBlkInfo(rollAppID string, lightBlk *tenderminttypes.LightBlock) error {
	//对lightBlk的数据进行基本校验
	err := lightBlk.ValidateBasic(rollAppID)
	if err != nil {
		return errorsmod.Wrapf(types.ErrValidateSubmitBlock,
			fmt.Sprintf("Validate light block error. err  = %s", err.Error()))

	}
	//校验承诺和签名信息
	if err = lightBlk.ValidatorSet.VerifyCommit(rollAppID, lightBlk.Commit.BlockID, lightBlk.Height, lightBlk.Commit); err != nil {
		return errorsmod.Wrapf(types.ErrValidateSubmitBlock,
			fmt.Sprintf("VerifyCommit block error. err  = %s", err.Error()))
	}
	return nil

}

func (k Keeper) IsRollappExist(ctx sdk.Context, rollappId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RollappKeyPrefix))

	b := store.Get(types.RollappKey(
		rollappId,
	))
	if b == nil {
		return false
	}
	return true
}
