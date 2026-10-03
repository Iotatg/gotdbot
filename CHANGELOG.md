# Changelog

All notable changes to this project are documented here.

## v0.14.0

Additive release. No existing API changed, so upgrading from v0.13.0 is a
version bump with no code changes.

### Added

- `Client.Invoke` / `Client.InvokeWithContext` — generic raw TDLib request
  escape hatch. Mirrors Pyrogram's `client.invoke`, for TDLib methods that have
  no generated wrapper yet. Returns the decoded concrete response type.
- `Client.InvokeJSON` / `Client.InvokeJSONWithContext` — raw escape hatch that
  accepts hand-written JSON, so TDLib methods newer than these bindings are
  usable immediately. `@extra` is injected automatically to correlate the
  response; TDLib errors are returned as `*Error`.
- `Client.UploadBytes(data []byte, fileName string) (*File, error)` — uploads
  in-memory bytes to Telegram and returns the resulting `File`.
- `Client.InputFileBytes(data []byte, fileName string) (InputFile, error)` —
  same, returning an `*InputFileId` ready for any `inputFile` field
  (`inputPhoto.photo`, `inputStickerSetItem.document`, and so on).
- `ClientOpts.UploadDir` — directory used to stage in-memory uploads. Ignored
  when `FilesDirectory` is set, since that is the directory TDLib was
  initialised with and the only one guaranteed to resolve relative file
  parameters correctly. Defaults to `./file`.

### Notes

TDLib has no `inputFileBytes` equivalent: file parameters are always paths on
disk that TDLib reads itself. `UploadBytes` therefore stages the bytes in a
temporary file, calls TDLib's `uploadFile`, and removes the staged file on every
exit path including errors. The returned value is an `*InputFileId` because
TDLib always addresses uploads by id.

Upload file names are sanitised at the write boundary, so a traversing name
cannot escape the upload directory.

The two pre-existing `unsafe.Pointer` vet notices in `internal/tdjson` come
from the CGO string bridging and are unchanged by this release.

## v0.13.0

- Callback `EditMessageMedia`, remaining keyboard builders, typed `OnCallback`
  and `OnRawUpdate`, missing `On*` handler aliases, field-only leftover message
  filters.

## v0.12.0

- User-scoped `AskFrom`, packed/signed callback data, inbound `RateLimit`
  middleware, and Mini App `initData` validation. Existing `Ask` and
  `CallbackButton` APIs unchanged.

## v0.11.0

- Member, join/invite, dialog, album, vote, keyboard-request and
  connect/disconnect helpers. Fixed invite link false positives,
  `PromoteUser` zero-value demotion, `CopyMediaGroup` opts mutation, history
  pagination, and `OnDisconnect` firing before the first ready state.