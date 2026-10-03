package gotdbot

import "testing"

// v0.14.1: AddHandlerGroup used to dereference the result of
// c.handlers.Load() without a nil check. The handler table is an
// atomic.Pointer that only NewClient initialises, so registering any handler on
// a zero-value Client panicked. RemoveHandlerGroup already guarded against nil;
// registration did not.
//
// The library's own tests never caught this because their helper builds the
// client by hand and stores an empty table first.

func TestRegisterOnZeroValueClientDoesNotPanic(t *testing.T) {
	c := &Client{}

	// Every registration entry point has to survive, because they all funnel
	// into AddHandlerGroup.
	c.OnCommand("ping", func(client *Client, m *Message) error { return nil })
	c.OnCommandGroup("pong", func(client *Client, m *Message) error { return nil }, 5)
	c.OnMessage(func(client *Client, m *Message) error { return nil }, nil)
	c.OnCallbackQuery(func(client *Client, q *UpdateNewCallbackQuery) error { return nil },
		func(q *UpdateNewCallbackQuery) bool { return true })
	c.OnChatMember(func(client *Client, u *UpdateChatMember) error { return nil },
		func(u *UpdateChatMember) bool { return true })
	c.OnRawUpdate(func(client *Client, u TlObject) error { return nil }, nil)
	c.AddHandler(NewCommandHandler("extra", func(client *Client, m *Message) error { return nil }))

	data := c.handlers.Load()
	if data == nil {
		t.Fatal("handler table is still nil after registration")
	}
	total := 0
	for _, list := range data.handlers {
		total += len(list)
	}
	if total != 7 {
		t.Errorf("registered %d handlers, want 7", total)
	}
}

func TestRegisterOnZeroValueClientKeepsGroupNumbers(t *testing.T) {
	c := &Client{}
	c.OnCommand("a", func(client *Client, m *Message) error { return nil })
	c.OnCommandGroup("b", func(client *Client, m *Message) error { return nil }, 3)

	data := c.handlers.Load()
	if got := len(data.handlers[0]); got != 1 {
		t.Errorf("group 0 has %d handlers, want 1", got)
	}
	if got := len(data.handlers[3]); got != 1 {
		t.Errorf("group 3 has %d handlers, want 1", got)
	}
}

func TestRemoveHandlerOnEmptyClientReturnsFalse(t *testing.T) {
	c := &Client{}
	if c.RemoveHandler(NewCommandHandler("gone", func(client *Client, m *Message) error { return nil })) {
		t.Error("removing from an empty client must report false, not panic")
	}
}

// A second registration on the same zero-value client must not reset what the
// first one stored.
func TestRepeatedRegistrationOnZeroValueClientAccumulates(t *testing.T) {
	c := &Client{}
	for i := 0; i < 3; i++ {
		c.OnCommand("x", func(client *Client, m *Message) error { return nil })
	}
	data := c.handlers.Load()
	if got := len(data.handlers[0]); got != 3 {
		t.Errorf("group 0 has %d handlers, want 3", got)
	}
}
