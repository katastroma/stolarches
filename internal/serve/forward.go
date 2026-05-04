//revive:disable:package-comments
package serve

import (
	"context"
	"log/slog"

	"github.com/katastroma/stolarches/internal/labeler"
	"github.com/katastroma/stolarches/internal/order"
)

func orderAndForward(ctx context.Context, log *slog.Logger, backend order.Backend, streamFn labeler.StreamFunc) {
	ctx = context.WithoutCancel(ctx)

	log.DebugContext(ctx, "ordering manifests")
	resources, err := backend.Order()
	if err != nil {
		fail(ctx, log, "ordering failed", err)
		return
	}
	log.DebugContext(ctx, "manifests ordered", "count", len(resources))

	log.DebugContext(ctx, "streaming to provisioner")
	if err = streamFn(ctx, resources); err != nil {
		fail(ctx, log, "streaming to provisioner failed", err)
		return
	}
	log.InfoContext(ctx, "streamed to provisioner")
}
