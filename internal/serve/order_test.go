package serve

import (
	"fmt"
	"log/slog"
	"testing"

	pb "github.com/katastroma/diataxis"
	"google.golang.org/grpc/metadata"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/katastroma/stolarches/internal/order"
	"github.com/katastroma/stolarches/internal/tests"
)

func TestReadOrdererType_Default(t *testing.T) {
	got := readOrdererType(t.Context())
	if got != defaultOrdererType {
		t.Errorf("expected default %v, got %v", defaultOrdererType, got)
	}
}

func TestReadOrdererType_FromMetadata(t *testing.T) {
	ctx := tests.OrdererContext(t, pb.OrdererType_ORDERER_TYPE_CLI_UTILS)

	got := readOrdererType(ctx)
	if got != pb.OrdererType_ORDERER_TYPE_CLI_UTILS {
		t.Errorf("expected CLI_UTILS, got %v", got)
	}
}

func TestReadOrdererType_UnknownFallsBackToDefault(t *testing.T) {
	ctx := tests.OrdererContextRaw(t, "ORDERER_TYPE_BOGUS")

	got := readOrdererType(ctx)
	if got != defaultOrdererType {
		t.Errorf("expected default %v for unknown type, got %v", defaultOrdererType, got)
	}
}

func TestReadOrdererType_MetadataWithoutKey(t *testing.T) {
	md := metadata.Pairs("other-key", "value")
	ctx := metadata.NewIncomingContext(t.Context(), md)

	got := readOrdererType(ctx)
	if got != defaultOrdererType {
		t.Errorf("expected default %v for missing key, got %v", defaultOrdererType, got)
	}
}

func TestOrder_BackendLookupError(t *testing.T) {
	var router order.Router

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererServer{
		Ctx: t.Context(),
	}

	if err := svc.Order(stream); err == nil {
		t.Fatal("expected error when backend not registered")
	}
}

func TestOrder_ReceiveError(t *testing.T) {
	backend := &tests.MockBackend{ReceiveErr: fmt.Errorf("receive failed")}

	var router order.Router
	router.Register(defaultOrdererType, backend)

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{{Data: []byte("data")}},
		Ctx:      t.Context(),
	}

	if err := svc.Order(stream); err == nil {
		t.Fatal("expected error when receive fails")
	}
}

func TestOrder_SendAndCloseError(t *testing.T) {
	backend := &tests.MockBackend{}

	var router order.Router
	router.Register(defaultOrdererType, backend)

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{{Data: []byte("data")}},
		SendErr:  fmt.Errorf("send failed"),
		Ctx:      t.Context(),
	}

	if err := svc.Order(stream); err == nil {
		t.Fatal("expected error when SendAndClose fails")
	}
}

func TestOrderStream(t *testing.T) {
	backend := &tests.MockBackend{OrderResult: tests.TestResources()}

	var router order.Router
	router.Register(defaultOrdererType, backend)

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{{Data: []byte("data")}},
		Ctx:      t.Context(),
	}

	if err := svc.OrderStream(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stream.Sent) != 1 {
		t.Fatalf("expected 1 response, got %d", len(stream.Sent))
	}
}

func TestOrderStream_BackendLookupError(t *testing.T) {
	var router order.Router

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererStreamServer{
		Ctx: t.Context(),
	}

	if err := svc.OrderStream(stream); err == nil {
		t.Fatal("expected error when backend not registered")
	}
}

func TestOrderStream_ReceiveError(t *testing.T) {
	backend := &tests.MockBackend{ReceiveErr: fmt.Errorf("receive failed")}

	var router order.Router
	router.Register(defaultOrdererType, backend)

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{{Data: []byte("data")}},
		Ctx:      t.Context(),
	}

	if err := svc.OrderStream(stream); err == nil {
		t.Fatal("expected error when receive fails")
	}
}

func TestOrderStream_OrderError(t *testing.T) {
	backend := &tests.MockBackend{OrderErr: fmt.Errorf("order failed")}

	var router order.Router
	router.Register(defaultOrdererType, backend)

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{{Data: []byte("data")}},
		Ctx:      t.Context(),
	}

	if err := svc.OrderStream(stream); err == nil {
		t.Fatal("expected error when order fails")
	}
}

func TestOrderStream_MarshalError(t *testing.T) {
	bad := []*unstructured.Unstructured{{
		Object: map[string]any{"kind": "Bad", "metadata": map[string]any{"name": make(chan int)}},
	}}
	backend := &tests.MockBackend{OrderResult: bad}

	var router order.Router
	router.Register(defaultOrdererType, backend)

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{{Data: []byte("data")}},
		Ctx:      t.Context(),
	}

	if err := svc.OrderStream(stream); err == nil {
		t.Fatal("expected error when marshal fails")
	}
}

func TestOrderStream_SendError(t *testing.T) {
	backend := &tests.MockBackend{OrderResult: tests.TestResources()}

	var router order.Router
	router.Register(defaultOrdererType, backend)

	svc := New(slog.Default(), &router, nil)
	stream := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{{Data: []byte("data")}},
		SendErr:  fmt.Errorf("send failed"),
		Ctx:      t.Context(),
	}

	if err := svc.OrderStream(stream); err == nil {
		t.Fatal("expected error when send fails")
	}
}

