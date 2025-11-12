package utils

import (
	"encoding/json"

	"github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var EventDataKey = "data"

func GenEventCompactAttr(eventType string, data interface{}) sdk.Event {
	byt, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	ev := sdk.NewEvent(eventType)
	ev.Attributes = append(ev.Attributes, types.EventAttribute{
		Key:   EventDataKey,
		Value: string(byt),
	})
	return ev
}

func GenEventCompactAttrWithBytes(eventType string, data []byte) sdk.Event {
	ev := sdk.NewEvent(eventType)
	ev.Attributes = append(ev.Attributes, types.EventAttribute{
		Key:   EventDataKey,
		Value: string(data),
	})
	return ev
}
