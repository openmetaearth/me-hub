package types

import didtypes "github.com/openmetaearth/me-hub/x/did/types"

type QuerySubAccountDid struct {
	SubAccount string `protobuf:"bytes,1,opt,name=sub_account,json=subAccount,proto3" json:"sub_account,omitempty"`
}

func (m *QuerySubAccountDid) Reset()         { *m = QuerySubAccountDid{} }
func (m *QuerySubAccountDid) String() string { return m.SubAccount }
func (*QuerySubAccountDid) ProtoMessage()    {}

type QuerySubAccountDidResponse struct {
	Info didtypes.DidInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info"`
}

func (m *QuerySubAccountDidResponse) Reset()         { *m = QuerySubAccountDidResponse{} }
func (m *QuerySubAccountDidResponse) String() string { return "" }
func (*QuerySubAccountDidResponse) ProtoMessage()    {}
