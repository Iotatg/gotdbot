# gotdbot v0.12.0 design

**Date:** 2026-09-18
**Status:** Approved
**Owner:** Iota coder
**Affects:** `github.com/Iotatg/gotdbot` (this repo)
**Depends on:** published `v0.11.0` (`2ef1500`), TDLib `v1.8.67`

---

## 1. Problem

`v0.11.0` already ships Pyrogram-style DX: handler aliases, plugins, middleware, chat-scoped `Ask`, `ParseCommand`, keyboard builders, flood-wait helpers, media groups.

AishaMusic and other bots still hit four library gaps:

1. **Group conversations are chat-scoped.** `Ask(chatId)` accepts any sender in the chat. Admin flows and multi-step prompts in groups need user-scoped wait.
2. **Callback data is raw strings.** AishaMusic uses prefixes (`play_skip`, `help_`) and `strings.Contains`. There is no pack/unpack helper and no authenticity check. Telegram limits callback data to 64 bytes.
3. **Inbound spam is unfiltered.** `floodwait.go` covers outbound TDLib retry. There is no inbound per-user rate-limit middleware for `Client.Use`.
4. **Mini App DX is thin.** Generated TDLib already has `WebAppButton`, `AnswerWebAppQuery`, `SetMenuButton`, `GetWebAppUrl`, `OpenWebApp`, and `Message.WebAppData()`. Missing: a handler alias, a menu-button helper, and HMAC validation of Mini App `initData`.

A TDLib bump, FSM/scenes, mock client, worker pool, Mini App hosting, and Bot-API rewrite are out of scope for this tag.

## 2. Goal

Ship an additive minor `v0.12.0` that stays a thin TDLib wrapper:

| Constant | Value |
|---|---|
| `Version` | `v0.12.0` |
| `TDLibVersion` | `v1.8.67` (unchanged) |

Must ship:

1. `AskFrom(chatId, userId, opts)` — user-scoped wait. Existing `Ask` stays chat-scoped.
2. `PackCallback` / `UnpackCallback` plus HMAC-signed variants. `CallbackButton(text, data string)` signature unchanged.
3. `RateLimit(n, window) Middleware` for `Client.Use`. Key = sender user id. Exceed = silent drop.
4. Mini App DX: `OnWebAppData`, `SetBotMenuButton` / `MenuButton`, `ValidateWebAppInitData`.

Success: existing `OnCommand` / prefix callbacks / `Ask` call sites keep compiling. Tests run without `libtdjson`.

## 3. Architecture

```text
Your Go app
    │
    ├─ AskFrom ──► WaitForChat(chatId) + sender filter
    ├─ PackCallback / PackCallbackSigned ──► CallbackButton data (≤64 bytes)
    ├─ Client.Use(RateLimit(n, window)) ──► handler dispatch
    └─ OnWebAppData / SetBotMenuButton / ValidateWebAppInitData
           │
           ▼
     gotdbot Client (v0.11.0 surface unchanged)
           │
           ▼
     purego → libtdjson v1.8.67 → Telegram
```

Four new DX units in the root `gotdbot` package. No generated-API regeneration. No new subpackages. No Mini App HTML/JS hosting.

## 4. Compatibility

Additive only. No breaking changes to `v0.11.0` APIs.

Unchanged:

- `Ask(chatId, opts)`
- `ListenMessage` (still wraps `Ask`)
- `CallbackButton(text, data string)`
- `WebAppButton` / `KeyboardWebAppButton`
- `Client.Use`, plugin loader, handler groups
- `floodwait.go` outbound helpers
- generated types/methods (`AnswerWebAppQuery`, `SetMenuButton`, …)

AishaMusic may adopt the new APIs in a later pin to `v0.12.0`. That pin is not part of this tag.

## 5. Components

### 5.1 AskFrom

```go
func (c *Client) AskFrom(chatId, userId int64, opts *WaitMessageOpts) (*Message, error)
```

Behavior:

- `userId == 0` returns `ErrInvalidUserID` immediately. It does not mean "any sender".
- Reuses `WaitMessageOpts`: `Filter`, `CancellationFilter`, `Timeout`. Nil opts default timeout is 1 minute (same as `Ask`).
- Waiter is registered with `WaitForChat(chatId, ...)`. Incoming `UpdateNewMessage` must:
  1. match `chatId` (already enforced by `WaitForChat`)
  2. `msg.SenderID() == userId`
  3. pass `Filter` if set
  4. not match `CancellationFilter` (if it matches → `ConversationCancelled`)
- Timeout → `ConversationTimeout` (existing).
- Concurrent waiters in the same chat keep today's processor behavior: every matching waiter is offered the update (non-blocking send). First-match-wins is not guaranteed if two waiters share a filter; callers should use distinct `userId` values.
- Implementation extracts an unexported `askFromFilter(chatId, userId int64, opts *WaitMessageOpts) func(*Client, TlObject) bool` so unit tests can assert sender matching without a live waiter loop. `userId == 0` in the filter means "any sender" and is used only by `Ask`.

Thin alias for symmetry with `ListenMessage`:

```go
func (c *Client) ListenMessageFrom(chatId, userId int64, filter func(*Message) bool, timeout time.Duration) (*Message, error)
```

This is `AskFrom` with `WaitMessageOpts{Filter, Timeout}`.

`Ask` and `ListenMessage` are not modified at the API level. `Ask` may share `askFromFilter` internally with `userId == 0`.

### 5.2 Callback pack

Unsigned:

```go
func PackCallback(action string, args ...string) (string, error)
func UnpackCallback(data string) (action string, args []string, err error)
```

Format: `action` or `action:arg1:arg2`. No escaping in this tag. Action and each arg must not contain `:`.

Signed:

```go
func PackCallbackSigned(secret []byte, action string, args ...string) (string, error)
func UnpackCallbackSigned(secret []byte, data string) (action string, args []string, err error)
```

HMAC:

- Algorithm: HMAC-SHA256 over the **unsigned payload** (`action` or `action:arg1:arg2`), then first 12 bytes encoded as 24 lowercase hex characters.
- Wire format: `<unsigned-payload>:<hmac24hex>`.
- Verify with `hmac.Equal` (constant-time).
- Empty `secret` → `ErrNoCallbackSecret`.

`ClientOpts` gains:

```go
CallbackSecret []byte
```

Optional. Client convenience methods:

```go
func (c *Client) PackCallbackSigned(action string, args ...string) (string, error)
func (c *Client) UnpackCallbackSigned(data string) (action string, args []string, err error)
```

These read `c.config.CallbackSecret`. If `config` is nil or secret empty → `ErrNoCallbackSecret`.

Button helper (does not change `CallbackButton`):

```go
func PackedCallbackButton(text, action string, args ...string) (InlineKeyboardButton, error)
```

Packs with `PackCallback`, then `CallbackButton(text, packed)`. On pack error, returns the zero button and the error. No panic.

Telegram limit: packed string must be ≤ 64 bytes (`ErrCallbackTooLong`). Signed payloads spend 25 bytes on `:` + hmac, leaving 39 bytes for `action` + args.

Unsigned `UnpackCallback` on signed data still splits fields; the last field will look like hmac hex. Callers who need authenticity must use `UnpackCallbackSigned`.

### 5.3 RateLimit middleware

```go
func RateLimit(n int, window time.Duration) Middleware
```

- `n <= 0` or `window <= 0` → no-op middleware (`next()` always).
- Key: sender user id.
  - `*UpdateNewMessage` → `Message.SenderID()`
  - `*UpdateNewCallbackQuery` → `SenderUserId`
  - any other update, or extracted id `0` → pass through (`next()`).
- Algorithm: **sliding window**. Per-user slice of timestamps. On each limited update, drop timestamps older than `now - window`. If remaining count `>= n`, return `nil` without calling `next()` (silent drop, no auto-reply, no required log). Otherwise append `now` and call `next()`.
- Storage: mutex-protected `map[int64][]time.Time` captured in the middleware closure. Per `RateLimit()` call, not shared across processes or Client instances unless the same middleware value is reused.
- Typical use: `c.Use(RateLimit(5, time.Second))`.

Does not replace `floodwait.go` (outbound). Does not answer callback queries on drop.

`processor` already offers updates to waiters **before** `runMiddlewares`. RateLimit therefore gates handler dispatch only. `Ask` / `AskFrom` waiters are not rate-limited.

### 5.4 Mini App DX

Already present in `v0.11.0` (do not duplicate):

- `WebAppButton(text, url)` / `KeyboardWebAppButton(text, url)`
- `Message.WebAppData()` / `filters/message.WebAppData`
- generated `AnswerWebAppQuery`, `SetMenuButton`, `GetMenuButton`, `GetWebAppUrl`, `OpenWebApp`

New in `v0.12.0`:

```go
func (c *Client) OnWebAppData(handler func(client *Client, message *Message) error, filter func(*Message) bool) *MessageHandler
```

Registers `OnMessage` with an inner filter that requires `msg.WebAppData() != nil`, then applies the caller filter. Existing `OnMessage` still works.

```go
func MenuButton(text, url string) *BotMenuButton
func (c *Client) SetBotMenuButton(text, url string, userId int64) error
```

`SetBotMenuButton` wraps generated `SetMenuButton`. `userId == 0` means all users (TDLib semantics). Empty `url` still sends `BotMenuButton{Text, Url}` (clears to default commands menu when both are empty).

```go
func ValidateWebAppInitData(botToken, initData string, maxAge time.Duration) (map[string]string, error)
func (c *Client) ValidateWebAppInitData(initData string, maxAge time.Duration) (map[string]string, error)
```

HMAC as Telegram Mini App docs (no TDLib):

1. Parse `initData` as a query string.
2. Take `hash`, drop it from the map.
3. Build `data_check_string`: remaining `key=value` pairs sorted by key, joined with `\n`.
4. `secret_key = HMAC-SHA256(key="WebAppData", data=botToken)`.
5. Compare `hex(HMAC-SHA256(key=secret_key, data=data_check_string))` to `hash` with `hmac.Equal`.
6. `maxAge == 0` skips `auth_date`. Otherwise `auth_date` older than `maxAge` → `ErrWebAppDataExpired`. Missing/unparseable `auth_date` when `maxAge > 0` → `ErrWebAppDataInvalid`.
7. Empty token, empty initData, missing hash, or bad hmac → `ErrWebAppDataInvalid`.

Client method uses `c.botToken`. Empty token → `ErrWebAppDataInvalid`.

Not this tag: hosting Mini App HTML/JS, Telegram.WebApp JS SDK, payment invoices inside Mini Apps.

### 5.5 New errors

| Error | When |
|---|---|
| `ErrInvalidUserID` | `AskFrom` / `ListenMessageFrom` with `userId == 0` |
| `ErrCallbackTooLong` | packed data > 64 bytes |
| `ErrCallbackInvalid` | empty action or empty data |
| `ErrCallbackColon` | `:` in action or an arg |
| `ErrCallbackBadHMAC` | signed unpack mismatch or missing hmac field |
| `ErrNoCallbackSecret` | signed helper with empty secret |
| `ErrWebAppDataInvalid` | Mini App initData missing/tampered/unparseable |
| `ErrWebAppDataExpired` | Mini App `auth_date` older than `maxAge` |

## 6. Data flow

### AskFrom

1. Handler calls `AskFrom(chatId, userId, opts)`.
2. If `userId == 0`, return `ErrInvalidUserID`.
3. Build a filter that type-asserts `*UpdateNewMessage`, checks sender, then existing Filter / CancellationFilter.
4. Block on `WaitForChat`.
5. Return `*Message`, `ConversationCancelled`, `ConversationTimeout`, or context error.

### Callback pack

1. Builder: `PackedCallbackButton("Skip", "play", "skip")` → data `play:skip`.
2. Handler: `action, args, err := UnpackCallback(cb.DataString())`.
3. Signed path uses `CallbackSecret` and rejects tampered data with `ErrCallbackBadHMAC`.

### RateLimit

1. `processor` offers matching waiters the update first (existing; unchanged).
2. Then `runMiddlewares(update, dispatch)` runs handlers.
3. `RateLimit` inspects the update, maybe drops it, otherwise `next()`.
4. Waiters (`Ask` / `AskFrom`) are not subject to RateLimit.

### Mini App

1. User opens a Web App from `WebAppButton` / menu button.
2. Keyboard Web App sends `MessageWebAppDataReceived` → `OnWebAppData`.
3. HTTPS Mini App backend calls `ValidateWebAppInitData(botToken, initData, maxAge)` before trusting `user`.
4. Bot answers queries with existing `AnswerWebAppQuery`.

## 7. Files

| File | Change |
|---|---|
| `consts.go` | `Version = "v0.12.0"` |
| `conversation.go` | `AskFrom`, `askFromFilter` |
| `listeners.go` | `ListenMessageFrom` |
| `callback_pack.go` | new: pack/unpack + signed |
| `builders.go` | `PackedCallbackButton`, `MenuButton` |
| `ratelimit.go` | new: `RateLimit` |
| `webapp.go` | new: `OnWebAppData`, `SetBotMenuButton`, `ValidateWebAppInitData` |
| `client_opts.go` | `CallbackSecret []byte` |
| `errors.go` | new sentinel errors |
| `callback_pack_test.go` | new table tests |
| `ratelimit_test.go` | new tests via `newTestClient` + `runMiddlewares` |
| `conversation_test.go` | `AskFrom` userId 0; `askFromFilter` sender match/reject |
| `webapp_test.go` | initData HMAC + `OnWebAppData` filter |
| `README.md` | Advanced section: AskFrom, callback pack, RateLimit, Mini Apps |

No edits to `gen_*.go`.

## 8. Testing

All tests must run without `libtdjson`.

Pack/unpack table cases:

- `play` → action `play`, no args
- `play:skip` → `play`, `["skip"]`
- colon in action/arg → `ErrCallbackColon`
- empty action → `ErrCallbackInvalid`
- payload of 65 bytes → `ErrCallbackTooLong`
- 64 bytes exactly → ok
- signed round-trip
- wrong secret / truncated hmac → `ErrCallbackBadHMAC`
- empty secret → `ErrNoCallbackSecret`
- `PackedCallbackButton` uses packed data on the callback button type

RateLimit:

- first `n` message updates from the same user call `next`
- `n+1` within window does not
- different user ids are independent
- non-message/callback updates always pass
- `n<=0` or `window<=0` always pass

AskFrom:

- `userId == 0` → `ErrInvalidUserID` without blocking
- constructed filter accepts matching sender and rejects other senders (unit-test the filter function; do not require a live waiter loop)

Mini App:

- valid initData round-trip
- tampered hash → `ErrWebAppDataInvalid`
- stale `auth_date` with `maxAge` → `ErrWebAppDataExpired`
- `OnWebAppData` matches `MessageWebAppDataReceived` and rejects plain text
- `MenuButton` fills `BotMenuButton`

Command to run (Go 1.27.1):

```bash
GOROOT=/usr/local/go-1.27.1 GOTOOLCHAIN=local /usr/local/go-1.27.1/bin/go test ./...
```

## 9. Out of scope (this tag)

- FSM / scenes
- mock / test Client that injects updates
- TDLib version bump or regenerating `gen_*.go`
- changing `Ask` to user-scoped
- changing `CallbackButton` signature
- percent-encoding or length-prefix for `:` in args
- RateLimit auto-reply or callback `Answer`
- shared/redis rate limiter
- outbound flood (already `floodwait.go`)
- hosting Mini App HTML/JS or the Telegram.WebApp JS SDK
- pinning AishaMusic to `v0.12.0` (follow-up)

## 10. Later tags (roadmap, not implemented now)

These are real gaps or likely future Telegram surface. Documented so v0.12.0 stays a slice, not a dump.

**Framework DX**

- FSM / scenes (aiogram-style states, per user+chat)
- mock / test Client that records calls and injects updates
- command prefixes beyond `/`, ignore-mention vs require-mention
- update context bag (`ctx.Value` style data on a handler context)
- router groups with scoped middleware (not only global `Use`)
- typed callback handlers (`OnCallback("play", ...)`)

**Mini Apps / Web Apps (beyond this tag)**

- `OnWebAppMessageSent` alias for `UpdateWebAppMessageSent`
- helpers around `GetMainWebApp` / `GetWebAppLinkUrl` / `CloseWebApp`
- Mini App file-download check (`CheckWebAppFileDownload`)
- example Mini App backend using `ValidateWebAppInitData`

**Product Telegram features (generated API exists; DX wrappers later)**

- Stars / payments / invoices / pre-checkout convenience
- gifts, giveaways, paid media, checklists
- forum topics and Direct Messages topics
- business connection helpers beyond `OnBusinessConnection`
- stories, boosts, checklists as first-class send helpers (some `Reply*` already exist)

**Ops**

- i18n catalog
- metrics / tracing middleware
- persistent session store
- worker pool / job scheduler

**TDLib**

- bump `libtdjson` and regenerate `gen_*.go` only if a real public binary newer than `v1.8.67` exists

**Never fake in gotdbot**

- hosting a Mini App frontend
- Bot-API HTTP rewrite
- paid third-party APIs
- a Telegram client UI

## 11. Follow-up after this tag

1. Publish `v0.12.0` on `Iotatg/gotdbot`.
2. AishaMusic: replace `replace => ./gotdbot-fork` with `github.com/Iotatg/gotdbot@v0.12.0` when ready.
3. Pick the next tag from section 10 (likely FSM or mock client).
