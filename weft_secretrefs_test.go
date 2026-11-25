package weftv1

import "testing"

func TestCubeFSMountRefGetters(t *testing.T) {
	var nilMount *CubeFSMount
	if nilMount.GetAccessKeyRef() != "" || nilMount.GetSecretKeyRef() != "" || nilMount.GetSecretVersion() != 0 {
		t.Error("nil CubeFSMount ref getter returned a non-zero value")
	}
	m := &CubeFSMount{AccessKeyRef: "a", SecretKeyRef: "b", SecretVersion: 3}
	if m.GetAccessKeyRef() != "a" || m.GetSecretKeyRef() != "b" || m.GetSecretVersion() != 3 {
		t.Error("CubeFSMount ref getters")
	}
}
