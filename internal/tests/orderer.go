//revive:disable:package-comments
package tests

import (
	"context"
	"io"

	pb "github.com/katastroma/diataxis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MockOrdererServer implements pb.OrdererService_OrderServer for testing.
type MockOrdererServer struct {
	// Requests are returned sequentially by Recv.
	Requests []*pb.OrderRequest
	// Responses collects messages passed to Send.
	Responses []*pb.OrderResponse
	// SendErr is returned by Send when set.
	SendErr error
	// RecvErr is returned by Recv when set.
	RecvErr error
	// Ctx is returned by Context.
	Ctx     context.Context
	recvIdx int
	grpc.ServerStream
}

// Recv returns the next pre-loaded request, the configured error, or io.EOF.
func (s *MockOrdererServer) Recv() (*pb.OrderRequest, error) {
	if s.RecvErr != nil {
		return nil, s.RecvErr
	}
	if s.recvIdx >= len(s.Requests) {
		return nil, io.EOF
	}
	req := s.Requests[s.recvIdx]
	s.recvIdx++
	return req, nil
}

// SendAndClose records the response or returns the configured error.
func (s *MockOrdererServer) SendAndClose(resp *pb.OrderResponse) error {
	if s.SendErr != nil {
		return s.SendErr
	}
	s.Responses = append(s.Responses, resp)
	return nil
}

// Context returns the configured context.
func (s *MockOrdererServer) Context() context.Context {
	return s.Ctx
}

// SetHeader is a no-op.
func (s *MockOrdererServer) SetHeader(_ metadata.MD) error { return nil }

// SendHeader is a no-op.
func (s *MockOrdererServer) SendHeader(_ metadata.MD) error { return nil }

// SetTrailer is a no-op.
func (s *MockOrdererServer) SetTrailer(_ metadata.MD) {}

// SendMsg is a no-op.
func (s *MockOrdererServer) SendMsg(_ any) error { return nil }

// RecvMsg signals end of stream.
func (s *MockOrdererServer) RecvMsg(_ any) error { return io.EOF }

// MockOrdererStreamServer implements pb.OrdererService_OrderStreamServer for testing.
type MockOrdererStreamServer struct {
	// Requests are returned sequentially by Recv.
	Requests []*pb.OrderStreamRequest
	// Sent collects messages passed to Send.
	Sent []*pb.OrderStreamResponse
	// SendErr is returned by Send when set.
	SendErr error
	// RecvErr is returned by Recv when set.
	RecvErr error
	// Ctx is returned by Context.
	Ctx     context.Context
	recvIdx int
	grpc.ServerStream
}

// Recv returns the next pre-loaded request, the configured error, or io.EOF.
func (s *MockOrdererStreamServer) Recv() (*pb.OrderStreamRequest, error) {
	if s.RecvErr != nil {
		return nil, s.RecvErr
	}
	if s.recvIdx >= len(s.Requests) {
		return nil, io.EOF
	}
	req := s.Requests[s.recvIdx]
	s.recvIdx++
	return req, nil
}

// Send records a sent response or returns the configured error.
func (s *MockOrdererStreamServer) Send(resp *pb.OrderStreamResponse) error {
	if s.SendErr != nil {
		return s.SendErr
	}
	s.Sent = append(s.Sent, resp)
	return nil
}

// Context returns the configured context.
func (s *MockOrdererStreamServer) Context() context.Context {
	return s.Ctx
}

// SetHeader is a no-op.
func (s *MockOrdererStreamServer) SetHeader(_ metadata.MD) error { return nil }

// SendHeader is a no-op.
func (s *MockOrdererStreamServer) SendHeader(_ metadata.MD) error { return nil }

// SetTrailer is a no-op.
func (s *MockOrdererStreamServer) SetTrailer(_ metadata.MD) {}

// SendMsg is a no-op.
func (s *MockOrdererStreamServer) SendMsg(_ any) error { return nil }

// RecvMsg signals end of stream.
func (s *MockOrdererStreamServer) RecvMsg(_ any) error { return io.EOF }
