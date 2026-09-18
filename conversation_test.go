package gotdbot

import (
	"errors"
	"testing"
)

func TestAskFromRejectsZeroUserID(t *testing.T) {
	c := newTestClient()
	msg, err := c.AskFrom(1, 0, nil)
	if msg != nil || !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("AskFrom zero user: msg=%v err=%v", msg, err)
	}
}

func TestListenMessageFromRejectsZeroUserID(t *testing.T) {
	c := newTestClient()
	msg, err := c.ListenMessageFrom(1, 0, nil, 0)
	if msg != nil || !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("ListenMessageFrom zero user: msg=%v err=%v", msg, err)
	}
}

func TestAskFromFilterMatchesSender(t *testing.T) {
	filter := askFromFilter(10, 42, nil)
	match := &UpdateNewMessage{Message: &Message{
		ChatId:   10,
		SenderId: &MessageSenderUser{UserId: 42},
	}}
	otherUser := &UpdateNewMessage{Message: &Message{
		ChatId:   10,
		SenderId: &MessageSenderUser{UserId: 7},
	}}
	otherChat := &UpdateNewMessage{Message: &Message{
		ChatId:   11,
		SenderId: &MessageSenderUser{UserId: 42},
	}}
	if !filter(nil, match) {
		t.Fatal("expected matching sender")
	}
	if filter(nil, otherUser) {
		t.Fatal("rejected wrong sender")
	}
	if filter(nil, otherChat) {
		t.Fatal("rejected wrong chat")
	}
	if filter(nil, &UpdateNewCallbackQuery{SenderUserId: 42}) {
		t.Fatal("rejected non-message update")
	}
}

func TestAskFromFilterUserIDZeroAnySender(t *testing.T) {
	filter := askFromFilter(10, 0, nil)
	msg := &UpdateNewMessage{Message: &Message{
		ChatId:   10,
		SenderId: &MessageSenderUser{UserId: 99},
	}}
	if !filter(nil, msg) {
		t.Fatal("userId 0 should accept any sender in chat")
	}
}

func TestAskFromFilterCustomAndCancel(t *testing.T) {
	opts := &WaitMessageOpts{
		Filter: func(m *Message) bool {
			return m.GetText() == "ok"
		},
		CancellationFilter: func(m *Message) bool {
			return m.GetText() == "cancel"
		},
	}
	filter := askFromFilter(1, 2, opts)
	base := func(text string) *UpdateNewMessage {
		return &UpdateNewMessage{Message: &Message{
			ChatId:   1,
			SenderId: &MessageSenderUser{UserId: 2},
			Content:  &MessageText{Text: &FormattedText{Text: text}},
		}}
	}
	if !filter(nil, base("ok")) {
		t.Fatal("custom filter should pass")
	}
	if filter(nil, base("nope")) {
		t.Fatal("custom filter should reject")
	}
	if !filter(nil, base("cancel")) {
		t.Fatal("cancellation filter should match")
	}
}
