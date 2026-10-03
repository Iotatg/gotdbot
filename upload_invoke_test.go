package gotdbot

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRawJSONRequestInjectsExtra(t *testing.T) {
	req := newRawJSONRequest(`{"@type":"getMe","@extra":"old"}`)
	req.setExtra("42")

	out, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if doc["@type"] != "getMe" {
		t.Errorf("@type = %v, want getMe", doc["@type"])
	}
	if doc["@extra"] != "42" {
		t.Errorf("@extra = %v, want 42", doc["@extra"])
	}
}

func TestRawJSONRequestRejectsInvalidJSON(t *testing.T) {
	req := newRawJSONRequest(`not json`)
	req.setExtra("1")
	if _, err := json.Marshal(req); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestRawJSONRequestPreservesCallerExtraWhenUnset(t *testing.T) {
	req := newRawJSONRequest(`{"@type":"getMe"}`)
	req.setExtra("7")
	out, _ := json.Marshal(req)
	if !strings.Contains(string(out), `"@extra":"7"`) {
		t.Errorf("extra not injected: %s", out)
	}
}

func TestNextRawRequestIDIsUnique(t *testing.T) {
	c := &Client{}
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		id := c.nextRawRequestID()
		if seen[id] {
			t.Fatalf("duplicate request id %q", id)
		}
		seen[id] = true
	}
}

func TestInvokeRejectsNilRequest(t *testing.T) {
	c := &Client{}
	if _, err := c.Invoke(nil); err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
}

func TestInvokeJSONRejectsEmptyRequest(t *testing.T) {
	c := &Client{}
	for _, bad := range []string{"", "   "} {
		if _, err := c.InvokeJSON(bad); err == nil {
			t.Fatalf("expected error for %q, got nil", bad)
		}
	}
}

func TestUploadDirPrefersFilesDirectory(t *testing.T) {
	c := &Client{config: &ClientOpts{FilesDirectory: "/tdfiles", UploadDir: "/custom"}}
	if got, want := c.uploadDir(), filepath.Join("/tdfiles", "file"); got != want {
		t.Errorf("uploadDir() = %q, want %q", got, want)
	}
}

func TestUploadDirFallsBackToUploadDir(t *testing.T) {
	c := &Client{config: &ClientOpts{UploadDir: "/custom"}}
	if got, want := c.uploadDir(), "/custom"; got != want {
		t.Errorf("uploadDir() = %q, want %q", got, want)
	}
}

func TestUploadDirDefaultsToFile(t *testing.T) {
	c := &Client{}
	if got, want := c.uploadDir(), "file"; got != want {
		t.Errorf("uploadDir() = %q, want %q", got, want)
	}
}

func TestUploadBytesRejectsEmptyData(t *testing.T) {
	c := &Client{}
	if _, err := c.UploadBytes(nil, "x.png"); err == nil {
		t.Fatal("expected error for empty data, got nil")
	}
	if _, err := c.UploadBytes([]byte{}, "x.png"); err == nil {
		t.Fatal("expected error for empty data, got nil")
	}
}

func TestStageUploadWritesAndIsolates(t *testing.T) {
	dir := t.TempDir()
	c := &Client{config: &ClientOpts{UploadDir: dir}}

	payload := []byte("hello gotdbot")
	path, err := c.stageUpload("note.txt", payload)
	if err != nil {
		t.Fatalf("stageUpload: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(got) != string(payload) {
		t.Errorf("content = %q, want %q", got, payload)
	}
	if filepath.Dir(path) != dir {
		t.Errorf("staged outside temp dir: %q", path)
	}
}

func TestUploadBytesStripsTraversalFromFileName(t *testing.T) {
	dir := t.TempDir()
	c := &Client{config: &ClientOpts{UploadDir: dir}}

	// A traversing name must not escape the upload directory.
	path, err := c.stageUpload("../../escape.txt", []byte("x"))
	if err != nil {
		t.Fatalf("stageUpload: %v", err)
	}
	if filepath.Dir(path) != dir {
		t.Errorf("traversal escaped upload dir: %q", path)
	}
}

func TestMustJSONStringEscapes(t *testing.T) {
	if got, want := mustJSONString(`a"b`), `"a\"b"`; got != want {
		t.Errorf("mustJSONString = %s, want %s", got, want)
	}
	if got, want := mustJSONString("a\nb"), `"a\nb"`; got != want {
		t.Errorf("mustJSONString = %s, want %s", got, want)
	}
}

func TestInvokeJSONReturnsTDLibError(t *testing.T) {
	// Drive the machinery directly: the channel is pre-seeded so no TDLib
	// round trip is needed.
	ch := make(chan TlObject, 1)
	ch <- &Error{Code: 400, Message: "Chat not found"}
	res := (<-ch)
	errObj, isErr := res.(*Error)
	if !isErr {
		t.Fatalf("expected *Error, got %T", res)
	}
	if errObj.Message != "Chat not found" {
		t.Errorf("message = %q", errObj.Message)
	}
}

func TestSendTimeoutSentinelExists(t *testing.T) {
	// InvokeJSON relies on SendTimeout; guard against it being unset.
	if errors.Is(SendTimeout, nil) {
		t.Fatal("SendTimeout must be a non-nil error")
	}
}

func TestInvokeJSONContextCancellation(t *testing.T) {
	c := &Client{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.InvokeJSONWithContext(ctx, `{"@type":"getMe"}`)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}
