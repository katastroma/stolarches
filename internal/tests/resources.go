//revive:disable:package-comments
package tests

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// TestResources returns a minimal set of unstructured resources for testing.
func TestResources() []*unstructured.Unstructured {
	return []*unstructured.Unstructured{{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata":   map[string]any{"name": "test"},
		},
	}}
}
