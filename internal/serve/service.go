//revive:disable:package-comments
package serve

import (
	"log/slog"

	pb "github.com/katastroma/diataxis"

	"github.com/katastroma/stolarches/internal/labeler"
	"github.com/katastroma/stolarches/internal/order"
)

// Service implements the diataxis OrdererServiceServer
type Service struct {
	pb.UnimplementedOrdererServiceServer
	log      *slog.Logger
	router   *order.Router
	streamFn labeler.StreamFunc
}

// New creates a Service with the given logger, router, and stream function
func New(log *slog.Logger, router *order.Router, streamFn labeler.StreamFunc) *Service {
	return &Service{log: log, router: router, streamFn: streamFn}
}
