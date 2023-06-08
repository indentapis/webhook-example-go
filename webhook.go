package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auditv1 "go.indent.com/indent-go/api/indent/audit/v1"
	"go.indent.com/indent-go/pkg/common"
)

type Result struct {
	Status *status.Status `json:"status,omitempty"`
}

func HandleRequest(ctx context.Context, req *auditv1.WriteRequest) (Result, error) {
	e := event(req.GetEvents(), common.EventRevoke, common.EventGrant)
	switch e.GetEvent() {
	case common.EventGrant:
		return handleGrant(ctx, e)
	case common.EventRevoke:
		return handleRevoke(ctx, e)
	default:
		return Result{
			Status: status.Newf(codes.InvalidArgument, "unknown event type %s", e.GetEvent()),
		}, nil
	}
}

func handleGrant(ctx context.Context, e *auditv1.Event) (Result, error) {
	// handle grant
	return Result{
		Status: status.New(codes.OK, ""),
	}, nil

}

func handleRevoke(ctx context.Context, e *auditv1.Event) (Result, error) {
	// handle revoke
	return Result{
		Status: status.New(codes.OK, ""),
	}, nil
}

func event(events []*auditv1.Event, eventTypes ...string) *auditv1.Event {
	for _, eventType := range eventTypes {
		for _, e := range events {
			if e.GetEvent() == eventType {
				return e
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
