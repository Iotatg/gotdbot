package gotdbot

import (
	"encoding/json"
	"testing"
)

// v0.14.2 - Send must not silently drop a request it cannot correlate.
//
// setExtra is unexported, so no request type defined outside this package can
// satisfy tlFunction. Before v0.14.2, SendWithContext sent such a request and
// returned (nil, nil): the call had already reached Telegram, and the caller was
// told nothing - no response, no error, and a return value that looks like
// success. Hand-written request types are how a caller reaches a TDLib method
// that has not been generated yet, so this was the escape hatch failing
// silently at exactly the moment it was needed.
//
// The whole end-to-end path cannot be exercised here, because tdjson.SendBytes
// goes through cgo and there is nothing to intercept. What is testable is the
// part that actually decides whether the request is usable: the document handed
// to TDLib has to carry an @extra the dispatcher can key on, with every other
// field intact.

// foreignRequest stands in for a hand-written request type. It can satisfy
// TlObject and nothing more, which is the whole point.
type foreignRequest struct {
	ChatID int64  `json:"chat_id"`
	Title  string `json:"custom_title"`
	Canary int64  `json:"canary"`
}

func (f *foreignRequest) GetType() string { return "setChatAdministratorCustomTitle" }

// foreignWithExtra is the awkward case: a caller who already carries an @extra
// field, as some generated types do.
type foreignWithExtra struct {
	Extra  string `json:"@extra,omitempty"`
	ChatID int64  `json:"chat_id"`
}

func (f *foreignWithExtra) GetType() string { return "getChat" }

func decodeFields(t *testing.T, raw []byte) map[string]json.RawMessage {
	t.Helper()
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatalf("marshalWithExtra produced invalid JSON: %v\n%s", err, raw)
	}
	return fields
}

func TestMarshalWithExtraGivesAForeignRequestACorrelationID(t *testing.T) {
	raw, err := marshalWithExtra(&foreignRequest{ChatID: -1004379943083, Title: "Full Admin", Canary: 1<<60 + 1}, "42")
	if err != nil {
		t.Fatalf("marshalWithExtra = %v", err)
	}

	fields := decodeFields(t, raw)
	extra, ok := fields["@extra"]
	if !ok {
		t.Fatalf("no @extra was added, so the response cannot be matched: %s", raw)
	}
	var got string
	if err := json.Unmarshal(extra, &got); err != nil {
		t.Fatalf("@extra is not a string: %v", err)
	}
	if got != "42" {
		t.Errorf("@extra = %q, want %q", got, "42")
	}
}

func TestMarshalWithExtraPreservesFieldValuesExactly(t *testing.T) {
	// Two separate claims, and only the second is about precision.
	//
	// A supergroup chat id is about 1e12, which float64 represents exactly, so on
	// its own it cannot detect a value being routed through float64 - an earlier
	// version of this test asserted that it could, and it could not. The Canary
	// field is the canary for the mechanism: 1<<60+1 is past 2^53, so any route
	// through an interface{} map, a hand-built value, or anything else that
	// re-encodes numbers as floats changes it. It is not a Telegram value; it is
	// there to fail loudly if the field-preserving path is ever swapped out.
	const (
		chatID = "-1004379943083"
		canary = "1152921504606846977" // 1<<60 + 1
	)

	raw, err := marshalWithExtra(&foreignRequest{
		ChatID: -1004379943083,
		Title:  "Full Admin",
		Canary: 1<<60 + 1,
	}, "7")
	if err != nil {
		t.Fatalf("marshalWithExtra = %v", err)
	}

	fields := decodeFields(t, raw)
	for _, c := range []struct{ name, want string }{
		{"chat_id", chatID},
		{"canary", canary},
	} {
		if got := string(fields[c.name]); got != c.want {
			t.Errorf("%s = %s, want %s", c.name, got, c.want)
		}
	}
	if got := string(fields["custom_title"]); got != `"Full Admin"` {
		t.Errorf("title = %s, want the string it was given", got)
	}
}

func TestMarshalWithExtraOverridesACallerSuppliedExtra(t *testing.T) {
	// If the caller's stale extra survived, the dispatcher would key on an ID
	// that is not the one it was stored under, and the response would be dropped
	// on the floor - a hang instead of a wrong answer.
	raw, err := marshalWithExtra(&foreignWithExtra{Extra: "stale", ChatID: 5}, "9")
	if err != nil {
		t.Fatalf("marshalWithExtra = %v", err)
	}

	var got string
	if err := json.Unmarshal(decodeFields(t, raw)["@extra"], &got); err != nil {
		t.Fatalf("@extra is not a string: %v", err)
	}
	if got != "9" {
		t.Errorf("@extra = %q, want the dispatcher's own %q", got, "9")
	}
}

func TestMarshalWithExtraRejectsARequestThatIsNotAnObject(t *testing.T) {
	// A request that marshalled to a bare string or null cannot be correlated,
	// and sending it would be the old silent failure again.
	bad := &notAnObject{}
	if _, err := marshalWithExtra(bad, "1"); err == nil {
		t.Error("marshalWithExtra accepted a non-object request")
	}
}

type notAnObject struct{}

func (n *notAnObject) GetType() string { return "nonsense" }

func (n *notAnObject) MarshalJSON() ([]byte, error) { return []byte(`"nope"`), nil }

// TestMarshalWithExtraMatchesAGeneratedRequestShape is the tie between the two
// paths: a hand-written request must reach TDLib looking exactly like a generated
// one, because nothing downstream can tell them apart.
func TestMarshalWithExtraMatchesAGeneratedRequestShape(t *testing.T) {
	raw, err := marshalWithExtra(&foreignRequest{ChatID: 7, Title: "x", Canary: 1<<60 + 1}, "1")
	if err != nil {
		t.Fatalf("marshalWithExtra = %v", err)
	}
	var got struct {
		Type   string `json:"@type"`
		Extra  string `json:"@extra"`
		ChatID int64  `json:"chat_id"`
		Title  string `json:"custom_title"`
		Canary int64  `json:"canary"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Type != "setChatAdministratorCustomTitle" {
		t.Errorf("@type = %q, want the type the request's GetType reported", got.Type)
	}
	if got.Extra != "1" || got.ChatID != 7 || got.Title != "x" || got.Canary != 1<<60+1 {
		t.Errorf("a field was lost or altered: %+v", got)
	}
}
