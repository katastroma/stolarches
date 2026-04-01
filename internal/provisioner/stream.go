//revive:disable:package-comments
package provisioner

import (
	"context"
	"log/slog"

	pb "github.com/katastroma/katartismos"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// StreamFunc streams ordered resources to the provisioner
type StreamFunc func(ctx context.Context, resources []*unstructured.Unstructured) error

// NewStreamFunc returns a StreamFunc that opens a Provisioner stream on conn
// and sends each resource as a separate message
func NewStreamFunc(log *slog.Logger, conn grpc.ClientConnInterface) StreamFunc {
	client := pb.NewProvisionerServiceClient(conn)
	return func(ctx context.Context, resources []*unstructured.Unstructured) error {
		return send(ctx, log, client, resources)
	}
}
