package gotdbot

import (
	"errors"
	"fmt"
)

var (
	// StopPropagation stops remaining handlers in the current group and later groups.
	StopPropagation = errors.New("group iteration ended")
	// EndGroups stops processing all remaining groups.
	EndGroups = StopPropagation
	// ContinueGroups skips remaining handlers in the current group and continues with the next group.
	ContinueGroups = errors.New("group iteration continued")
	// ContinueHandlers continues processing the remaining handlers in the same group.
	ContinueHandlers = errors.New("handler iteration continued")

	ConversationCancelled = errors.New("conversation cancelled")
	ConversationTimeout   = errors.New("conversation wait timeout")

	SendTimeout         = errors.New("TDLib send timeout")
	WaitPremiumPurchase = errors.New("account requires Telegram Premium")

	ErrNoFormattedText  = errors.New("message does not contain formatted text or caption")
	ErrNotMediaGroup    = errors.New("the message does not belong to a media group")
	ErrDownloadStopped  = errors.New("file download stopped before completion")
	ErrChatNotResolved  = errors.New("TDLib resolved no chat for this id")
	ErrInvalidMessageID = errors.New("message id must be greater than zero")
	ErrInvalidUserID    = errors.New("user id must be greater than zero")

	ErrCallbackTooLong  = errors.New("callback data exceeds 64 bytes")
	ErrCallbackInvalid  = errors.New("callback data is invalid")
	ErrCallbackColon    = errors.New("callback action and args must not contain ':'")
	ErrCallbackBadHMAC  = errors.New("callback hmac is invalid")
	ErrNoCallbackSecret = errors.New("callback secret is empty")

	ErrWebAppDataInvalid = errors.New("web app init data is invalid")
	ErrWebAppDataExpired = errors.New("web app init data has expired")
)

// UnexpectedResponse reports that TDLib answered with something other than the type
// the call promised.
//
// The generated methods assert the response type directly and would panic here.
// Where a method cannot, because it has a fallback path to take, this is the error
// it returns instead - a panic in a library loses the caller's error and, with it,
// the reason the call failed at all.
type UnexpectedResponse struct {
	// Method is the TDLib method that was called.
	Method string
	// Expected is the response type the method documents.
	Expected string
	// GotType is the @type that actually came back, if it had one.
	GotType string
}

func (e *UnexpectedResponse) Error() string {
	if e.GotType == "" {
		return fmt.Sprintf("gotdbot: %s answered with an untyped response, want %s",
			e.Method, e.Expected)
	}
	return fmt.Sprintf("gotdbot: %s answered with %s, want %s",
		e.Method, e.GotType, e.Expected)
}

// unexpectedResponse builds an UnexpectedResponse, reading the @type off whatever
// came back without asserting it.
func unexpectedResponse(method, expected string, got TlObject) error {
	resp := &UnexpectedResponse{Method: method, Expected: expected}
	if got != nil {
		resp.GotType = got.GetType()
	}
	return resp
}
