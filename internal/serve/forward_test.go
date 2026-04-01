package serve

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/katastroma/stolarches/internal/tests"
)

func TestOrderAndForward(t *testing.T) {
	var forwarded []*unstructured.Unstructured
	backend := &tests.MockBackend{OrderResult: tests.TestResources()}

	orderAndForward(t.Context(), slog.Default(), backend, func(_ context.Context, resources []*unstructured.Unstructured) error {
		forwarded = resources
		return nil
	})

	if len(forwarded) != 1 {
		t.Fatalf("expected 1 resource forwarded, got %d", len(forwarded))
	}

	if forwarded[0].GetKind() != "ConfigMap" {
		t.Errorf("expected ConfigMap, got %s", forwarded[0].GetKind())
	}
}

func TestOrderAndForward_OrderError(t *testing.T) {
	backend := &tests.MockBackend{OrderErr: fmt.Errorf("order failed")}

	orderAndForward(t.Context(), slog.Default(), backend, func(_ context.Context, _ []*unstructured.Unstructured) error {
		t.Fatal("streamFn should not be called when order fails")
		return nil
	})
}

func TestOrderAndForward_ForwardError(t *testing.T) {
	backend := &tests.MockBackend{OrderResult: tests.TestResources()}

	orderAndForward(t.Context(), slog.Default(), backend, func(_ context.Context, _ []*unstructured.Unstructured) error {
		return fmt.Errorf("forward failed")
	})
}
