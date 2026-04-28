package cli

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openmetaearth/me-hub/utils"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" // nolint:staticcheck
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
)

func Debug() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me-debug",
		Short: "Tool for helping with debugging your application",
		RunE:  client.ValidateCmd,
	}
	cmd.AddCommand(
		ToStringCmd(),
		ToBytes32Cmd(),
		ModuleAddressCmd(),
		ChecksumEthAddressCmd(),
		ConvertTxDataToHashCmd(),
		DecodeSimulateTxCmd(),
		VerifyTxCmd(),
		PubkeyCmd(),
		AddrCmd(),
		debug.RawBytesCmd(),
		ConvertTronAddrCmd(),
		ShowAddressCmd(),
	)
	cmd.PersistentFlags().StringP(tmcli.OutputFlag, "o", "json", "Output format (text|json)")
	return cmd
}

func ToStringCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "2str [hex/base64/base58] [data]",
		Short: "Decode to string tools",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var decodeString []byte
			switch args[0] {
			case "hex":
				decodeString, err = hex.DecodeString(strings.TrimPrefix(args[1], "0x"))
				if err != nil {
					return err
				}
			case "base64":
				decodeString, err = base64.StdEncoding.DecodeString(args[1])
				if err != nil {
					return err
				}
			case "base58":
				decodeString, _, err = base58.CheckDecode(args[1])
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid encode type: %s", args[0])
			}
			cmd.Println(string(decodeString))
			return nil
		},
	}
}

func ToBytes32Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "2bytes32",
		Short: "String to bytes32 hex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args[0]) > 32 {
				return fmt.Errorf("input data length greater than 32")
			}
			var byte32 [32]byte
			copy(byte32[:], args[0])
			cmd.Println(hex.EncodeToString(byte32[:]))
			return nil
		},
	}
}

// ModuleAddressCmd
func ModuleAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module-addr <module name>",
		Short: "Get module address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(authtypes.NewModuleAddress(args[0]).String())
			return nil
		},
	}
	return cmd
}

func VerifyTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify-tx [base64TxData]",
		Short:   "Verify tx",
		Example: fmt.Sprintf("%s debug verify-tx 'CucHC...==='", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			sdkTx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
			if err != nil {
				return err
			}

			builder, err := clientCtx.TxConfig.WrapTxBuilder(sdkTx)
			if err != nil {
				return err
			}
			stdTx := builder.GetTx()

			sigTx, ok := sdkTx.(authsigning.SigVerifiableTx)
			if !ok {
				return errors.New("invalid transaction type")
			}
			// stdSigs contains the sequence number, account number, and signatures.
			// When simulating, this would just be a 0-length slice.
			sigs, err := sigTx.GetSignaturesV2()
			if err != nil {
				return fmt.Errorf("get signature error %s", err.Error())
			}
			signerAddrs := sigTx.GetSigners()

			// check that signer length and signature length are the same
			if len(sigs) != len(signerAddrs) {
				return fmt.Errorf("invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
			}
			status, err := clientCtx.Client.Status(cmd.Context())
			if err != nil {
				return err
			}
			chainId := status.NodeInfo.Network
			queryClient := authtypes.NewQueryClient(clientCtx)
			for i, sig := range sigs {
				accountResponse, err := queryClient.Account(cmd.Context(), &authtypes.QueryAccountRequest{Address: signerAddrs[i].String()})
				if err != nil {
					return err
				}
				var acc authtypes.AccountI
				err = clientCtx.InterfaceRegistry.UnpackAny(accountResponse.GetAccount(), &acc)
				if err != nil {
					return err
				}
				// retrieve pubkey
				pubKey := acc.GetPubKey()
				sequence := sig.Sequence
				signerData := authsigning.SignerData{
					ChainID:       chainId,
					AccountNumber: acc.GetAccountNumber(),
					Sequence:      sequence,
				}

				bz := legacytx.StdSignBytes(
					chainId, acc.GetAccountNumber(), sequence, stdTx.GetTimeoutHeight(),
					legacytx.StdFee{Amount: stdTx.GetFee(), Gas: stdTx.GetGas()},
					sdkTx.GetMsgs(), stdTx.GetMemo(), nil,
				)
				if err = clientCtx.PrintString(string(bz) + "\n"); err != nil {
					return err
				}

				if err = authsigning.VerifySignature(pubKey, signerData, sig.Data, clientCtx.TxConfig.SignModeHandler(), sdkTx); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}

func ConvertTxDataToHashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx-hash [base64TxData]",
		Short:   "Convert base64 tx data to txHash",
		Example: fmt.Sprintf("%s debug tx-hash 'CucHC...==='", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBytes, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}
			hashBytes := sha256.Sum256(txBytes)
			cmd.Println(fmt.Sprintf("%X", hashBytes))
			return nil
		},
	}
	return cmd
}

func DecodeSimulateTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulate-tx [base64/hex]",
		Short: "decode base64 or hex tx data to json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var txBytes []byte
			if useHex, _ := cmd.Flags().GetBool("hex"); useHex {
				txBytes, err = hexutil.Decode(args[0])
			} else {
				txBytes, err = base64.StdEncoding.DecodeString(args[0])
			}
			if err != nil {
				return err
			}

			simulateReq := new(tx.SimulateRequest)
			if err := proto.Unmarshal(txBytes, simulateReq); err != nil {
				return err
			}
			return clientCtx.PrintProto(simulateReq)
		},
	}
	cmd.Flags().BoolP("hex", "x", false, "Treat input as hexadecimal instead of base64")
	return cmd
}

func ChecksumEthAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checksum [eth address]",
		Short: "Checksum eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			if !gethcommon.IsHexAddress(args[0]) {
				return fmt.Errorf("not hex address: %s", args[0])
			}
			return clientCtx.PrintString(fmt.Sprintf("%s\n", gethcommon.HexToAddress(args[0]).Hex()))
		},
	}
	return cmd
}

func PubkeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pubkey [pubkey]",
		Short: "Decode a pubkey from proto JSON",
		Example: fmt.Sprintf(`
$ %s debug pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AurroA7jvfPd1AadmmOvWM2rJSwipXfRf8yD6pLbA2DJ"}'
$ %s debug pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"eKlxn6Xoe9LNmD53omoNQrVrws5KT73hfmqeCSqL87A="}'
`, version.AppName, version.AppName),
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			var pubkey cryptotypes.PubKey
			if len(args) == 0 {
				serverCtx := server.GetServerContextFromCmd(cmd)
				serverCfg := serverCtx.Config
				privValidator := privval.LoadFilePV(serverCfg.PrivValidatorKeyFile(), serverCfg.PrivValidatorStateFile())
				valPubKey, err := privValidator.GetPubKey()
				if err != nil {
					return err
				}
				pubkey, err = cryptocodec.FromTmPubKeyInterface(valPubKey)
				if err != nil {
					return err
				}
			} else {
				if err = clientCtx.Codec.UnmarshalInterfaceJSON([]byte(args[0]), &pubkey); err != nil {
					if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.ConsPK, args[0]); err == nil { // nolint:staticcheck
					} else if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.AccPK, args[0]); err == nil { // nolint:staticcheck
					} else if pubkey, err = legacybech32.UnmarshalPubKey(legacybech32.ValPK, args[0]); err == nil { // nolint:staticcheck
					} else {
						return fmt.Errorf("pubkey '%s' invalid", args[0])
					}
				}
			}
			pubkeyJson, err := clientCtx.Codec.MarshalInterfaceJSON(pubkey)
			if err != nil {
				return err
			}
			var data []byte
			switch pubkey.Type() {
			case "ed25519":
				data, err = json.MarshalIndent(map[string]interface{}{
					"address":        strings.ToUpper(hex.EncodeToString(pubkey.Address().Bytes())),
					"val_cons_pub":   json.RawMessage(pubkeyJson),
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			case "secp256k1":
				data, err = json.MarshalIndent(map[string]interface{}{
					"acc_address":    sdk.AccAddress(pubkey.Address().Bytes()).String(),
					"val_address":    sdk.ValAddress(pubkey.Address().Bytes()).String(),
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			case "eth_secp256k1":
				data, err = json.MarshalIndent(map[string]interface{}{
					"eip55_address":  gethcommon.BytesToAddress(pubkey.Address()).String(),
					"acc_address":    sdk.AccAddress(pubkey.Address().Bytes()).String(),
					"val_address":    sdk.ValAddress(pubkey.Address().Bytes()).String(),
					"pub_key_hex":    hex.EncodeToString(pubkey.Bytes()),
					"pub_key_base64": base64.StdEncoding.EncodeToString(pubkey.Bytes()),
				}, "", "  ")
			default:
				return fmt.Errorf("invalid public key type: %s", pubkey.Type())
			}
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(data))
		},
	}
}

func AddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr [address]",
		Short: "Convert an address between hex and bech32",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			bech32prefix, err := cmd.Flags().GetString("prefix")
			if err != nil {
				return err
			}
			// try hex, then bech32
			addrString := args[0]
			var addr []byte
			addr, err = hexutil.Decode(addrString)
			if err != nil {
				_, addr, err = bech32.DecodeAndConvert(addrString)
				if err != nil {
					return errors.New("expected hex or bech32")
				}
			}

			meAddress, err := bech32.ConvertAndEncode(bech32prefix, addr)
			if err != nil {
				return err
			}

			// Tron address prefix is 0x41
			const tronAddressPrefix = 0x41
			tronAddress := base58.CheckEncode(addr, tronAddressPrefix)

			raw, err := json.Marshal(map[string]interface{}{
				"base64": addr,
				"hex":    utils.ToChecksummed(addr),
				"bech32": meAddress,
				"tron":   tronAddress,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintRaw(raw)
		},
	}
	cmd.Flags().StringP("prefix", "p", "me", "Bech32 Prefix to encode to")
	return cmd
}

func ConvertTronAddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tron-addr-convert [tron-base58-address]",
		Short: "Convert a Tron Base58 address to its corresponding Ethereum and Cosmos addresses",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			tronAddr := args[0]

			// Tron addresses are Base58Check encoded. The decoded payload is 21 bytes:
			// 1-byte prefix (0x41) + 20-byte address.
			// The btcsuite/btcutil/base58 library's CheckDecode function handles this by
			// returning the 20-byte address as the payload and the 0x41 prefix as the version.
			addressBytes, version, err := base58.CheckDecode(tronAddr)
			if err != nil {
				return fmt.Errorf("failed to decode base58 check address: %w", err)
			}

			// Validate the Tron address prefix (0x41) and length (20 bytes for the address itself).
			const tronAddressPrefix = 0x41
			if version != tronAddressPrefix || len(addressBytes) != 20 {
				return fmt.Errorf("invalid tron address format: incorrect version or length")
			}

			// Convert to Ethereum address
			ethAddress := gethcommon.BytesToAddress(addressBytes)

			// Convert to Cosmos address (using the default 'me' prefix from the context)
			cosmosAddress, err := bech32.ConvertAndEncode("me", addressBytes)
			if err != nil {
				return fmt.Errorf("failed to encode cosmos address: %w", err)
			}

			output := map[string]string{
				"tron_address": tronAddr,
				"eth_address":  ethAddress.Hex(), // Checksummed address
				"me_address":   cosmosAddress,
			}

			result, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return err
			}

			return clientCtx.PrintString(string(result))
		},
	}
	return cmd
}

func ShowAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address [private-key-hex]",
		Short: "Derive public key and addresses (Tron, ME, ETH) from a private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			// Decode the private key from hex
			privKey, err := crypto.HexToECDSA(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode private key: %w", err)
			}

			// Get the public key
			pubKey := privKey.PublicKey
			pubKeyBytes := crypto.FromECDSAPub(&pubKey)

			// Get the address from the public key
			addrBytes := crypto.PubkeyToAddress(pubKey).Bytes()

			// Convert to ME address
			meAddress, err := bech32.ConvertAndEncode("me", addrBytes)
			if err != nil {
				return fmt.Errorf("failed to encode ME address: %w", err)
			}

			// Convert to Tron address
			const tronAddressPrefix = 0x41
			tronAddress := base58.CheckEncode(addrBytes, tronAddressPrefix)

			// Convert to Ethereum address
			ethAddress := gethcommon.BytesToAddress(addrBytes)

			output := map[string]string{
				"public_key_hex": hex.EncodeToString(pubKeyBytes),
				"eth_address":    ethAddress.Hex(),
				"me_address":     meAddress,
				"tron_address":   tronAddress,
			}

			result, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return err
			}

			return clientCtx.PrintString(string(result))
		},
	}
	return cmd
}
