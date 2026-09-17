package gotdbot

func (c *Client) GetMediaGroup(chatId, messageId int64) ([]Message, error) {
	if messageId <= 0 {
		return nil, ErrInvalidMessageID
	}

	ids := make([]int64, 0, 19)
	for id := messageId - 9; id <= messageId+9; id++ {
		if id > 0 {
			ids = append(ids, id)
		}
	}

	res, err := c.GetMessages(chatId, ids)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrNotMediaGroup
	}

	var albumID int64
	found := false
	for i := range res.Messages {
		msg := &res.Messages[i]
		if msg.Id == messageId {
			albumID = msg.MediaAlbumId
			found = true
			break
		}
	}
	if !found || albumID == 0 {
		return nil, ErrNotMediaGroup
	}

	group := make([]Message, 0, 10)
	for i := range res.Messages {
		msg := res.Messages[i]
		if msg.Id != 0 && msg.MediaAlbumId == albumID {
			group = append(group, msg)
		}
	}
	if len(group) == 0 {
		return nil, ErrNotMediaGroup
	}
	return group, nil
}

func (m *Message) GetMediaGroup(c *Client) ([]Message, error) {
	if m == nil {
		return nil, ErrInvalidMessageID
	}
	return c.GetMediaGroup(m.ChatId, m.Id)
}

func mediaGroupIDs(messages []Message) []int64 {
	ids := make([]int64, 0, len(messages))
	for i := range messages {
		if messages[i].Id != 0 {
			ids = append(ids, messages[i].Id)
		}
	}
	return ids
}

func copyForwardOpts(opts *ForwardMessagesOpts, sendCopy bool) *ForwardMessagesOpts {
	out := &ForwardMessagesOpts{}
	if opts != nil {
		*out = *opts
	}
	if sendCopy {
		out.SendCopy = true
	}
	return out
}

func (c *Client) CopyMediaGroup(chatId, fromChatId, messageId int64, opts *ForwardMessagesOpts) (*Messages, error) {
	group, err := c.GetMediaGroup(fromChatId, messageId)
	if err != nil {
		return nil, err
	}
	return c.ForwardMessages(chatId, fromChatId, mediaGroupIDs(group), copyForwardOpts(opts, true))
}

func (c *Client) ForwardMediaGroup(chatId, fromChatId, messageId int64, opts *ForwardMessagesOpts) (*Messages, error) {
	group, err := c.GetMediaGroup(fromChatId, messageId)
	if err != nil {
		return nil, err
	}
	return c.ForwardMessages(chatId, fromChatId, mediaGroupIDs(group), opts)
}

func (m *Message) CopyMediaGroup(c *Client, chatId int64, opts *ForwardMessagesOpts) (*Messages, error) {
	if m == nil {
		return nil, ErrInvalidMessageID
	}
	return c.CopyMediaGroup(chatId, m.ChatId, m.Id, opts)
}

func (m *Message) ForwardMediaGroup(c *Client, chatId int64, opts *ForwardMessagesOpts) (*Messages, error) {
	if m == nil {
		return nil, ErrInvalidMessageID
	}
	return c.ForwardMediaGroup(chatId, m.ChatId, m.Id, opts)
}
