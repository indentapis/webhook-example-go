package lambda

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/grpc/codes"

	auditv1 "go.indent.com/indent-go/api/indent/audit/v1"
	"go.indent.com/indent-go/pkg/event"
	"go.indent.com/webhook-go"
)

const (
	// EnvWebhookSecret is the environment variable name for the webhook secret.
	EnvWebhookSecret = "INDENT_WEBHOOK_SECRET"
)

// ServeApply is a helper function to serve an ApplyUpdateRequest.
func ServeApply(applyHandler webhook.ApplyHandler) {
	logger := webhook.Logger()

	secret := []byte(os.Getenv(EnvWebhookSecret))
	handler := func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		logger.Info("Received apply request, attempting to verify and decode")
		applyReq := new(auditv1.WriteRequest)
		if err := event.VerifyHeaders(secret, req.Headers, []byte(req.Body)); err != nil {
			return apiGatewayResp(logger, codes.InvalidArgument, "failed to verify message", err)
		} else if err = webhook.Decode([]byte(req.Body), applyReq); err != nil {
			return apiGatewayResp(logger, codes.InvalidArgument, "failed to decode message", err)
		}

		logger.Info("Performing apply update and encoding result")
		var data []byte
		if applyResp, err := webhook.ApplyUpdate(ctx, applyHandler, applyReq); err != nil {
			return apiGatewayResp(logger, codes.Internal, "failed to apply update", err)
		} else if data, err = webhook.Encode(applyResp); err != nil {
			return apiGatewayResp(logger, codes.Internal, "failed to encode response", err)
		}

		logger.Info("Successfully applied update")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(data),
		}, nil
	}
	lambda.Start(handler)
}
