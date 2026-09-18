package gotdbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (c *Client) OnWebAppData(handler func(client *Client, message *Message) error, filter func(*Message) bool) *MessageHandler {
	return c.OnMessage(handler, func(msg *Message) bool {
		if msg.WebAppData() == nil {
			return false
		}
		if filter != nil && !filter(msg) {
			return false
		}
		return true
	})
}

func (c *Client) SetBotMenuButton(text, url string, userId int64) error {
	return c.SetMenuButton(MenuButton(text, url), userId)
}

func (c *Client) ValidateWebAppInitData(initData string, maxAge time.Duration) (map[string]string, error) {
	token := ""
	if c != nil {
		token = c.botToken
	}
	return ValidateWebAppInitData(token, initData, maxAge)
}

func ValidateWebAppInitData(botToken, initData string, maxAge time.Duration) (map[string]string, error) {
	if botToken == "" || initData == "" {
		return nil, ErrWebAppDataInvalid
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, ErrWebAppDataInvalid
	}

	hash := values.Get("hash")
	if hash == "" {
		return nil, ErrWebAppDataInvalid
	}
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	fields := make(map[string]string, len(keys))
	for _, key := range keys {
		value := values.Get(key)
		fields[key] = value
		pairs = append(pairs, key+"="+value)
	}
	dataCheckString := strings.Join(pairs, "\n")

	secretMac := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = secretMac.Write([]byte(botToken))
	secretKey := secretMac.Sum(nil)

	checkMac := hmac.New(sha256.New, secretKey)
	_, _ = checkMac.Write([]byte(dataCheckString))
	wantHash := hex.EncodeToString(checkMac.Sum(nil))
	if !hmac.Equal([]byte(wantHash), []byte(hash)) {
		return nil, ErrWebAppDataInvalid
	}

	if maxAge > 0 {
		authDateRaw, ok := fields["auth_date"]
		if !ok {
			return nil, ErrWebAppDataInvalid
		}
		authDate, err := strconv.ParseInt(authDateRaw, 10, 64)
		if err != nil {
			return nil, ErrWebAppDataInvalid
		}
		if time.Since(time.Unix(authDate, 0)) > maxAge {
			return nil, ErrWebAppDataExpired
		}
	}

	return fields, nil
}
