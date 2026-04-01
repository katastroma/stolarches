//revive:disable:package-comments
package helm

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	helmutil "helm.sh/helm/v4/pkg/release/v1/util"
)

func decodeManifests(dec runtime.Serializer, sorted []helmutil.Manifest) ([]*unstructured.Unstructured, error) {
	resources := make([]*unstructured.Unstructured, 0, len(sorted))

	for _, m := range sorted {
		obj := &unstructured.Unstructured{}
		if _, _, err := dec.Decode([]byte(m.Content), nil, obj); err != nil {
			return nil, fmt.Errorf("decoding sorted manifest: %w", err)
		}

		resources = append(resources, obj)
	}

	return resources, nil
}
