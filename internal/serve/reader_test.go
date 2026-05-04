package serve

import (
	"io"
	"testing"

	pb "github.com/katastroma/diataxis"

	"github.com/katastroma/stolarches/internal/tests"
)

func TestStreamReader_SingleMessage(t *testing.T) {
	mock := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{{Data: []byte("hello")}},
		Ctx:      t.Context(),
	}
	r := &streamReader{stream: mock}

	buf := make([]byte, 16)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected %q, got %q", "hello", string(buf[:n]))
	}
}

func TestStreamReader_MultipleMessages(t *testing.T) {
	mock := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{
			{Data: []byte("ab")},
			{Data: []byte("cd")},
		},
		Ctx: t.Context(),
	}
	r := &streamReader{stream: mock}

	result, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "abcd" {
		t.Errorf("expected %q, got %q", "abcd", string(result))
	}
}

func TestStreamReader_PartialRead(t *testing.T) {
	mock := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{{Data: []byte("abcdef")}},
		Ctx:      t.Context(),
	}
	r := &streamReader{stream: mock}

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on first read: %v", err)
	}
	if string(buf[:n]) != "abc" {
		t.Errorf("expected %q, got %q", "abc", string(buf[:n]))
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on second read: %v", err)
	}
	if string(buf[:n]) != "def" {
		t.Errorf("expected %q, got %q", "def", string(buf[:n]))
	}
}

func TestStreamReader_EOF(t *testing.T) {
	mock := &tests.MockOrdererServer{
		Requests: []*pb.OrderRequest{},
		Ctx:      t.Context(),
	}
	r := &streamReader{stream: mock}

	buf := make([]byte, 16)
	_, err := r.Read(buf)
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

func TestOrderStreamReader_SingleMessage(t *testing.T) {
	mock := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{{Data: []byte("hello")}},
		Ctx:      t.Context(),
	}
	r := &orderStreamReader{stream: mock}

	buf := make([]byte, 16)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected %q, got %q", "hello", string(buf[:n]))
	}
}

func TestOrderStreamReader_MultipleMessages(t *testing.T) {
	mock := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{
			{Data: []byte("ab")},
			{Data: []byte("cd")},
		},
		Ctx: t.Context(),
	}
	r := &orderStreamReader{stream: mock}

	result, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "abcd" {
		t.Errorf("expected %q, got %q", "abcd", string(result))
	}
}

func TestOrderStreamReader_PartialRead(t *testing.T) {
	mock := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{{Data: []byte("abcdef")}},
		Ctx:      t.Context(),
	}
	r := &orderStreamReader{stream: mock}

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on first read: %v", err)
	}
	if string(buf[:n]) != "abc" {
		t.Errorf("expected %q, got %q", "abc", string(buf[:n]))
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on second read: %v", err)
	}
	if string(buf[:n]) != "def" {
		t.Errorf("expected %q, got %q", "def", string(buf[:n]))
	}
}

func TestOrderStreamReader_EOF(t *testing.T) {
	mock := &tests.MockOrdererStreamServer{
		Requests: []*pb.OrderStreamRequest{},
		Ctx:      t.Context(),
	}
	r := &orderStreamReader{stream: mock}

	buf := make([]byte, 16)
	_, err := r.Read(buf)
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}
