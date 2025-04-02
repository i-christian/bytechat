package cookies

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to generate a valid AES key
func generateAESKey() []byte {
	key := make([]byte, 32) // AES-256 key
	_, _ = rand.Read(key)
	return key
}

// Test Write function
func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()
	cookie := http.Cookie{Name: "test", Value: "value"}

	err := Write(w, cookie)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	resp := w.Result()
	defer resp.Body.Close()

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	decodedValue, err := base64.URLEncoding.DecodeString(cookies[0].Value)
	if err != nil || string(decodedValue) != "value" {
		t.Fatalf("expected cookie value 'value', got %s", string(decodedValue))
	}
}

// Test Write with oversized cookie
func TestWriteTooLong(t *testing.T) {
	w := httptest.NewRecorder()
	longValue := make([]byte, 5000)
	cookie := http.Cookie{Name: "test", Value: string(longValue)}

	err := Write(w, cookie)
	if !errors.Is(err, ErrValueTooLong) {
		t.Fatalf("expected ErrValueTooLong, got %v", err)
	}
}

// Test Read function
func TestRead(t *testing.T) {
	w := httptest.NewRecorder()
	cookie := http.Cookie{Name: "test", Value: "value"}
	_ = Write(w, cookie)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	value, err := Read(req, "test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "value" {
		t.Fatalf("expected 'value', got %s", value)
	}
}

// Test Read with invalid base64
func TestReadInvalidBase64(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "test", Value: "invalid$$$"})

	_, err := Read(req, "test")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

// Test WriteEncrypted function
func TestWriteEncrypted(t *testing.T) {
	w := httptest.NewRecorder()
	secretKey := generateAESKey()
	cookie := http.Cookie{Name: "test", Value: "secret"}

	err := WriteEncrypted(w, cookie, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
}

// Test WriteEncrypted with an invalid key length
func TestWriteEncryptedInvalidKey(t *testing.T) {
	w := httptest.NewRecorder()
	invalidKey := []byte("short") // Not 16, 24, or 32 bytes

	cookie := http.Cookie{Name: "test", Value: "secret"}
	err := WriteEncrypted(w, cookie, invalidKey)
	if err == nil {
		t.Fatal("expected error for invalid key length, got nil")
	}
}

// Test ReadEncrypted function
func TestReadEncrypted(t *testing.T) {
	w := httptest.NewRecorder()
	secretKey := generateAESKey()
	cookie := http.Cookie{Name: "test", Value: "secret"}

	_ = WriteEncrypted(w, cookie, secretKey)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	value, err := ReadEncrypted(req, "test", secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "secret" {
		t.Fatalf("expected 'secret', got %s", value)
	}
}

// Test ReadEncrypted with missing nonce
func TestReadEncryptedMissingNonce(t *testing.T) {
	w := httptest.NewRecorder()
	secretKey := generateAESKey()

	cookie := http.Cookie{Name: "test", Value: "secret"}
	_ = WriteEncrypted(w, cookie, secretKey)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "test", Value: "short"})

	_, err := ReadEncrypted(req, "test", secretKey)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue for short nonce, got %v", err)
	}
}

// Test ReadEncrypted with wrong key
func TestReadEncryptedWrongKey(t *testing.T) {
	w := httptest.NewRecorder()
	secretKey := generateAESKey()
	wrongKey := generateAESKey()
	cookie := http.Cookie{Name: "test", Value: "secret"}

	_ = WriteEncrypted(w, cookie, secretKey)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	_, err := ReadEncrypted(req, "test", wrongKey)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

// Test ReadEncrypted with mismatched cookie name
func TestReadEncryptedWrongName(t *testing.T) {
	w := httptest.NewRecorder()
	secretKey := generateAESKey()
	cookie := http.Cookie{Name: "wrongname", Value: "secret"}

	_ = WriteEncrypted(w, cookie, secretKey)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(&http.Cookie{Name: "test", Value: c.Value})
	}

	_, err := ReadEncrypted(req, "test", secretKey)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue for mismatched name, got %v", err)
	}
}

// Test ReadEncrypted with corrupted data
func TestReadEncryptedCorruptedData(t *testing.T) {
	w := httptest.NewRecorder()
	secretKey := generateAESKey()
	cookie := http.Cookie{Name: "test", Value: "secret"}

	_ = WriteEncrypted(w, cookie, secretKey)

	req := httptest.NewRequest("GET", "/", nil)
	for range w.Result().Cookies() {
		req.AddCookie(&http.Cookie{Name: "test", Value: "corrupteddata"})
	}

	_, err := ReadEncrypted(req, "test", secretKey)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}
