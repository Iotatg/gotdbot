package gotdbot

func FileFromPath(path string) InputFile {
	return &InputFileLocal{Path: path}
}

func FileFromID(id int32) InputFile {
	return &InputFileId{Id: id}
}

func FileFromRemote(id string) InputFile {
	return &InputFileRemote{Id: id}
}

func InlineKeyboard(rows ...[]InlineKeyboardButton) *ReplyMarkupInlineKeyboard {
	return &ReplyMarkupInlineKeyboard{Rows: rows}
}

func InlineRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	return buttons
}

func CallbackButton(text, data string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeCallback{Data: []byte(data)},
	}
}

func URLButton(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeUrl{Url: url},
	}
}

func WebAppButton(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeWebApp{Url: url},
	}
}

func UserButton(text string, userId int64) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeUser{UserId: userId},
	}
}

func CopyTextButton(text, copyText string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeCopyText{Text: copyText},
	}
}

func SwitchInlineButton(text, query string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeSwitchInline{
			Query:      query,
			TargetChat: &TargetChatCurrent{},
		},
	}
}

func SwitchInlineChosenButton(text, query string, allowUsers, allowBots, allowGroups, allowChannels bool) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeSwitchInline{
			Query: query,
			TargetChat: &TargetChatChosen{
				Types: &TargetChatTypes{
					AllowBotChats:     allowBots,
					AllowUserChats:    allowUsers,
					AllowGroupChats:   allowGroups,
					AllowChannelChats: allowChannels,
				},
			},
		},
	}
}

func LoginURLButton(text, url string, id int64, forwardText string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeLoginUrl{
			Url:         url,
			Id:          id,
			ForwardText: forwardText,
		},
	}
}

func GameButton(text string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeCallbackGame{},
	}
}

func ReplyKeyboard(rows ...[]KeyboardButton) *ReplyMarkupShowKeyboard {
	return &ReplyMarkupShowKeyboard{
		Rows:           rows,
		ResizeKeyboard: true,
	}
}

func KeyboardRow(buttons ...KeyboardButton) []KeyboardButton {
	return buttons
}

func TextButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
		Type: &KeyboardButtonTypeText{},
	}
}

func RequestContactButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
		Type: &KeyboardButtonTypeRequestPhoneNumber{},
	}
}

func RequestLocationButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
		Type: &KeyboardButtonTypeRequestLocation{},
	}
}

func KeyboardWebAppButton(text, url string) KeyboardButton {
	return KeyboardButton{
		Text: text,
		Type: &KeyboardButtonTypeWebApp{Url: url},
	}
}

func RemoveKeyboard() *ReplyMarkupRemoveKeyboard {
	return &ReplyMarkupRemoveKeyboard{}
}

func ForceReplyMarkup(placeholder string, personal bool) *ReplyMarkupForceReply {
	return &ReplyMarkupForceReply{
		InputFieldPlaceholder: placeholder,
		IsPersonal:            personal,
	}
}

func DisabledButton(text string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Type: &InlineKeyboardButtonTypeDisabled{},
	}
}

type RequestUsersButtonOpts struct {
	MaxQuantity           int32
	RequestName           bool
	RequestPhoto          bool
	RequestUsername       bool
	RestrictUserIsBot     bool
	RestrictUserIsPremium bool
	UserIsBot             bool
	UserIsPremium         bool
}

func RequestUsersButton(text string, id int32, opts *RequestUsersButtonOpts) KeyboardButton {
	btn := KeyboardButtonTypeRequestUsers{Id: id, MaxQuantity: 1}
	if opts != nil {
		btn.MaxQuantity = opts.MaxQuantity
		btn.RequestName = opts.RequestName
		btn.RequestPhoto = opts.RequestPhoto
		btn.RequestUsername = opts.RequestUsername
		btn.RestrictUserIsBot = opts.RestrictUserIsBot
		btn.RestrictUserIsPremium = opts.RestrictUserIsPremium
		btn.UserIsBot = opts.UserIsBot
		btn.UserIsPremium = opts.UserIsPremium
	}
	if btn.MaxQuantity <= 0 {
		btn.MaxQuantity = 1
	}
	return KeyboardButton{Text: text, Type: &btn}
}

type RequestChatButtonOpts struct {
	BotAdministratorRights  *ChatAdministratorRights
	BotIsMember             bool
	ChatHasUsername         bool
	ChatIsChannel           bool
	ChatIsCreated           bool
	ChatIsForum             bool
	RequestPhoto            bool
	RequestTitle            bool
	RequestUsername         bool
	RestrictChatHasUsername bool
	RestrictChatIsForum     bool
	UserAdministratorRights *ChatAdministratorRights
}

func RequestChatButton(text string, id int32, opts *RequestChatButtonOpts) KeyboardButton {
	btn := KeyboardButtonTypeRequestChat{Id: id}
	if opts != nil {
		btn.BotAdministratorRights = opts.BotAdministratorRights
		btn.BotIsMember = opts.BotIsMember
		btn.ChatHasUsername = opts.ChatHasUsername
		btn.ChatIsChannel = opts.ChatIsChannel
		btn.ChatIsCreated = opts.ChatIsCreated
		btn.ChatIsForum = opts.ChatIsForum
		btn.RequestPhoto = opts.RequestPhoto
		btn.RequestTitle = opts.RequestTitle
		btn.RequestUsername = opts.RequestUsername
		btn.RestrictChatHasUsername = opts.RestrictChatHasUsername
		btn.RestrictChatIsForum = opts.RestrictChatIsForum
		btn.UserAdministratorRights = opts.UserAdministratorRights
	}
	return KeyboardButton{Text: text, Type: &btn}
}

func RequestPollButton(text string, quiz bool) KeyboardButton {
	return KeyboardButton{
		Text: text,
		Type: &KeyboardButtonTypeRequestPoll{
			ForceQuiz:    quiz,
			ForceRegular: !quiz,
		},
	}
}

func PackedCallbackButton(text, action string, args ...string) (InlineKeyboardButton, error) {
	packed, err := PackCallback(action, args...)
	if err != nil {
		return InlineKeyboardButton{}, err
	}
	return CallbackButton(text, packed), nil
}

func MenuButton(text, url string) *BotMenuButton {
	return &BotMenuButton{Text: text, Url: url}
}
