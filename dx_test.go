package gotdbot_test

import (
	"errors"
	"testing"

	"github.com/Iotatg/gotdbot"
	"github.com/Iotatg/gotdbot/filters/message"
)

func TestMessageContentAccessors(t *testing.T) {
	photo := &gotdbot.Message{
		Content: &gotdbot.MessagePhoto{Photo: &gotdbot.Photo{}},
	}
	if !photo.HasPhoto() || !photo.HasMedia() || photo.HasVideo() {
		t.Fatalf("photo accessors failed")
	}

	text := &gotdbot.Message{
		Content: &gotdbot.MessageText{Text: &gotdbot.FormattedText{Text: "hello"}},
	}
	if !text.HasText() || text.HasMedia() || text.HasPhoto() {
		t.Fatalf("text accessors failed")
	}

	album := &gotdbot.Message{
		MediaAlbumId: 99,
		Content:      &gotdbot.MessagePhoto{Photo: &gotdbot.Photo{}, HasSpoiler: true},
	}
	if !album.IsAlbum() || !album.HasSpoiler() {
		t.Fatalf("album/spoiler accessors failed")
	}

	fwd := &gotdbot.Message{ForwardInfo: &gotdbot.MessageForwardInfo{}}
	if !fwd.IsForwarded() {
		t.Fatalf("forwarded accessor failed")
	}

	forum := &gotdbot.Message{TopicId: &gotdbot.MessageTopicForum{ForumTopicId: 7}}
	if !forum.IsForum() || forum.ForumTopicID() != 7 || forum.IsDirect() {
		t.Fatalf("forum accessors failed")
	}
}

func TestMessageFilters(t *testing.T) {
	photo := &gotdbot.Message{Content: &gotdbot.MessagePhoto{Photo: &gotdbot.Photo{}}}
	if !message.Photo(photo) || !message.Media(photo) || message.Video(photo) {
		t.Fatalf("photo filter failed")
	}

	cmd := &gotdbot.Message{
		Content: &gotdbot.MessageText{
			Text: &gotdbot.FormattedText{
				Text: "/start@MyBot extra",
				Entities: []gotdbot.TextEntity{
					{Offset: 0, Length: 6, Type: &gotdbot.TextEntityTypeBotCommand{}},
				},
			},
		},
	}
	if !message.Command(cmd) || !message.Commands("start")(cmd) || message.Commands("help")(cmd) {
		t.Fatalf("command filter failed")
	}

	from := &gotdbot.Message{SenderId: &gotdbot.MessageSenderUser{UserId: 42}, ChatId: 100}
	if !message.User(42)(from) || message.User(1)(from) || !message.Chat(100)(from) {
		t.Fatalf("user/chat filter failed")
	}

	svc := &gotdbot.Message{Content: &gotdbot.MessagePinMessage{MessageId: 1}}
	if !message.Service(svc) || !message.PinnedMessage(svc) {
		t.Fatalf("service filter failed")
	}
}

func TestBuilders(t *testing.T) {
	kb := gotdbot.InlineKeyboard(
		gotdbot.InlineRow(
			gotdbot.CallbackButton("A", "a"),
			gotdbot.URLButton("B", "https://example.com"),
		),
	)
	if len(kb.Rows) != 1 || len(kb.Rows[0]) != 2 {
		t.Fatalf("inline keyboard rows")
	}
	cb, ok := kb.Rows[0][0].Type.(*gotdbot.InlineKeyboardButtonTypeCallback)
	if !ok || string(cb.Data) != "a" {
		t.Fatalf("callback button")
	}

	rk := gotdbot.ReplyKeyboard(gotdbot.KeyboardRow(gotdbot.TextButton("Hi"), gotdbot.RequestContactButton("Phone")))
	if len(rk.Rows) != 1 || !rk.ResizeKeyboard {
		t.Fatalf("reply keyboard")
	}

	f := gotdbot.FileFromPath("/tmp/a.jpg")
	local, ok := f.(*gotdbot.InputFileLocal)
	if !ok || local.Path != "/tmp/a.jpg" {
		t.Fatalf("file from path")
	}

	remote := gotdbot.FileFromRemote("AgAD")
	if _, ok := remote.(*gotdbot.InputFileRemote); !ok {
		t.Fatalf("file from remote")
	}
}

func TestUserChatHelpers(t *testing.T) {
	u := &gotdbot.User{
		Id:        9,
		FirstName: "Iota",
		LastName:  "coder",
		Usernames: &gotdbot.Usernames{ActiveUsernames: []string{"iotacoder"}},
		Type:      &gotdbot.UserTypeBot{},
	}
	if u.FullName() != "Iota coder" || u.Username() != "iotacoder" || !u.IsBot() {
		t.Fatalf("user helpers")
	}
	if u.MentionHTML() == "" || u.MentionMarkdown() == "" {
		t.Fatalf("user mention")
	}

	ch := &gotdbot.Chat{Type: &gotdbot.ChatTypeSupergroup{IsChannel: true}}
	if !ch.IsChannelChat() || ch.IsSupergroup() || ch.IsPrivate() {
		t.Fatalf("channel chat helpers")
	}
	sg := &gotdbot.Chat{Type: &gotdbot.ChatTypeSupergroup{IsChannel: false}}
	if !sg.IsSupergroup() || !sg.IsGroupChat() || sg.IsChannelChat() {
		t.Fatalf("supergroup helpers")
	}
}

func TestStopPropagationAlias(t *testing.T) {
	if !errors.Is(gotdbot.StopPropagation, gotdbot.EndGroups) {
		t.Fatalf("StopPropagation should alias EndGroups")
	}
}

func TestSenderHelpers(t *testing.T) {
	s := gotdbot.UserSender(5)
	u, ok := s.(*gotdbot.MessageSenderUser)
	if !ok || u.UserId != 5 {
		t.Fatalf("user sender")
	}
}

func TestParseCommand(t *testing.T) {
	p := gotdbot.ParseCommand(`/start@MyBot extra "quoted arg"`, "/")
	if p == nil || p.Name != "start" || p.Mention != "MyBot" {
		t.Fatalf("command parse: %+v", p)
	}
	if len(p.Args) != 2 || p.Args[0] != "extra" || p.Args[1] != "quoted arg" {
		t.Fatalf("quoted args: %#v", p.Args)
	}
	if gotdbot.ParseCommand("hello", "/") != nil {
		t.Fatalf("non-command should be nil")
	}
	bang := gotdbot.ParseCommand("!help me", "!")
	if bang == nil || bang.Name != "help" || len(bang.Args) != 1 || bang.Args[0] != "me" {
		t.Fatalf("custom prefix: %+v", bang)
	}

	msg := &gotdbot.Message{
		Content: &gotdbot.MessageText{Text: &gotdbot.FormattedText{Text: "/ping@Bot a b"}},
	}
	if msg.Command() != "ping" || msg.CommandMention() != "Bot" {
		t.Fatalf("message command helpers")
	}
	if got := msg.CommandArgsQuoted(); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("message quoted args: %#v", got)
	}
}

func TestInputMediaBuilders(t *testing.T) {
	photo := gotdbot.InputMediaPhotoPath("/tmp/a.jpg")
	p, ok := photo.(*gotdbot.InputMessagePhoto)
	if !ok || p.Photo == nil {
		t.Fatalf("photo media")
	}
	video := gotdbot.InputMediaVideoPath("/tmp/a.mp4")
	if _, ok := video.(*gotdbot.InputMessageVideo); !ok {
		t.Fatalf("video media")
	}
	doc := gotdbot.InputMediaDocumentPath("/tmp/a.bin")
	if _, ok := doc.(*gotdbot.InputMessageDocument); !ok {
		t.Fatalf("document media")
	}
}

func TestFloodWaitHelpers(t *testing.T) {
	err := &gotdbot.Error{Code: 429, Message: "Too Many Requests: retry after 7"}
	if !gotdbot.IsFloodWait(err) || gotdbot.FloodWaitSeconds(err) != 7 {
		t.Fatalf("flood wait parse")
	}
	if gotdbot.IsFloodWait(errors.New("nope")) || gotdbot.FloodWaitSeconds(nil) != 0 {
		t.Fatalf("non flood")
	}
}

func TestParseModeNormalize(t *testing.T) {
	if gotdbot.NormalizeParseMode("html") != gotdbot.ParseModeHTML {
		t.Fatalf("html")
	}
	if gotdbot.NormalizeParseMode("mdv2") != gotdbot.ParseModeMarkdownV2 {
		t.Fatalf("mdv2")
	}
	if gotdbot.NormalizeParseMode("plain") != gotdbot.ParseModeNone {
		t.Fatalf("none")
	}
}

func TestFileProgressPercent(t *testing.T) {
	f := &gotdbot.File{Size: 200, Local: &gotdbot.LocalFile{DownloadedSize: 50}}
	if f.ProgressPercent() != 25 {
		t.Fatalf("progress percent: %v", f.ProgressPercent())
	}
}

func TestCommandsWithPrefixFilter(t *testing.T) {
	msg := &gotdbot.Message{
		Content: &gotdbot.MessageText{Text: &gotdbot.FormattedText{Text: "!ban user"}},
	}
	if !message.CommandsWithPrefix([]string{"!"}, "ban")(msg) {
		t.Fatalf("prefix command filter")
	}
	if message.CommandsWithPrefix([]string{"!"}, "kick")(msg) {
		t.Fatalf("wrong command")
	}
}
