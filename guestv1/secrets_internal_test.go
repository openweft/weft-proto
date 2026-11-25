package guestv1

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
)

// A nil message answers its zero value from every getter: generated code
// relies on it, and callers hold *FetchRequest values that may be nil.
func TestNilGettersAnswerZero(t *testing.T) {
	var req *FetchRequest
	if req.GetSecretId() != "" || req.GetBootToken() != nil || req.GetKnownVersion() != 0 {
		t.Error("nil FetchRequest getter returned a non-zero value")
	}
	var resp *FetchResponse
	if resp.GetValue() != nil || resp.GetVersion() != 0 || resp.GetUnchanged() {
		t.Error("nil FetchResponse getter returned a non-zero value")
	}
	var cm *CubeFSMount
	if cm.GetAccessKeyRef() != "" || cm.GetSecretKeyRef() != "" || cm.GetSecretVersion() != 0 {
		t.Error("nil CubeFSMount ref getter returned a non-zero value")
	}
	var wg *WireGuard
	if wg.GetPrivateKeyRef() != "" {
		t.Error("nil WireGuard ref getter returned a non-zero value")
	}
}

func TestGettersReturnSetValues(t *testing.T) {
	req := &FetchRequest{SecretId: "a", BootToken: []byte{1}, KnownVersion: 4}
	if req.GetSecretId() != "a" || len(req.GetBootToken()) != 1 || req.GetKnownVersion() != 4 {
		t.Error("FetchRequest getters")
	}
	resp := &FetchResponse{Value: []byte{2}, Version: 5, Unchanged: true}
	if len(resp.GetValue()) != 1 || resp.GetVersion() != 5 || !resp.GetUnchanged() {
		t.Error("FetchResponse getters")
	}
	cm := &CubeFSMount{AccessKeyRef: "x", SecretKeyRef: "y", SecretVersion: 6}
	if cm.GetAccessKeyRef() != "x" || cm.GetSecretKeyRef() != "y" || cm.GetSecretVersion() != 6 {
		t.Error("CubeFSMount getters")
	}
	if (&WireGuard{PrivateKeyRef: "z"}).GetPrivateKeyRef() != "z" {
		t.Error("WireGuard getter")
	}
}

// Initialising a file descriptor twice must be a no-op, not a second
// registration (which panics).
func TestSecretsInitIsIdempotent(t *testing.T) {
	file_secrets_proto_init()
	file_secrets_proto_init()
}

type failingConn struct{}

func (failingConn) Invoke(context.Context, string, any, any, ...grpc.CallOption) error {
	return errors.New("boom")
}
func (failingConn) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, errors.New("boom")
}

func TestClientFetchPropagatesTransportError(t *testing.T) {
	resp, err := NewGuestSecretsClient(failingConn{}).Fetch(context.Background(), &FetchRequest{})
	if err == nil || resp != nil {
		t.Fatalf("got (%v, %v), want (nil, error)", resp, err)
	}
}

type fakeSecrets struct {
	UnimplementedGuestSecretsServer
	called bool
}

func (f *fakeSecrets) Fetch(context.Context, *FetchRequest) (*FetchResponse, error) {
	f.called = true
	return &FetchResponse{Version: 1}, nil
}

func TestHandlerDecodeErrorNeverReachesTheService(t *testing.T) {
	srv := &fakeSecrets{}
	_, err := _GuestSecrets_Fetch_Handler(srv, context.Background(),
		func(any) error { return errors.New("bad frame") }, nil)
	if err == nil || srv.called {
		t.Fatalf("decode error ignored: err=%v called=%v", err, srv.called)
	}
}

func TestHandlerRunsTheInterceptorWhenSet(t *testing.T) {
	srv := &fakeSecrets{}
	var seen string
	icpt := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
		seen = info.FullMethod
		return h(ctx, req)
	}
	if _, err := _GuestSecrets_Fetch_Handler(srv, context.Background(), func(any) error { return nil }, icpt); err != nil {
		t.Fatal(err)
	}
	if seen != GuestSecrets_Fetch_FullMethodName || !srv.called {
		t.Fatalf("interceptor not run: seen=%q called=%v", seen, srv.called)
	}
	srv.called = false
	if _, err := _GuestSecrets_Fetch_Handler(srv, context.Background(), func(any) error { return nil }, nil); err != nil || !srv.called {
		t.Fatalf("nil interceptor path: err=%v called=%v", err, srv.called)
	}
}
