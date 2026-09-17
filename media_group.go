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
