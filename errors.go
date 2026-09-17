package gotdbot

import (
	"errors"
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
	ErrInvalidMessageID = errors.New("message id must be greater than zero")
)
