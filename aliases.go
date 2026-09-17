package gotdbot

import (
	"strings"
)

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
	CanBeEdited bool
}

func DefaultAdminRights() *ChatAdministratorRights {
	return &ChatAdministratorRights{
		CanChangeInfo:       true,
		CanDeleteMessages:   true,
		CanInviteUsers:      true,
		CanPinMessages:      true,
		CanManageVideoChats: true,
		CanRestrictMembers:  true,
		CanManageChat:       true,
		CanManageTopics:     true,
		CanPostStories:      true,
		CanEditStories:      true,
		CanDeleteStories:    true,
	}
}

func (c *Client) PromoteUser(chatId, userId int64, opts *PromoteOpts) error {
	rights := DefaultAdminRights()
	canBeEdited := true
	if opts != nil {
		if opts.Rights != nil {
			rights = opts.Rights
		}
		canBeEdited = opts.CanBeEdited
	}
	return c.SetChatMemberStatus(chatId, UserSender(userId), &ChatMemberStatusAdministrator{
		CanBeEdited: canBeEdited,
		Rights:      rights,
	})
}

func (c *Client) GetChatMemberByUser(chatId, userId int64) (*ChatMember, error) {
	return c.GetChatMember(chatId, UserSender(userId))
}

func (c *Client) ResolveUsername(username string) (*Chat, error) {
	username = strings.TrimPrefix(strings.TrimSpace(username), "@")
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
	if limit <= 0 {
		limit = 100
	}
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
	if limit <= 0 {
		limit = 100
	}
	cursor := fromMessageId
	for {
		page, err := c.GetChatHistoryPage(chatId, cursor, limit)
		if err != nil {
			return err
		}
		if len(page) == 0 {
			return nil
		}
		for _, msg := range page {
			if !fn(msg) {
				return nil
			}
			cursor = msg.Id
		}
		if int32(len(page)) < limit {
			return nil
		}
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
