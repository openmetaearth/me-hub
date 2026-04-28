package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/app/upgrades/v2_0_13"
	bsctypes "github.com/openmetaearth/me-hub/x/bsc/types"
	trontypes "github.com/openmetaearth/me-hub/x/tron/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/evmos/ethermint/crypto/hd"
	"github.com/spf13/cobra"
)

func GenRelayersCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "gen-relayers [address1,address2,...] [coins]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithKeyringOptions(hd.EthSecp256k1Option())
			cdc := clientCtx.Codec
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("parse coins: %w", err)
			}
			if !coins.IsValid() {
				return fmt.Errorf("invalid coins")
			}

			rawAddrList := strings.Split(args[0], ",")
			var proposalRelayers []string
			var addrs []sdk.AccAddress
			seen := make(map[string]struct{})
			for _, s := range rawAddrList {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				acc, err := sdk.AccAddressFromBech32(s)
				if err != nil {
					return fmt.Errorf("invalid bech32 address %s: %w", s, err)
				}

				addrStr := acc.String()
				if _, ok := seen[addrStr]; ok {
					return fmt.Errorf("duplicate address provided: %s", addrStr)
				}
				seen[addrStr] = struct{}{}
				addrs = append(addrs, acc)
				proposalRelayers = append(proposalRelayers, addrStr)
			}
			if len(addrs) == 0 {
				return fmt.Errorf("no valid addresses provided")
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)
			existingAccs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("unpack accounts: %w", err)
			}

			for _, a := range addrs {
				if existingAccs.Contains(a) {
					return fmt.Errorf("address already exists in genesis: %s", a.String())
				}
			}

			for _, a := range addrs {
				base := authtypes.NewBaseAccount(a, nil, 0, 0)
				if err := base.Validate(); err != nil {
					return fmt.Errorf("validate new account %s: %w", a.String(), err)
				}
				existingAccs = append(existingAccs, base)
			}
			existingAccs = authtypes.SanitizeGenesisAccounts(existingAccs)

			packed, err := authtypes.PackAccounts(existingAccs)
			if err != nil {
				return fmt.Errorf("pack accounts: %w", err)
			}
			authGenState.Accounts = packed
			authBz, err := cdc.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("marshal auth genesis: %w", err)
			}
			appState[authtypes.ModuleName] = authBz

			bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
			for _, a := range addrs {
				bal := banktypes.Balance{
					Address: a.String(),
					Coins:   coins.Sort(),
				}
				bankGenState.Balances = append(bankGenState.Balances, bal)
				bankGenState.Supply = bankGenState.Supply.Add(bal.Coins...)
			}
			delegateAmount := sdk.NewInt(1 * 1e8)
			bondedAmount := delegateAmount.Mul(sdk.NewInt(int64(len(addrs))))

			bal1 := banktypes.Balance{
				Address: authtypes.NewModuleAddress(bsctypes.ModuleName).String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, bondedAmount)).Sort(),
			}
			bankGenState.Balances = append(bankGenState.Balances, bal1)
			bankGenState.Supply = bankGenState.Supply.Add(bal1.Coins...)

			bal2 := banktypes.Balance{
				Address: authtypes.NewModuleAddress(trontypes.ModuleName).String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(params.BaseDenom, bondedAmount)).Sort(),
			}
			bankGenState.Balances = append(bankGenState.Balances, bal2)
			bankGenState.Supply = bankGenState.Supply.Add(bal2.Coins...)

			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

			bankBz, err := cdc.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("marshal bank genesis: %w", err)
			}
			appState[banktypes.ModuleName] = bankBz

			{
				bscGenState := v2_0_13.GenGravityGenesis(0, proposalRelayers, bsctypes.DefaultGenesisState(), delegateAmount, bsctypes.ModuleName)
				bscGenStateBz, err := cdc.MarshalJSON(bscGenState)
				if err != nil {
					return fmt.Errorf("marshal bsc genesis: %w", err)
				}
				appState[bsctypes.ModuleName] = bscGenStateBz
			}

			{
				tronGenState := v2_0_13.GenGravityGenesis(0, proposalRelayers, trontypes.DefaultGenesisState(), delegateAmount, trontypes.ModuleName)
				tronGenStateBz, err := cdc.MarshalJSON(tronGenState)
				if err != nil {
					return fmt.Errorf("marshal tron genesis: %w", err)
				}
				appState[trontypes.ModuleName] = tronGenStateBz
			}
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("marshal app state: %w", err)
			}
			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
