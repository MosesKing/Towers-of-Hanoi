package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
	"google.golang.org/protobuf/types/known/structpb"

	// Import the unstructured package
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured"
)

type Move struct {
	Disc int
	From string
	To   string
}

func SolveTowerOfHanoi(n int, from, to, via string, moves *[]Move) {
	if n > 0 {
		SolveTowerOfHanoi(n-1, from, via, to, moves)
		*moves = append(*moves, Move{Disc: n, From: from, To: to})
		SolveTowerOfHanoi(n-1, via, to, from, moves)
	}
}

func main() {
	// Placeholder for function server setup using Crossplane SDK
}

func RunFunction(ctx context.Context, req *v1beta1.RunFunctionRequest) (*v1beta1.RunFunctionResponse, error) {
	log := logging.NewNopLogger()
	rsp := response.To(req, response.DefaultTTL)

	discsValue, exists := req.Input.Fields["discs"]
	if !exists {
		log.Info("Discs field is missing from input")
		return nil, fmt.Errorf("discs field is missing from input")
	}

	discs, err := strconv.Atoi(discsValue.GetStringValue())
	if err != nil {
		log.Info("Error parsing discs input", "error", err)
		return nil, err
	}

	var moves []Move
	SolveTowerOfHanoi(discs, "A", "C", "B", &moves)

	// Create a new WrapperClient
	client := unstructured.NewClient(nil)

	for i, move := range moves {
		moveDescription := fmt.Sprintf("Move disk %d from %s to %s", move.Disc, move.From, move.To)
		cmData := make(map[string]interface{})
		cmData["move"] = moveDescription
		dataStruct, err := structpb.NewStruct(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":   fmt.Sprintf("hanoi-move-%d", i+1),
				"labels": map[string]string{"tower": "hanoi"},
			},
			"data": cmData,
		})
		if err != nil {
			log.Info("Failed to create struct for ConfigMap", "error", err)
			continue
		}

		// Use the WrapperClient to create the object
		obj := &unstructured.Unstructured{Object: dataStruct.AsMap()}
		if err := client.Create(ctx, obj); err != nil {
			log.Info("Failed to create resource", "error", err)
			continue
		}

		// Convert the unstructured.Unstructured object to *resource.Composite
		composite := &resource.Composite{Resource: obj}

		// Add the resource to the response using SetDesiredCompositeResource
		if err := response.SetDesiredCompositeResource(rsp, composite); err != nil {
			log.Info("Failed to set desired composite resource", "error", err)
			continue
		}
	}

	return rsp, nil
}
