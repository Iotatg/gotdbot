package gotdbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func signWebAppInitData(token string, fields url.Values) string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		if key == "hash" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, key+"="+fields.Get(key))
	}
	secretMac := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = secretMac.Write([]byte(token))
	checkMac := hmac.New(sha256.New, secretMac.Sum(nil))
	_, _ = checkMac.Write([]byte(strings.Join(pairs, "\n")))
	fields.Set("hash", hex.EncodeToString(checkMac.Sum(nil)))
	return fields.Encode()
}

func TestValidateWebAppInitData(t *testing.T) {
	token := "123:ABC"
	fields := url.Values{}
	fields.Set("query_id", "AAE")
	fields.Set("user", `{"id":1}`)
	fields.Set("auth_date", strconv.FormatInt(time.Now().Unix(), 10))
	initData := signWebAppInitData(token, fields)

	got, err := ValidateWebAppInitData(token, initData, time.Hour)
	if err != nil {
		t.Fatalf("valid init data: %v", err)
	}
	if got["query_id"] != "AAE" || got["user"] != `{"id":1}` {
		t.Fatalf("fields: %#v", got)
	}

	if _, err := ValidateWebAppInitData(token, initData+"x", time.Hour); !errors.Is(err, ErrWebAppDataInvalid) {
		t.Fatalf("tampered: %v", err)
	}
	if _, err := ValidateWebAppInitData("", initData, 0); !errors.Is(err, ErrWebAppDataInvalid) {
		t.Fatalf("empty token: %v", err)
	}
	if _, err := ValidateWebAppInitData(token, "", 0); !errors.Is(err, ErrWebAppDataInvalid) {
		t.Fatalf("empty data: %v", err)
	}
}

func TestValidateWebAppInitDataExpired(t *testing.T) {
	token := "123:ABC"
	fields := url.Values{}
	fields.Set("auth_date", strconv.FormatInt(time.Now().Add(-2*time.Hour).Unix(), 10))
	initData := signWebAppInitData(token, fields)
	if _, err := ValidateWebAppInitData(token, initData, time.Minute); !errors.Is(err, ErrWebAppDataExpired) {
		t.Fatalf("expired: %v", err)
	}
	if _, err := ValidateWebAppInitData(token, initData, 0); err != nil {
		t.Fatalf("maxAge 0 should skip: %v", err)
	}
}

func TestValidateWebAppInitDataMissingAuthDate(t *testing.T) {
	token := "123:ABC"
	fields := url.Values{}
	fields.Set("query_id", "AAE")
	initData := signWebAppInitData(token, fields)
	if _, err := ValidateWebAppInitData(token, initData, time.Hour); !errors.Is(err, ErrWebAppDataInvalid) {
		t.Fatalf("missing auth_date: %v", err)
	}
}

func TestClientValidateWebAppInitDataEmptyToken(t *testing.T) {
	c := newTestClient()
	if _, err := c.ValidateWebAppInitData("a=b", 0); !errors.Is(err, ErrWebAppDataInvalid) {
		t.Fatalf("empty bot token: %v", err)
	}
}

func TestOnWebAppDataFilter(t *testing.T) {
	c := newTestClient()
	h := c.OnWebAppData(func(client *Client, message *Message) error { return nil }, nil)
	web := &UpdateNewMessage{Message: &Message{
		Content: &MessageWebAppDataReceived{ButtonText: "Open", Data: "ok"},
	}}
	text := &UpdateNewMessage{Message: &Message{
		Content: &MessageText{Text: &FormattedText{Text: "hi"}},
	}}
	if !h.CheckUpdate(c, web) {
		t.Fatal("web app data should match")
	}
	if h.CheckUpdate(c, text) {
		t.Fatal("plain text should not match")
	}
}

func TestOnWebAppDataCallerFilter(t *testing.T) {
	c := newTestClient()
	h := c.OnWebAppData(func(client *Client, message *Message) error { return nil }, func(m *Message) bool {
		return m.WebAppData().Data == "keep"
	})
	keep := &UpdateNewMessage{Message: &Message{
		Content: &MessageWebAppDataReceived{Data: "keep"},
	}}
	drop := &UpdateNewMessage{Message: &Message{
		Content: &MessageWebAppDataReceived{Data: "drop"},
	}}
	if !h.CheckUpdate(c, keep) {
		t.Fatal("caller filter should keep")
	}
	if h.CheckUpdate(c, drop) {
		t.Fatal("caller filter should drop")
	}
}

func TestMenuButton(t *testing.T) {
	btn := MenuButton("Open", "https://example.com/app")
	if btn == nil || btn.Text != "Open" || btn.Url != "https://example.com/app" {
		t.Fatalf("menu button: %+v", btn)
	}
}
