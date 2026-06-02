package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func (s *KeeperTestSuite) TestExportImportGenesisFixedDepositRuntimeState() {
	s.SetupTest()

	regionId := types.MeEarthRegionId
	term := int64(30)
	cfg := types.FixedDepositCfg{
		RegionId: regionId,
		Term:     term,
		Rate:     sdk.MustNewDecFromStr("0.1234"),
		Status:   types.RegionFixedDepositCfgStatusActive,
	}
	count := types.FixedDepositCountOfCfg{
		RegionId: regionId,
		Term:     term,
		Count:    2,
	}
	total := types.FixedDepositTotal{
		Amount: sdk.NewCoin(params.BaseDenom, sdk.NewInt(123456)),
	}

	s.Keeper().SetFixedDepositCfg(s.Ctx, cfg)
	s.Keeper().SetFixedDepositCountOfCfg(s.Ctx, regionId, term, count.Count)
	s.Keeper().SetFixedDepositTotalAmount(s.Ctx, total)

	genesis := s.Keeper().ExportGenesis(s.Ctx)
	s.Require().Equal([]types.FixedDepositCfg{cfg}, genesis.FixedDepositCfgList)
	s.Require().Equal([]types.FixedDepositCountOfCfg{count}, genesis.FixedDepositCountOfCfgList)
	s.Require().NotNil(genesis.FixedDepositTotalAmount)
	s.Require().Equal(total, *genesis.FixedDepositTotalAmount)

	s.Keeper().RemoveFixedDepositCfg(s.Ctx, regionId, term)
	s.Keeper().SetFixedDepositCountOfCfg(s.Ctx, regionId, term, 0)
	s.Keeper().SetFixedDepositTotalAmount(s.Ctx, types.FixedDepositTotal{
		Amount: sdk.NewCoin(params.BaseDenom, sdk.ZeroInt()),
	})

	s.Keeper().InitGenesis(s.Ctx, genesis)

	importedCfg, found := s.Keeper().GetFixedDepositCfg(s.Ctx, regionId, term)
	s.Require().True(found)
	s.Require().Equal(cfg, importedCfg)
	s.Require().Equal(count.Count, s.Keeper().GetFixedDepositCountOfCfg(s.Ctx, regionId, term))

	importedTotal, found := s.Keeper().GetFixedDepositTotalAmount(s.Ctx)
	s.Require().True(found)
	s.Require().Equal(total, importedTotal)
}
