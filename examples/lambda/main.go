package main

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	indentv1 "go.indent.com/indent-go/api/indent/v1"
	"go.indent.com/indent-go/pkg/common"
	"go.indent.com/webhook-go"
	"go.indent.com/webhook-go/lambda"
)

func main() {
	handler := webhook.ApplyHandler{
		Supported:   []string{common.KindRole},
		ApplyGrant:  handleGrant,
		ApplyRevoke: handleRevoke,
	}
	lambda.ServeApply(handler)
}

func handleGrant(ctx context.Context, req *webhook.ApplyRequest) (*indentv1.ApplyUpdateResponse, error) {
	req.Logger.Info("Handling grant")
	// TODO: add req.Petitioner to req.Target
	return &indentv1.ApplyUpdateResponse{
		Status: &status.Status{Code: int32(codes.OK)},
	}, nil
}

func handleRevoke(ctx context.Context, req *webhook.ApplyRequest) (*indentv1.ApplyUpdateResponse, error) {
	req.Logger.Info("Handling revoke")
	// TODO: remove req.Petitioner from req.Target
	return &indentv1.ApplyUpdateResponse{
		Status: &status.Status{Code: int32(codes.OK)},
	}, nil
}
