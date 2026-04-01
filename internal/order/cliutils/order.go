//revive:disable:package-comments
package cliutils

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"sigs.k8s.io/cli-utils/pkg/ordering"
)

// Backend implements order.Backend using the cli-utils GVK sorting algorithm.
type Backend struct {
	resources []*unstructured.Unstructured
}

// New creates a cli-utils ordering backend.
func New() *Backend {
	return &Backend{}
}

// Receive parses the YAML manifest blob into unstructured resources.
func (b *Backend) Receive(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading manifests: %w", err)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), len(data))

	for {
		obj := &unstructured.Unstructured{}
		if err := decoder.Decode(obj); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("decoding manifest: %w", err)
		}

		if len(obj.Object) == 0 {
			continue
		}

		b.resources = append(b.resources, obj)
	}

	return nil
}

// Order sorts the received resources by GVK priority.
func (b *Backend) Order() ([]*unstructured.Unstructured, error) {
	sort.Sort(ordering.SortableUnstructureds(b.resources))
	return b.resources, nil
}
