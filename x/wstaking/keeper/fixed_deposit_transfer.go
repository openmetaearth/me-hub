package keeper

import (
	"errors"
	"fmt"

	"github.com/openmetaearth/me-hub/x/wstaking/types"
)

func fixedDepositConfigByTerm(configs []types.FixedDepositCfg) map[int64]types.FixedDepositCfg {
	configByTerm := make(map[int64]types.FixedDepositCfg, len(configs))
	for _, cfg := range configs {
		configByTerm[cfg.Term] = cfg
	}
	return configByTerm
}

func validateFixedDepositTransferConfig(configs map[int64]types.FixedDepositCfg, fixed types.FixedDeposit) error {
	cfg, exists := configs[fixed.Term]
	if !exists {
		return fmt.Errorf("deposit cfg not found.fixed.Rate=%s,fixed.Term=%v", fixed.Rate.String(), fixed.Term)
	}
	if cfg.Status == types.RegionFixedDepositCfgStatusInactive {
		return errors.New("fixed deposit cfg status is inactive")
	}
	if !cfg.Rate.Equal(fixed.Rate) {
		return fmt.Errorf("deposit cfg not same.rate=%s,fixed.Rate=%s,exists=%v,fixed.Term=%v", cfg.Rate.String(), fixed.Rate.String(), exists, fixed.Term)
	}
	return nil
}
