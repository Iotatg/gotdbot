package gotdbot

import (
	"sync"
	"time"
)

func RateLimit(n int, window time.Duration) Middleware {
	if n <= 0 || window <= 0 {
		return func(client *Client, update TlObject, next func() error) error {
			return next()
		}
	}

	var mu sync.Mutex
	hits := make(map[int64][]time.Time)

	return func(client *Client, update TlObject, next func() error) error {
		userID := rateLimitUserID(update)
		if userID == 0 {
			return next()
		}

		now := time.Now()
		cutoff := now.Add(-window)

		mu.Lock()
		stamps := hits[userID]
		kept := stamps[:0]
		for _, ts := range stamps {
			if ts.After(cutoff) {
				kept = append(kept, ts)
			}
		}
		if len(kept) >= n {
			hits[userID] = kept
			mu.Unlock()
			return nil
		}
		hits[userID] = append(kept, now)
		mu.Unlock()
		return next()
	}
}

func rateLimitUserID(update TlObject) int64 {
	switch u := update.(type) {
	case *UpdateNewMessage:
		if u.Message == nil {
			return 0
		}
		return u.Message.SenderID()
	case *UpdateNewCallbackQuery:
		return u.SenderUserId
	default:
		return 0
	}
}
