//revive:disable:package-comments
package serve

import (
	"log/slog"

	pb "github.com/katastroma/diataxis"

	"github.com/katastroma/stolarches/internal/order"
	"github.com/katastroma/stolarches/internal/provisioner"
)

// Service implements the diataxis OrdererServiceServer
type Service struct {
	pb.UnimplementedOrdererServiceServer
	log      *slog.Logger
	router   *order.Router
	streamFn provisioner.StreamFunc
}

// New creates a Service with the given logger, router, and stream function
func New(log *slog.Logger, router *order.Router, streamFn provisioner.StreamFunc) *Service {
	return &Service{log: log, router: router, streamFn: streamFn}
}
