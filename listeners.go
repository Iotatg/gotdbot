package gotdbot

import "time"

func (c *Client) Listen(filter func(client *Client, update TlObject) bool, timeout time.Duration) (TlObject, error) {
	return c.WaitFor(filter, timeout)
}

func (c *Client) ListenChat(chatId int64, filter func(client *Client, update TlObject) bool, timeout time.Duration) (TlObject, error) {
	return c.WaitForChat(chatId, filter, timeout)
}

func (c *Client) ListenMessage(chatId int64, filter func(*Message) bool, timeout time.Duration) (*Message, error) {
	return c.Ask(chatId, &WaitMessageOpts{Filter: filter, Timeout: timeout})
}

func (c *Client) ListenMessageFrom(chatId, userId int64, filter func(*Message) bool, timeout time.Duration) (*Message, error) {
	return c.AskFrom(chatId, userId, &WaitMessageOpts{Filter: filter, Timeout: timeout})
}

func (c *Client) ListenCallback(filter func(*UpdateNewCallbackQuery) bool, timeout time.Duration) (*UpdateNewCallbackQuery, error) {
	u, err := c.WaitFor(func(client *Client, update TlObject) bool {
		cb, ok := update.(*UpdateNewCallbackQuery)
		if !ok {
			return false
		}
		if filter != nil && !filter(cb) {
			return false
		}
		return true
	}, timeout)
	if err != nil {
		return nil, err
	}
	return u.(*UpdateNewCallbackQuery), nil
}

func (c *Client) ListenInlineQuery(filter func(*UpdateNewInlineQuery) bool, timeout time.Duration) (*UpdateNewInlineQuery, error) {
	u, err := c.WaitFor(func(client *Client, update TlObject) bool {
		q, ok := update.(*UpdateNewInlineQuery)
		if !ok {
			return false
		}
		if filter != nil && !filter(q) {
			return false
		}
		return true
	}, timeout)
	if err != nil {
		return nil, err
	}
	return u.(*UpdateNewInlineQuery), nil
}
