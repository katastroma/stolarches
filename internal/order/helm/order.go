//revive:disable:package-comments
package helm

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	helmutil "helm.sh/helm/v4/pkg/release/v1/util"
)

const manifestKey = "manifests"

// Backend implements order.Backend using helm's InstallOrder sorting.
type Backend struct {
	raw string
}

// New creates a helm ordering backend.
func New() *Backend {
	return &Backend{}
}

// Receive reads the YAML manifest blob as a raw string.
func (b *Backend) Receive(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading manifests: %w", err)
	}

	b.raw = string(data)
	return nil
}

// Order sorts the received manifests using helm's InstallOrder and returns
// them as unstructured resources.
func (b *Backend) Order() ([]*unstructured.Unstructured, error) {
	_, sorted, err := helmutil.SortManifests(
		map[string]string{manifestKey: b.raw},
		nil,
		helmutil.InstallOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("sorting manifests: %w", err)
	}

	dec := k8syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	return decodeManifests(dec, sorted)
}
