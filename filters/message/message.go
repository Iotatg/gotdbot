package message

import (
	"regexp"
	"strings"

	"github.com/Iotatg/gotdbot"
	"github.com/Iotatg/gotdbot/filters"
)

// Prefix checks if the message text starts with the given prefix.
func Prefix(prefix string) filters.Message {
	return func(msg *gotdbot.Message) bool {
		return msg.StartsWith(prefix)
	}
}

// Suffix checks if the message text ends with the given suffix.
func Suffix(suffix string) filters.Message {
	return func(msg *gotdbot.Message) bool {
		return msg.EndsWith(suffix)
	}
}

// Equal checks if the message text equals the given match string.
func Equal(match string) filters.Message {
	return func(msg *gotdbot.Message) bool {
		return msg.GetText() == match
	}
}

// Regex checks if the message text matches the given regex pattern.
func Regex(pattern string) filters.Message {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return func(msg *gotdbot.Message) bool {
			return false
		}
	}
	return func(msg *gotdbot.Message) bool {
		return reg.MatchString(msg.GetText())
	}
}

// Private filters messages sent in private chats.
func Private(msg *gotdbot.Message) bool {
	return msg.IsPrivate()
}

// Group filters messages sent in group or supergroup chats.
func Group(msg *gotdbot.Message) bool {
	return msg.IsGroup() || msg.IsSupergroupOrChannel()
}

// Channel filters messages sent in channel chats.
func Channel(msg *gotdbot.Message) bool {
	return msg.IsChannel()
}

// Incoming filters incoming messages.
func Incoming(msg *gotdbot.Message) bool {
	return !msg.IsOutgoingMessage()
}

// Outgoing filters outgoing messages.
func Outgoing(msg *gotdbot.Message) bool {
	return msg.IsOutgoingMessage()
}

// Edited filters edited messages.
func Edited(msg *gotdbot.Message) bool {
	return msg.IsEditedMessage()
}

func Command(msg *gotdbot.Message) bool {
	return msg.IsCommand()
}

func Text(msg *gotdbot.Message) bool            { return msg.HasText() }
func Caption(msg *gotdbot.Message) bool         { return msg.HasCaption() }
func Photo(msg *gotdbot.Message) bool           { return msg.HasPhoto() }
func Video(msg *gotdbot.Message) bool           { return msg.HasVideo() }
func Animation(msg *gotdbot.Message) bool       { return msg.HasAnimation() }
func Audio(msg *gotdbot.Message) bool           { return msg.HasAudio() }
func Document(msg *gotdbot.Message) bool        { return msg.HasDocument() }
func Sticker(msg *gotdbot.Message) bool         { return msg.HasSticker() }
func Voice(msg *gotdbot.Message) bool           { return msg.HasVoice() }
func VideoNote(msg *gotdbot.Message) bool       { return msg.HasVideoNote() }
func Contact(msg *gotdbot.Message) bool         { return msg.HasContact() }
func Location(msg *gotdbot.Message) bool        { return msg.HasLocation() }
func Venue(msg *gotdbot.Message) bool           { return msg.HasVenue() }
func Poll(msg *gotdbot.Message) bool            { return msg.HasPoll() }
func Dice(msg *gotdbot.Message) bool            { return msg.HasDice() }
func Game(msg *gotdbot.Message) bool            { return msg.HasGame() }
func Story(msg *gotdbot.Message) bool           { return msg.HasStory() }
func Invoice(msg *gotdbot.Message) bool         { return msg.HasInvoice() }
func PaidMedia(msg *gotdbot.Message) bool       { return msg.HasPaidMedia() }
func Checklist(msg *gotdbot.Message) bool       { return msg.HasChecklist() }
func Gift(msg *gotdbot.Message) bool            { return msg.HasGift() }
func Giveaway(msg *gotdbot.Message) bool        { return msg.HasGiveaway() }
func GiveawayWinners(msg *gotdbot.Message) bool { return msg.GiveawayWinners() != nil }
func UsersShared(msg *gotdbot.Message) bool     { return msg.UsersShared() != nil }
func ChatShared(msg *gotdbot.Message) bool      { return msg.ChatShared() != nil }
func Media(msg *gotdbot.Message) bool           { return msg.HasMedia() }
func MediaGroup(msg *gotdbot.Message) bool      { return msg.IsAlbum() }
func MediaSpoiler(msg *gotdbot.Message) bool    { return msg.HasSpoiler() }
func Forwarded(msg *gotdbot.Message) bool       { return msg.IsForwarded() }
func Reply(msg *gotdbot.Message) bool           { return msg.IsReply() }
func ViaBot(msg *gotdbot.Message) bool          { return msg.IsViaBot() }
func Scheduled(msg *gotdbot.Message) bool       { return msg.IsScheduled() }
func Service(msg *gotdbot.Message) bool         { return msg.IsService() }
func Forum(msg *gotdbot.Message) bool           { return msg.IsForum() }
func Direct(msg *gotdbot.Message) bool          { return msg.IsDirect() }
func ChannelPost(msg *gotdbot.Message) bool     { return msg.IsChannelPostMessage() }
func InlineKeyboard(msg *gotdbot.Message) bool {
	if msg == nil || msg.ReplyMarkup == nil {
		return false
	}
	_, ok := msg.ReplyMarkup.(*gotdbot.ReplyMarkupInlineKeyboard)
	return ok
}
func ReplyKeyboard(msg *gotdbot.Message) bool {
	if msg == nil || msg.ReplyMarkup == nil {
		return false
	}
	_, ok := msg.ReplyMarkup.(*gotdbot.ReplyMarkupShowKeyboard)
	return ok
}

func Commands(names ...string) filters.Message {
	return CommandsWithPrefix([]string{"/"}, names...)
}

func CommandsWithPrefix(prefixes []string, names ...string) filters.Message {
	set := make(map[string]struct{}, len(names))
	for _, n := range names {
		set[strings.ToLower(strings.TrimPrefix(n, "/"))] = struct{}{}
	}
	return func(msg *gotdbot.Message) bool {
		parsed := gotdbot.ParseCommand(msg.GetText(), prefixes...)
		if parsed == nil {
			return false
		}
		_, ok := set[strings.ToLower(parsed.Name)]
		return ok
	}
}

func User(ids ...int64) filters.Message {
	set := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return func(msg *gotdbot.Message) bool {
		_, ok := set[msg.SenderID()]
		return ok
	}
}

func Chat(ids ...int64) filters.Message {
	set := make(map[int64]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return func(msg *gotdbot.Message) bool {
		_, ok := set[msg.ChatID()]
		return ok
	}
}

func Topic(ids ...int32) filters.Message {
	set := make(map[int32]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return func(msg *gotdbot.Message) bool {
		_, ok := set[msg.ForumTopicID()]
		return ok
	}
}

func NewChatMembers(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessageChatAddMembers)
	return ok
}

func LeftChatMember(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessageChatDeleteMember)
	return ok
}

func PinnedMessage(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessagePinMessage)
	return ok
}

func NewChatTitle(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessageChatChangeTitle)
	return ok
}

func NewChatPhoto(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessageChatChangePhoto)
	return ok
}

func DeleteChatPhoto(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessageChatDeletePhoto)
	return ok
}

func VideoChatStarted(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessageVideoChatStarted)
	return ok
}

func VideoChatEnded(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessageVideoChatEnded)
	return ok
}

func SuccessfulPayment(msg *gotdbot.Message) bool {
	_, ok := msg.Content.(*gotdbot.MessagePaymentSuccessful)
	return ok
}

func WebAppData(msg *gotdbot.Message) bool {
	return msg.WebAppData() != nil
}

func Me(msg *gotdbot.Message) bool {
	return msg != nil && msg.IsOutgoingMessage()
}

func Mentioned(msg *gotdbot.Message) bool {
	return msg.IsMentioned()
}

func LiveLocation(msg *gotdbot.Message) bool {
	return msg.HasLiveLocation()
}

func WebPage(msg *gotdbot.Message) bool {
	return msg.HasLinkPreview()
}

func Quote(msg *gotdbot.Message) bool {
	return msg.HasQuote()
}

func Business(msg *gotdbot.Message) bool {
	return msg.IsBusiness()
}

func Ephemeral(msg *gotdbot.Message) bool {
	return msg.IsEphemeral()
}

// FromUserID checks if the message sender's user ID matches the given ID.
func FromUserID(id int64) filters.Message {
	return func(msg *gotdbot.Message) bool {
		return msg.SenderID() == id
	}
}

// ChatID checks if the message chat ID matches the given ID.
func ChatID(id int64) filters.Message {
	return func(msg *gotdbot.Message) bool {
		return msg.ChatID() == id
	}
}

// And combines multiple filters, all must pass.
func And(fs ...filters.Message) filters.Message {
	return func(msg *gotdbot.Message) bool {
		for _, f := range fs {
			if !f(msg) {
				return false
			}
		}
		return true
	}
}

// Or combines multiple filters, at least one must pass.
func Or(fs ...filters.Message) filters.Message {
	return func(msg *gotdbot.Message) bool {
		for _, f := range fs {
			if f(msg) {
				return true
			}
		}
		return false
	}
}

// Not negates a filter.
func Not(f filters.Message) filters.Message {
	return func(msg *gotdbot.Message) bool {
		return !f(msg)
	}
}
