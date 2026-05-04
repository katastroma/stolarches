package labeler_test

import (
	"fmt"
	"log/slog"
	"testing"

	"google.golang.org/grpc"

	"github.com/katastroma/stolarches/internal/labeler"
	"github.com/katastroma/stolarches/internal/tests"
)

func TestNewStreamFunc(t *testing.T) {
	cs := &tests.MockClientStream{Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}
	streamFn := labeler.NewStreamFunc(slog.Default(), conn)

	if err := streamFn(t.Context(), tests.TestResources()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewStreamFunc_StreamError(t *testing.T) {
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}
	streamFn := labeler.NewStreamFunc(slog.Default(), conn)

	if err := streamFn(t.Context(), tests.TestResources()); err == nil {
		t.Fatal("expected error when stream fails")
	}
}
