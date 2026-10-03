package gotdbot

import "testing"

// v0.14.1 helpers: Message.Payload, SenderUserID, SenderChatID,
// ChatMember.MemberUserID and IsMemberStatus.
//
// Each of these existed because a port hit a case where the existing API could
// not express what was needed, so the code was hand-rolled at the call site.
// Payload in particular exists because Args and ArgsList normalise away
// whitespace via strings.Fields and then re-join with single spaces.

// textMessage builds a Message whose content is plain text.
func textMessage(text string) *Message {
	return &Message{Content: &MessageText{Text: &FormattedText{Text: text}}}
}

// Payload is the equivalent of Python's
// `message.text.split(maxsplit=1)[1].strip()`.
func TestPayloadMatchesPythonSplit(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"/welcome", ""},
		{"/welcome   ", ""},
		{"/welcome hello there", "hello there"},
		{"/welcome   spaced   out  ", "spaced   out"},
		{"/welcome multi\nline", "multi\nline"},
		{"/welcome Ünïcödé 用户 🌸", "Ünïcödé 用户 🌸"},
		{"  /welcome   x  ", "x"},
		{"/st r #ff0000 s2 3", "r #ff0000 s2 3"},
		{"/petname", ""},
		{"/welcome\ta\tb", "a\tb"},
	}
	for _, c := range cases {
		if got := textMessage(c.text).Payload(); got != c.want {
			t.Errorf("Payload(%q) = %q, want %q", c.text, got, c.want)
		}
	}
}

// Payload must not corrupt multi-byte characters, which a byte-wise scan for
// spaces could do.
func TestPayloadDoesNotSplitMultiByteRunes(t *testing.T) {
	const text = "/cmd 🌸🌸🌸 end"
	if got, want := textMessage(text).Payload(), "🌸🌸🌸 end"; got != want {
		t.Errorf("Payload(%q) = %q, want %q", text, got, want)
	}
}

// Args normalises whitespace; Payload must not. This is the whole reason
// Payload exists.
func TestPayloadPreservesInteriorSpacingUnlikeArgs(t *testing.T) {
	const text = "/cmd   a    b"
	m := textMessage(text)
	if got, want := m.Payload(), "a    b"; got != want {
		t.Errorf("Payload() = %q, want %q", got, want)
	}
	if got, want := m.Args(), "a b"; got != want {
		t.Errorf("Args() = %q, want %q (Args is expected to normalise)", got, want)
	}
}

func TestPayloadOnNonTextMessage(t *testing.T) {
	m := &Message{}
	if got := m.Payload(); got != "" {
		t.Errorf("Payload() on a contentless message = %q, want %q", got, "")
	}
}

func TestSenderUserID(t *testing.T) {
	cases := []struct {
		name   string
		sender MessageSender
		wantID int64
		wantOK bool
	}{
		{"pointer user", &MessageSenderUser{UserId: 42}, 42, true},
		{"value user", MessageSenderUser{UserId: 42}, 42, true},
		{"pointer chat", &MessageSenderChat{ChatId: -100123}, -100123, false},
		{"value chat", MessageSenderChat{ChatId: -100123}, -100123, false},
		{"nil", nil, 0, false},
	}
	for _, c := range cases {
		gotID, gotOK := SenderUserID(c.sender)
		if gotID != c.wantID || gotOK != c.wantOK {
			t.Errorf("%s: SenderUserID() = (%d, %v), want (%d, %v)",
				c.name, gotID, gotOK, c.wantID, c.wantOK)
		}
	}
}

func TestSenderChatID(t *testing.T) {
	cases := []struct {
		name   string
		sender MessageSender
		wantID int64
		wantOK bool
	}{
		{"pointer chat", &MessageSenderChat{ChatId: -100123}, -100123, true},
		{"value chat", MessageSenderChat{ChatId: -100123}, -100123, true},
		{"pointer user", &MessageSenderUser{UserId: 42}, 42, false},
		{"value user", MessageSenderUser{UserId: 42}, 42, false},
		{"nil", nil, 0, false},
	}
	for _, c := range cases {
		gotID, gotOK := SenderChatID(c.sender)
		if gotID != c.wantID || gotOK != c.wantOK {
			t.Errorf("%s: SenderChatID() = (%d, %v), want (%d, %v)",
				c.name, gotID, gotOK, c.wantID, c.wantOK)
		}
	}
}

func TestMemberUserID(t *testing.T) {
	m := &ChatMember{MemberId: &MessageSenderUser{UserId: 7}}
	if got, ok := m.MemberUserID(); got != 7 || !ok {
		t.Errorf("MemberUserID() = (%d, %v), want (7, true)", got, ok)
	}

	anon := &ChatMember{MemberId: &MessageSenderChat{ChatId: -100}}
	if got, ok := anon.MemberUserID(); ok || got != -100 {
		t.Errorf("MemberUserID() on a chat sender = (%d, %v), want (-100, false)", got, ok)
	}

	var nilMember *ChatMember
	if got, ok := nilMember.MemberUserID(); got != 0 || ok {
		t.Errorf("MemberUserID() on nil = (%d, %v), want (0, false)", got, ok)
	}
}

// IsMemberStatus is the test the welcome handler needed to detect a join: the
// new status is a plain member and the old one was not.
func TestIsMemberStatus(t *testing.T) {
	if !IsMemberStatus(ChatMemberStatusMember{}) {
		t.Error("ChatMemberStatusMember must report true")
	}
	for _, s := range []ChatMemberStatus{
		ChatMemberStatusLeft{},
		ChatMemberStatusBanned{},
		ChatMemberStatusAdministrator{},
		ChatMemberStatusCreator{},
		ChatMemberStatusRestricted{},
		nil,
	} {
		if IsMemberStatus(s) {
			t.Errorf("IsMemberStatus(%T) = true, want false", s)
		}
	}
}

// The join case IsMemberStatus exists for, spelled out. A join is
// "new status is a plain member AND the old status was not", so a promotion
// from member to administrator must not register as a join.
func TestIsMemberStatusDetectsJoin(t *testing.T) {
	isJoin := func(old, new ChatMemberStatus) bool {
		return IsMemberStatus(new) && !IsMemberStatus(old)
	}

	if !isJoin(ChatMemberStatusLeft{}, ChatMemberStatusMember{}) {
		t.Error("left -> member must be detected as a join")
	}
	if !isJoin(ChatMemberStatusBanned{}, ChatMemberStatusMember{}) {
		t.Error("banned -> member must be detected as a join")
	}
	if isJoin(ChatMemberStatusMember{}, ChatMemberStatusMember{}) {
		t.Error("member -> member must not be detected as a join")
	}
	if isJoin(ChatMemberStatusLeft{}, ChatMemberStatusAdministrator{}) {
		t.Error("left -> administrator must not be detected as a join")
	}
	if isJoin(ChatMemberStatusMember{}, ChatMemberStatusAdministrator{}) {
		t.Error("a promotion must not be detected as a join")
	}
}
