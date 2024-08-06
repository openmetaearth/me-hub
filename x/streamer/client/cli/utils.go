package cli

import (
	"fmt"

	"github.com/osmosis-labs/osmosis/v15/osmoutils"
	"github.com/st-chain/me-hub/x/streamer/types"
)

func parseRecords(gaugesRaw, weightsRaw string) ([]types.DistrRecord, error) {
	gaugeIds, err := osmoutils.ParseUint64SliceFromString(gaugesRaw, ",")
	if err != nil {
		return nil, err
	}

	weights, err := osmoutils.ParseSdkIntFromString(weightsRaw, ",")
	if err != nil {
		return nil, err
	}

	if len(gaugeIds) != len(weights) {
		return nil, fmt.Errorf("the length of gauge ids and weights not matched")
	}

	if len(gaugeIds) == 0 {
		return nil, fmt.Errorf("records is empty")
	}

	var records []types.DistrRecord
	for i, gaugeId := range gaugeIds {
		records = append(records, types.DistrRecord{
			GaugeId: gaugeId,
			Weight:  weights[i],
		})
	}
	return records, nil
}
