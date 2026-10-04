package gotdbot

import (
	"encoding/json"
	"strings"
	"testing"
)

// v0.14.2 - pinned messages and chat-type predicates.

func TestPinnedChatMessagesRequestShape(t *testing.T) {
	// The request has to be indistinguishable from a generated one, because
	// SendWithContext routes on @extra and TDLib routes on @type.
	raw, err := json.Marshal(&GetPinnedChatMessages{ChatId: -1004379943083, Limit: 100})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got struct {
		Type   string `json:"@type"`
		ChatID int64  `json:"chat_id"`
		Limit  int32  `json:"limit"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Type != "getPinnedChatMessages" {
		t.Errorf("@type = %q", got.Type)
	}
	if got.ChatID != -1004379943083 || got.Limit != 100 {
		t.Errorf("chat_id/limit = %d/%d", got.ChatID, got.Limit)
	}
}

func TestPinnedChatMessagesResultDecodes(t *testing.T) {
	// TDLib sends total_count and messages; if the field names are wrong the list
	// silently comes back empty and "unpin the latest" unpins nothing.
	const payload = `{"@type":"pinnedChatMessages","total_count":2,
		"messages":[{"id":11,"chat_id":7},{"id":40,"chat_id":7}]}`

	var got PinnedChatMessages
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.TotalCount != 2 {
		t.Errorf("total_count = %d, want 2", got.TotalCount)
	}
	if len(got.Messages) != 2 {
		t.Fatalf("decoded %d messages, want 2", len(got.Messages))
	}
	if got.Messages[0].Id != 11 || got.Messages[1].Id != 40 {
		t.Errorf("message ids = %d, %d; want 11, 40",
			got.Messages[0].Id, got.Messages[1].Id)
	}
}

func TestPinnedChatMessagesMarshalsBackWithItsType(t *testing.T) {
	raw, err := json.Marshal(&PinnedChatMessages{TotalCount: 0})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(raw), `"@type":"pinnedChatMessages"`) {
		t.Errorf("the result lost its @type: %s", raw)
	}
}

func TestLatestPinnedChoiceIgnoresResultOrder(t *testing.T) {
	// Telegram does not promise which end of the list a client gets first, and a
	// wrong assumption here would unpin an old announcement instead of the new one.
	ascending := []Message{{Id: 11}, {Id: 40}, {Id: 12}}
	descending := []Message{{Id: 40}, {Id: 11}, {Id: 12}}

	a, ok := latestPinnedID(ascending)
	if !ok || a != 40 {
		t.Errorf("ascending order gave %d, want 40", a)
	}
	d, ok := latestPinnedID(descending)
	if !ok || d != 40 {
		t.Errorf("descending order gave %d, want 40", d)
	}
}

func TestLatestPinnedChoiceReportsNothingPinned(t *testing.T) {
	// "Nothing is pinned" is an ordinary answer, not an error, and must not come
	// back as message id 0: that is a value UnpinChatMessage would happily send.
	if id, ok := latestPinnedID(nil); ok {
		t.Errorf("an empty list produced id %d and ok=true", id)
	}
	if id, ok := latestPinnedID([]Message{}); ok {
		t.Errorf("an empty list produced id %d and ok=true", id)
	}
}

func TestUnexpectedResponseNamesBothTypes(t *testing.T) {
	err := unexpectedResponse("getPinnedChatMessages", "pinnedChatMessages", &Ok{})
	msg := err.Error()
	for _, want := range []string{"getPinnedChatMessages", "pinnedChatMessages"} {
		if !strings.Contains(msg, want) {
			t.Errorf("the error does not mention %s: %s", want, msg)
		}
	}
}

func TestUnexpectedResponseHandlesNil(t *testing.T) {
	if err := unexpectedResponse("getPinnedChatMessages", "pinnedChatMessages", nil); err == nil {
		t.Error("a nil response produced no error")
	}
}

func TestChatTypePredicatesAreDistinct(t *testing.T) {
	group := &Chat{Type: &ChatTypeBasicGroup{}}
	super := &Chat{Type: &ChatTypeSupergroup{}}
	priv := &Chat{Type: &ChatTypePrivate{}}
	secret := &Chat{Type: &ChatTypeSecret{}}

	if !group.ChatTypeIsBasicGroup() || group.ChatTypeIsSupergroup() {
		t.Error("a basic group is not being recognised as one")
	}
	if !super.ChatTypeIsSupergroup() || super.ChatTypeIsBasicGroup() {
		t.Error("a supergroup is not being recognised as one")
	}
	if !priv.ChatTypeIsPrivate() || priv.ChatTypeIsSupergroup() {
		t.Error("a private chat is not being recognised as one")
	}
	if !secret.ChatTypeIsSecret() || secret.ChatTypeIsPrivate() {
		t.Error("a secret chat is not being recognised as one")
	}
}

func TestChatTypePredicatesTolerateNilType(t *testing.T) {
	// A chat decoded from a response that carried no type must answer false, not
	// panic, and not be mistaken for a group.
	c := &Chat{}
	if c.ChatTypeIsSupergroup() || c.ChatTypeIsBasicGroup() ||
		c.ChatTypeIsPrivate() || c.ChatTypeIsSecret() ||
		c.ChatTypeIsSupergroupOrChannel() {
		t.Error("a chat with no type matched a predicate")
	}
}

func TestChatTypePredicatesReadTheTypeNotTheID(t *testing.T) {
	// The whole point of these: a -100... id says nothing about which of
	// supergroup or private-chat-with-a-huge-id this is, and a private chat id is
	// positive while a supergroup's is not. A predicate that looked at the id
	// instead of the type would be the bug this file exists to stop repeating.
	super := &Chat{Id: -1004379943083, Type: &ChatTypeSupergroup{}}
	if !super.ChatTypeIsSupergroupOrChannel() {
		t.Error("a supergroup was not recognised")
	}
	if super.ChatTypeIsBasicGroup() {
		t.Error("a supergroup was mistaken for a basic group")
	}
}
