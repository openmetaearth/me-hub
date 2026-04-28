package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/openmetaearth/me-hub/x/gravity/types"
)

// RouterKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type RouterKeeper struct {
	router Router
}

// NewRouterKeeper returns a new instance of the cross chain keeper
func NewRouterKeeper(rtr Router) RouterKeeper {
	rtr.Seal()

	return RouterKeeper{
		router: rtr,
	}
}

// Router returns the gov Keeper's Router
func (k RouterKeeper) Router() Router {
	return k.router
}

// Logger returns a module-specific logger.
func (k RouterKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

type ModuleHandler struct {
	QueryServer types.QueryServer
	MsgServer   types.MsgServer
}

func NewModuleHandler(keeper Keeper) *ModuleHandler {
	return &ModuleHandler{
		QueryServer: NewQueryServerImpl(keeper),
		MsgServer:   NewMsgServerImpl(keeper),
	}
}

var _ Router = (*router)(nil)

// Router implements a crosschain Server Handler router.
type Router interface {
	AddRoute(name string, moduleHandler *ModuleHandler) (router Router)
	HasRoute(name string) bool
	GetRoute(name string) (moduleHandler *ModuleHandler)
	Seal()
}

type router struct {
	routes map[string]*ModuleHandler
	sealed bool
}

// NewRouter creates a new Router interface instance
func NewRouter() Router {
	return &router{
		routes: make(map[string]*ModuleHandler),
	}
}

// Seal seals the router which prohibits any subsequent route handlers to be
// added. Seal will panic if called more than once.
func (rtr *router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}
	rtr.sealed = true
}

// AddRoute adds a governance handler for a given path. It returns the Router
// so AddRoute calls can be linked. It will panic if the router is sealed.
func (rtr *router) AddRoute(path string, moduleHandler *ModuleHandler) Router {
	if rtr.sealed {
		panic("router sealed; cannot add route handler")
	}

	if !sdk.IsAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.HasRoute(path) {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = moduleHandler
	return rtr
}

// HasRoute returns true if the router has a path registered or false otherwise.
func (rtr *router) HasRoute(path string) bool {
	return rtr.routes[path] != nil
}

// GetRoute returns a Handler for a given path.
func (rtr *router) GetRoute(path string) *ModuleHandler {
	if !rtr.HasRoute(path) {
		panic(fmt.Sprintf("route \"%s\" does not exist", path))
	}

	return rtr.routes[path]
}
