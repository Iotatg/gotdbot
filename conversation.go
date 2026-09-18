package gotdbot

import (
	"time"
)

// WaitMessageOpts holds optional parameters for Ask.
type WaitMessageOpts struct {
	Filter             func(*Message) bool
	CancellationFilter func(*Message) bool
	Timeout            time.Duration
}

// Ask waits for a new message in the specified chat.
func (c *Client) Ask(chatId int64, opts *WaitMessageOpts) (*Message, error) {
	return c.ask(chatId, 0, opts)
}

// AskFrom waits for a new message in the specified chat from a specific user.
func (c *Client) AskFrom(chatId, userId int64, opts *WaitMessageOpts) (*Message, error) {
	if userId == 0 {
		return nil, ErrInvalidUserID
	}
	return c.ask(chatId, userId, opts)
}

func (c *Client) ask(chatId, userId int64, opts *WaitMessageOpts) (*Message, error) {
	if opts == nil {
		opts = &WaitMessageOpts{Timeout: 1 * time.Minute}
	}

	u, err := c.WaitForChat(chatId, askFromFilter(chatId, userId, opts), opts.Timeout)
	if err != nil {
		return nil, err
	}

	msg := u.(*UpdateNewMessage).Message
	if opts.CancellationFilter != nil && opts.CancellationFilter(msg) {
		return nil, ConversationCancelled
	}

	return msg, nil
}

func askFromFilter(chatId, userId int64, opts *WaitMessageOpts) func(*Client, TlObject) bool {
	var filter func(*Message) bool
	var cancellationFilter func(*Message) bool
	if opts != nil {
		filter = opts.Filter
		cancellationFilter = opts.CancellationFilter
	}
	return func(client *Client, update TlObject) bool {
		u, ok := update.(*UpdateNewMessage)
		if !ok || u.Message == nil {
			return false
		}

		msg := u.Message
		if msg.ChatId != chatId {
			return false
		}
		if userId != 0 && msg.SenderID() != userId {
			return false
		}

		if cancellationFilter != nil && cancellationFilter(msg) {
			return true
		}

		if filter != nil && !filter(msg) {
			return false
		}
		return true
	}
}
