package gotdbot

import (
	"errors"
	"strings"
	"testing"
)

func TestPackUnpackCallback(t *testing.T) {
	packed, err := PackCallback("play")
	if err != nil || packed != "play" {
		t.Fatalf("pack play: %q %v", packed, err)
	}
	action, args, err := UnpackCallback(packed)
	if err != nil || action != "play" || args != nil {
		t.Fatalf("unpack play: %q %#v %v", action, args, err)
	}

	packed, err = PackCallback("play", "skip")
	if err != nil || packed != "play:skip" {
		t.Fatalf("pack skip: %q %v", packed, err)
	}
	action, args, err = UnpackCallback(packed)
	if err != nil || action != "play" || len(args) != 1 || args[0] != "skip" {
		t.Fatalf("unpack skip: %q %#v %v", action, args, err)
	}
}

func TestPackCallbackErrors(t *testing.T) {
	if _, err := PackCallback(""); !errors.Is(err, ErrCallbackInvalid) {
		t.Fatalf("empty action: %v", err)
	}
	if _, err := PackCallback("play:x"); !errors.Is(err, ErrCallbackColon) {
		t.Fatalf("colon action: %v", err)
	}
	if _, err := PackCallback("play", "a:b"); !errors.Is(err, ErrCallbackColon) {
		t.Fatalf("colon arg: %v", err)
	}
	if _, _, err := UnpackCallback(""); !errors.Is(err, ErrCallbackInvalid) {
		t.Fatalf("empty unpack: %v", err)
	}
}

func TestPackCallbackLength(t *testing.T) {
	ok := strings.Repeat("a", 64)
	packed, err := PackCallback(ok)
	if err != nil || packed != ok {
		t.Fatalf("64 bytes should pass: %v", err)
	}
	if _, err := PackCallback(strings.Repeat("a", 65)); !errors.Is(err, ErrCallbackTooLong) {
		t.Fatalf("65 bytes: %v", err)
	}
}

func TestPackCallbackSignedRoundTrip(t *testing.T) {
	secret := []byte("super-secret")
	packed, err := PackCallbackSigned(secret, "play", "skip")
	if err != nil {
		t.Fatalf("signed pack: %v", err)
	}
	action, args, err := UnpackCallbackSigned(secret, packed)
	if err != nil || action != "play" || len(args) != 1 || args[0] != "skip" {
		t.Fatalf("signed unpack: %q %#v %v", action, args, err)
	}
	if _, _, err := UnpackCallbackSigned([]byte("other"), packed); !errors.Is(err, ErrCallbackBadHMAC) {
		t.Fatalf("wrong secret: %v", err)
	}
	if _, _, err := UnpackCallbackSigned(secret, "play:skip"); !errors.Is(err, ErrCallbackBadHMAC) {
		t.Fatalf("missing hmac: %v", err)
	}
	if _, _, err := UnpackCallbackSigned(secret, packed[:len(packed)-2]); !errors.Is(err, ErrCallbackBadHMAC) {
		t.Fatalf("truncated hmac: %v", err)
	}
	action, args, err = UnpackCallback(packed)
	if err != nil || action != "play" || len(args) != 2 {
		t.Fatalf("unsigned unpack of signed still splits: %q %#v %v", action, args, err)
	}
	if _, err := PackCallbackSigned(nil, "play"); !errors.Is(err, ErrNoCallbackSecret) {
		t.Fatalf("empty secret pack: %v", err)
	}
	if _, _, err := UnpackCallbackSigned(nil, packed); !errors.Is(err, ErrNoCallbackSecret) {
		t.Fatalf("empty secret unpack: %v", err)
	}
}

func TestPackCallbackSignedTooLong(t *testing.T) {
	secret := []byte("k")
	if _, err := PackCallbackSigned(secret, strings.Repeat("a", 40)); !errors.Is(err, ErrCallbackTooLong) {
		t.Fatalf("signed too long: %v", err)
	}
}

func TestPackedCallbackButton(t *testing.T) {
	btn, err := PackedCallbackButton("Skip", "play", "skip")
	if err != nil {
		t.Fatalf("button: %v", err)
	}
	cb, ok := btn.Type.(*InlineKeyboardButtonTypeCallback)
	if !ok || string(cb.Data) != "play:skip" || btn.Text != "Skip" {
		t.Fatalf("button payload: %+v", btn)
	}
	if _, err := PackedCallbackButton("X", "a:b"); !errors.Is(err, ErrCallbackColon) {
		t.Fatalf("button error: %v", err)
	}
}

func TestClientCallbackSecretHelpers(t *testing.T) {
	c := newTestClient()
	if _, err := c.PackCallbackSigned("play"); !errors.Is(err, ErrNoCallbackSecret) {
		t.Fatalf("nil config secret: %v", err)
	}
	c.config = &ClientOpts{CallbackSecret: []byte("s")}
	packed, err := c.PackCallbackSigned("play", "1")
	if err != nil {
		t.Fatalf("client pack: %v", err)
	}
	action, args, err := c.UnpackCallbackSigned(packed)
	if err != nil || action != "play" || len(args) != 1 || args[0] != "1" {
		t.Fatalf("client unpack: %q %#v %v", action, args, err)
	}
}
