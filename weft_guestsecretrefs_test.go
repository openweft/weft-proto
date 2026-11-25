package weftv1

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

func TestGuestSecretRefMessages(t *testing.T) {
	var nilRef *GuestSecretRef
	if nilRef.GetSecretId() != "" || nilRef.GetRef() != "" || nilRef.GetVersion() != 0 {
		t.Error("nil GuestSecretRef getter returned a non-zero value")
	}
	r := &GuestSecretRef{SecretId: "k", Ref: "file:/x", Version: 2}
	if r.GetSecretId() != "k" || r.GetRef() != "file:/x" || r.GetVersion() != 2 {
		t.Error("GuestSecretRef getters")
	}
	set := &SetGuestSecretRefRequest{SecretId: "a", Ref: "b"}
	if set.GetSecretId() != "a" || set.GetRef() != "b" {
		t.Error("SetGuestSecretRefRequest getters")
	}
	var nilSet *SetGuestSecretRefRequest
	if nilSet.GetSecretId() != "" || nilSet.GetRef() != "" {
		t.Error("nil SetGuestSecretRefRequest getter")
	}
	if (&SetGuestSecretRefResponse{Version: 5}).GetVersion() != 5 || (*SetGuestSecretRefResponse)(nil).GetVersion() != 0 {
		t.Error("SetGuestSecretRefResponse.GetVersion")
	}
	if (&DeleteGuestSecretRefRequest{SecretId: "a"}).GetSecretId() != "a" || (*DeleteGuestSecretRefRequest)(nil).GetSecretId() != "" {
		t.Error("DeleteGuestSecretRefRequest.GetSecretId")
	}
	l := &ListGuestSecretRefsResponse{Refs: []*GuestSecretRef{r}}
	if len(l.GetRefs()) != 1 || (*ListGuestSecretRefsResponse)(nil).GetRefs() != nil {
		t.Error("ListGuestSecretRefsResponse.GetRefs")
	}
}

// ProtoReflect on a live message caches its message info in the message's
// own state; a nil message must still reflect as invalid.
func TestGuestSecretRefProtoReflectCachesMessageInfo(t *testing.T) {
	a, b, c, d, e, f, g := &GuestSecretRef{}, &SetGuestSecretRefRequest{}, &SetGuestSecretRefResponse{},
		&DeleteGuestSecretRefRequest{}, &DeleteGuestSecretRefResponse{},
		&ListGuestSecretRefsRequest{}, &ListGuestSecretRefsResponse{}
	live := []struct {
		m protoreflect.ProtoMessage
		p protoimpl.Pointer
	}{
		{a, protoimpl.Pointer(a)}, {b, protoimpl.Pointer(b)}, {c, protoimpl.Pointer(c)},
		{d, protoimpl.Pointer(d)}, {e, protoimpl.Pointer(e)}, {f, protoimpl.Pointer(f)}, {g, protoimpl.Pointer(g)},
	}
	for _, x := range live {
		x.m.ProtoReflect()
		if protoimpl.X.MessageStateOf(x.p).LoadMessageInfo() == nil {
			t.Errorf("%T: message info not cached", x.m)
		}
	}
	nils := []protoreflect.ProtoMessage{
		(*GuestSecretRef)(nil), (*SetGuestSecretRefRequest)(nil), (*SetGuestSecretRefResponse)(nil),
		(*DeleteGuestSecretRefRequest)(nil), (*DeleteGuestSecretRefResponse)(nil),
		(*ListGuestSecretRefsRequest)(nil), (*ListGuestSecretRefsResponse)(nil),
	}
	for _, m := range nils {
		if m.ProtoReflect().IsValid() {
			t.Errorf("%T: a nil message reflected as valid", m)
		}
	}
}

type failingSecretRefConn struct{}

func (failingSecretRefConn) Invoke(context.Context, string, any, any, ...grpc.CallOption) error {
	return errors.New("boom")
}
func (failingSecretRefConn) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, errors.New("boom")
}

func TestGuestSecretRefClientsPropagateTransportErrors(t *testing.T) {
	c := NewWeftAgentClient(failingSecretRefConn{})
	ctx := context.Background()
	if r, err := c.SetGuestSecretRef(ctx, &SetGuestSecretRefRequest{}); err == nil || r != nil {
		t.Errorf("SetGuestSecretRef: got (%v, %v)", r, err)
	}
	if r, err := c.DeleteGuestSecretRef(ctx, &DeleteGuestSecretRefRequest{}); err == nil || r != nil {
		t.Errorf("DeleteGuestSecretRef: got (%v, %v)", r, err)
	}
	if r, err := c.ListGuestSecretRefs(ctx, &ListGuestSecretRefsRequest{}); err == nil || r != nil {
		t.Errorf("ListGuestSecretRefs: got (%v, %v)", r, err)
	}
}

type fakeSecretRefAgent struct {
	UnimplementedWeftAgentServer
	calls []string
}

func (f *fakeSecretRefAgent) SetGuestSecretRef(context.Context, *SetGuestSecretRefRequest) (*SetGuestSecretRefResponse, error) {
	f.calls = append(f.calls, "set")
	return &SetGuestSecretRefResponse{Version: 1}, nil
}
func (f *fakeSecretRefAgent) DeleteGuestSecretRef(context.Context, *DeleteGuestSecretRefRequest) (*DeleteGuestSecretRefResponse, error) {
	f.calls = append(f.calls, "delete")
	return &DeleteGuestSecretRefResponse{}, nil
}
func (f *fakeSecretRefAgent) ListGuestSecretRefs(context.Context, *ListGuestSecretRefsRequest) (*ListGuestSecretRefsResponse, error) {
	f.calls = append(f.calls, "list")
	return &ListGuestSecretRefsResponse{}, nil
}

type secretRefHandler func(any, context.Context, func(any) error, grpc.UnaryServerInterceptor) (any, error)

func TestGuestSecretRefHandlers(t *testing.T) {
	cases := []struct {
		name, method, call string
		h                  secretRefHandler
	}{
		{"set", WeftAgent_SetGuestSecretRef_FullMethodName, "set", _WeftAgent_SetGuestSecretRef_Handler},
		{"delete", WeftAgent_DeleteGuestSecretRef_FullMethodName, "delete", _WeftAgent_DeleteGuestSecretRef_Handler},
		{"list", WeftAgent_ListGuestSecretRefs_FullMethodName, "list", _WeftAgent_ListGuestSecretRefs_Handler},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := &fakeSecretRefAgent{}
			// A decode error never reaches the service.
			if _, err := tc.h(srv, context.Background(), func(any) error { return errors.New("bad frame") }, nil); err == nil || len(srv.calls) != 0 {
				t.Fatalf("decode error ignored: err=%v calls=%v", err, srv.calls)
			}
			// The interceptor, when set, wraps the call.
			var seen string
			icpt := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) {
				seen = info.FullMethod
				return h(ctx, req)
			}
			if _, err := tc.h(srv, context.Background(), func(any) error { return nil }, icpt); err != nil {
				t.Fatal(err)
			}
			if seen != tc.method || len(srv.calls) != 1 || srv.calls[0] != tc.call {
				t.Fatalf("interceptor not run: seen=%q calls=%v", seen, srv.calls)
			}
			// No interceptor: the service is called directly.
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("nil interceptor path panicked: %v", r)
					}
				}()
				out, err := tc.h(srv, context.Background(), func(any) error { return nil }, nil)
				if err != nil || out == nil || len(srv.calls) != 2 {
					t.Errorf("nil interceptor path: out=%v err=%v calls=%v", out, err, srv.calls)
				}
			}()
		})
	}
}
