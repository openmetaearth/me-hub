package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/gogoproto/proto"
	"sigs.k8s.io/yaml"
)

func marshalToYAML(v proto.Message) string {
	bz, err := codec.ProtoMarshalJSON(v, nil)
	if err != nil {
		panic(err)
	}

	out, err := yaml.JSONToYAML(bz)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (v ValidatorV1) String() string {
	return marshalToYAML(&v)
}

func (v ValidatorV2Panic) String() string {
	return marshalToYAML(&v)
}

func (v ValidatorV2) String() string {
	return marshalToYAML(&v)
}
