package guestv1

import (
	"bytes"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The wire numbers below are the contract the host and the guest agree on;
// renumbering any of them silently breaks a deployed pair.
func TestSecretsWireNumbers(t *testing.T) {
	cases := []struct {
		msg    proto.Message
		field  protoreflect.Name
		number protoreflect.FieldNumber
	}{
		{&FetchRequest{}, "secret_id", 1},
		{&FetchRequest{}, "boot_token", 2},
		{&FetchRequest{}, "known_version", 3},
		{&FetchResponse{}, "value", 1},
		{&FetchResponse{}, "version", 2},
		{&FetchResponse{}, "unchanged", 3},
		{&CubeFSMount{}, "access_key_ref", 7},
		{&CubeFSMount{}, "secret_key_ref", 8},
		{&CubeFSMount{}, "secret_version", 9},
		{&WireGuard{}, "private_key_ref", 6},
	}
	for _, c := range cases {
		fd := c.msg.ProtoReflect().Descriptor().Fields().ByName(c.field)
		if fd == nil {
			t.Errorf("%T has no field %s", c.msg, c.field)
			continue
		}
		if fd.Number() != c.number {
			t.Errorf("%T.%s = %d, want %d", c.msg, c.field, fd.Number(), c.number)
		}
	}
}

func TestFetchRoundTrip(t *testing.T) {
	in := &FetchRequest{SecretId: "cubefs-ws.key", BootToken: bytes.Repeat([]byte{7}, 32), KnownVersion: 3}
	b, err := proto.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out FetchRequest
	if err := proto.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(in, &out) {
		t.Fatalf("round trip: got %v want %v", &out, in)
	}
}

// A spec written before the ref fields existed must still decode, with the
// refs empty: the change is additive.
func TestPreRefSpecStillDecodes(t *testing.T) {
	old, err := proto.Marshal(&WireGuard{Interface: "wg0", PrivateKey: "k"})
	if err != nil {
		t.Fatal(err)
	}
	var wg WireGuard
	if err := proto.Unmarshal(old, &wg); err != nil {
		t.Fatal(err)
	}
	if wg.GetPrivateKeyRef() != "" || wg.GetPrivateKey() != "k" {
		t.Fatalf("got %+v", &wg)
	}
}
