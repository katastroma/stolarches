package helm

import (
	"fmt"
	"io"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	helmutil "helm.sh/helm/v4/pkg/release/v1/util"
)

type failingDecoder struct{}

func (failingDecoder) Decode([]byte, *schema.GroupVersionKind, runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	return nil, nil, fmt.Errorf("decode failed")
}

func (failingDecoder) Encode(runtime.Object, io.Writer) error {
	return fmt.Errorf("encode not supported")
}

func (failingDecoder) Identifier() runtime.Identifier {
	return "failing"
}

func TestDecodeManifests(t *testing.T) {
	dec := k8syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	sorted := []helmutil.Manifest{
		{Content: "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n"},
	}

	resources, err := decodeManifests(dec, sorted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}

	if resources[0].GetKind() != "ConfigMap" {
		t.Errorf("expected ConfigMap, got %s", resources[0].GetKind())
	}
}

func TestDecodeManifests_DecodeError(t *testing.T) {
	sorted := []helmutil.Manifest{
		{Content: "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n"},
	}

	if _, err := decodeManifests(failingDecoder{}, sorted); err == nil {
		t.Fatal("expected error when decoder fails")
	}
}
