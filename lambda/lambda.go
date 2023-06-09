package lambda

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	indentv1 "go.indent.com/indent-go/api/indent/v1"
	"go.indent.com/webhook-go"
)

// apiGatewayResp returns an API Gateway formatted response.
func apiGatewayResp(logger *zap.Logger, code codes.Code, msg string, err error) (events.APIGatewayProxyResponse, error) {
	logger.Error(msg, zap.Error(err))
	msg += ": " + err.Error()
	applyResp := &indentv1.ApplyUpdateResponse{
		Status: &status.Status{
			Code:    int32(code),
			Message: msg,
		},
	}

	var data []byte
	if data, err = webhook.Encode(applyResp); err != nil {
		logger.Error("Failed to encode response", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(data),
	}, nil
}
