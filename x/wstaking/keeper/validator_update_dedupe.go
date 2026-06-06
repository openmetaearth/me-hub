package keeper

import (
	"encoding/hex"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
)

func deduplicateValidatorUpdates(validatorUpdates []abci.ValidatorUpdate) []abci.ValidatorUpdate {
	if len(validatorUpdates) < 2 {
		return validatorUpdates
	}

	indexByPubKey := make(map[string]int, len(validatorUpdates))
	deduplicated := make([]abci.ValidatorUpdate, 0, len(validatorUpdates))
	for _, update := range validatorUpdates {
		pubKey := validatorUpdatePubKeyID(update)
		if index, found := indexByPubKey[pubKey]; found {
			deduplicated[index] = update
			continue
		}

		indexByPubKey[pubKey] = len(deduplicated)
		deduplicated = append(deduplicated, update)
	}
	return deduplicated
}

func validatorUpdatePubKeyID(update abci.ValidatorUpdate) string {
	if pubKey := update.PubKey.GetEd25519(); len(pubKey) > 0 {
		return "ed25519:" + hex.EncodeToString(pubKey)
	}
	if pubKey := update.PubKey.GetSecp256K1(); len(pubKey) > 0 {
		return "secp256k1:" + hex.EncodeToString(pubKey)
	}
	return fmt.Sprintf("%v", update.PubKey)
}
