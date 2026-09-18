package gotdbot

import (
	"testing"
	"time"
)

func messageUpdate(userID int64) *UpdateNewMessage {
	return &UpdateNewMessage{Message: &Message{
		SenderId: &MessageSenderUser{UserId: userID},
	}}
}

func TestRateLimitAllowsThenDrops(t *testing.T) {
	c := newTestClient()
	c.Use(RateLimit(2, time.Second))
	var n int
	next := func() error {
		n++
		return nil
	}
	u := messageUpdate(9)
	if err := c.runMiddlewares(u, next); err != nil {
		t.Fatalf("first: %v", err)
	}
	if err := c.runMiddlewares(u, next); err != nil {
		t.Fatalf("second: %v", err)
	}
	if err := c.runMiddlewares(u, next); err != nil {
		t.Fatalf("third: %v", err)
	}
	if n != 2 {
		t.Fatalf("next called %d times", n)
	}
}

func TestRateLimitIndependentUsers(t *testing.T) {
	c := newTestClient()
	c.Use(RateLimit(1, time.Second))
	var n int
	next := func() error {
		n++
		return nil
	}
	if err := c.runMiddlewares(messageUpdate(1), next); err != nil {
		t.Fatal(err)
	}
	if err := c.runMiddlewares(messageUpdate(2), next); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("users should be independent, got %d", n)
	}
}

func TestRateLimitCallbackAndPassthrough(t *testing.T) {
	c := newTestClient()
	c.Use(RateLimit(1, time.Second))
	var n int
	next := func() error {
		n++
		return nil
	}
	cb := &UpdateNewCallbackQuery{SenderUserId: 3}
	if err := c.runMiddlewares(cb, next); err != nil {
		t.Fatal(err)
	}
	if err := c.runMiddlewares(cb, next); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("callback limit: %d", n)
	}
	if err := c.runMiddlewares(&UpdateUserStatus{}, next); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("other updates should pass: %d", n)
	}
}

func TestRateLimitNoop(t *testing.T) {
	c := newTestClient()
	c.Use(RateLimit(0, time.Second))
	c.Use(RateLimit(3, 0))
	var n int
	next := func() error {
		n++
		return nil
	}
	u := messageUpdate(1)
	_ = c.runMiddlewares(u, next)
	_ = c.runMiddlewares(u, next)
	if n != 2 {
		t.Fatalf("noop should always pass, got %d", n)
	}
}

func TestRateLimitUserIDZeroPasses(t *testing.T) {
	c := newTestClient()
	c.Use(RateLimit(1, time.Second))
	var n int
	next := func() error {
		n++
		return nil
	}
	_ = c.runMiddlewares(&UpdateNewMessage{Message: &Message{}}, next)
	_ = c.runMiddlewares(&UpdateNewMessage{Message: &Message{}}, next)
	if n != 2 {
		t.Fatalf("zero sender should pass: %d", n)
	}
}
