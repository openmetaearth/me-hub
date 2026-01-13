package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

func (m Migrator) Migrate(_ sdk.Context) error {
	_ = m.keeper
	return nil
}
