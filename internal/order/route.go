//revive:disable:package-comments
package order

import (
	"fmt"
	"io"

	pb "github.com/katastroma/diataxis"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Backend receives a YAML manifest blob and produces a sorted version.
type Backend interface {
	// Receive reads the manifest stream into the backend's internal format.
	Receive(io.Reader) error
	// Order sorts the received manifests into a safe apply sequence.
	Order() ([]*unstructured.Unstructured, error)
}

// Router dispatches ordering to the backend registered for a given type.
type Router struct {
	backends map[pb.OrdererType]Backend
}

// Register adds an ordering backend for an orderer type.
func (r *Router) Register(ordererType pb.OrdererType, backend Backend) {
	if r.backends == nil {
		r.backends = make(map[pb.OrdererType]Backend)
	}
	r.backends[ordererType] = backend
}

// Lookup returns the backend registered for the given type.
func (r *Router) Lookup(ordererType pb.OrdererType) (Backend, error) {
	backend, ok := r.backends[ordererType]
	if !ok {
		return nil, fmt.Errorf("no orderer registered for %s", ordererType)
	}

	return backend, nil
}
