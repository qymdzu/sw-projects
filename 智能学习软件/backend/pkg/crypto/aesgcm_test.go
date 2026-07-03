package crypto

import (
	"bytes"
	"errors"
	"testing"
)

func TestAESGCM_RoundTrip(t *testing.T) {
	a, err := NewAESGCM("unit-test-secret-please-ignore")
	if err != nil {
		t.Fatalf("NewAESGCM: %v", err)
	}
	plain := []byte("sk-proj-1234567890ABCDEFGHIJKLMNOP")
	ct, nonce, err := a.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if bytes.Equal(ct, plain) {
		t.Fatal("ciphertext must differ from plaintext")
	}
	got, err := a.Decrypt(ct, nonce)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("decrypt mismatch: got=%s want=%s", got, plain)
	}
}

func TestAESGCM_EmptySecret(t *testing.T) {
	if _, err := NewAESGCM(""); err == nil {
		t.Fatal("expected error for empty secret")
	}
}

func TestAESGCM_BadNonceSize(t *testing.T) {
	a, err := NewAESGCM("unit-test-secret")
	if err != nil {
		t.Fatalf("NewAESGCM: %v", err)
	}
	_, err = a.Decrypt([]byte("x"), []byte("short"))
	if !errors.Is(err, ErrCipher) {
		t.Fatalf("want ErrCipher, got %v", err)
	}
}

func TestAESGCM_DifferentNonces(t *testing.T) {
	a, _ := NewAESGCM("unit-test-secret")
	plain := []byte("same-plaintext-twice")
	ct1, n1, _ := a.Encrypt(plain)
	ct2, n2, _ := a.Encrypt(plain)
	if bytes.Equal(n1, n2) {
		t.Fatal("两次加密 nonce 居然相同")
	}
	if bytes.Equal(ct1, ct2) {
		t.Fatal("两次加密密文居然相同（nonce 必然失效）")
	}
}