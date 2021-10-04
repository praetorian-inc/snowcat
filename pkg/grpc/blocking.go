package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// BlockingDial is a helper method to dial the given address, using optional TLS credentials,
// and blocking until the returned connection is ready. If the given credentials are nil, the
// connection will be insecure (plain-text).
func BlockingDial(ctx context.Context, network, address string, creds credentials.TransportCredentials,
	opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// grpc.Dial doesn't provide any information on permanent connection errors (like
	// TLS handshake failures). So in order to provide good error messages, we need a
	// custom dialer that can provide that info. That means we manage the TLS handshake.
	result := make(chan interface{}, 1)

	writeResult := func(res interface{}) {
		// non-blocking write: we only need the first result
		select {
		case result <- res:
		default:
		}
	}

	dialer := func(ctx context.Context, address string) (net.Conn, error) {
		// NB: We *could* handle the TLS handshake ourselves, in the custom
		// dialer (instead of customizing both the dialer and the credentials).
		// But that requires using WithInsecure dial option (so that the gRPC
		// library doesn't *also* try to do a handshake). And that would mean
		// that the library would send the wrong ":scheme" metaheader to
		// servers: it would send "http" instead of "https" because it is
		// unaware that TLS is actually in use.
		conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
		if err != nil {
			writeResult(err)
		}
		return conn, err
	}

	// Even with grpc.FailOnNonTempDialError, this call will usually timeout in
	// the face of TLS handshake errors. So we can't rely on grpc.WithBlock() to
	// know when we're done. So we run it in a goroutine and then use result
	// channel to either get the connection or fail-fast.
	go func() {
		// We put grpc.FailOnNonTempDialError *before* the explicitly provided
		// options so that it could be overridden.
		opts = append([]grpc.DialOption{grpc.FailOnNonTempDialError(true)}, opts...)
		// But we don't want caller to be able to override these two, so we put
		// them *after* the explicitly provided options.
		opts = append(opts, grpc.WithBlock(), grpc.WithContextDialer(dialer))

		if creds == nil {
			opts = append(opts, grpc.WithInsecure())
		} else {
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
		conn, err := grpc.DialContext(ctx, address, opts...)
		var res interface{}
		if err != nil {
			res = err
		} else {
			res = conn
		}
		writeResult(res)
	}()

	select {
	case res := <-result:
		if conn, ok := res.(*grpc.ClientConn); ok {
			return conn, nil
		}
		return nil, res.(error)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
