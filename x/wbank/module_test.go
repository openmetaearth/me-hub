package wbank

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/stretchr/testify/require"
	golanggrpc "google.golang.org/grpc"
)

type grpcServerMock struct{}

func (s grpcServerMock) RegisterService(sd *golanggrpc.ServiceDesc, ss interface{}) {}

type configuratorMock struct {
	msgServer                 grpcServerMock
	queryServer               grpcServerMock
	capturedMigrationVersions []uint64
}

func newConfiguratorMock() *configuratorMock {
	msgServer := grpcServerMock{}
	queryServer := grpcServerMock{}

	return &configuratorMock{
		msgServer:   msgServer,
		queryServer: queryServer,
	}
}

func (c *configuratorMock) MsgServer() grpc.Server {
	return c.msgServer
}

func (c *configuratorMock) QueryServer() grpc.Server {
	return c.queryServer
}

func (c *configuratorMock) RegisterMigration(
	moduleName string, forVersion uint64, handler module.MigrationHandler,
) error {
	c.capturedMigrationVersions = append(c.capturedMigrationVersions, forVersion)
	return nil
}

// The test checks the migration registration of the original bank.
// Since we override the "Register Services" we want to be sure that after the update of the SDK,
// The original bank won't have unexpected migrations.
func TestAppModuleOriginalBank_RegisterServices(t *testing.T) {
	bankModule := bank.NewAppModule(&codec.AminoCodec{}, bankkeeper.BaseKeeper{}, keeper.AccountKeeper{}, nil)
	configurator := newConfiguratorMock()
	bankModule.RegisterServices(configurator)
	require.Equal(t, []uint64{1, 2, 3}, configurator.capturedMigrationVersions)
	require.Equal(t, uint64(4), bankModule.ConsensusVersion())
}
