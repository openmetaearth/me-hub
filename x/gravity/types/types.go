package types

import (
	"bytes"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/openmetaearth/me-hub/app/params"
	"github.com/openmetaearth/me-hub/utils"
	"math"
	"math/big"
	"sort"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// --- ERC20Token --- //
func NewERC20Token(amount sdkmath.Int, contract string) ERC20Token {
	return ERC20Token{Amount: amount, Contract: contract}
}

// ValidateBasic permforms stateless validation
func (m *ERC20Token) ValidateBasic() error {
	if err := utils.ValidateEthereumAddress(m.Contract); err != nil {
		return errorsmod.Wrap(err, "invalid contract address")
	}
	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(ErrInvalid, "amount")
	}
	return nil
}

// --- BRIDGE VALIDATOR(S) --- //

// ValidateBasic performs stateless checks on validity
func (m *BridgeValidator) ValidateBasic() error {
	if m.Power == 0 {
		return errorsmod.Wrap(ErrEmpty, "power")
	}
	if err := utils.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return errorsmod.Wrap(ErrInvalid, "external address")
	}
	return nil
}

// BridgeValidators is the sorted set of validator data for Ethereum bridge MultiSig set
type BridgeValidators []BridgeValidator

func (b BridgeValidators) Len() int {
	return len(b)
}

func (b BridgeValidators) Less(i, j int) bool {
	if b[i].Power == b[j].Power {
		// Secondary sort on eth address in case powers are equal
		return bytes.Compare([]byte(b[i].ExternalAddress), []byte(b[j].ExternalAddress)) == -1
	}
	return b[i].Power > b[j].Power
}

func (b BridgeValidators) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// PowerDiff returns the difference in power between two bridge validator sets
// note this is Gravity bridge power *not* Cosmos voting power. Cosmos voting
// power is based on the absolute number of tokens in the staking pool at any given
// time Gravity bridge power is normalized using the equation.
//
// validators cosmos voting power / total cosmos voting power in this block = gravity bridge power / u32_max
//
// As an example if someone has 52% of the Cosmos voting power when a validator set is created their Gravity
// bridge voting power is u32_max * .52
//
// Normalized voting power dramatically reduces how often we have to produce new validator set updates. For example
// if the total on chain voting power increases by 1% due to inflation, we shouldn't have to generate a new validator
// set, after all the validators retained their relative percentages during inflation and normalized Gravity bridge power
// shows no difference.
func (b BridgeValidators) PowerDiff(c BridgeValidators) float64 {
	powers := map[string]int64{}
	// loop over b and initialize the map with their powers
	for _, bv := range b {
		powers[bv.ExternalAddress] = int64(bv.Power)
	}

	// subtract c powers from powers in the map, initializing
	// uninitialized keys with negative numbers
	for _, bv := range c {
		if val, ok := powers[bv.ExternalAddress]; ok {
			powers[bv.ExternalAddress] = val - int64(bv.Power)
		} else {
			powers[bv.ExternalAddress] = -int64(bv.Power)
		}
	}

	var delta float64
	for _, v := range powers {
		// NOTE: we care about the absolute value of the changes
		delta += math.Abs(float64(v))
	}

	return math.Abs(delta / float64(PowerBase))
}

// TotalPower returns the total power in the bridge validator set
func (b BridgeValidators) TotalPower() (out uint64) {
	for _, v := range b {
		out += v.Power
	}
	return
}

// HasDuplicates returns true if there are duplicates in the set
func (b BridgeValidators) HasDuplicates() bool {
	m := make(map[string]struct{}, len(b))
	for i := range b {
		m[b[i].ExternalAddress] = struct{}{}
	}
	return len(m) != len(b)
}

// GetPowers returns only the power values for all members
func (b BridgeValidators) GetPowers() []uint64 {
	r := make([]uint64, len(b))
	for i := range b {
		r[i] = b[i].Power
	}
	return r
}

// ValidateBasic performs stateless checks
func (b BridgeValidators) ValidateBasic() error {
	if len(b) == 0 {
		return ErrEmpty
	}
	for i := range b {
		if err := b[i].ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "member %d", i)
		}
	}
	if b.HasDuplicates() {
		return errorsmod.Wrap(ErrDuplicate, "address")
	}
	return nil
}

func (b BridgeValidators) Equal(o BridgeValidators) bool {
	if len(b) != len(o) {
		return false
	}

	for i, bv := range b {
		ov := o[i]
		if bv.Power != ov.Power || bv.ExternalAddress != ov.ExternalAddress {
			return false
		}
	}

	return true
}

// --- RelayerSet(S) --- //

func CurrentRelayerSet(nonce, height uint64, members BridgeValidators) *RelayerSet {
	sort.Sort(members)
	return &RelayerSet{
		Nonce:   nonce,
		Members: members,
		Height:  height,
	}
}

// GetCheckpoint returns the checkpoint
func (r *RelayerSet) GetCheckpoint(gravityIDStr string) ([]byte, error) {
	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := utils.StrToByte32(gravityIDStr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}
	checkpoint, err := utils.StrToByte32("checkpoint")
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse checkpoint")
	}

	memberAddresses := make([]gethcommon.Address, len(r.Members))
	convertedPowers := make([]*big.Int, len(r.Members))
	for i, m := range r.Members {
		//memberAddresses[i] = gethcommon.HexToAddress(m.ExternalAddress)
		memberAddresses[i] = toHexAddr(gravityIDStr, m.ExternalAddress)
		convertedPowers[i] = big.NewInt(int64(m.Power))
	}
	// the word 'checkpoint' needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	packBytes, packErr := relayerSetCheckpointABI.Pack("checkpoint", gravityID, checkpoint, big.NewInt(int64(r.Nonce)), memberAddresses, convertedPowers)

	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if packErr != nil {
		return nil, errorsmod.Wrap(packErr, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	hash := crypto.Keccak256Hash(packBytes[4:])
	return hash.Bytes(), nil
}

func (m *RelayerSet) Equal(o *RelayerSet) (bool, error) {
	if m.Height != o.Height {
		return false, errorsmod.Wrap(ErrInvalid, "relayer set heights mismatch")
	}

	if m.Nonce != o.Nonce {
		return false, errorsmod.Wrap(ErrInvalid, "relayer set nonce mismatch")
	}

	if !BridgeValidators(m.Members).Equal(o.Members) {
		return false, errorsmod.Wrap(ErrInvalid, "relayer set members mismatch")
	}

	return true, nil
}

type RelayerSets []*RelayerSet

func (v RelayerSets) Len() int {
	return len(v)
}

func (v RelayerSets) Less(i, j int) bool {
	return v[i].Nonce > v[j].Nonce
}

func (v RelayerSets) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// --- OutgoingTxBatch --- //

type OutgoingTxBatches []*OutgoingTxBatch

func (v OutgoingTxBatches) Len() int {
	return len(v)
}

func (v OutgoingTxBatches) Less(i, j int) bool {
	return v[i].BatchNonce > v[j].BatchNonce
}

func (v OutgoingTxBatches) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// GetFees returns the total fees contained within a given batch
func (m *OutgoingTxBatch) GetFees() sdkmath.Int {
	sum := sdkmath.ZeroInt()
	for _, t := range m.Transactions {
		sum = sum.Add(t.Fee.Amount)
	}
	return sum
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (m *OutgoingTxBatch) GetCheckpoint(gravityIDString string) ([]byte, error) {
	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := utils.StrToByte32(gravityIDString)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}

	// Create the methodName argument which salts the signature
	batchMethodName, err := utils.StrToByte32("transactionBatch")
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse transactionBatch")
	}

	// Run through the elements of the batch and serialize them
	txAmounts := make([]*big.Int, len(m.Transactions))
	txDestinations := make([]gethcommon.Address, len(m.Transactions))
	txFees := make([]*big.Int, len(m.Transactions))
	for i, tx := range m.Transactions {
		txAmounts[i] = tx.Token.Amount.BigInt()
		//txDestinations[i] = gethcommon.HexToAddress(tx.DestAddress)
		txDestinations[i] = toHexAddr(gravityIDString, tx.DestAddress)
		txFees[i] = tx.Fee.Amount.BigInt()
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := outgoingBatchTxCheckpointABI.Pack("submitBatch",
		gravityID,
		batchMethodName,
		txAmounts,
		txDestinations,
		txFees,
		big.NewInt(int64(m.BatchNonce)),
		toHexAddr(gravityIDString, m.TokenContract),
		big.NewInt(int64(m.BatchTimeout)),
		toHexAddr(gravityIDString, m.FeeReceive),
	)
	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, errorsmod.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	return crypto.Keccak256Hash(abiEncodedBatch[4:]).Bytes(), nil
}

func (m *OutgoingTxBatch) TotalAmount() sdk.Int {
	totalAmount := sdk.ZeroInt()
	for _, tx := range m.Transactions {
		totalAmount = totalAmount.Add(tx.Token.Amount)
	}
	return totalAmount
}

// --- Relayer(S) --- //

func (m *Relayer) GetRelayer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.RelayerAddress)
}

func (m *Relayer) GetSlashAmount(slashFraction sdk.Dec) sdk.Coin {
	slashAmount := sdk.NewDecFromInt(m.DelegateAmount).Mul(slashFraction).MulInt64(m.SlashTimes).TruncateInt()
	slashAmount = sdkmath.MinInt(slashAmount, m.DelegateAmount)
	slashAmount = sdkmath.MaxInt(slashAmount, sdkmath.ZeroInt())
	return sdk.NewCoin(params.BaseDenom, slashAmount)
}

func (m *Relayer) GetPower() sdkmath.Int {
	return m.DelegateAmount.Quo(sdk.DefaultPowerReduction)
}

type Relayers []Relayer

func (v Relayers) Len() int {
	return len(v)
}

func (v Relayers) Less(i, j int) bool {
	return v[i].DelegateAmount.Sub(v[j].DelegateAmount).IsPositive()
}

func (v Relayers) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func MinBatchFeeToBaseFees(ms []MinBatchFee) map[string]sdkmath.Int {
	kv := make(map[string]sdkmath.Int, len(ms))
	for _, m := range ms {
		if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
			continue
		}
		kv[m.TokenContract] = m.BaseFee
	}
	return kv
}

type OutgoingTransferTxs []*OutgoingTransferTx

func (bs OutgoingTransferTxs) TotalFee() sdkmath.Int {
	totalFee := sdkmath.NewInt(0)
	for _, tx := range bs {
		totalFee = totalFee.Add(tx.Fee.Amount)
	}
	return totalFee
}

func toHexAddr(gravityId, addr string) gethcommon.Address {
	if strings.Contains(gravityId, "tron") {
		tronAddr, err := address.Base58ToAddress(addr)
		if err != nil {
			return gethcommon.Address{}
		}
		return gethcommon.BytesToAddress(tronAddr.Bytes()[1:])
	}
	return gethcommon.HexToAddress(addr)
}
