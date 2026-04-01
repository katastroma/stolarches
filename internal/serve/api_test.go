package serve_test

import (
	"context"
	"log/slog"
	"testing"

	pb "github.com/katastroma/diataxis"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/katastroma/stolarches/internal/order"
	"github.com/katastroma/stolarches/internal/serve"
	"github.com/katastroma/stolarches/internal/tests"
)

func TestNew(t *testing.T) {
	svc := serve.New(slog.Default(), nil, nil)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestOrder(t *testing.T) {
	var forwarded []*unstructured.Unstructured

	var router order.Router
	router.Register(pb.OrdererType_ORDERER_TYPE_HELM, &tests.MockBackend{
		OrderResult: tests.TestResources(),
	})

	svc := serve.New(
		slog.Default(),
		&router,
		func(_ context.Context, resources []*unstructured.Unstructured) error {
			forwarded = resources
			return nil
		},
	)

	stream := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{{Manifest: []byte("data")}},
		Ctx:      t.Context(),
	}

	if err := svc.Order(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(forwarded) != 1 {
		t.Fatalf("expected 1 resource forwarded, got %d", len(forwarded))
	}

	if forwarded[0].GetKind() != "ConfigMap" {
		t.Errorf("expected ConfigMap, got %s", forwarded[0].GetKind())
	}

	if len(stream.Responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(stream.Responses))
	}
}
