//revive:disable:package-comments
package tests

import (
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/katastroma/stolarches/internal/order"
)

var _ order.Backend = (*MockBackend)(nil)

// MockBackend implements order.Backend for testing.
type MockBackend struct {
	// ReceiveErr is returned by Receive when set.
	ReceiveErr error
	// OrderResult is returned by Order on success.
	OrderResult []*unstructured.Unstructured
	// OrderErr is returned by Order when set.
	OrderErr error
}

// Receive returns the configured error.
func (b *MockBackend) Receive(_ io.Reader) error {
	return b.ReceiveErr
}

// Order returns the configured result or error.
func (b *MockBackend) Order() ([]*unstructured.Unstructured, error) {
	if b.OrderErr != nil {
		return nil, b.OrderErr
	}
	return b.OrderResult, nil
}
