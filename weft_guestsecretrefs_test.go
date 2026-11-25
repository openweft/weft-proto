package weftv1

import "testing"

func TestGuestSecretRefMessages(t *testing.T) {
	var nilRef *GuestSecretRef
	if nilRef.GetSecretId() != "" || nilRef.GetRef() != "" || nilRef.GetVersion() != 0 {
		t.Error("nil GuestSecretRef getter returned a non-zero value")
	}
	r := &GuestSecretRef{SecretId: "k", Ref: "file:/x", Version: 2}
	if r.GetSecretId() != "k" || r.GetRef() != "file:/x" || r.GetVersion() != 2 {
		t.Error("GuestSecretRef getters")
	}
	if (&SetGuestSecretRefRequest{SecretId: "a", Ref: "b"}).GetSecretId() != "a" {
		t.Error("SetGuestSecretRefRequest.GetSecretId")
	}
	if (&SetGuestSecretRefRequest{SecretId: "a", Ref: "b"}).GetRef() != "b" {
		t.Error("SetGuestSecretRefRequest.GetRef")
	}
	if (&SetGuestSecretRefResponse{Version: 5}).GetVersion() != 5 {
		t.Error("SetGuestSecretRefResponse.GetVersion")
	}
	if (&DeleteGuestSecretRefRequest{SecretId: "a"}).GetSecretId() != "a" {
		t.Error("DeleteGuestSecretRefRequest.GetSecretId")
	}
	l := &ListGuestSecretRefsResponse{Refs: []*GuestSecretRef{r}}
	if len(l.GetRefs()) != 1 {
		t.Error("ListGuestSecretRefsResponse.GetRefs")
	}
	if _, ok := any(UnimplementedWeftAgentServer{}).(WeftAgentServer); !ok {
		t.Error("UnimplementedWeftAgentServer must satisfy WeftAgentServer")
	}
}
