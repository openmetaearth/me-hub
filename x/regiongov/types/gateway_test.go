package types

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	gateway "github.com/cosmos/gogogateway"
	"github.com/cosmos/gogoproto/proto"
	legacyproto "github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	stdproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	regiongovapi "github.com/openmetaearth/me-hub/api/metaearth/regiongov"
)

func TestStandardAPIWireCompatibility(t *testing.T) {
	for _, tc := range []struct {
		name     string
		standard stdproto.Message
		gogo     proto.Message
		json     string
	}{
		{
			name:     "proposal",
			standard: &regiongovapi.Proposal{}, gogo: &Proposal{},
			json: `{"id":"7","kind":"PROPOSAL_KIND_TREASURY_WITHDRAW","status":"PROPOSAL_STATUS_PASSED",` +
				`"submitTime":"2026-11-11T00:00:00Z","votingEndTime":"2026-11-11T01:00:00Z","votingPeriod":"3600",` +
				`"electorate":"ELECTORATE_KIND_REGION_EQUITY","daoVote":"DAO_VOTE_YES",` +
				`"treasury":{"withdrawals":[{"regionId":"chn","receiver":"receiver","amount":[{"denom":"umec","amount":"10"}]}]}}`,
		},
		{
			name:     "pending action",
			standard: &regiongovapi.PendingAction{}, gogo: &PendingAction{},
			json: `{"id":"8","kind":"PENDING_ACTION_KIND_FREEZE_EQUITY","status":"PENDING_ACTION_STATUS_APPROVED",` +
				`"regionId":"chn","addresses":["equity"],"submitTime":"2026-11-11T00:00:00Z"}`,
		},
		{
			name:     "vote transaction",
			standard: &regiongovapi.MsgVote{}, gogo: &MsgVote{},
			json: `{"voter":"voter","proposalId":"7","option":"VOTE_OPTION_YES"}`,
		},
		{
			name:     "proposal transaction",
			standard: &regiongovapi.MsgSubmitProposal{}, gogo: &MsgSubmitProposal{},
			json: `{"proposer":"proposer","kind":"PROPOSAL_KIND_UPDATE_MIN_VOTING_PERIOD","votingPeriod":"120",` +
				`"updateMinVotingPeriod":{"minVotingPeriod":"300"}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, protojson.Unmarshal([]byte(tc.json), tc.standard))
			wire, err := stdproto.Marshal(tc.standard)
			require.NoError(t, err)
			require.NoError(t, proto.Unmarshal(wire, tc.gogo))
			wire, err = proto.Marshal(tc.gogo)
			require.NoError(t, err)
			actual := tc.standard.ProtoReflect().New().Interface()
			require.NoError(t, stdproto.Unmarshal(wire, actual))
			require.True(t, stdproto.Equal(tc.standard, actual))
			body, err := protojson.Marshal(actual)
			require.NoError(t, err)
			require.JSONEq(t, tc.json, string(body))
		})
	}
}

func TestStandardAPIDescriptors(t *testing.T) {
	// Discover files from the gogo registry so a new regiongov proto file or enum
	// cannot silently be omitted from standard API generation.
	for path, compressed := range proto.AllFileDescriptors() {
		if !strings.HasPrefix(path, "metaearth/regiongov/") {
			continue
		}
		t.Run(path, func(t *testing.T) {
			reader, err := gzip.NewReader(bytes.NewReader(compressed))
			require.NoError(t, err)
			raw, err := io.ReadAll(reader)
			require.NoError(t, err)
			require.NoError(t, reader.Close())
			var original descriptorpb.FileDescriptorProto
			require.NoError(t, stdproto.Unmarshal(raw, &original))

			file, err := protoregistry.GlobalFiles.FindFileByPath(path)
			require.NoError(t, err, "run make proto-gen-api")
			standard := protodesc.ToFileDescriptorProto(file)
			// Buf managed mode changes the Go import path and other file options,
			// but must not change the protobuf messages, fields, enums or services.
			original.Options, standard.Options = nil, nil
			require.True(t, stdproto.Equal(&original, standard), "API schema differs from gogo schema; run make proto-gen-api")

			for _, enum := range original.EnumType {
				fullName := file.Package().Append(protoreflect.Name(enum.GetName()))
				registered, err := protoregistry.GlobalTypes.FindEnumByName(fullName)
				require.NoError(t, err)
				require.Equal(t, fullName, registered.Descriptor().FullName())
				want := make(map[string]int32, len(enum.Value))
				for _, value := range enum.Value {
					want[value.GetName()] = value.GetNumber()
				}
				require.Equal(t, want, legacyproto.EnumValueMap(string(fullName)), "Gateway must discover every enum")
			}
		})
	}
}

func TestGatewayAllEnumFields(t *testing.T) {
	for _, tc := range []struct {
		message proto.Message
		field   string
		enum    string
	}{
		{&Proposal{}, "kind", "ProposalKind"},
		{&Proposal{}, "status", "ProposalStatus"},
		{&Proposal{}, "electorate", "ElectorateKind"},
		{&Proposal{}, "dao_vote", "DaoVote"},
		{&Vote{}, "option", "VoteOption"},
		{&PendingAction{}, "kind", "PendingActionKind"},
		{&PendingAction{}, "status", "PendingActionStatus"},
	} {
		for name, number := range proto.EnumValueMap("metaearth.regiongov." + tc.enum) {
			for _, input := range []string{name, strconv.Itoa(int(number))} {
				t.Run(tc.enum+"/"+input, func(t *testing.T) {
					message := proto.Clone(tc.message)
					require.NoError(t, runtime.PopulateFieldFromPath(message, tc.field, input))
					body, err := (&gateway.JSONPb{EmitDefaults: true, OrigName: true}).Marshal(message)
					require.NoError(t, err)
					var fields map[string]json.RawMessage
					require.NoError(t, json.Unmarshal(body, &fields))
					require.JSONEq(t, strconv.Quote(name), string(fields[tc.field]))
				})
			}
		}
	}
}

func TestGatewayQueryEnums(t *testing.T) {
	type testCase struct {
		name      string
		path      string
		want      proto.Message
		wantError string
	}
	tests := []testCase{
		{
			name: "reported proposal query",
			path: "/metaearth/regiongov/proposals?status=PROPOSAL_STATUS_VOTING&kind=PROPOSAL_KIND_UNSPECIFIED",
			want: &QueryProposalsRequest{Status: ProposalStatusVoting},
		},
		{
			name: "numeric proposal query",
			path: "/metaearth/regiongov/proposals?status=1&kind=0",
			want: &QueryProposalsRequest{Status: ProposalStatusVoting},
		},
		{
			name: "named status and numeric kind",
			path: "/metaearth/regiongov/proposals?status=PROPOSAL_STATUS_PASSED&kind=2",
			want: &QueryProposalsRequest{Status: ProposalStatusPassed, Kind: ProposalKindReplaceRegionOperator},
		},
		{
			name: "numeric status and named kind",
			path: "/metaearth/regiongov/proposals?status=4&kind=PROPOSAL_KIND_UPDATE_MIN_VOTING_PERIOD",
			want: &QueryProposalsRequest{Status: ProposalStatusFailed, Kind: ProposalKindUpdateMinVotingPeriod},
		},
		{
			name: "omitted proposal filters",
			path: "/metaearth/regiongov/proposals",
			want: &QueryProposalsRequest{},
		},
		{
			name: "proposal pagination",
			path: "/metaearth/regiongov/proposals?status=PROPOSAL_STATUS_VOTING&kind=PROPOSAL_KIND_TREASURY_WITHDRAW&pagination.offset=2&pagination.limit=10&pagination.count_total=true&pagination.reverse=true",
			want: &QueryProposalsRequest{
				Status: ProposalStatusVoting, Kind: ProposalKindTreasuryWithdraw,
				Pagination: &query.PageRequest{Offset: 2, Limit: 10, CountTotal: true, Reverse: true},
			},
		},
		{
			name: "omitted pending action filter",
			path: "/metaearth/regiongov/pending_actions",
			want: &QueryPendingActionsRequest{},
		},
		{
			name: "pending action pagination",
			path: "/metaearth/regiongov/pending_actions?status=PENDING_ACTION_STATUS_APPROVED&pagination.key=AQI%3D&pagination.limit=5",
			want: &QueryPendingActionsRequest{
				Status:     PendingActionStatusApproved,
				Pagination: &query.PageRequest{Key: []byte{1, 2}, Limit: 5},
			},
		},
	}

	// Exercise every enum value in both its Swagger name and numeric form.
	for value, name := range ProposalStatus_name {
		for _, input := range []string{name, strconv.Itoa(int(value))} {
			tests = append(tests, testCase{
				name: "proposal status " + input,
				path: "/metaearth/regiongov/proposals?status=" + input,
				want: &QueryProposalsRequest{Status: ProposalStatus(value)},
			})
		}
	}
	for value, name := range ProposalKind_name {
		for _, input := range []string{name, strconv.Itoa(int(value))} {
			tests = append(tests, testCase{
				name: "proposal kind " + input,
				path: "/metaearth/regiongov/proposals?kind=" + input,
				want: &QueryProposalsRequest{Kind: ProposalKind(value)},
			})
		}
	}
	for value, name := range PendingActionStatus_name {
		for _, input := range []string{name, strconv.Itoa(int(value))} {
			tests = append(tests, testCase{
				name: "pending action status " + input,
				path: "/metaearth/regiongov/pending_actions?status=" + input,
				want: &QueryPendingActionsRequest{Status: PendingActionStatus(value)},
			})
		}
	}
	for _, field := range []struct {
		path  string
		param string
	}{
		{"/metaearth/regiongov/proposals", "status"},
		{"/metaearth/regiongov/proposals", "kind"},
		{"/metaearth/regiongov/pending_actions", "status"},
	} {
		for _, input := range []string{"NOT_AN_ENUM", "99", "-1"} {
			tests = append(tests, testCase{
				name:      field.path + " invalid " + field.param + " " + input,
				path:      field.path + "?" + field.param + "=" + input,
				wantError: input + " is not a valid",
			})
		}
	}

	for _, registration := range []struct {
		name     string
		register func(context.Context, *runtime.ServeMux, QueryServer) error
	}{
		{"server", RegisterQueryHandlerServer},
		{"client", func(ctx context.Context, mux *runtime.ServeMux, server QueryServer) error {
			return RegisterQueryHandlerClient(ctx, mux, gatewayQueryClient{server: server})
		}},
	} {
		t.Run(registration.name, func(t *testing.T) {
			for _, tc := range tests {
				t.Run(tc.name, func(t *testing.T) {
					server := &gatewayQueryServer{}
					mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &gateway.JSONPb{
						EmitDefaults: true, OrigName: true,
					}))
					require.NoError(t, registration.register(context.Background(), mux, server))
					response := httptest.NewRecorder()
					mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tc.path, nil))

					if tc.wantError != "" {
						require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
						var body struct {
							Code    int    `json:"code"`
							Message string `json:"message"`
						}
						require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
						require.Equal(t, int(codes.InvalidArgument), body.Code)
						require.Contains(t, body.Message, tc.wantError)
						require.Nil(t, server.request, "invalid filters must not reach the query server")
						return
					}

					require.Equal(t, http.StatusOK, response.Code, response.Body.String())
					require.True(t, json.Valid(response.Body.Bytes()), response.Body.String())
					require.Equal(t, tc.want, server.request)
				})
			}
		})
	}
}

type gatewayQueryServer struct {
	UnimplementedQueryServer
	request proto.Message
}

func (s *gatewayQueryServer) Proposals(_ context.Context, req *QueryProposalsRequest) (*QueryProposalsResponse, error) {
	s.request = req
	return &QueryProposalsResponse{}, nil
}

func (s *gatewayQueryServer) PendingActions(_ context.Context, req *QueryPendingActionsRequest) (*QueryPendingActionsResponse, error) {
	s.request = req
	return &QueryPendingActionsResponse{}, nil
}

type gatewayQueryClient struct {
	QueryClient
	server QueryServer
}

func (c gatewayQueryClient) Proposals(ctx context.Context, req *QueryProposalsRequest, _ ...grpc.CallOption) (*QueryProposalsResponse, error) {
	return c.server.Proposals(ctx, req)
}

func (c gatewayQueryClient) PendingActions(ctx context.Context, req *QueryPendingActionsRequest, _ ...grpc.CallOption) (*QueryPendingActionsResponse, error) {
	return c.server.PendingActions(ctx, req)
}
