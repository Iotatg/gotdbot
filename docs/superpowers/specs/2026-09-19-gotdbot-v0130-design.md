# gotdbot v0.13.0 design

**Date:** 2026-09-19
**Status:** Approved
**Owner:** Iota coder
**Affects:** `github.com/Iotatg/gotdbot` (this repo)
**Depends on:** published `v0.12.0` (`7b48c3f`), TDLib `v1.8.67`

---

## 1. Problem

`v0.12.0` already ships AskFrom, callback pack/HMAC, RateLimit, and Mini App DX.

A Kurigram v2.2.26 / Pyrogram-latest comparison of dispatcher, keyboards, callback bound methods, client options, filters, and handler aliases still leaves small additive gaps that AishaMusic-style bots hit:

1. **Callback media edits need a Client round-trip.** `UpdateNewCallbackQuery` has `EditMessageText` / `EditMessageCaption` / `EditMessageReplyMarkup`. `Client.EditMessageMedia` and `Message.EditMedia` exist. There is no callback bound helper.
2. **Keyboard builders miss three Telegram button types.** Generated types exist: `InlineKeyboardButtonTypeBuy`, `InlineKeyboardButtonTypeCallbackWithPassword`, `KeyboardButtonTypeRequestManagedBot`. `SwitchInlineButton` already targets `TargetChatCurrent` but has no current-chat alias name.
3. **Handler aliases are incomplete versus Kurigram decorators.** Missing: edited/deleted business messages, purchased paid media, managed bot, message reaction count. There is no catch-all `OnRawUpdate`. There is no typed `OnCallback("play", ...)`.
4. **A few field-only message filters are missing.** `me`, `mentioned`, `live_location`, `web_page`, `quote`, `business`, `ephemeral`.

`bot` and `admin` filters need extra TDLib fetches (`GetUser` / `GetChat`). `Message.Click` is user-client keyboard clicking. Those stay out of this tag.

FSM/scenes, mock client, Stars/payments DX, Mini App hosting, Bot-API rewrite, and pinning AishaMusic are out of scope.

## 2. Goal

Ship an additive minor `v0.13.0` that stays a thin TDLib wrapper:

| Constant | Value |
|---|---|
| `Version` | `v0.13.0` |
| `TDLibVersion` | `v1.8.67` (unchanged) |

Must ship:

1. `UpdateNewCallbackQuery.EditMessageMedia` bound helper.
2. Keyboard builders: `PayButton`, `PasswordCallbackButton`, `RequestManagedBotButton`, `SwitchInlineCurrentButton`.
3. Handler DX: typed `OnCallback` / `OnCallbackGroup`; missing `On*` aliases; `OnRawUpdate` / `OnRawUpdateGroup`.
4. Field-only leftover filters in `filters/message`: `Me`, `Mentioned`, `LiveLocation`, `WebPage`, `Quote`, `Business`, `Ephemeral`. Matching `Message` accessors.

Success: existing `OnCommand` / `OnCallbackQuery` / `CallbackButton` / `Ask` call sites keep compiling. Tests run without `libtdjson`.

## 3. Architecture

```text
Your Go app
    │
    ├─ PayButton / PasswordCallbackButton / RequestManagedBotButton
    ├─ q.EditMessageMedia(c, content, opts)
    ├─ OnCallback("play", handler, filter)
    ├─ OnRawUpdate(handler, filter)
    └─ filters/message.{Me,Mentioned,LiveLocation,WebPage,Quote,Business,Ephemeral}
           │
           ▼
     gotdbot Client (v0.12.0 surface unchanged)
           │
           ▼
     purego → libtdjson v1.8.67 → Telegram
```

Four DX units. No generated-API regeneration. No new subpackages. No Mini App hosting. No Bot-API rewrite.

## 4. Compatibility

Additive only. No breaking changes to `v0.12.0` APIs.

Unchanged:

- `CallbackButton(text, data string)`
- `SwitchInlineButton` (already `TargetChatCurrent`)
- `OnCallbackQuery` / `OnCallbackQueryGroup`
- `PackCallback` / `UnpackCallback` / signed variants
- `Client.Use`, plugin loader, handler groups, `RateLimit`
- generated types/methods (`EditMessageMedia`, `OnUpdateManagedBot`, …)

AishaMusic pin to this tag is a follow-up, not part of this tag.

## 5. Components

### 5.1 Callback bound EditMessageMedia

```go
func (t *UpdateNewCallbackQuery) EditMessageMedia(c *Client, content InputMessageContent, opts *EditMessageMediaOpts) (*Message, error)
```

Wraps `c.EditMessageMedia(t.ChatId, content, t.MessageId, opts)`. Nil `opts` is forwarded as-is (generated method already accepts nil). Does not answer the callback query.

Existing `EditMessageText` / `EditMessageCaption` / `EditMessageReplyMarkup` stay as they are.

### 5.2 Keyboard builders

```go
func PayButton(text string) InlineKeyboardButton
func PasswordCallbackButton(text, data string) InlineKeyboardButton
func RequestManagedBotButton(text string, id int32, suggestedName, suggestedUsername string) KeyboardButton
func SwitchInlineCurrentButton(text, query string) InlineKeyboardButton
```

Behavior:

- `PayButton` sets `Type` to `&InlineKeyboardButtonTypeBuy{}`. Telegram requires this as the first button of an invoice keyboard; the builder does not enforce layout.
- `PasswordCallbackButton` sets `Type` to `&InlineKeyboardButtonTypeCallbackWithPassword{Data: []byte(data)}`. Same data rules as `CallbackButton` (raw string, no pack).
- `RequestManagedBotButton` sets `Type` to `&KeyboardButtonTypeRequestManagedBot{Id, SuggestedName, SuggestedUsername}`. `id == 0` still builds; TDLib validates.
- `SwitchInlineCurrentButton` is an alias of `SwitchInlineButton` (already `TargetChatCurrent`). Do not change `SwitchInlineButton`.
- Empty `text` still builds (same as `CallbackButton`). No panics. No new errors.

### 5.3 Typed OnCallback

```go
func (c *Client) OnCallback(action string, handler func(client *Client, update *UpdateNewCallbackQuery) error, filter func(*UpdateNewCallbackQuery) bool)
func (c *Client) OnCallbackGroup(action string, handler func(client *Client, update *UpdateNewCallbackQuery) error, filter func(*UpdateNewCallbackQuery) bool, group int)
```

Registers via existing `OnCallbackQueryGroup`. Inner filter:

1. `UnpackCallback(u.DataString())`. On error → no match (`ErrCallbackInvalid` / empty data).
2. If `action != ""` and unpacked action != `action` → no match.
3. If `action == ""`, any unpackable payload matches.
4. If caller `filter != nil`, AND it.

Handler signature stays `OnCallbackQuery`'s signature. Callers that need args call `UnpackCallback` (or signed unpack) inside the handler.

Unsigned unpack on HMAC-signed data still splits fields; the last field may be hmac hex. `OnCallback("play")` still matches `play:skip:<hmac>`. Authenticity is `UnpackCallbackSigned`, not this helper.

Empty `action` does not mean "raw unparsed data". Unparseable data never matches.

### 5.4 Missing On* aliases and OnRawUpdate

Thin wrappers over generated methods:

| Alias | Generated |
|---|---|
| `OnEditedBusinessMessage` | `OnUpdateBusinessMessageEdited` |
| `OnDeletedBusinessMessages` | `OnUpdateBusinessMessagesDeleted` |
| `OnPurchasedPaidMedia` | `OnUpdatePaidMediaPurchased` |
| `OnManagedBot` | `OnUpdateManagedBot` |
| `OnMessageReactionCount` | `OnUpdateMessageReactions` |

```go
func (c *Client) OnRawUpdate(handler func(client *Client, update TlObject) error, filter func(TlObject) bool)
func (c *Client) OnRawUpdateGroup(handler func(client *Client, update TlObject) error, filter func(TlObject) bool, group int)
```

New unexported `rawUpdateHandler` implementing `Handler`:

- `CheckUpdate`: `filter == nil` or `filter(update)`.
- `HandleUpdate`: `handler(client, update)`.

`OnRawUpdate` uses group 0. It matches every update unless filtered. First matching handler still stops later groups unless the handler returns `ContinueHandlers` / `ContinueGroups`. Register `OnRawUpdate` in a later group, or return `ContinueHandlers`, if typed handlers in the same or earlier groups must still run. This matches existing group dispatch; it is not a parallel pyrogram-style dual match.

`OnConnect` / `OnDisconnect` / `OnError` stay lifecycle hooks, not raw-update handlers.

### 5.5 Leftover filters and accessors

`Message` accessors (nil-safe, no Client, no network):

```go
func (m *Message) LiveLocation() *MessageLiveLocation
func (m *Message) HasLiveLocation() bool
func (m *Message) LinkPreview() *LinkPreview
func (m *Message) HasLinkPreview() bool
func (m *Message) Quote() *TextQuote
func (m *Message) HasQuote() bool
func (m *Message) IsMentioned() bool
func (m *Message) IsBusiness() bool
func (m *Message) IsEphemeral() bool
```

Definitions:

- `LiveLocation()`: content type-assert `*MessageLiveLocation`.
- `LinkPreview()`: `*MessageText` with `LinkPreview != nil`.
- `Quote()`: `ReplyTo` type-assert `*MessageReplyToMessage` with `Quote != nil`.
- `IsMentioned()`: `ContainsUnreadMention`.
- `IsBusiness()`: `SenderBusinessBotUserId != 0`.
- `IsEphemeral()`: `EphemeralMessageId != 0`.

`filters/message` wrappers (same style as `Photo` / `Forwarded`):

```go
func Me(msg *gotdbot.Message) bool            // msg.IsOutgoingMessage()
func Mentioned(msg *gotdbot.Message) bool     // msg.IsMentioned()
func LiveLocation(msg *gotdbot.Message) bool // msg.HasLiveLocation()
func WebPage(msg *gotdbot.Message) bool      // msg.HasLinkPreview()
func Quote(msg *gotdbot.Message) bool        // msg.HasQuote()
func Business(msg *gotdbot.Message) bool     // msg.IsBusiness()
func Ephemeral(msg *gotdbot.Message) bool    // msg.IsEphemeral()
```

`Me` is outgoing, matching Pyrogram `filters.me` (`sender.is_self or outgoing`). Bots typically have outgoing service/self messages; this does not call `GetMe`.

Not this tag: `bot` (needs `User.Type` via `GetUser`), `admin` (needs chat member rights).

## 6. Data flow

### EditMessageMedia

1. Callback handler receives `*UpdateNewCallbackQuery`.
2. `q.EditMessageMedia(c, content, opts)` uses `q.ChatId` and `q.MessageId`.
3. Generated `Client.EditMessageMedia` talks to TDLib.

### Builders

1. App builds markup with `PayButton` / `PasswordCallbackButton` / `RequestManagedBotButton`.
2. Markup is attached to `Send*` / `Edit*` as today.
3. TDLib validates invoice/password/managed-bot placement.

### OnCallback

1. Update hits processor waiters, then middleware, then handlers.
2. Inner filter unpacks callback data and matches `action`.
3. Handler runs with the original `*UpdateNewCallbackQuery`.

### OnRawUpdate

1. Same dispatch path as any `Handler`.
2. If `CheckUpdate` passes, `HandleUpdate` runs.
3. Default stop-propagation applies.

### Filters

Constructed `*Message` fields only. No `Client`, no `libtdjson`.

## 7. Errors

No new sentinel errors.

Existing pack/unpack errors still apply only when callers use pack helpers. `OnCallback` treats unpack failure as a non-match, not as a returned error.

Empty builder text / `id == 0` for managed-bot still build.

## 8. Files

| File | Change |
|---|---|
| `consts.go` | `Version = "v0.13.0"` |
| `builders.go` | `PayButton`, `PasswordCallbackButton`, `RequestManagedBotButton`, `SwitchInlineCurrentButton` |
| `callback_helpers.go` | `EditMessageMedia` bound helper |
| `handler_aliases.go` | missing `On*` aliases, `OnCallback` / `OnCallbackGroup`, `OnRawUpdate` / `OnRawUpdateGroup`, `rawUpdateHandler` |
| `content_helpers.go` | `LiveLocation` / `LinkPreview` accessors |
| `message_helpers.go` | `Quote` / `IsMentioned` / `IsBusiness` / `IsEphemeral` accessors |
| `filters/message/message.go` | leftover filters |
| `dx_test.go` | accessors + leftover filters |
| `handler_dx_test.go` | `OnCallback`, `OnRawUpdate`, alias registration |
| `builders_test.go` | new: button type assertions |
| `README.md` | Advanced: OnCallback, OnRawUpdate, extra builders/filters |

No edits to `gen_*.go`.

## 9. Testing

All tests must run without `libtdjson`.

Builders:

- `PayButton` type is `*InlineKeyboardButtonTypeBuy`
- `PasswordCallbackButton` type is `*InlineKeyboardButtonTypeCallbackWithPassword` with the given data
- `RequestManagedBotButton` type is `*KeyboardButtonTypeRequestManagedBot` with id/name/username
- `SwitchInlineCurrentButton` equals `SwitchInlineButton` (`TargetChatCurrent`)

Callback bound:

- Method is callable in tests without sending: construct `UpdateNewCallbackQuery{ChatId, MessageId}` and assert those fields are what the wrapper would pass. Do not call TDLib. A compile-time assignment / one-liner wrapper test is enough.

OnCallback (`newTestClient`):

- `OnCallback("play", ...)` matches `DataString()=="play"` and `"play:skip"`
- rejects `"help"` and empty/invalid data
- `OnCallback("", ...)` matches any unpackable payload, not empty data
- extra filter can reject a matching action

OnRawUpdate:

- nil filter matches any `TlObject`
- filter can reject
- registers in the requested group

Aliases:

- each alias registers a handler on `newTestClient` (non-zero handler count in group 0)

Filters / accessors table on constructed `*Message`:

- outgoing → `Me` true
- `ContainsUnreadMention` → `Mentioned` true
- `*MessageLiveLocation` → `LiveLocation` true; `*MessageLocation` false
- `*MessageText` with `LinkPreview` → `WebPage` true; text without preview false
- `ReplyTo=*MessageReplyToMessage{Quote:...}` → `Quote` true
- `SenderBusinessBotUserId != 0` → `Business` true
- `EphemeralMessageId != 0` → `Ephemeral` true

Command to run (Go 1.27.1):

```bash
GOROOT=/usr/local/go-1.27.1 GOTOOLCHAIN=local /usr/local/go-1.27.1/bin/go test ./...
```

## 10. Out of scope (this tag)

- FSM / scenes
- mock / test Client that injects updates
- `Message.Click`
- `bot` / `admin` filters (extra TDLib fetch)
- changing `SwitchInlineButton` semantics
- HMAC-aware `OnCallback` (callers use `UnpackCallbackSigned`)
- RateLimit / AskFrom / Mini App changes
- Stars / payments / invoices convenience beyond `PayButton`
- TDLib bump or regenerating `gen_*.go`
- pinning AishaMusic to `v0.13.0`

## 11. Later tags (roadmap, not implemented now)

**Framework DX**

- FSM / scenes (aiogram-style states, per user+chat)
- mock / test Client that records calls and injects updates
- command prefixes beyond `/`, ignore-mention vs require-mention
- update context bag (`ctx.Value` style data on a handler context)
- router groups with scoped middleware (not only global `Use`)
- HMAC-aware typed callback handlers
- `bot` / `admin` filters backed by cached GetUser/GetChat
- `Message.Click` for user clients

**Mini Apps / Web Apps**

- `OnWebAppMessageSent` alias
- helpers around `GetMainWebApp` / `GetWebAppLinkUrl` / `CloseWebApp`
- example Mini App backend using `ValidateWebAppInitData`

**Product Telegram features**

- Stars / payments / invoices / pre-checkout convenience
- gifts, giveaways, paid media, checklists
- forum topics and Direct Messages topics
- business connection helpers beyond aliases
- stories, boosts as first-class send helpers

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

## 12. Follow-up after this tag

1. Publish `v0.13.0` on `Iotatg/gotdbot`.
2. AishaMusic: replace `replace => ./gotdbot-fork` with `github.com/Iotatg/gotdbot@v0.13.0` when ready.
