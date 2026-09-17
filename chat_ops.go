package gotdbot

import (
	"strings"
)

type GetChatMembersOpts struct {
	Query        string
	Filter       ChatMembersFilter
	MemberFilter SupergroupMembersFilter
	Limit        int32
	Offset       int32
}

type CreateInviteOpts struct {
	Name               string
	ExpirationDate     int32
	MemberLimit        int32
	CreatesJoinRequest bool
}

type SearchMessagesIterOpts struct {
	MaxDate int32
	MinDate int32
	Filter  *SearchMessagesOpts
}

func (c *Client) GetChatMembers(chatId int64, opts *GetChatMembersOpts) ([]ChatMember, error) {
	if opts == nil {
		opts = &GetChatMembersOpts{}
	}
	limit := clampLimit(opts.Limit, 200, 200)
	if opts.Query != "" {
		res, err := c.SearchChatMembers(chatId, limit, opts.Query, &SearchChatMembersOpts{Filter: opts.Filter})
		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, nil
		}
		return res.Members, nil
	}

	chat, err := c.GetChat(chatId)
	if err != nil {
		return nil, err
	}
	if chat == nil || chat.Type == nil {
		return nil, nil
	}
	if sg, ok := chat.Type.(*ChatTypeSupergroup); ok {
		res, err := c.GetSupergroupMembers(limit, opts.Offset, sg.SupergroupId, &GetSupergroupMembersOpts{Filter: opts.MemberFilter})
		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, nil
		}
		return res.Members, nil
	}
	if bg, ok := chat.Type.(*ChatTypeBasicGroup); ok {
		info, err := c.GetBasicGroupFullInfo(bg.BasicGroupId)
		if err != nil {
			return nil, err
		}
		if info == nil {
			return nil, nil
		}
		return sliceMembers(info.Members, opts.Offset, limit), nil
	}
	return nil, nil
}

func (c *Client) IterChatMembers(chatId int64, opts *GetChatMembersOpts, fn func(ChatMember) bool) error {
	if fn == nil {
		return nil
	}
	if opts == nil {
		opts = &GetChatMembersOpts{}
	}
	limit := clampLimit(opts.Limit, 200, 200)

	chat, err := c.GetChat(chatId)
	if err != nil {
		return err
	}
	if chat != nil {
		if bg, ok := chat.Type.(*ChatTypeBasicGroup); ok && opts.Query == "" {
			info, err := c.GetBasicGroupFullInfo(bg.BasicGroupId)
			if err != nil {
				return err
			}
			if info == nil {
				return nil
			}
			for _, member := range info.Members {
				if !fn(member) {
					return nil
				}
			}
			return nil
		}
	}

	offset := opts.Offset
	for {
		pageOpts := *opts
		pageOpts.Limit = limit
		pageOpts.Offset = offset
		members, err := c.GetChatMembers(chatId, &pageOpts)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return nil
		}
		for _, member := range members {
			if !fn(member) {
				return nil
			}
		}
		if int32(len(members)) < limit || opts.Query != "" {
			return nil
		}
		offset += int32(len(members))
	}
}

func (c *Client) GetChatMembersCount(chatId int64) (int32, error) {
	chat, err := c.GetChat(chatId)
	if err != nil {
		return 0, err
	}
	if chat == nil || chat.Type == nil {
		return 0, nil
	}
	switch t := chat.Type.(type) {
	case *ChatTypeSupergroup:
		sg, err := c.GetSupergroup(t.SupergroupId)
		if err != nil {
			return 0, err
		}
		if sg == nil {
			return 0, nil
		}
		return sg.MemberCount, nil
	case *ChatTypeBasicGroup:
		bg, err := c.GetBasicGroup(t.BasicGroupId)
		if err != nil {
			return 0, err
		}
		if bg == nil {
			return 0, nil
		}
		return bg.MemberCount, nil
	default:
		return 0, nil
	}
}

func (c *Client) JoinChatByUsername(username string) (ChatJoinResult, error) {
	username = strings.TrimSpace(username)
	if isInviteLink(username) {
		return c.JoinChatByInviteLink(normalizeInviteLink(username))
	}
	username = publicUsername(username)
	if username == "" {
		return nil, nil
	}
	chat, err := c.ResolveUsername(username)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, nil
	}
	return c.JoinChat(chat.Id)
}

func (c *Client) JoinInvite(inviteLink string) (ChatJoinResult, error) {
	return c.JoinChatByInviteLink(normalizeInviteLink(inviteLink))
}

func (c *Client) ApproveJoinRequest(chatId, userId int64) error {
	return c.ProcessChatJoinRequest(chatId, userId, &ProcessChatJoinRequestOpts{Approve: true})
}

func (c *Client) DeclineJoinRequest(chatId, userId int64) error {
	return c.ProcessChatJoinRequest(chatId, userId, &ProcessChatJoinRequestOpts{Approve: false})
}

func (c *Client) ExportChatInvite(chatId int64) (*ChatInviteLink, error) {
	return c.ReplacePrimaryChatInviteLink(chatId)
}

func (c *Client) CreateInvite(chatId int64, opts *CreateInviteOpts) (*ChatInviteLink, error) {
	if opts == nil {
		opts = &CreateInviteOpts{}
	}
	return c.CreateChatInviteLink(chatId, opts.ExpirationDate, opts.MemberLimit, opts.Name, &CreateChatInviteLinkOpts{
		CreatesJoinRequest: opts.CreatesJoinRequest,
	})
}

func (c *Client) RevokeInvite(chatId int64, inviteLink string) (*ChatInviteLinks, error) {
	return c.RevokeChatInviteLink(chatId, inviteLink)
}

func (c *Client) BlockUser(userId int64) error {
	return c.SetMessageSenderBlockList(UserSender(userId), &SetMessageSenderBlockListOpts{BlockList: &BlockListMain{}})
}

func (c *Client) UnblockUser(userId int64) error {
	return c.SetMessageSenderBlockList(UserSender(userId), nil)
}

func (c *Client) SetAdminTitle(chatId, userId int64, title string) error {
	return c.SetChatMemberTag(chatId, title, userId)
}

func (c *Client) GetDialogs(limit int32, chatList ChatList) ([]*Chat, error) {
	if limit <= 0 {
		limit = 100
	}
	_ = c.LoadChats(limit, &LoadChatsOpts{ChatList: chatList})
	res, err := c.GetChats(limit, &GetChatsOpts{ChatList: chatList})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	out := make([]*Chat, 0, len(res.ChatIds))
	for _, id := range res.ChatIds {
		chat, err := c.GetChat(id)
		if err != nil {
			return nil, err
		}
		if chat != nil {
			out = append(out, chat)
		}
	}
	return out, nil
}

func (c *Client) IterDialogs(limit int32, chatList ChatList, fn func(*Chat) bool) error {
	if fn == nil {
		return nil
	}
	page := limit
	if page <= 0 {
		page = 100
	}
	seen := make(map[int64]struct{})
	for {
		_ = c.LoadChats(page, &LoadChatsOpts{ChatList: chatList})
		res, err := c.GetChats(int32(len(seen))+page, &GetChatsOpts{ChatList: chatList})
		if err != nil {
			return err
		}
		if res == nil || len(res.ChatIds) == 0 {
			return nil
		}
		added := 0
		for _, id := range res.ChatIds {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			added++
			chat, err := c.GetChat(id)
			if err != nil {
				return err
			}
			if chat != nil && !fn(chat) {
				return nil
			}
			if limit > 0 && int32(len(seen)) >= limit {
				return nil
			}
		}
		if added == 0 {
			return nil
		}
	}
}

func clampLimit(limit, def, max int32) int32 {
	if limit <= 0 {
		limit = def
	}
	if max > 0 && limit > max {
		return max
	}
	return limit
}

func sliceMembers(members []ChatMember, offset, limit int32) []ChatMember {
	if offset < 0 {
		offset = 0
	}
	if int(offset) >= len(members) {
		return nil
	}
	members = members[offset:]
	if limit > 0 && int(limit) < len(members) {
		members = members[:limit]
	}
	return members
}

func (c *Client) IterSearchMessages(query string, limit int32, opts *SearchMessagesIterOpts, fn func(Message) bool) error {
	if fn == nil {
		return nil
	}
	limit = clampLimit(limit, 100, 100)
	var maxDate, minDate int32
	var searchOpts *SearchMessagesOpts
	if opts != nil {
		maxDate = opts.MaxDate
		minDate = opts.MinDate
		searchOpts = opts.Filter
	}
	offset := ""
	for {
		res, err := c.SearchMessages(limit, maxDate, minDate, offset, query, searchOpts)
		if err != nil {
			return err
		}
		if res == nil || len(res.Messages) == 0 {
			return nil
		}
		for _, msg := range res.Messages {
			if !fn(msg) {
				return nil
			}
		}
		if res.NextOffset == "" || res.NextOffset == offset {
			return nil
		}
		offset = res.NextOffset
	}
}

func (c *Client) SendMessageDraft(chatId, draftId int64, text string, forumTopicId int32, opts *SendTextMessageDraftOpts) error {
	if opts == nil {
		opts = &SendTextMessageDraftOpts{}
	}
	if text != "" && opts.Text == nil {
		opts.Text = &FormattedText{Text: text}
	}
	return c.SendTextMessageDraft(chatId, draftId, forumTopicId, opts)
}

func (c *Client) VotePoll(chatId, messageId int64, optionIds []int32) error {
	return c.SetPollAnswer(chatId, messageId, optionIds)
}

func (c *Client) RetractPollVote(chatId, messageId int64) error {
	return c.SetPollAnswer(chatId, messageId, nil)
}

func isInviteLink(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "+") {
		return len(s) > 1
	}
	hosts := []string{"t.me/", "telegram.me/", "telegram.dog/"}
	for _, host := range hosts {
		if strings.Contains(s, host+"+") || strings.Contains(s, host+"joinchat/") {
			return true
		}
	}
	return false
}

func normalizeInviteLink(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "+") {
		return "https://t.me/" + s
	}
	return s
}

func publicUsername(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "@")
	lower := strings.ToLower(s)
	prefixes := []string{
		"https://t.me/",
		"http://t.me/",
		"https://telegram.me/",
		"http://telegram.me/",
		"https://telegram.dog/",
		"http://telegram.dog/",
		"t.me/",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			s = s[len(prefix):]
			break
		}
	}
	s = strings.TrimPrefix(s, "@")
	if i := strings.IndexAny(s, "/?"); i >= 0 {
		s = s[:i]
	}
	return s
}
