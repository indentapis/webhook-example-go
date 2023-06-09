package webhook

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	auditv1 "go.indent.com/indent-go/api/indent/audit/v1"
	indentv1 "go.indent.com/indent-go/api/indent/v1"
	"go.indent.com/indent-go/pkg/common"
)

const (
	// numPetitionEventResources is the number of Resources in an Event for a Petition.
	numPetitionEventResources = 2
)

type ApplyRequest struct {
	Events     []*auditv1.Event
	Event      *auditv1.Event
	Petitioner *auditv1.Resource
	Target     *auditv1.Resource

	Logger *zap.Logger
}

type ApplyHandler struct {
	Supported               []string
	ApplyGrant, ApplyRevoke func(ctx context.Context, req *ApplyRequest) (*indentv1.ApplyUpdateResponse, error)
}

func ApplyUpdate(ctx context.Context, handler ApplyHandler, req *auditv1.WriteRequest) (*indentv1.ApplyUpdateResponse, error) {
	supportedField := zap.Strings("supported", handler.Supported)
	applyReq, err := applyRequest(Logger().With(supportedField), handler.Supported, req)
	if err != nil {
		return &indentv1.ApplyUpdateResponse{
			Status: &status.Status{
				Code:    int32(codes.InvalidArgument),
				Message: err.Error(),
			},
		}, nil
	}

	applyReq.Logger.Info("Performing apply", zap.String("event", applyReq.Event.GetEvent()))
	switch applyReq.Event.GetEvent() {
	case common.EventGrant:
		return handler.ApplyGrant(ctx, applyReq)
	case common.EventRevoke:
		return handler.ApplyRevoke(ctx, applyReq)
	default:
		return &indentv1.ApplyUpdateResponse{
			Status: &status.Status{
				Code:    int32(codes.InvalidArgument),
				Message: "invalid event type",
			},
		}, nil
	}
}

func applyRequest(logger *zap.Logger, supported []string, req *auditv1.WriteRequest) (*ApplyRequest, error) {
	e := applyEvent(req.GetEvents())
	if len(e.GetResources()) != numPetitionEventResources {
		return nil, fmt.Errorf("expected %d resources in event, got %d", numPetitionEventResources, len(e.GetResources()))
	}
	petitioner, target := e.GetResources()[0], e.GetResources()[1]

	if !isSupported(supported, target) {
		logger.Error("Unsupported resource kind", zap.String("kind", target.GetKind()))
		return nil, fmt.Errorf("unsupported resource kind '%s'", target.GetKind())
	}
	logger = logger.With(zap.Object("petitioner", petitioner), zap.Object("target", target))
	return &ApplyRequest{
		Events:     req.GetEvents(),
		Event:      e,
		Petitioner: petitioner,
		Target:     target,
		Logger:     logger,
	}, nil
}

func isSupported(kinds []string, target *auditv1.Resource) bool {
	targetKind := target.GetKind()
	for _, k := range kinds {
		if k == targetKind {
			return true
		}
	}
	return false
}

func applyEvent(events []*auditv1.Event) *auditv1.Event {
	for _, eventType := range []string{common.EventRevoke, common.EventGrant} {
		for _, event := range events {
			if event.GetEvent() == eventType {
				return event
			}
		}
	}
	return nil
}
