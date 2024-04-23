package main

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/function-hanoi/input/v1beta1"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/response"
)

func TestRunFunction(t *testing.T) {
	// Helper to create request with JSON input as a protobuf struct
	createReq := func(input *v1beta1.Input) *fnv1beta1.RunFunctionRequest {
		inputMap := map[string]interface{}{
			"discs": input.Discs,
		}
		pbStruct, err := structpb.NewStruct(inputMap)
		if err != nil {
			t.Fatalf("failed to convert input map to protobuf Struct: %v", err)
		}
		return &fnv1beta1.RunFunctionRequest{
			Meta:  &fnv1beta1.RequestMeta{Tag: "test"},
			Input: pbStruct,
		}
	}

	type args struct {
		ctx context.Context
		req *fnv1beta1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1beta1.RunFunctionResponse
		err bool
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ValidInput": {
			reason: "The Function should correctly calculate the moves for 3 discs",
			args: args{
				ctx: context.Background(),
				req: createReq(&v1beta1.Input{Discs: 3}),
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{Tag: "test", Ttl: durationpb.New(response.DefaultTTL)},
					Results: []*fnv1beta1.Result{
						{
							Severity: fnv1beta1.Severity_SEVERITY_NORMAL,
							Message:  `[{"disc":1,"from":"A","to":"C"},{"disc":2,"from":"A","to":"B"},{"disc":1,"from":"C","to":"B"},{"disc":3,"from":"A","to":"C"},{"disc":1,"from":"B","to":"A"},{"disc":2,"from":"B","to":"C"},{"disc":1,"from":"A","to":"C"}]`,
						},
					},
				},
				err: false,
			},
		},
		"InvalidInput": {
			reason: "The Function should return an error for invalid input (negative number of discs)",
			args: args{
				ctx: context.Background(),
				req: createReq(&v1beta1.Input{Discs: -1}),
			},
			want: want{
				err: true,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := &Function{log: logging.NewNopLogger()}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)

			if !tc.want.err && err != nil {
				t.Errorf("%s\nUnexpected error: %v", tc.reason, err)
			}
			if tc.want.err && err == nil {
				t.Errorf("%s\nExpected error but got none", tc.reason)
				t.Logf("Expected error but received none. Return data: %+v", rsp)
			}

			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform()); diff != "" && !tc.want.err {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}
		})
	}
}
