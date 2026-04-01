//revive:disable:package-comments
package tests

import (
	"context"
	"testing"

	pb "github.com/katastroma/diataxis"
	"google.golang.org/grpc/metadata"
)

// OrdererContext returns a context with orderer type gRPC metadata.
func OrdererContext(t *testing.T, ordererType pb.OrdererType) context.Context {
	t.Helper()
	md := metadata.Pairs(pb.OrdererTypeMetadataKey, ordererType.String())
	return metadata.NewIncomingContext(t.Context(), md)
}

// OrdererContextRaw returns a context with a raw orderer type string in gRPC metadata.
func OrdererContextRaw(t *testing.T, value string) context.Context {
	t.Helper()
	md := metadata.Pairs(pb.OrdererTypeMetadataKey, value)
	return metadata.NewIncomingContext(t.Context(), md)
}
