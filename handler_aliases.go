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

func (c *Client) OnStory(handler func(client *Client, update *UpdateStory) error, filter func(u *UpdateStory) bool) {
	c.OnUpdateStory(handler, filter)
}

func (c *Client) OnUserStatus(handler func(client *Client, update *UpdateUserStatus) error, filter func(u *UpdateUserStatus) bool) {
	c.OnUpdateUserStatus(handler, filter)
}

func (c *Client) OnMessageReaction(handler func(client *Client, update *UpdateMessageReaction) error, filter func(u *UpdateMessageReaction) bool) {
	c.OnUpdateMessageReaction(handler, filter)
}

func (c *Client) OnGuest(handler func(client *Client, update *UpdateNewGuestQuery) error, filter func(u *UpdateNewGuestQuery) bool) {
	c.OnUpdateNewGuestQuery(handler, filter)
}

func (c *Client) OnEditedBusinessMessage(handler func(client *Client, update *UpdateBusinessMessageEdited) error, filter func(u *UpdateBusinessMessageEdited) bool) {
	c.OnUpdateBusinessMessageEdited(handler, filter)
}

func (c *Client) OnDeletedBusinessMessages(handler func(client *Client, update *UpdateBusinessMessagesDeleted) error, filter func(u *UpdateBusinessMessagesDeleted) bool) {
	c.OnUpdateBusinessMessagesDeleted(handler, filter)
}

func (c *Client) OnPurchasedPaidMedia(handler func(client *Client, update *UpdatePaidMediaPurchased) error, filter func(u *UpdatePaidMediaPurchased) bool) {
	c.OnUpdatePaidMediaPurchased(handler, filter)
}

func (c *Client) OnManagedBot(handler func(client *Client, update *UpdateManagedBot) error, filter func(u *UpdateManagedBot) bool) {
	c.OnUpdateManagedBot(handler, filter)
}

func (c *Client) OnMessageReactionCount(handler func(client *Client, update *UpdateMessageReactions) error, filter func(u *UpdateMessageReactions) bool) {
	c.OnUpdateMessageReactions(handler, filter)
}

func (c *Client) OnCallback(action string, handler func(client *Client, update *UpdateNewCallbackQuery) error, filter func(u *UpdateNewCallbackQuery) bool) {
	c.OnCallbackGroup(action, handler, filter, 0)
}

func (c *Client) OnCallbackGroup(action string, handler func(client *Client, update *UpdateNewCallbackQuery) error, filter func(u *UpdateNewCallbackQuery) bool, group int) {
	c.OnCallbackQueryGroup(handler, func(u *UpdateNewCallbackQuery) bool {
		unpacked, _, err := UnpackCallback(u.DataString())
		if err != nil {
			return false
		}
		if action != "" && unpacked != action {
			return false
		}
		if filter != nil && !filter(u) {
			return false
		}
		return true
	}, group)
}

type rawUpdateHandler struct {
	filter   func(TlObject) bool
	response func(client *Client, update TlObject) error
}

func (h *rawUpdateHandler) CheckUpdate(client *Client, update TlObject) bool {
	if h.filter != nil && !h.filter(update) {
		return false
	}
	return true
}

func (h *rawUpdateHandler) HandleUpdate(client *Client, update TlObject) error {
	return h.response(client, update)
}

func (c *Client) OnRawUpdate(handler func(client *Client, update TlObject) error, filter func(TlObject) bool) {
	c.OnRawUpdateGroup(handler, filter, 0)
}

func (c *Client) OnRawUpdateGroup(handler func(client *Client, update TlObject) error, filter func(TlObject) bool, group int) {
	c.AddHandlerGroup(&rawUpdateHandler{filter: filter, response: handler}, group)
}

func (c *Client) OnError(handler func(client *Client, update TlObject, err error) error) {
	c.errorMu.Lock()
	c.errorHandler = handler
	c.errorMu.Unlock()
}

func (c *Client) OnConnect(handler func(client *Client) error) {
	if handler == nil {
		return
	}
	c.lifecycleMu.Lock()
	c.connectHandlers = append(c.connectHandlers, handler)
	c.lifecycleMu.Unlock()
}

func (c *Client) OnDisconnect(handler func(client *Client) error) {
	if handler == nil {
		return
	}
	c.lifecycleMu.Lock()
	c.disconnectHandlers = append(c.disconnectHandlers, handler)
	c.lifecycleMu.Unlock()
}
