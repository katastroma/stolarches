//revive:disable:package-comments
package serve

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/katastroma/diataxis"
	"google.golang.org/grpc/metadata"
)

// const defaultOrdererType = pb.OrdererType_ORDERER_TYPE_CLI_UTILS

const defaultOrdererType = pb.OrdererType_ORDERER_TYPE_HELM

func fail(ctx context.Context, log *slog.Logger, msg string, err error) {
	log.ErrorContext(ctx, msg, "error", err)
}

// Order receives a YAML manifest blob from the stream, orders the manifests,
// and streams them to the provisioner.
func (s *Service) Order(stream pb.OrdererService_OrderServer) error {
	ctx := stream.Context()
	s.log.InfoContext(ctx, "order requested")

	ordererType := readOrdererType(ctx)
	log := s.log.With("orderer", ordererType.String())
	log.InfoContext(ctx, "orderer type resolved")

	log.DebugContext(ctx, "looking up backend")
	backend, err := s.router.Lookup(ordererType)
	if err != nil {
		fail(ctx, log, "backend lookup failed", err)
		return fmt.Errorf("backend lookup: %w", err)
	}
	log.DebugContext(ctx, "backend found")

	log.DebugContext(ctx, "receiving manifests")
	if err = backend.Receive(&streamReader{stream: stream}); err != nil {
		fail(ctx, log, "receiving failed", err)
		return fmt.Errorf("receiving: %w", err)
	}
	log.DebugContext(ctx, "manifests received")

	if err = stream.SendAndClose(&pb.OrderResponse{}); err != nil {
		fail(ctx, log, "sending response failed", err)
		return fmt.Errorf("sending response: %w", err)
	}

	go orderAndForward(ctx, log, backend, s.streamFn)

	return nil
}

func readOrdererType(ctx context.Context) pb.OrdererType {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return defaultOrdererType
	}

	values := md.Get(pb.OrdererTypeMetadataKey)
	if len(values) == 0 {
		return defaultOrdererType
	}

	ordererType, ok := pb.OrdererType_value[values[0]]
	if !ok {
		return defaultOrdererType
	}

	return pb.OrdererType(ordererType)
}
