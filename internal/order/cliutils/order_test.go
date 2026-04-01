package cliutils_test

import (
	"strings"
	"testing"

	"github.com/katastroma/stolarches/internal/order/cliutils"
	"github.com/katastroma/stolarches/internal/tests"
)

func TestReceiveAndOrder(t *testing.T) {
	b := cliutils.New()

	manifest := "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: config\n---\napiVersion: v1\nkind: Namespace\nmetadata:\n  name: test\n"

	if err := b.Receive(strings.NewReader(manifest)); err != nil {
		t.Fatalf("receive: %v", err)
	}

	resources, err := b.Order()
	if err != nil {
		t.Fatalf("order: %v", err)
	}

	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	if resources[0].GetKind() != "Namespace" {
		t.Errorf("expected Namespace first, got %s", resources[0].GetKind())
	}

	if resources[1].GetKind() != "ConfigMap" {
		t.Errorf("expected ConfigMap second, got %s", resources[1].GetKind())
	}
}

func TestReceive_ReadError(t *testing.T) {
	b := cliutils.New()

	if err := b.Receive(tests.ErrReader{}); err == nil {
		t.Fatal("expected error when reader fails")
	}
}

func TestReceive_InvalidYAML(t *testing.T) {
	b := cliutils.New()

	if err := b.Receive(strings.NewReader("not: valid: yaml: {{{")); err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestReceive_EmptyDocuments(t *testing.T) {
	b := cliutils.New()

	if err := b.Receive(strings.NewReader("---\n---\n")); err != nil {
		t.Fatalf("receive: %v", err)
	}

	resources, err := b.Order()
	if err != nil {
		t.Fatalf("order: %v", err)
	}

	if len(resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(resources))
	}
}
