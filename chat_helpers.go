package gotdbot

import "strings"

func (u *User) FullName() string {
	if u == nil {
		return ""
	}
	name := strings.TrimSpace(u.FirstName + " " + u.LastName)
	return name
}

func (u *User) Username() string {
	if u == nil || u.Usernames == nil || len(u.Usernames.ActiveUsernames) == 0 {
		return ""
	}
	return u.Usernames.ActiveUsernames[0]
}

func (u *User) IsBot() bool {
	if u == nil || u.Type == nil {
		return false
	}
	_, ok := u.Type.(*UserTypeBot)
	return ok
}

func (u *User) MentionHTML() string {
	if u == nil {
		return ""
	}
	return Mention(u.FullName(), u.Id, true, true)
}

func (u *User) MentionMarkdown() string {
	if u == nil {
		return ""
	}
	return Mention(u.FullName(), u.Id, false, true)
}

func (ch *Chat) IsPrivate() bool {
	if ch == nil || ch.Type == nil {
		return false
	}
	_, ok := ch.Type.(*ChatTypePrivate)
	return ok
}

func (ch *Chat) IsBasicGroup() bool {
	if ch == nil || ch.Type == nil {
		return false
	}
	_, ok := ch.Type.(*ChatTypeBasicGroup)
	return ok
}

func (ch *Chat) IsSupergroup() bool {
	if ch == nil || ch.Type == nil {
		return false
	}
	t, ok := ch.Type.(*ChatTypeSupergroup)
	return ok && !t.IsChannel
}

func (ch *Chat) IsChannelChat() bool {
	if ch == nil || ch.Type == nil {
		return false
	}
	t, ok := ch.Type.(*ChatTypeSupergroup)
	return ok && t.IsChannel
}

func (ch *Chat) IsSecret() bool {
	if ch == nil || ch.Type == nil {
		return false
	}
	_, ok := ch.Type.(*ChatTypeSecret)
	return ok
}

func (ch *Chat) IsGroupChat() bool {
	return ch.IsBasicGroup() || ch.IsSupergroup()
}

func UserSender(userId int64) MessageSender {
	return &MessageSenderUser{UserId: userId}
}

func ChatSender(chatId int64) MessageSender {
	return &MessageSenderChat{ChatId: chatId}
}
