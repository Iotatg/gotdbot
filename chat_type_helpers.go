package gotdbot

// The authoritative chat-type predicates.
//
// Message.IsGroup, Message.IsSupergroupOrChannel and Message.IsChannel all work
// off the chat id alone, because a Message carries no chat type. That works for
// separating private chats from groups, and it is genuinely ambiguous in exactly
// one direction: a supergroup and a channel both have a -100... id, so no id
// arithmetic can tell them apart. Message.IsChannel has therefore always been an
// alias for IsSupergroupOrChannel, and true for every supergroup in Telegram.
//
// Its doc comment claimed otherwise, which is worse than having no comment: the
// name reads like an exact answer and the answer was wrong. Read *Chat.Type when
// the difference matters, through the predicates below.

// ChatTypeIsBasicGroup reports whether the chat is a basic group.
func (c *Chat) ChatTypeIsBasicGroup() bool {
	_, ok := c.Type.(*ChatTypeBasicGroup)
	return ok
}

// ChatTypeIsSupergroup reports whether the chat is a supergroup.
func (c *Chat) ChatTypeIsSupergroup() bool {
	_, ok := c.Type.(*ChatTypeSupergroup)
	return ok
}

// ChatTypeIsPrivate reports whether the chat is a private chat with a user.
func (c *Chat) ChatTypeIsPrivate() bool {
	_, ok := c.Type.(*ChatTypePrivate)
	return ok
}

// ChatTypeIsSecret reports whether the chat is a secret chat.
func (c *Chat) ChatTypeIsSecret() bool {
	_, ok := c.Type.(*ChatTypeSecret)
	return ok
}

// ChatTypeIsSupergroupOrChannel reports whether the chat is a supergroup or a
// broadcast channel.
//
// TDLib reports both as chatTypeSupergroup, and there is no chat type that means
// "channel" alone - so this is the most specific answer the chat type can give,
// and it is also why a separate ChatTypeIsChannel does not exist here. That is not
// an oversight to be tidied up later; it is what the type is.
func (c *Chat) ChatTypeIsSupergroupOrChannel() bool {
	_, ok := c.Type.(*ChatTypeSupergroup)
	return ok
}

// MessageChatType returns the message's chat type, reading it from Telegram.
//
// Every id-based helper on Message is a guess in one direction or the other;
// this is the one that is not. It costs a round trip, which is why the guesses
// exist.
func (m *Message) MessageChatType(c *Client) (ChatType, error) {
	chat, err := c.GetChat(m.ChatId)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, ErrChatNotResolved
	}
	return chat.Type, nil
}

// MessageChatTypeIsBasicGroup reports whether the message is from a basic group,
// authoritatively.
func (m *Message) MessageChatTypeIsBasicGroup(c *Client) (bool, error) {
	chatType, err := m.MessageChatType(c)
	if err != nil {
		return false, err
	}
	_, ok := chatType.(*ChatTypeBasicGroup)
	return ok, nil
}

// MessageChatTypeIsSupergroupOrChannel reports whether the message is from a
// supergroup or a channel, authoritatively.
func (m *Message) MessageChatTypeIsSupergroupOrChannel(c *Client) (bool, error) {
	chatType, err := m.MessageChatType(c)
	if err != nil {
		return false, err
	}
	switch chatType.(type) {
	case *ChatTypeSupergroup, *ChatTypeBasicGroup:
		return true, nil
	}
	return false, nil
}
