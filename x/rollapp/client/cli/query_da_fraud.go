package cli

/*
#cgo LDFLAGS: -L../../../../cgoCelestia -lcgoCelestia -Wl,-rpath,number
#cgo CFLAGS: -I../../../../cgoCelestia
#include "cgoCelestialib.h"
*/
import "C"
import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/st-chain/me-hub/x/rollapp/types"
	"unsafe"

	//"plugin"
	"strconv"
)

func ConvertString(goStr string) *C.char {
	cStr := C.CString(goStr)
	//	defer C.free(unsafe.Pointer(cStr))
	return cStr
}

func CmdGetSubmitBlockDaInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSubmitBlockDaInfo [rollappID] [startHeight] [numberBlocks]",
		Short: "get submit block and da data",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rollappId := args[0]

			startHeight, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			numberBlocks, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			argRollappId := args[0]
			if argRollappId == "" {
				return fmt.Errorf("rollappID can not be empty")
			}
			req := &types.MsgGetBlockDaInfoRequest{
				Creator:      clientCtx.GetFromAddress().String(),
				RollappId:    rollappId,
				StartHeight:  startHeight,
				NumberBlocks: uint32(numberBlocks),
			}

			res, err := queryClient.GetSubmitBlockDaInfo(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmGetAppendingDaFraudChallenge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getAppendingDaFraudChallenge",
		Short: "get appending da fraud challenge",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			req := &types.MsgGetDaFraudChallengeRequest{Creator: clientCtx.GetFromAddress().String()}
			res, err := queryClient.GetAppendingDaFraudChallenge(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdVerifyCommitmentProof() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verifyCommitmentProof [commitmentProof] [daRoot] [namespace]",
		Short: "verifyCommitmentProof",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			/*
				commitProof, err := hex.DecodeString(args[0])
				if err != nil {
					return err
				}

				daRoot, err := hex.DecodeString(args[1])
				if err != nil {
					return err
				}

				namespace, err := hex.DecodeString(args[2])
				if err != nil {
					return err
				}

			*/

			/*
					p, err := plugin.Open("celestiaPlugin.so")
					if err != nil {
						return fmt.Errorf("load plugin error.err = %s", err.Error())
					}
					pVal, err := p.Lookup("VerifyDACommitmentProof")
					if err != nil {
						return fmt.Errorf(" Lookup function error.err = %s", err.Error())
					}
					daVerify, ok := pVal.(func([]byte, []byte, []byte) (int, error))


				if !ok {
					return fmt.Errorf("plugin's verify function type is not match")
				}

			*/
			pCmtProof_c := ConvertString(args[0])
			pDaRoot_c := ConvertString(args[1])
			pNamespace_c := ConvertString(args[2])
			res := C.VerifyDACommitmentProof_c(pCmtProof_c, pDaRoot_c, pNamespace_c)
			strErr := C.GoString(res.Err)
			status := int(res.Status)

			defer func() { //释放对应内存
				C.free(unsafe.Pointer(res.Err))
				C.free(unsafe.Pointer(pCmtProof_c))
				C.free(unsafe.Pointer(pDaRoot_c))
				C.free(unsafe.Pointer(pNamespace_c))
			}()
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("verify result = %d,err = %s", status, strErr))

		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
