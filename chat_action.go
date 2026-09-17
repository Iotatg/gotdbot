package gotdbot

func (c *Client) SendAction(chatId int64, action ChatAction) error {
	return c.SendChatAction("", chatId, &SendChatActionOpts{Action: action})
}

func (c *Client) SendTyping(chatId int64) error {
	return c.SendAction(chatId, &ChatActionTyping{})
}

func (c *Client) SendUploadPhoto(chatId int64, progress int32) error {
	return c.SendAction(chatId, &ChatActionUploadingPhoto{Progress: progress})
}

func (c *Client) SendUploadVideo(chatId int64, progress int32) error {
	return c.SendAction(chatId, &ChatActionUploadingVideo{Progress: progress})
}

func (c *Client) SendUploadDocument(chatId int64, progress int32) error {
	return c.SendAction(chatId, &ChatActionUploadingDocument{Progress: progress})
}

func (c *Client) SendUploadVoice(chatId int64, progress int32) error {
	return c.SendAction(chatId, &ChatActionUploadingVoiceNote{Progress: progress})
}

func (c *Client) SendUploadVideoNote(chatId int64, progress int32) error {
	return c.SendAction(chatId, &ChatActionUploadingVideoNote{Progress: progress})
}

func (c *Client) SendRecordVideo(chatId int64) error {
	return c.SendAction(chatId, &ChatActionRecordingVideo{})
}

func (c *Client) SendRecordVoice(chatId int64) error {
	return c.SendAction(chatId, &ChatActionRecordingVoiceNote{})
}

func (c *Client) SendRecordVideoNote(chatId int64) error {
	return c.SendAction(chatId, &ChatActionRecordingVideoNote{})
}

func (c *Client) SendChooseSticker(chatId int64) error {
	return c.SendAction(chatId, &ChatActionChoosingSticker{})
}

func (c *Client) SendChooseLocation(chatId int64) error {
	return c.SendAction(chatId, &ChatActionChoosingLocation{})
}

func (c *Client) SendChooseContact(chatId int64) error {
	return c.SendAction(chatId, &ChatActionChoosingContact{})
}

func (c *Client) CancelChatAction(chatId int64) error {
	return c.SendAction(chatId, &ChatActionCancel{})
}
