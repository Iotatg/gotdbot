package gotdbot

func (c *Client) OnCallbackQuery(handler func(client *Client, update *UpdateNewCallbackQuery) error, filter func(u *UpdateNewCallbackQuery) bool) {
	c.OnUpdateNewCallbackQuery(handler, filter)
}

func (c *Client) OnCallbackQueryGroup(handler func(client *Client, update *UpdateNewCallbackQuery) error, filter func(u *UpdateNewCallbackQuery) bool, group int) {
	c.AddUpdateNewCallbackQueryHandlerGroup(handler, filter, group)
}

func (c *Client) OnInlineQuery(handler func(client *Client, update *UpdateNewInlineQuery) error, filter func(u *UpdateNewInlineQuery) bool) {
	c.OnUpdateNewInlineQuery(handler, filter)
}

func (c *Client) OnChosenInlineResult(handler func(client *Client, update *UpdateNewChosenInlineResult) error, filter func(u *UpdateNewChosenInlineResult) bool) {
	c.OnUpdateNewChosenInlineResult(handler, filter)
}

func (c *Client) OnEditedMessage(handler func(client *Client, update *UpdateMessageEdited) error, filter func(u *UpdateMessageEdited) bool) {
	c.OnUpdateMessageEdited(handler, filter)
}

func (c *Client) OnDeletedMessages(handler func(client *Client, update *UpdateDeleteMessages) error, filter func(u *UpdateDeleteMessages) bool) {
	c.OnUpdateDeleteMessages(handler, filter)
}

func (c *Client) OnChatMember(handler func(client *Client, update *UpdateChatMember) error, filter func(u *UpdateChatMember) bool) {
	c.OnUpdateChatMember(handler, filter)
}

func (c *Client) OnChatJoinRequest(handler func(client *Client, update *UpdateNewChatJoinRequest) error, filter func(u *UpdateNewChatJoinRequest) bool) {
	c.OnUpdateNewChatJoinRequest(handler, filter)
}

func (c *Client) OnPoll(handler func(client *Client, update *UpdatePoll) error, filter func(u *UpdatePoll) bool) {
	c.OnUpdatePoll(handler, filter)
}

func (c *Client) OnPollAnswer(handler func(client *Client, update *UpdatePollAnswer) error, filter func(u *UpdatePollAnswer) bool) {
	c.OnUpdatePollAnswer(handler, filter)
}

func (c *Client) OnPreCheckoutQuery(handler func(client *Client, update *UpdateNewPreCheckoutQuery) error, filter func(u *UpdateNewPreCheckoutQuery) bool) {
	c.OnUpdateNewPreCheckoutQuery(handler, filter)
}

func (c *Client) OnShippingQuery(handler func(client *Client, update *UpdateNewShippingQuery) error, filter func(u *UpdateNewShippingQuery) bool) {
	c.OnUpdateNewShippingQuery(handler, filter)
}

func (c *Client) OnBusinessConnection(handler func(client *Client, update *UpdateBusinessConnection) error, filter func(u *UpdateBusinessConnection) bool) {
	c.OnUpdateBusinessConnection(handler, filter)
}

func (c *Client) OnBusinessMessage(handler func(client *Client, update *UpdateNewBusinessMessage) error, filter func(u *UpdateNewBusinessMessage) bool) {
	c.OnUpdateNewBusinessMessage(handler, filter)
}

func (c *Client) OnChatBoost(handler func(client *Client, update *UpdateChatBoost) error, filter func(u *UpdateChatBoost) bool) {
	c.OnUpdateChatBoost(handler, filter)
}
