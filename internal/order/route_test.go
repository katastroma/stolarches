package order_test

import (
	"testing"

	"github.com/katastroma/diataxis"

	"github.com/katastroma/stolarches/internal/order"
	"github.com/katastroma/stolarches/internal/tests"
)

func TestRouter_Lookup(t *testing.T) {
	var r order.Router
	r.Register(diataxis.OrdererType_ORDERER_TYPE_CLI_UTILS, &tests.MockBackend{})

	backend, err := r.Lookup(diataxis.OrdererType_ORDERER_TYPE_CLI_UTILS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if backend == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestRouter_Lookup_NotRegistered(t *testing.T) {
	var r order.Router
	r.Register(diataxis.OrdererType_ORDERER_TYPE_CLI_UTILS, &tests.MockBackend{})

	if _, err := r.Lookup(diataxis.OrdererType_ORDERER_TYPE_HELM); err == nil {
		t.Fatal("expected error when no backend registered for type")
	}
}

func TestRouter_Lookup_Empty(t *testing.T) {
	var r order.Router

	if _, err := r.Lookup(diataxis.OrdererType_ORDERER_TYPE_CLI_UTILS); err == nil {
		t.Fatal("expected error with no registered backends")
	}
}
