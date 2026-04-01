//revive:disable:package-comments
package provisioner

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/katastroma/katartismos"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func send(
	ctx context.Context,
	log *slog.Logger,
	client pb.ProvisionerServiceClient,
	resources []*unstructured.Unstructured,
) error {
	stream, err := client.Provision(ctx)
	if err != nil {
		log.ErrorContext(ctx, "opening provisioner stream failed", "error", err)
		return fmt.Errorf("opening provisioner stream: %w", err)
	}

	log.DebugContext(ctx, "streaming resources to provisioner", "count", len(resources))

	for i, obj := range resources {
		kind := obj.GetKind()
		name := obj.GetName()

		var data []byte
		data, err = obj.MarshalJSON()
		if err != nil {
			log.ErrorContext(ctx, "marshaling resource failed", "kind", kind, "name", name, "error", err)
			return fmt.Errorf("marshaling %s/%s: %w", kind, name, err)
		}

		log.DebugContext(ctx, "sending resource", "index", i, "kind", kind, "name", name, "bytes", len(data))

		if err = stream.Send(&pb.ProvisionRequest{Manifest: data}); err != nil {
			log.ErrorContext(ctx, "sending resource failed", "kind", kind, "name", name, "error", err)
			return fmt.Errorf("sending %s/%s to provisioner: %w", kind, name, err)
		}
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		log.ErrorContext(ctx, "provisioner failed", "error", err)
		return fmt.Errorf("provisioner: %w", err)
	}

	return nil
}
