package gotdbot

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// v0.14.3 - the raw escape hatch actually reaching TDLib.
//
// InvokeJSON built its own pending-request entry under an @extra it had allocated,
// then sent the caller's JSON string with tdjson.SendBytes directly. rawJSONRequest
// has a MarshalJSON whose whole purpose is to inject that @extra, and nothing ever
// called it - the bytes sent were the caller's, unmodified.
//
// TDLib therefore received no @extra and echoed none. The dispatcher sees an empty
// extra and hands the reply to the update channel, so the caller waiting on the
// pending entry waited out the full 30 seconds and got SendTimeout. Every
// InvokeJSON call failed, for every method, including the ones it was written for.
//
// Nothing caught it because nothing could: tdjson.SendBytes is cgo, so a unit
// test cannot see the bytes that were sent, and the only visible symptom was a
// timeout that looked exactly like TDLib being unresponsive.

// TestRawJSONRequestMarshalsTheExtraItWasGiven is the property the old path
// skipped. With it, Send takes the generated branch for rawJSONRequest and the
// document TDLib receives carries @extra, so the reply can be matched.
func TestRawJSONRequestMarshalsTheExtraItWasGiven(t *testing.T) {
	req := newRawJSONRequest(`{"@type":"getMe"}`)
	req.setExtra("raw7")

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(raw), `"@extra":"raw7"`) {
		t.Errorf("the document TDLib would receive carries no @extra: %s", raw)
	}

	var got struct {
		Type  string `json:"@type"`
		Extra string `json:"@extra"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Extra != "raw7" {
		t.Errorf("@extra = %q, want raw7", got.Extra)
	}
	// The caller's own fields must survive the injection.
	if got.Type != "getMe" {
		t.Errorf("@type = %q, want getMe", got.Type)
	}
}

// TestRawJSONRequestOverridesACallersExtra pins that the allocator wins. A stale
// extra left in place would key the pending entry on an id the dispatcher never
// sees, which is the same lost-reply failure again.
func TestRawJSONRequestOverridesACallersExtra(t *testing.T) {
	req := newRawJSONRequest(`{"@type":"getMe","@extra":"stale"}`)
	req.setExtra("raw8")

	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got struct {
		Extra string `json:"@extra"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Extra != "raw8" {
		t.Errorf("@extra = %q, want raw8", got.Extra)
	}
}

// TestOnlyClientSendsToTDLib is the guard the old code needed and did not have.
//
// One function writes to TDLib. Anything else that does so is re-implementing the
// pending-request, timeout and retry machinery, and the failure mode is the one
// that just happened: a send that looks right and never comes back. If a second
// call site is ever genuinely needed, this test should be the thing that has to be
// argued with.
func TestOnlyClientSendsToTDLib(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	fset := token.NewFileSet()

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(".", name)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok || sel.Sel.Name != "SendBytes" {
					return true
				}
				if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "tdjson" {
					return true
				}
				if name != "client.go" || fn.Name.Name != "SendWithContext" {
					t.Errorf("%s:%s writes to TDLib directly. Only "+
						"Client.SendWithContext may, so that @extra is assigned and "+
						"the reply can be matched.",
						name, fn.Name.Name)
				}
				return true
			})
		}
	}
}

// TestSendWithContextDoesNotSendOnADeadContext keeps a cancelled caller from
// mutating a chat. The pre-check is what makes that true, and the old raw path had
// it while Send did not - so delegating to Send without it would have been a
// regression.
func TestSendWithContextDoesNotSendOnADeadContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := &Client{}

	if _, err := c.SendWithContext(ctx, &GetMe{}); err == nil {
		t.Error("SendWithContext returned no error for a cancelled context")
	}
}
