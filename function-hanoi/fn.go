package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/function-sdk-go/logging"
	"github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/response"
	"google.golang.org/protobuf/types/known/structpb"
)

type Move struct {
	Disc int    `json:"disc"`
	From string `json:"from"`
	To   string `json:"to"`
}

func SolveTowerOfHanoi(n int, from, to, via string, moves *[]Move) {
	if n > 0 {
		SolveTowerOfHanoi(n-1, from, via, to, moves)
		*moves = append(*moves, Move{Disc: n, From: from, To: to})
		SolveTowerOfHanoi(n-1, via, to, from, moves)
	}
}

type Function struct {
	v1beta1.UnimplementedFunctionRunnerServiceServer
	log logging.Logger
}

func (f *Function) RunFunction(ctx context.Context, req *v1beta1.RunFunctionRequest) (*v1beta1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())
	rsp := response.To(req, response.DefaultTTL)

	discsValue, exists := req.Input.Fields["discs"]
	if !exists {
		f.log.Info("Discs field is missing from input")
		return nil, errors.Errorf("discs field is missing from input")
	}

	var discs int
	var err error

	switch v := discsValue.GetKind().(type) {
	case *structpb.Value_NumberValue:
		discs = int(v.NumberValue)
	case *structpb.Value_StringValue:
		discs, err = strconv.Atoi(v.StringValue)
		if err != nil {
			f.log.Info("Error parsing discs field from string to integer", "error", err)
			return nil, fmt.Errorf("discs field is not a valid integer: %v", err)
		}
	default:
		f.log.Info("Discs field is of an unexpected type")
		return nil, errors.Errorf("unexpected type for discs field: %T", discsValue.GetKind())
	}

	var moves []Move
	SolveTowerOfHanoi(discs, "A", "C", "B", &moves)

	resources := make([]*v1beta1.Resource, len(moves))
	for i, move := range moves {
		moveDescription := fmt.Sprintf("Move disk %d from %s to %s", move.Disc, move.From, move.To)
		// Use Struct to create the ConfigMap equivalent structure
		cm := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"apiVersion": structpb.NewStringValue("v1"),
				"kind":       structpb.NewStringValue("ConfigMap"),
				"metadata": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"name": structpb.NewStringValue(fmt.Sprintf("tower-hanoi-move-%d", i+1)),
					},
				}),
				"data": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"move": structpb.NewStringValue(moveDescription),
					},
				}),
			},
		}
		resources[i] = &v1beta1.Resource{Spec: cm}
	}

	rsp.Resources = resources
	return rsp, nil
}
