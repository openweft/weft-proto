package weftv1

import (
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"
)

// The bucket secret is write-only: BucketInfo.secret_access_key keeps its
// tag (7, never reused) but is marked deprecated so no client builds on it.
func TestBucketInfoSecretIsDeprecatedAndKeepsTag(t *testing.T) {
	fd := (&BucketInfo{}).ProtoReflect().Descriptor().Fields().ByName("secret_access_key")
	if fd == nil {
		t.Fatal("BucketInfo.secret_access_key is gone: its tag must stay allocated")
	}
	if fd.Number() != 7 {
		t.Errorf("secret_access_key tag = %d, want 7", fd.Number())
	}
	opts, ok := fd.Options().(*descriptorpb.FieldOptions)
	if !ok || !opts.GetDeprecated() {
		t.Error("BucketInfo.secret_access_key must be marked deprecated")
	}
}
