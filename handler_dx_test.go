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
