package cli

import (
	"cosmossdk.io/math"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/st-chain/me-hub/app/params"
	gomath "math"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/st-chain/me-hub/x/wstaking/types"
)

// default values
var (
	DefaultTokens                  = sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	defaultAmount                  = DefaultTokens.String() + sdk.DefaultBondDenom
	defaultCommissionRate          = "0.1"
	defaultCommissionMaxRate       = "0.2"
	defaultCommissionMaxChangeRate = "0.01"
	defaultMinSelfDelegation       = "1"
)

// NewTxCmd returns a root CLI command handler for all x/staking transaction commands.
func NewTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		NewCreateValidatorCmd(),
		NewCreateExperienceNodeCmd(),
		CmdNewRegion(),
		CmdRemoveRegion(),
	)

	return stakingTxCmd
}

// NewCreateValidatorCmd returns a CLI command handler for creating a MsgCreateValidator transaction.
func NewCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create a new validator initialized with a self-delegation to it",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetPublicKey())
	cmd.Flags().AddFlagSet(FlagSetAmount())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())
	cmd.Flags().AddFlagSet(FlagSetCommissionCreate())
	cmd.Flags().AddFlagSet(FlagSetMinSelfDelegation())

	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	_ = cmd.MarkFlagRequired(FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagRegionId)
	return cmd
}

func NewCreateExperienceNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-experience-node",
		Short: "create new validator initialized with a self-delegation to it",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			txf, msgCreateValidator, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			msgCreateRegion := types.NewMsgNewRegion(
				clientCtx.GetFromAddress().String(),
				strings.ToLower(types.ExperienceRegion),
				types.ExperienceRegion,
				msgCreateValidator.ValidatorAddress)
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msgCreateValidator, msgCreateRegion)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetPublicKey())
	cmd.Flags().AddFlagSet(FlagSetAmount())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())
	cmd.Flags().AddFlagSet(FlagSetCommissionCreate())
	//cmd.Flags().AddFlagSet(FlagSetMinSelfStake())

	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")
	cmd.Flags().String(FlagValidatorAddress, "", "validator address(prefix is me)")
	cmd.Flags().String(FlagRegionId, "", "region id")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagPubKey)
	_ = cmd.MarkFlagRequired(FlagMoniker)
	_ = cmd.MarkFlagRequired(FlagValidatorAddress)
	_ = cmd.MarkFlagRequired(FlagRegionId)
	return cmd
}

func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *stakingtypes.MsgCreateValidator, error) {
	fAmount, _ := fs.GetString(FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	globalDao := clientCtx.GetFromAddress()
	pkStr, err := fs.GetString(FlagPubKey)
	if err != nil {
		return txf, nil, err
	}

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(FlagMoniker)
	identity, _ := fs.GetString(FlagIdentity)
	website, _ := fs.GetString(FlagWebsite)
	security, _ := fs.GetString(FlagSecurityContact)
	details, _ := fs.GetString(FlagDetails)
	description := stakingtypes.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)
	regionId, _ := fs.GetString(FlagRegionId)
	description.RegionId = regionId

	// get the initial validator commission parameters
	rateStr, _ := fs.GetString(FlagCommissionRate)
	maxRateStr, _ := fs.GetString(FlagCommissionMaxRate)
	maxChangeRateStr, _ := fs.GetString(FlagCommissionMaxChangeRate)

	commissionRates, err := buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txf, nil, err
	}

	minSelfDelegation := math.NewInt(int64(gomath.Pow10(params.BaseDenomUnit)))

	var pkAny *codectypes.Any
	if pk != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pk); err != nil {
			return txf, nil, err
		}
	}

	validatorAddress, err := fs.GetString(FlagValidatorAddress)

	msg := &stakingtypes.MsgCreateValidator{
		Description:       description,
		DelegatorAddress:  globalDao.String(),
		ValidatorAddress:  sdk.ValAddress(sdk.MustAccAddressFromBech32(validatorAddress)).String(),
		Pubkey:            pkAny,
		Value:             amount,
		Commission:        commissionRates,
		MinSelfDelegation: minSelfDelegation,
	}
	if err := msg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	genOnly, _ := fs.GetBool(flags.FlagGenerateOnly)
	if genOnly {
		ip, _ := fs.GetString(FlagIP)
		p2pPort, _ := fs.GetUint(FlagP2PPort)
		nodeID, _ := fs.GetString(FlagNodeID)

		if nodeID != "" && ip != "" && p2pPort > 0 {
			txf = txf.WithMemo(fmt.Sprintf("%s@%s:%d", nodeID, ip, p2pPort))
		}
	}

	return txf, msg, nil
}

// Return the flagset, particular flags, and a description of defaults
// this is anticipated to be used with the gen-tx
func CreateValidatorMsgFlagSet(ipDefault string) (fs *flag.FlagSet, defaultsDesc string) {
	fsCreateValidator := flag.NewFlagSet("", flag.ContinueOnError)
	fsCreateValidator.String(FlagIP, ipDefault, "The node's public P2P IP")
	fsCreateValidator.Uint(FlagP2PPort, 26656, "The node's public P2P port")
	fsCreateValidator.String(FlagNodeID, "", "The node's NodeID")
	fsCreateValidator.String(FlagMoniker, "", "The validator's (optional) moniker")
	fsCreateValidator.String(FlagWebsite, "", "The validator's (optional) website")
	fsCreateValidator.String(FlagSecurityContact, "", "The validator's (optional) security contact email")
	fsCreateValidator.String(FlagDetails, "", "The validator's (optional) details")
	fsCreateValidator.String(FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	fsCreateValidator.String(FlagRegionId, "", "Region id")
	fsCreateValidator.String(FlagValidatorAddress, "", "Region id")
	fsCreateValidator.AddFlagSet(FlagSetCommissionCreate())
	fsCreateValidator.AddFlagSet(FlagSetMinSelfDelegation())
	fsCreateValidator.AddFlagSet(FlagSetAmount())
	fsCreateValidator.AddFlagSet(FlagSetPublicKey())

	defaultsDesc = fmt.Sprintf(`
	delegation amount:           %s
	commission rate:             %s
	commission max rate:         %s
	commission max change rate:  %s
	minimum self delegation:     %s
`, defaultAmount, defaultCommissionRate,
		defaultCommissionMaxRate, defaultCommissionMaxChangeRate,
		defaultMinSelfDelegation)

	return fsCreateValidator, defaultsDesc
}

type TxCreateValidatorConfig struct {
	ChainID string
	NodeID  string
	Moniker string

	Amount string

	CommissionRate          string
	CommissionMaxRate       string
	CommissionMaxChangeRate string
	MinSelfDelegation       string

	PubKey cryptotypes.PubKey

	IP               string
	P2PPort          uint
	Website          string
	SecurityContact  string
	Details          string
	Identity         string
	RegionId         string
	ValidatorAddress string
}

func PrepareConfigForTxCreateValidator(flagSet *flag.FlagSet, moniker, nodeID, chainID string, valPubKey cryptotypes.PubKey) (TxCreateValidatorConfig, error) {
	c := TxCreateValidatorConfig{}

	ip, err := flagSet.GetString(FlagIP)
	if err != nil {
		return c, err
	}

	if ip == "" {
		_, _ = fmt.Fprintf(os.Stderr, "failed to retrieve an external IP; the tx's memo field will be unset")
	}

	p2pPort, err := flagSet.GetUint(FlagP2PPort)
	if err != nil {
		return c, err
	}

	website, err := flagSet.GetString(FlagWebsite)
	if err != nil {
		return c, err
	}

	securityContact, err := flagSet.GetString(FlagSecurityContact)
	if err != nil {
		return c, err
	}

	details, err := flagSet.GetString(FlagDetails)
	if err != nil {
		return c, err
	}

	identity, err := flagSet.GetString(FlagIdentity)
	if err != nil {
		return c, err
	}

	c.Amount, err = flagSet.GetString(FlagAmount)
	if err != nil {
		return c, err
	}

	c.CommissionRate, err = flagSet.GetString(FlagCommissionRate)
	if err != nil {
		return c, err
	}

	c.CommissionMaxRate, err = flagSet.GetString(FlagCommissionMaxRate)
	if err != nil {
		return c, err
	}

	c.CommissionMaxChangeRate, err = flagSet.GetString(FlagCommissionMaxChangeRate)
	if err != nil {
		return c, err
	}

	c.MinSelfDelegation, err = flagSet.GetString(FlagMinSelfDelegation)
	if err != nil {
		return c, err
	}

	c.RegionId, err = flagSet.GetString(FlagRegionId)
	if err != nil {
		return c, err
	}

	c.ValidatorAddress, err = flagSet.GetString(FlagValidatorAddress)
	if err != nil {
		return c, err
	}

	c.IP = ip
	c.P2PPort = p2pPort
	c.Website = website
	c.SecurityContact = securityContact
	c.Identity = identity
	c.NodeID = nodeID
	c.PubKey = valPubKey
	c.Website = website
	c.SecurityContact = securityContact
	c.Details = details
	c.Identity = identity
	c.ChainID = chainID
	c.Moniker = moniker

	if c.Amount == "" {
		c.Amount = defaultAmount
	}

	if c.CommissionRate == "" {
		c.CommissionRate = defaultCommissionRate
	}

	if c.CommissionMaxRate == "" {
		c.CommissionMaxRate = defaultCommissionMaxRate
	}

	if c.CommissionMaxChangeRate == "" {
		c.CommissionMaxChangeRate = defaultCommissionMaxChangeRate
	}

	if c.MinSelfDelegation == "" {
		c.MinSelfDelegation = defaultMinSelfDelegation
	}

	return c, nil
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(clientCtx client.Context, config TxCreateValidatorConfig, txBldr tx.Factory, generateOnly bool) (tx.Factory, sdk.Msg, error) {
	amounstStr := config.Amount
	amount, err := sdk.ParseCoinNormalized(amounstStr)
	if err != nil {
		return txBldr, nil, err
	}

	globalDao := clientCtx.GetFromAddress()
	description := stakingtypes.NewDescription(
		config.Moniker,
		config.Identity,
		config.Website,
		config.SecurityContact,
		config.Details,
	)
	description.RegionId = config.RegionId

	// get the initial validator commission parameters
	rateStr := config.CommissionRate
	maxRateStr := config.CommissionMaxRate
	maxChangeRateStr := config.CommissionMaxChangeRate
	commissionRates, err := buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txBldr, nil, err
	}

	// get the initial validator min self delegation
	msbStr := config.MinSelfDelegation
	minSelfDelegation, ok := sdk.NewIntFromString(msbStr)

	if !ok {
		return txBldr, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	var pkAny *codectypes.Any
	if config.PubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(config.PubKey); err != nil {
			return txBldr, nil, err
		}
	}

	msg := &stakingtypes.MsgCreateValidator{
		Description:       description,
		DelegatorAddress:  globalDao.String(),
		ValidatorAddress:  sdk.ValAddress(sdk.MustAccAddressFromBech32(config.ValidatorAddress)).String(),
		Pubkey:            pkAny,
		Value:             amount,
		Commission:        commissionRates,
		MinSelfDelegation: minSelfDelegation,
	}
	if err := msg.ValidateBasic(); err != nil {
		return txBldr, nil, err
	}

	if generateOnly {
		ip := config.IP
		p2pPort := config.P2PPort
		nodeID := config.NodeID

		if nodeID != "" && ip != "" && p2pPort > 0 {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:%d", nodeID, ip, p2pPort))
		}
	}

	return txBldr, msg, nil
}
