package cmd

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/errors"
	cmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/dao/types"
	didtypes "github.com/st-chain/me-hub/x/did/types"
	kyctypes "github.com/st-chain/me-hub/x/kyc/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
)

func SetDAOCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gentx_DAO --pubkey [DAO_ADDRESS]",
		Short: "create new genesis DAO and DID",
		Args:  cobra.NoArgs,
		//FIXME:
		Long: fmt.Sprintf(`Generate a genesis sequencer, by providing the public key of the sequencer and the rollapp address of the sequencer.
Example:
$ %s gentx \'%s dymint show-sequencer\' --home=/path/to/home/dir --keyring-backend=os --from sequencer-account
	`, version.AppName, version.AppName,
		),

		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// daoAddr := args[0]
			// _, err = sdk.AccAddressFromBech32(daoAddr)
			// if err != nil {
			// 	return fmt.Errorf("invalid DAO account address (%s)", err.Error())
			// }
			pkStr, err := cmd.Flags().GetString(stakingcli.FlagPubKey)
			if err != nil {
				return err
			}

			var pk cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
				return err
			}
			daoAddr, err := bech32.ConvertAndEncode("me", pk.Address())
			if err != nil {
				return err
			}
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			genDoc, err := cmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return errors.Wrap(err, "failed to read genesis doc from file")
			}

			// create the app state
			appGenesisState, err := genutiltypes.GenesisStateFromGenDoc(*genDoc)
			if err != nil {
				return err
			}

			appGenesisState, err = AddDAOToGenesis(clientCtx.Codec, appGenesisState, daoAddr)
			if err != nil {
				return err
			}

			appGenesisState, err = SetGenesisIssuerToGenesis(clientCtx.Codec, appGenesisState, daoAddr, pkStr)
			if err != nil {
				return err
			}
			appState, err := json.MarshalIndent(appGenesisState, "", "  ")
			if err != nil {
				return err
			}

			genDoc.AppState = appState
			err = genutil.ExportGenesisFile(genDoc, config.GenesisFile())

			return err
		},
	}

	cmd.Flags().AddFlagSet(stakingcli.FlagSetPublicKey())
	//cmd.Flags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test|memory)")

	_ = cmd.MarkFlagRequired(stakingcli.FlagPubKey)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func AddDAOToGenesis(
	cdc codec.JSONCodec, appGenesisState map[string]json.RawMessage, addr string,
) (map[string]json.RawMessage, error) {

	var genState types.GenesisState
	cdc.MustUnmarshalJSON(appGenesisState[types.ModuleName], &genState)

	genState.DaoAddresses = types.DaoAddresses{
		GlobalDao:      addr,
		MeidDao:        addr,
		DevOperator:    addr,
		AirdropAddress: addr,
	}
	appGenesisState[types.ModuleName] = cdc.MustMarshalJSON(&genState)

	return appGenesisState, nil
}
func SetGenesisIssuerToGenesis(
	cdc codec.JSONCodec, appGenesisState map[string]json.RawMessage, addr string, pkStr string,
) (map[string]json.RawMessage, error) {

	var genState kyctypes.GenesisState
	cdc.MustUnmarshalJSON(appGenesisState[kyctypes.ModuleName], &genState)

	genState.Issuers = []didtypes.DidInfo{
		{
			Did:    "1000000000000001",
			Pubkey: pkStr,
			Status: didtypes.DID_STATUS_ACTIVE,
		},
	}
	appGenesisState[kyctypes.ModuleName] = cdc.MustMarshalJSON(&genState)

	return appGenesisState, nil
}

// jq '.app_state["dao"]["dao_addresses"]["global_dao"] = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
// jq '.app_state["dao"]["dao_addresses"]["meid_dao"] = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
// jq '.app_state["dao"]["dao_addresses"]["dev_operator"] = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
// jq '.app_state["dao"]["dao_addresses"]["airdrop_address"] = "me139mq752delxv78jvtmwxhasyrycufsvr0mue6u"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
// jq '.app_state["kyc"]["issuer"]["did"] = "1000000000000001"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
// jq '.app_state["kyc"]["issuer"]["pubkey"] = "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"Aggm+J77xeXPyJMOnpdtEu+nmCG/ia9zudrm3kGs722z\"}"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
// jq '.app_state["kyc"]["issuer"]["status"] = "DID_STATUS_ACTIVE"' "$GENESIS_FILE" > "$tmp" && mv "$tmp" "$GENESIS_FILE"
