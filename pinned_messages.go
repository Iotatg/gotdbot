package gotdbot

import "encoding/json"

// v0.14.2 - reading pinned messages, which v0.14.1 could not do at all.
//
// getPinnedChatMessages had no request type and no method. UnpinAllChatMessages
// already had both, so it is not re-added here.
//
// The missing read is the more consequential one. "Unpin the most recent pinned
// message" is a thing clients are asked to do - Pyrogram documents message_id 0
// as exactly that - and the only way to honour it is to ask Telegram which
// message is the most recent. Without getPinnedChatMessages the caller has to
// choose between unpinning nothing and unpinAllChatMessages, which is a much
// larger action and not the one that was requested.

// GetPinnedChatMessages returns the list of messages pinned in a chat; requires
// can_pin_messages member right for basic groups and supergroups, or can_edit_messages
// administrator right for channels.
type GetPinnedChatMessages struct {
	Extra string `json:"@extra,omitempty"` // @extra field
	// Identifier of the chat
	ChatId int64 `json:"chat_id"`
	// The maximum number of messages to be returned; 1-100
	Limit int32 `json:"limit"`
}

func (t *GetPinnedChatMessages) setExtra(extra string) { t.Extra = extra }

func (t GetPinnedChatMessages) GetType() string {
	return "getPinnedChatMessages"
}

func (t GetPinnedChatMessages) MarshalJSON() ([]byte, error) {
	type Alias GetPinnedChatMessages
	return json.Marshal(&struct {
		TypeStr string `json:"@type"`
		Extra   string `json:"@extra,omitempty"`
		*Alias
	}{
		TypeStr: "getPinnedChatMessages",
		Extra:   t.Extra,
		Alias:   (*Alias)(&t),
	})
}

// PinnedChatMessages contains a list of pinned messages.
type PinnedChatMessages struct {
	Extra      string `json:"@extra,omitempty"` // @extra field
	TotalCount int32  `json:"total_count"`
	// The list of pinned messages
	Messages []Message `json:"messages"`
}

func (t PinnedChatMessages) GetType() string {
	return "pinnedChatMessages"
}

func (t PinnedChatMessages) pinnedChatMessages() {}

func (t PinnedChatMessages) MarshalJSON() ([]byte, error) {
	type Alias PinnedChatMessages
	return json.Marshal(&struct {
		TypeStr string `json:"@type"`
		Extra   string `json:"@extra,omitempty"`
		*Alias
	}{
		TypeStr: "pinnedChatMessages",
		Extra:   t.Extra,
		Alias:   (*Alias)(&t),
	})
}

// defaultPinnedLookupLimit is what LatestPinnedMessage asks for.
//
// The list order is not something to rely on, so the newest pinned message is
// picked by message id instead of by position - within one chat, message ids
// only ever increase. Asking for enough of them makes that choice independent of
// however TDLib decided to order the result.
const defaultPinnedLookupLimit = 100

// GetPinnedChatMessages returns the pinned messages of a chat, at most limit of
// them. A limit of 0 is passed through unchanged, which TDLib reads as 1.
func (c *Client) GetPinnedChatMessages(chatId int64, limit int32) (*PinnedChatMessages, error) {
	req := &GetPinnedChatMessages{
		ChatId: chatId,
		Limit:  limit,
	}
	resp, err := c.Send(req)
	if err != nil {
		return nil, err
	}
	pinned, ok := resp.(*PinnedChatMessages)
	if !ok {
		// The generated methods assert the type directly. This one cannot,
		// because a response that did not decode into the expected type would
		// otherwise panic, and a panic here loses the caller's error.
		return nil, unexpectedResponse("getPinnedChatMessages", "pinnedChatMessages", resp)
	}
	return pinned, nil
}

// LatestPinnedMessage returns the most recently pinned message in a chat, or nil if
// nothing is pinned.
//
// This is what makes "unpin the latest" expressible without guessing: the caller
// gets a real message id to pass to UnpinChatMessage.
func (c *Client) LatestPinnedMessage(chatId int64) (*Message, error) {
	pinned, err := c.GetPinnedChatMessages(chatId, defaultPinnedLookupLimit)
	if err != nil {
		return nil, err
	}

	id, ok := latestPinnedID(pinned.Messages)
	if !ok {
		return nil, nil
	}
	for i := range pinned.Messages {
		if pinned.Messages[i].Id == id {
			return &pinned.Messages[i], nil
		}
	}
	return nil, nil
}

// latestPinnedID picks the newest message out of a pinned list.
//
// Message ids only ever increase within one chat, so the highest id is the most
// recent. Order is not used, because nothing promises the list arrives oldest
// first: getting this backwards would unpin a months-old announcement while
// reporting the newest as gone.
//
// Returns ok=false for an empty list, which is distinct from a message id of 0 -
// 0 is a value UnpinChatMessage would send, and it resolves to nothing.
func latestPinnedID(messages []Message) (int64, bool) {
	var best int64
	var found bool
	for i := range messages {
		if !found || messages[i].Id > best {
			best = messages[i].Id
			found = true
		}
	}
	return best, found
}
