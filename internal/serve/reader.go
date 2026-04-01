//revive:disable:package-comments
package serve

import pb "github.com/katastroma/diataxis"

// streamReader adapts a diataxis Order stream as an io.Reader. Each
// Read call consumes from the current message buffer, receiving the next
// OrderRequest when the buffer is exhausted.
type streamReader struct {
	stream pb.OrdererService_OrderServer
	buf    []byte
}

// Read fills p from the stream, receiving new messages as needed.
func (r *streamReader) Read(p []byte) (int, error) {
	if len(r.buf) == 0 {
		req, err := r.stream.Recv()
		if err != nil {
			return 0, err
		}
		r.buf = req.GetManifest()
	}

	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}
