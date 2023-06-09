package webhook

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	marshaller jsonpb.Marshaler
)

// Encode returns a JSON encoded version of msg.
func Encode(msg proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	if err := marshaller.Marshal(&buf, msg); err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}
	return buf.Bytes(), nil
}

// Decode decodes JSON data into msg.
func Decode(data []byte, msg proto.Message) error {
	return jsonpb.Unmarshal(bytes.NewReader(data), msg)
}
