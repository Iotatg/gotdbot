package gotdbot

import (
	"errors"
	"testing"
)

// v0.14.2 - reply-target helpers.
//
// The reply target's ids are already on the message, so all of this is testable
// without a client - which is the point: "is this a reply, and to what" is a
// question every moderation handler asks and none of them had an answer to.

func TestMessageIsReplyReadsTheTargetOffTheMessage(t *testing.T) {
	m := &Message{
		ChatId: -1004379943083,
		Id:     900,
		ReplyTo: &MessageReplyToMessage{
			ChatId:    -1004379943083,
			MessageId: 512,
		},
	}

	chatID, messageID, ok := m.MessageIsReply()
	if !ok {
		t.Fatal("a reply was not recognised as one")
	}
	if chatID != -1004379943083 || messageID != 512 {
		t.Errorf("got %d/%d, want -1004379943083/512", chatID, messageID)
	}
}

func TestMessageIsReplyUsesTheTargetsOwnChat(t *testing.T) {
	// A cross-chat reply - a comment reply, a forum topic - lives in a different
	// chat from the reply. Looking the message up in the replying message's chat
	// returns nothing, so the reply target's own chat id has to win when present.
	m := &Message{
		ChatId: -1004379943083,
		Id:     900,
		ReplyTo: &MessageReplyToMessage{
			ChatId:    -1009999999999,
			MessageId: 7,
		},
	}

	chatID, messageID, ok := m.MessageIsReply()
	if !ok {
		t.Fatal("a cross-chat reply was not recognised")
	}
	if chatID != -1009999999999 || messageID != 7 {
		t.Errorf("got %d/%d, want -1009999999999/7", chatID, messageID)
	}
}

func TestMessageIsReplyFallsBackToTheMessagesChat(t *testing.T) {
	// TDLib leaves chat_id 0 for a reply inside the same chat, so the replying
	// message's own chat is the only thing to go on.
	m := &Message{
		ChatId:  -1004379943083,
		ReplyTo: &MessageReplyToMessage{MessageId: 512},
	}

	chatID, messageID, ok := m.MessageIsReply()
	if !ok {
		t.Fatal("a same-chat reply was not recognised")
	}
	if chatID != -1004379943083 || messageID != 512 {
		t.Errorf("got %d/%d, want -1004379943083/512", chatID, messageID)
	}
}

func TestMessageIsReplyRejectsEverythingElse(t *testing.T) {
	cases := []struct {
		name string
		m    *Message
	}{
		{"no reply at all", &Message{ChatId: 7}},
		{"a reply to a story", &Message{ChatId: 7, ReplyTo: &MessageReplyToStory{StoryId: 3}}},
		// TDLib documents message_id 0 as "in an unknown chat". Message 0 is not a
		// message, and GetMessage(0) would fail with a confusing error instead.
		{"a reply target of zero", &Message{ChatId: 7, ReplyTo: &MessageReplyToMessage{MessageId: 0}}},
		{"a negative reply target", &Message{ChatId: 7, ReplyTo: &MessageReplyToMessage{MessageId: -1}}},
		{"a typed nil reply target", &Message{ChatId: 7, ReplyTo: (*MessageReplyToMessage)(nil)}},
		{"no message at all", nil},
	}
	for _, c := range cases {
		if _, _, ok := c.m.MessageIsReply(); ok {
			t.Errorf("%s was accepted as a reply", c.name)
		}
	}
}

func TestReplyTargetUserIDRequiresARealUser(t *testing.T) {
	// A message sent on behalf of a chat, or by a channel, has no user behind it.
	// 0 is a value a caller could pass to a moderation request as though it were
	// somebody, so this returns false rather than 0.
	cases := []struct {
		name   string
		msg    *Message
		wantID int64
		wantOK bool
	}{
		{"a user sender", &Message{SenderId: &MessageSenderUser{UserId: 777}}, 777, true},
		{"a chat sender", &Message{SenderId: &MessageSenderChat{ChatId: 7}}, 0, false},
		{"no sender at all", &Message{}, 0, false},
		{"a nil message", nil, 0, false},
	}
	for _, c := range cases {
		id, ok := replyTargetUserID(c.msg)
		if id != c.wantID || ok != c.wantOK {
			t.Errorf("%s gave %d/%v, want %d/%v", c.name, id, ok, c.wantID, c.wantOK)
		}
	}
}

func TestRepliedMessageReportsAMissingTargetAsNoReply(t *testing.T) {
	// Cannot reach a client, so this checks the branch that needs none: the
	// message is not a reply, and the caller must get ErrNoReply rather than a
	// panic or a nil-dereference on the way to a round trip.
	m := &Message{ChatId: 7}
	if _, err := m.RepliedMessage(nil); !errors.Is(err, ErrNoReply) {
		t.Errorf("err = %v, want ErrNoReply", err)
	}
	if _, _, err := (&Message{}).ReplyToSenderID(nil); err != nil {
		t.Errorf("ReplyToSenderID on a non-reply returned %v, want no error", err)
	}
}
