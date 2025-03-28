package types

import (
	"testing"
)

func TestSecretEncryption(t *testing.T) {
	secret := NewSecret("secret")
	res := secret.Encryption("data")
	if res != "mBjjMGulrCZ7XyZ5/kq9N+bNe1Q=" {
		t.Error("secret encryption error")
	}
}
