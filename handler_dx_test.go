package gotdbot

import (
	"errors"
	"testing"
)

func newTestClient() *Client {
	c := &Client{}
	c.handlers.Store(&handlersData{handlers: make(map[int][]Handler)})
	return c
}

func TestPluginLoadExcludeAndRemove(t *testing.T) {
	c := newTestClient()
	keep := NewPlugin("keep").OnCommand("start", func(client *Client, message *Message) error { return nil }, 2)
	skip := NewPlugin("skip").OnCommand("help", func(client *Client, message *Message) error { return nil }, 2)
	c.LoadPlugins([]*Plugin{keep, skip}, &PluginLoadOpts{Exclude: []string{"skip"}})

	data := c.handlers.Load()
	if len(data.handlers[2]) != 1 {
		t.Fatalf("expected 1 loaded handler, got %d", len(data.handlers[2]))
	}
	if !c.RemoveHandlerGroup(data.handlers[2][0], 2) {
		t.Fatalf("remove handler")
	}
	if len(c.handlers.Load().handlers[2]) != 0 {
		t.Fatalf("handler still present")
	}
}

func TestPluginInclude(t *testing.T) {
	c := newTestClient()
	a := NewPlugin("a").OnMessage(func(client *Client, message *Message) error { return nil }, nil, 3)
	b := NewPlugin("b").OnMessage(func(client *Client, message *Message) error { return nil }, nil, 3)
	c.LoadPlugins([]*Plugin{a, b}, &PluginLoadOpts{Include: []string{"b"}})
	if len(c.handlers.Load().handlers[3]) != 1 {
		t.Fatalf("include filter")
	}
}

func TestMiddlewareOrder(t *testing.T) {
	c := newTestClient()
	var order []int
	c.Use(func(client *Client, update TlObject, next func() error) error {
		order = append(order, 1)
		err := next()
		order = append(order, 3)
		return err
	})
	c.Use(func(client *Client, update TlObject, next func() error) error {
		order = append(order, 2)
		return next()
	})
	if err := c.runMiddlewares(nil, func() error {
		order = append(order, 4)
		return nil
	}); err != nil {
		t.Fatalf("middleware: %v", err)
	}
	want := []int{1, 2, 4, 3}
	if len(order) != len(want) {
		t.Fatalf("order %#v", order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order %#v", order)
		}
	}
}

func TestMiddlewareCanStop(t *testing.T) {
	c := newTestClient()
	called := false
	c.Use(func(client *Client, update TlObject, next func() error) error {
		return EndGroups
	})
	err := c.runMiddlewares(nil, func() error {
		called = true
		return nil
	})
	if !errors.Is(err, EndGroups) || called {
		t.Fatalf("middleware should stop chain")
	}
}

func TestGetMediaGroupInvalidID(t *testing.T) {
	c := newTestClient()
	_, err := c.GetMediaGroup(1, 0)
	if !errors.Is(err, ErrInvalidMessageID) {
		t.Fatalf("expected invalid id, got %v", err)
	}
}

func TestHandlerAliasesRegister(t *testing.T) {
	c := newTestClient()
	c.OnStory(func(client *Client, update *UpdateStory) error { return nil }, nil)
	c.OnUserStatus(func(client *Client, update *UpdateUserStatus) error { return nil }, nil)
	c.OnMessageReaction(func(client *Client, update *UpdateMessageReaction) error { return nil }, nil)
	c.OnGuest(func(client *Client, update *UpdateNewGuestQuery) error { return nil }, nil)
	c.OnError(func(client *Client, update TlObject, err error) error { return err })
	if c.errorHandler == nil {
		t.Fatalf("OnError not stored")
	}
}

func TestConnectDisconnectHandlers(t *testing.T) {
	c := newTestClient()
	connects := 0
	disconnects := 0
	c.OnConnect(func(client *Client) error {
		connects++
		return nil
	})
	c.OnDisconnect(func(client *Client) error {
		disconnects++
		return nil
	})
	if err := c.connectionStateHandler(c, &UpdateConnectionState{State: &ConnectionStateWaitingForNetwork{}}); err != nil {
		t.Fatalf("waiting before ready: %v", err)
	}
	if connects != 0 || disconnects != 0 {
		t.Fatalf("disconnect before first connect: c=%d d=%d", connects, disconnects)
	}
	if err := c.connectionStateHandler(c, &UpdateConnectionState{State: &ConnectionStateReady{}}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := c.connectionStateHandler(c, &UpdateConnectionState{State: &ConnectionStateReady{}}); err != nil {
		t.Fatalf("second ready: %v", err)
	}
	if connects != 1 {
		t.Fatalf("connect called %d times", connects)
	}
	if err := c.connectionStateHandler(c, &UpdateConnectionState{State: &ConnectionStateWaitingForNetwork{}}); err != nil {
		t.Fatalf("disconnect: %v", err)
	}
	if disconnects != 1 {
		t.Fatalf("disconnect called %d times", disconnects)
	}
	if err := c.connectionStateHandler(c, &UpdateConnectionState{State: &ConnectionStateWaitingForNetwork{}}); err != nil {
		t.Fatalf("second waiting: %v", err)
	}
	if disconnects != 1 {
		t.Fatalf("duplicate disconnect: %d", disconnects)
	}
}

func TestMediaGroupIDs(t *testing.T) {
	ids := mediaGroupIDs([]Message{{Id: 1}, {Id: 0}, {Id: 3}})
	if len(ids) != 2 || ids[0] != 1 || ids[1] != 3 {
		t.Fatalf("ids %#v", ids)
	}
}

func TestInviteLinkHelper(t *testing.T) {
	valid := []string{
		"https://t.me/+abc",
		"https://t.me/joinchat/xyz",
		"https://telegram.me/+hash",
		"+invitehash",
		"t.me/+abc",
	}
	for _, s := range valid {
		if !isInviteLink(s) {
			t.Fatalf("expected invite link: %s", s)
		}
	}
	invalid := []string{"iota", "joinchat", "https://t.me/iota", "@iota", "https://t.me/joinchatbot"}
	for _, s := range invalid {
		if isInviteLink(s) {
			t.Fatalf("false invite link: %s", s)
		}
	}
	if got := normalizeInviteLink("+abc"); got != "https://t.me/+abc" {
		t.Fatalf("normalize: %s", got)
	}
	if got := publicUsername("https://t.me/IotaCoder/123"); got != "IotaCoder" {
		t.Fatalf("public username: %s", got)
	}
}

func TestClampLimit(t *testing.T) {
	if clampLimit(0, 100, 100) != 100 || clampLimit(-1, 50, 100) != 50 || clampLimit(5, 100, 100) != 5 || clampLimit(500, 100, 100) != 100 {
		t.Fatalf("clampLimit")
	}
}

func TestPromoteSettingsDefaultCanBeEdited(t *testing.T) {
	rights, canBeEdited := promoteSettings(&PromoteOpts{})
	if rights == nil || !canBeEdited {
		t.Fatalf("empty PromoteOpts should keep default canBeEdited")
	}
	f := false
	_, canBeEdited = promoteSettings(&PromoteOpts{CanBeEdited: &f})
	if canBeEdited {
		t.Fatalf("explicit false should stick")
	}
}

func TestCopyForwardOptsDoesNotMutate(t *testing.T) {
	orig := &ForwardMessagesOpts{RemoveCaption: true}
	got := copyForwardOpts(orig, true)
	if !got.SendCopy || !got.RemoveCaption {
		t.Fatalf("copy opts: %+v", got)
	}
	if orig.SendCopy {
		t.Fatalf("caller opts mutated")
	}
}

func TestSliceMembers(t *testing.T) {
	members := []ChatMember{{}, {}, {}, {}}
	if got := sliceMembers(members, 1, 2); len(got) != 2 {
		t.Fatalf("slice len %d", len(got))
	}
	if got := sliceMembers(members, 10, 2); got != nil {
		t.Fatalf("offset past end")
	}
}
