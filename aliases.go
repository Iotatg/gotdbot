package gotdbot

func (c *Client) BanUser(chatId, userId int64, bannedUntilDate int32, revokeMessages bool) error {
	return c.BanChatMember(bannedUntilDate, chatId, UserSender(userId), &BanChatMemberOpts{RevokeMessages: revokeMessages})
}

func (c *Client) UnbanUser(chatId, userId int64) error {
	return c.SetChatMemberStatus(chatId, UserSender(userId), &ChatMemberStatusLeft{})
}

func (c *Client) KickUser(chatId, userId int64) error {
	return c.SetChatMemberStatus(chatId, UserSender(userId), &ChatMemberStatusLeft{})
}

func (c *Client) RestrictUser(chatId, userId int64, permissions *ChatPermissions, restrictedUntilDate int32) error {
	return c.SetChatMemberStatus(chatId, UserSender(userId), &ChatMemberStatusRestricted{
		IsMember:            true,
		Permissions:         permissions,
		RestrictedUntilDate: restrictedUntilDate,
	})
}

type PromoteOpts struct {
	Rights      *ChatAdministratorRights
	CanBeEdited *bool
}

func DefaultAdminRights() *ChatAdministratorRights {
	return &ChatAdministratorRights{
		CanChangeInfo:          true,
		CanDeleteMessages:      true,
		CanInviteUsers:         true,
		CanPinMessages:         true,
		CanManageVideoChats:    true,
		CanRestrictMembers:     true,
		CanManageChat:          true,
		CanManageTopics:        true,
		CanPostStories:         true,
		CanEditStories:         true,
		CanDeleteStories:       true,
		CanSendWelcomeMessages: true,
	}
}

func promoteSettings(opts *PromoteOpts) (*ChatAdministratorRights, bool) {
	rights := DefaultAdminRights()
	canBeEdited := true
	if opts != nil {
		if opts.Rights != nil {
			rights = opts.Rights
		}
		if opts.CanBeEdited != nil {
			canBeEdited = *opts.CanBeEdited
		}
	}
	return rights, canBeEdited
}

func (c *Client) PromoteUser(chatId, userId int64, opts *PromoteOpts) error {
	rights, canBeEdited := promoteSettings(opts)
	return c.SetChatMemberStatus(chatId, UserSender(userId), &ChatMemberStatusAdministrator{
		CanBeEdited: canBeEdited,
		Rights:      rights,
	})
}

func (c *Client) GetChatMemberByUser(chatId, userId int64) (*ChatMember, error) {
	return c.GetChatMember(chatId, UserSender(userId))
}

func (c *Client) ResolveUsername(username string) (*Chat, error) {
	username = publicUsername(username)
	if username == "" {
		return nil, nil
	}
	return c.SearchPublicChat(username)
}

func (c *Client) SendMediaGroup(chatId int64, contents []InputMessageContent, opts *SendMessageAlbumOpts) (*Messages, error) {
	return c.SendMessageAlbum(chatId, contents, opts)
}

func (c *Client) DownloadMedia(msg *Message, synchronous bool) (*File, error) {
	if msg == nil {
		return nil, nil
	}
	return msg.Download(c, 1, 0, 0, synchronous)
}

func (c *Client) GetChatHistoryPage(chatId int64, fromMessageId int64, limit int32) ([]Message, error) {
	limit = clampLimit(limit, 100, 100)
	res, err := c.GetChatHistory(chatId, fromMessageId, limit, 0, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.Messages, nil
}

func (c *Client) IterChatHistory(chatId int64, fromMessageId int64, limit int32, fn func(msg Message) bool) error {
	if fn == nil {
		return nil
	}
	limit = clampLimit(limit, 100, 100)
	cursor := fromMessageId
	skipCursor := false
	for {
		page, err := c.GetChatHistoryPage(chatId, cursor, limit)
		if err != nil {
			return err
		}
		if len(page) == 0 {
			return nil
		}
		next := cursor
		yielded := 0
		for _, msg := range page {
			if msg.Id == 0 || (skipCursor && msg.Id == cursor) {
				continue
			}
			if !fn(msg) {
				return nil
			}
			next = msg.Id
			yielded++
		}
		if yielded == 0 || next == cursor || int32(len(page)) < limit {
			return nil
		}
		cursor = next
		skipCursor = true
	}
}

func (m *Message) ReplyMarkupKeyboard() ReplyMarkup {
	if m == nil {
		return nil
	}
	return m.ReplyMarkup
}

func (m *Message) Promote(c *Client, opts *PromoteOpts) error {
	return c.PromoteUser(m.ChatId, m.SenderID(), opts)
}

func (m *Message) DownloadSync(c *Client) (*File, error) {
	return m.Download(c, 1, 0, 0, true)
}

func (m *Message) LocalPath(c *Client) (string, error) {
	f, err := m.DownloadSync(c)
	if err != nil || f == nil || f.Local == nil {
		return "", err
	}
	return f.Local.Path, nil
}
