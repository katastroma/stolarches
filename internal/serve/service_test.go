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
	var router order.Router
	router.Register(pb.OrdererType_ORDERER_TYPE_HELM, &tests.MockBackend{
		OrderResult: tests.TestResources(),
	})

	svc := serve.New(
		slog.Default(),
		&router,
		func(_ context.Context, _ []*unstructured.Unstructured) error { return nil },
	)

	stream := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{{Manifest: []byte("data")}},
		Ctx:      t.Context(),
	}

	if err := svc.Order(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stream.Responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(stream.Responses))
	}
}
