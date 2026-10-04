package gotdbot

import "errors"

// v0.14.2 - reaching the message being replied to.
//
// A Message carries ReplyTo, and the variant that matters - messageReplyToMessage
// - already names the message: both its chat id and its message id. Nothing about
// a reply target is hidden behind a round trip.
//
// What there was not, until now, was any way to ask. ReplyTo is a
// MessageReplyTo interface with two implementations, one of which is a reply to a
// story and carries no message at all, so every caller who needed "who is this a
// reply to" had to write the type switch themselves and then decide what to do
// about the story case. Moderation commands are almost entirely about that
// question, which makes this the most-used thing in this file and the least
// obvious to get right.
//
// ErrNoReply is what "not a reply" looks like as an error, for callers who would
// rather branch on one than unpack three values.

// ErrNoReply is returned when a message is not a reply to another message.
var ErrNoReply = errors.New("the message is not a reply to another message")

// messageReplyToMessage unwraps a reply target into the variant that names a
// message.
//
// The story variant names no message, so it is not a reply as far as this file is
// concerned. MessageId may also be 0, which TDLib documents as "the replied
// message is in an unknown chat" - an unresolvable target, not message 0.
func messageReplyToMessage(replyTo MessageReplyTo) (*MessageReplyToMessage, bool) {
	in, ok := replyTo.(*MessageReplyToMessage)
	if !ok || in == nil || in.MessageId <= 0 {
		return nil, false
	}
	return in, true
}

// MessageIsReply reports whether this message replies to a specific message, and
// gives the chat it lives in along with its id.
//
// It costs nothing: both ids are already on the message. The chat id matters
// because a reply can cross chats - a comment reply in a channel, a reply in a
// forum topic - and looking the message up in the replying message's chat finds
// nothing.
func (m *Message) MessageIsReply() (chatId int64, messageId int64, ok bool) {
	if m == nil || m.ReplyTo == nil {
		return 0, 0, false
	}
	in, ok := messageReplyToMessage(m.ReplyTo)
	if !ok {
		return 0, 0, false
	}
	// A reply within the same chat carries no separate chat id.
	if in.ChatId == 0 {
		return m.ChatId, in.MessageId, true
	}
	return in.ChatId, in.MessageId, true
}

// RepliedMessage returns the message this one replies to, or ErrNoReply.
//
// ErrNoReply covers "not a reply" and "the reply target no longer exists", which
// are the same thing to a caller: there is no message to act on.
func (m *Message) RepliedMessage(c *Client) (*Message, error) {
	chatID, messageID, ok := m.MessageIsReply()
	if !ok {
		return nil, ErrNoReply
	}
	msg, err := c.GetMessage(chatID, messageID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, ErrNoReply
	}
	return msg, nil
}

// ReplyToSenderID returns the user id of whoever wrote the message this one replies
// to.
//
// Moderation commands are almost all "act on who I replied to", and each one was
// left to make this round trip itself to reach a single integer. ok is false when
// the message is not a reply, when the reply target has been deleted, and when the
// target was sent by a chat rather than a user - the cases where there is no user
// to act on, and where a caller that treated 0 as an id would be acting on
// nobody.
func (m *Message) ReplyToSenderID(c *Client) (id int64, ok bool, err error) {
	replied, err := m.RepliedMessage(c)
	if err != nil {
		if err == ErrNoReply {
			return 0, false, nil
		}
		return 0, false, err
	}
	if id, ok = replyTargetUserID(replied); !ok {
		return 0, false, nil
	}
	return id, true, nil
}

// replyTargetUserID reads the user id out of an already-fetched reply target.
//
// Split out because the rule is worth stating on its own and because it is
// otherwise only reachable with a live client: a message sent on behalf of a chat,
// or by a channel, has no user behind it, and 0 is a value a caller could then
// pass to a moderation request as though it were somebody.
func replyTargetUserID(replied *Message) (int64, bool) {
	if replied == nil {
		return 0, false
	}
	// The sender's type, not its id. Message.SenderID deliberately answers with a
	// chat id for a chat sender, so an id-only test would let a message sent by a
	// channel through as though a user had sent it - and the id it returned would
	// be the channel's, which is a real id and a plausible-looking moderation
	// target.
	sender, ok := replied.SenderId.(*MessageSenderUser)
	if !ok || sender == nil || sender.UserId <= 0 {
		return 0, false
	}
	return sender.UserId, true
}
