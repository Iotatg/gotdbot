package gotdbot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Iotatg/gotdbot/internal/tdjson"
)

// Invoke sends a raw TDLib request and returns the decoded response.
//
// It is the generic escape hatch for TDLib methods that have no generated
// wrapper yet, mirroring Pyrogram's client.invoke. Prefer the generated
// helpers when they exist; use Invoke when they do not.
//
//	raw, err := c.Invoke(&GetMessage{ChatId: chat, MessageId: id})
//	msg, ok := raw.(*Message)
func (c *Client) Invoke(req TlObject) (TlObject, error) {
	return c.InvokeWithContext(context.Background(), req)
}

// InvokeWithContext is Invoke with caller-supplied context.
func (c *Client) InvokeWithContext(ctx context.Context, req TlObject) (TlObject, error) {
	if req == nil {
		return nil, fmt.Errorf("gotdbot: Invoke requires a non-nil request")
	}
	return c.SendWithContext(ctx, req)
}

// InvokeJSON sends a raw JSON-encoded TDLib request and returns the decoded
// response. It accepts methods that are not present in this package at all,
// which makes it possible to use brand-new TDLib functionality before the
// bindings are regenerated.
//
//	resp, err := c.InvokeJSON(`{"@type":"getMe","@extra":"1"}`)
func (c *Client) InvokeJSON(request string) (TlObject, error) {
	return c.InvokeJSONWithContext(context.Background(), request)
}

// InvokeJSONWithContext is InvokeJSON with caller-supplied context.
func (c *Client) InvokeJSONWithContext(ctx context.Context, request string) (TlObject, error) {
	trimmed := strings.TrimSpace(request)
	if trimmed == "" {
		return nil, fmt.Errorf("gotdbot: InvokeJSON requires a non-empty request")
	}

	fn := newRawJSONRequest(trimmed)
	// The dispatcher keys pending responses off @extra, so a function request
	// is required to receive the reply.
	fn.setExtra(c.nextRawRequestID())
	ch := make(chan TlObject, 1)
	c.pendingRequests.Store(fn.Extra, ch)

	select {
	case <-ctx.Done():
		c.pendingRequests.Delete(fn.Extra)
		go func() { <-ch }()
		return nil, ctx.Err()
	default:
	}

	tdjson.SendBytes(c.clientID, []byte(trimmed))

	select {
	case res := <-ch:
		c.pendingRequests.Delete(fn.Extra)
		if errObj, isErr := res.(*Error); isErr {
			return nil, errObj
		}
		return res, nil
	case <-ctx.Done():
		c.pendingRequests.Delete(fn.Extra)
		go func() { <-ch }()
		return nil, ctx.Err()
	case <-time.After(invokeJSONTimeout):
		c.pendingRequests.Delete(fn.Extra)
		go func() { <-ch }()
		return nil, SendTimeout
	}
}

// rawJSONRequest adapts a pre-serialized request to the tlFunction contract
// so it can share the pending-request machinery with generated functions.
type rawJSONRequest struct {
	Extra string `json:"@extra,omitempty"`
	raw   string
}

func newRawJSONRequest(raw string) *rawJSONRequest { return &rawJSONRequest{raw: raw} }

func (t *rawJSONRequest) setExtra(extra string) { t.Extra = extra }

func (t *rawJSONRequest) GetType() string { return "rawJsonRequest" }

// MarshalJSON re-injects @extra into the caller's JSON document so TDLib
// echoes it back on the matching response.
func (t *rawJSONRequest) MarshalJSON() ([]byte, error) {
	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(t.raw), &doc); err != nil {
		return nil, err
	}
	extra, err := json.Marshal(t.Extra)
	if err != nil {
		return nil, err
	}
	doc["@extra"] = extra
	return json.Marshal(doc)
}

var rawRequestCounter atomic.Uint64

func (c *Client) nextRawRequestID() string {
	return fmt.Sprintf("raw%d", rawRequestCounter.Add(1))
}

// invokeJSONTimeout bounds a raw JSON request that produced no response.
const invokeJSONTimeout = 30 * time.Second

// ── In-memory uploads ───────────────────────────────────────────────────────

// UploadBytes uploads data to Telegram and returns the resulting File.
//
// TDLib has no inputFileBytes equivalent: file parameters are always paths on
// disk, and TDLib itself reads the content. So the bytes are staged in a
// temporary file inside the client UploadDir and then handed to TDLib's
// uploadFile, which returns a File id usable in any inputFile slot.
//
// The staged file is always removed, including on error.
func (c *Client) UploadBytes(data []byte, fileName string) (*File, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("gotdbot: UploadBytes requires non-empty data")
	}
	if fileName == "" {
		fileName = "upload.bin"
	}

	path, err := c.stageUpload(fileName, data)
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)

	// TDLib resolves relative file parameters against its file directory; the
	// client mirrors that layout so the reference below is always valid.
	ref := filepath.ToSlash(path)
	if !filepath.IsAbs(ref) {
		ref = "file/" + ref
	}

	raw, err := c.InvokeJSONWithContext(
		context.Background(),
		fmt.Sprintf(
			`{"@type":"uploadFile","file":%s,"file_name":%s,"file_size":%d}`,
			mustJSONString(ref), mustJSONString(fileName), len(data),
		),
	)
	if err != nil {
		return nil, err
	}
	file, ok := raw.(*File)
	if !ok || file == nil {
		return nil, fmt.Errorf("gotdbot: uploadFile returned %T, want *File", raw)
	}
	return file, nil
}

// InputFileBytes uploads data and returns an InputFile ready to drop into any
// inputFile field (for example inputStickerSetItem.document or
// inputPhoto.photo).
//
// The name mirrors the Pyrogram concept this replaces; the returned value is
// an *InputFileId because TDLib always addresses uploads by id.
func (c *Client) InputFileBytes(data []byte, fileName string) (InputFile, error) {
	file, err := c.UploadBytes(data, fileName)
	if err != nil {
		return nil, err
	}
	return InputFileId{Id: file.Id}, nil
}

// stageUpload writes data into the client upload directory and returns the
// created path. The name is sanitised here rather than in the callers so the
// disk write itself can never escape the upload directory.
func (c *Client) stageUpload(fileName string, data []byte) (string, error) {
	dir := c.uploadDir()
	// filepath.Base on a slash-rooted path collapses any traversal segments.
	safeName := filepath.Base(filepath.Clean("/" + fileName))
	if safeName == "." || safeName == string(filepath.Separator) || safeName == "" {
		return "", fmt.Errorf("gotdbot: invalid upload file name %q", fileName)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("gotdbot: cannot create upload dir %q: %w", dir, err)
	}
	path := filepath.Join(dir, safeName)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("gotdbot: cannot stage upload: %w", err)
	}
	return path, nil
}

// uploadDir resolves where staged uploads live. FilesDirectory wins because it
// is the directory TDLib was actually initialised with, which is the only one
// guaranteed to resolve relative file parameters correctly.
func (c *Client) uploadDir() string {
	if c == nil || c.config == nil {
		return "file"
	}
	if c.config.FilesDirectory != "" {
		return filepath.Join(c.config.FilesDirectory, "file")
	}
	if c.config.UploadDir != "" {
		return c.config.UploadDir
	}
	return "file"
}

func mustJSONString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return `""`
	}
	return string(b)
}
