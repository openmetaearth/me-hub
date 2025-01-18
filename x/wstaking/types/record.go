package types

import "github.com/cosmos/cosmos-sdk/codec"

// MustMarshalRecord returns the delegation bytes. Panics if fails
func MustMarshalRecord(cdc codec.BinaryCodec, record Record) []byte {
	return cdc.MustMarshal(&record)
}

// MustUnmarshalRecord return the unmarshaled delegation from bytes.
// Panics if fails.
func MustUnmarshalRecord(cdc codec.BinaryCodec, value []byte) Record {
	var record Record
	err := cdc.Unmarshal(value, &record)
	if err != nil {
		panic(err)
	}

	return record
}

// MustMarshalReviewRecord returns the delegation bytes. Panics if fails
func MustMarshalReviewRecord(cdc codec.BinaryCodec, rr ReviewRecord) []byte {
	return cdc.MustMarshal(&rr)
}

// MustUnmarshalReviewRecord return the unmarshaled delegation from bytes.
// Panics if fails.
func MustUnmarshalReviewRecord(cdc codec.BinaryCodec, value []byte) ReviewRecord {
	var rr ReviewRecord
	err := cdc.Unmarshal(value, &rr)
	if err != nil {
		panic(err)
	}

	return rr
}
