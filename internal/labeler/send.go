//revive:disable:package-comments
package labeler

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/katastroma/akrostolion"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func send(
	ctx context.Context,
	log *slog.Logger,
	client pb.LabelerServiceClient,
	resources []*unstructured.Unstructured,
) error {
	stream, err := client.Label(ctx)
	if err != nil {
		log.ErrorContext(ctx, "opening labeler stream failed", "error", err)
		return fmt.Errorf("opening labeler stream: %w", err)
	}

	log.DebugContext(ctx, "streaming resources to labeler", "count", len(resources))

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

		if err = stream.Send(&pb.LabelRequest{Manifest: data}); err != nil {
			log.ErrorContext(ctx, "sending resource failed", "kind", kind, "name", name, "error", err)
			return fmt.Errorf("sending %s/%s to labeler: %w", kind, name, err)
		}
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		log.ErrorContext(ctx, "labeler failed", "error", err)
		return fmt.Errorf("labeler: %w", err)
	}

	return nil
}
