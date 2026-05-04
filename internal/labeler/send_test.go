package labeler

import (
	"fmt"
	"log/slog"
	"testing"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/katastroma/stolarches/internal/tests"
)

func TestSend(t *testing.T) {
	cs := &tests.MockClientStream{Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}
	streamFn := NewStreamFunc(slog.Default(), conn)

	if err := streamFn(t.Context(), tests.TestResources()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cs.SendMsgCount == 0 {
		t.Fatal("expected at least one message sent")
	}
}

func TestSend_MultipleResources(t *testing.T) {
	cs := &tests.MockClientStream{Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}
	streamFn := NewStreamFunc(slog.Default(), conn)

	resources := []*unstructured.Unstructured{
		{Object: map[string]any{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]any{"name": "ns"}}},
		{Object: map[string]any{"apiVersion": "v1", "kind": "ConfigMap", "metadata": map[string]any{"name": "cm"}}},
		{Object: map[string]any{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"name": "deploy"}}},
	}

	if err := streamFn(t.Context(), resources); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cs.SendMsgCount != 3 {
		t.Errorf("expected 3 messages sent, got %d", cs.SendMsgCount)
	}
}

func TestSend_OpenError(t *testing.T) {
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}
	streamFn := NewStreamFunc(slog.Default(), conn)

	if err := streamFn(t.Context(), tests.TestResources()); err == nil {
		t.Fatal("expected error when stream fails to open")
	}
}

func TestSend_MarshalError(t *testing.T) {
	cs := &tests.MockClientStream{Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}
	streamFn := NewStreamFunc(slog.Default(), conn)

	bad := []*unstructured.Unstructured{{
		Object: map[string]any{"kind": "Bad", "metadata": map[string]any{"name": make(chan int)}},
	}}

	if err := streamFn(t.Context(), bad); err == nil {
		t.Fatal("expected error when marshal fails")
	}
}

func TestSend_SendError(t *testing.T) {
	cs := &tests.MockClientStream{SendMsgErr: fmt.Errorf("send failed"), Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}
	streamFn := NewStreamFunc(slog.Default(), conn)

	if err := streamFn(t.Context(), tests.TestResources()); err == nil {
		t.Fatal("expected error when send fails")
	}
}

func TestSend_CloseAndRecvError(t *testing.T) {
	cs := &tests.MockClientStream{RecvMsgErr: fmt.Errorf("labeler failed"), Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}
	streamFn := NewStreamFunc(slog.Default(), conn)

	if err := streamFn(t.Context(), tests.TestResources()); err == nil {
		t.Fatal("expected error when labeler fails")
	}
}
