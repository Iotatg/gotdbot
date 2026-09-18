package gotdbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const (
	callbackDataMaxBytes = 64
	callbackHMACBytes    = 12
)

func PackCallback(action string, args ...string) (string, error) {
	payload, err := packCallbackPayload(action, args)
	if err != nil {
		return "", err
	}
	if len(payload) > callbackDataMaxBytes {
		return "", ErrCallbackTooLong
	}
	return payload, nil
}

func UnpackCallback(data string) (string, []string, error) {
	if data == "" {
		return "", nil, ErrCallbackInvalid
	}
	parts := strings.Split(data, ":")
	if parts[0] == "" {
		return "", nil, ErrCallbackInvalid
	}
	if len(parts) == 1 {
		return parts[0], nil, nil
	}
	return parts[0], parts[1:], nil
}

func PackCallbackSigned(secret []byte, action string, args ...string) (string, error) {
	if len(secret) == 0 {
		return "", ErrNoCallbackSecret
	}
	payload, err := packCallbackPayload(action, args)
	if err != nil {
		return "", err
	}
	mac := callbackHMAC(secret, payload)
	packed := payload + ":" + mac
	if len(packed) > callbackDataMaxBytes {
		return "", ErrCallbackTooLong
	}
	return packed, nil
}

func UnpackCallbackSigned(secret []byte, data string) (string, []string, error) {
	if len(secret) == 0 {
		return "", nil, ErrNoCallbackSecret
	}
	if data == "" {
		return "", nil, ErrCallbackInvalid
	}
	idx := strings.LastIndex(data, ":")
	if idx <= 0 || idx == len(data)-1 {
		return "", nil, ErrCallbackBadHMAC
	}
	payload := data[:idx]
	gotMAC := data[idx+1:]
	wantMAC := callbackHMAC(secret, payload)
	if !hmac.Equal([]byte(gotMAC), []byte(wantMAC)) {
		return "", nil, ErrCallbackBadHMAC
	}
	return UnpackCallback(payload)
}

func (c *Client) PackCallbackSigned(action string, args ...string) (string, error) {
	return PackCallbackSigned(c.callbackSecret(), action, args...)
}

func (c *Client) UnpackCallbackSigned(data string) (string, []string, error) {
	return UnpackCallbackSigned(c.callbackSecret(), data)
}

func (c *Client) callbackSecret() []byte {
	if c == nil || c.config == nil {
		return nil
	}
	return c.config.CallbackSecret
}

func packCallbackPayload(action string, args []string) (string, error) {
	if action == "" {
		return "", ErrCallbackInvalid
	}
	if strings.Contains(action, ":") {
		return "", ErrCallbackColon
	}
	parts := make([]string, 0, 1+len(args))
	parts = append(parts, action)
	for _, arg := range args {
		if strings.Contains(arg, ":") {
			return "", ErrCallbackColon
		}
		parts = append(parts, arg)
	}
	return strings.Join(parts, ":"), nil
}

func callbackHMAC(secret []byte, payload string) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payload))
	sum := mac.Sum(nil)
	if len(sum) > callbackHMACBytes {
		sum = sum[:callbackHMACBytes]
	}
	return hex.EncodeToString(sum)
}
