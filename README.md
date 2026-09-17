<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=pulse&color=gradient&customColorList=6,11,19,26&height=220&section=header&text=gotdbot&fontSize=80&fontColor=ffffff&animation=scaleIn&fontAlignY=48&stroke=00E5FF&strokeWidth=1" width="100%" alt="gotdbot pulse banner"/>
</p>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=6,11,19,26&height=120&section=header&text=&fontSize=1" width="100%" alt="wave transition"/>
</p>

<p align="center">
  <img src="https://readme-typing-svg.demolab.com?font=Fira+Code&weight=700&size=24&duration=2800&pause=800&color=00E5FF&center=true&vCenter=true&repeat=true&width=780&height=50&lines=Pure+Go+wrapper+for+Telegram+TDLib;No+CGO+%C2%B7+powered+by+purego;Type-safe+generated+API;Conversations%2C+filters+%26+helpers;Go+1.27.1+%C2%B7+maintained+by+Iota+coder" alt="gotdbot typing banner"/>
</p>

<p align="center">
  <a href="https://github.com/Iotatg/gotdbot/stargazers"><img src="https://img.shields.io/github/stars/Iotatg/gotdbot?style=for-the-badge&logo=github&color=FFD700" alt="Stars"/></a>
  <a href="https://github.com/Iotatg/gotdbot/network/members"><img src="https://img.shields.io/github/forks/Iotatg/gotdbot?style=for-the-badge&logo=github&color=00C853" alt="Forks"/></a>
  <a href="https://pkg.go.dev/github.com/Iotatg/gotdbot"><img src="https://img.shields.io/badge/Go-Reference-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Reference"/></a>
  <a href="https://github.com/Iotatg/gotdbot/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-MIT-2EA44F?style=for-the-badge" alt="MIT License"/></a>
  <img src="https://img.shields.io/badge/Go-1.27.1+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.27.1"/>
  <img src="https://img.shields.io/badge/CGO-not%20required-7C4DFF?style=for-the-badge" alt="No CGO"/>
  <img src="https://hitscounter.dev/api/hit?url=https%3A%2F%2Fgithub.com%2FIotatg%2Fgotdbot&label=Visitors&icon=eye&color=%2300E5FF&style=for-the-badge" alt="Visitors"/>
</p>

<p align="center">
  <img src="https://skillicons.dev/icons?i=go,linux,github,bots&theme=dark" alt="stack icons"/>
</p>

<p align="center">
  <b>gotdbot</b> is a pure Go wrapper for <a href="https://github.com/tdlib/td">TDLib</a>.<br/>
  Idiomatic APIs, generated types, and zero CGO — built for Telegram bots and user clients.
</p>

<p align="center">
  <a href="#installation">Install</a> ·
  <a href="#quick-start">Quick Start</a> ·
  <a href="#features">Features</a> ·
  <a href="#architecture">Architecture</a> ·
  <a href="#examples">Examples</a> ·
  <a href="#license">License</a>
</p>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=soft&color=gradient&customColorList=6,11,19&height=70&text=Features&fontSize=28&fontColor=ffffff&animation=twinkling&fontAlignY=52" width="100%" alt="features divider"/>
</p>

## Features

<table>
  <tr>
    <td width="50%" valign="top">
      <b>Pure Go</b><br/>
      Loads <code>libtdjson</code> with <code>purego</code>. No C compiler, no CGO flags.
    </td>
    <td width="50%" valign="top">
      <b>Type-safe API</b><br/>
      100% generated Go structs and methods for every TDLib type and function.
    </td>
  </tr>
  <tr>
    <td width="50%" valign="top">
      <b>High performance</b><br/>
      Direct JSON binding to TDLib with minimal overhead on the hot path.
    </td>
    <td width="50%" valign="top">
      <b>Composable filters</b><br/>
      Match private chats, incoming messages, callbacks, and more with combinators.
    </td>
  </tr>
  <tr>
    <td width="50%" valign="top">
      <b>Conversations</b><br/>
      Built-in <code>Ask</code> / wait helpers for sequential multi-step flows.
    </td>
    <td width="50%" valign="top">
      <b>Bots and users</b><br/>
      Bot tokens, phone login, QR code, and 2FA on the same client.
    </td>
  </tr>
</table>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=soft&color=gradient&customColorList=11,19,26&height=70&text=Architecture&fontSize=28&fontColor=ffffff&animation=blinking&fontAlignY=52" width="100%" alt="architecture divider"/>
</p>

## Architecture

```mermaid
flowchart LR
    A["Your Go app"] --> B["gotdbot Client"]
    B --> C["purego loader"]
    C --> D["libtdjson"]
    D --> E["Telegram"]
```

Updates flow back the same path. Handlers, filters, and conversations sit on the client — you never talk to TDLib JSON by hand.

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=soft&color=gradient&customColorList=6,26,11&height=70&text=Install&fontSize=28&fontColor=ffffff&animation=fadeIn&fontAlignY=52" width="100%" alt="install divider"/>
</p>

## Installation

```bash
go get github.com/Iotatg/gotdbot
```

### Requirements

| Tool | Version |
|------|---------|
| Go | 1.27.1 or newer |
| TDLib | compiled `libtdjson` shared library |

<details>
<summary><b>Quick TDLib setup</b></summary>

Download precompiled binaries from [tdlib-build](https://github.com/FallenProjects/tdlib-build/releases):

```bash
go run github.com/Iotatg/gotdbot/scripts/tools
```

</details>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=soft&color=gradient&customColorList=19,6,11&height=70&text=Quick%20Start&fontSize=28&fontColor=ffffff&animation=scaleIn&fontAlignY=52" width="100%" alt="quick start divider"/>
</p>

## Quick Start

```go
package main

import (
	"log"

	"github.com/Iotatg/gotdbot"
)

func main() {
	bot, err := gotdbot.NewClient(12345, "YOUR_API_HASH", "YOUR_BOT_TOKEN", nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	bot.OnCommand("start", func(c *gotdbot.Client, u *gotdbot.Message) error {
		_, err := u.ReplyText(c, "Hello! I am powered by gotdbot.", nil)
		return err
	})

	bot.OnMessage(func(c *gotdbot.Client, u *gotdbot.Message) error {
		_, err := c.ForwardMessages(u.ChatId, u.ChatId, []int64{u.Id}, &gotdbot.ForwardMessagesOpts{SendCopy: true})
		return err
	}, nil)

	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	bot.Idle()
}
```

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=soft&color=gradient&customColorList=26,6,19&height=70&text=Advanced&fontSize=28&fontColor=ffffff&animation=twinkling&fontAlignY=52" width="100%" alt="advanced divider"/>
</p>

## Advanced Usage

<details>
<summary><b>Composable filters</b></summary>

```go
import (
	"github.com/Iotatg/gotdbot/filters/message"
)

bot.OnMessage(handlePrivate, message.And(message.Private, message.Incoming))
```

Filter packages:

- `github.com/Iotatg/gotdbot/filters/message` — `gotdbot.Message`
- `github.com/Iotatg/gotdbot/filters/callbackquery` — `gotdbot.UpdateNewCallbackQuery`

</details>

<details>
<summary><b>Conversations</b></summary>

```go
bot.OnCommand("rename", func(c *gotdbot.Client, u *gotdbot.Message) error {
	u.ReplyText(c, "What is your new name?", nil)

	res, err := c.Ask(u.ChatId, &gotdbot.WaitMessageOpts{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return err
	}

	u.ReplyText(c, "Nice to meet you, "+res.GetText(), nil)
	return nil
})
```

</details>

<details>
<summary><b>Message helpers</b></summary>

```go
msg.ReplyPhoto(c, gotdbot.InputFileLocal{Path: "image.png"}, nil)
msg.Delete(c, true)
msg.Pin(c, false, false)
msg.React(c, []gotdbot.ReactionType{&gotdbot.ReactionTypeEmoji{Emoji: "👍"}}, nil)
```

</details>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=soft&color=gradient&customColorList=11,6,26&height=70&text=Examples&fontSize=28&fontColor=ffffff&animation=blinking&fontAlignY=52" width="100%" alt="examples divider"/>
</p>

## Examples

| Example | What it shows |
|---------|----------------|
| [Echo Bot](examples/echoBot) | Commands, replies, basic handlers |
| [Conversation Bot](examples/conversation) | `Ask` / wait API |
| [Multi Bot](examples/echoMultiBot) | Multiple clients in one process |

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=soft&color=gradient&customColorList=6,11,19&height=70&text=Contributing&fontSize=28&fontColor=ffffff&animation=fadeIn&fontAlignY=52" width="100%" alt="contributing divider"/>
</p>

## Contributing

Default branch is `main`.

1. Fork the project
2. Create a branch: `git checkout -b feat/amazing-feature`
3. Commit your changes
4. Push: `git push origin feat/amazing-feature`
5. Open a pull request against `main`

<p align="center">
  <img src="https://github-readme-stats.vercel.app/api/pin/?username=Iotatg&repo=gotdbot&theme=tokyonight&hide_border=true" alt="gotdbot stats"/>
  <img src="https://github-readme-stats.vercel.app/api/top-langs/?username=Iotatg&layout=compact&theme=tokyonight&hide_border=true&langs_count=6" alt="top languages"/>
</p>

<p align="center">
  <img src="https://starchart.cc/Iotatg/gotdbot.svg?variant=adaptive" alt="starchart" width="80%"/>
</p>

## License

MIT License. Copyright (c) 2026 [Iota coder](https://github.com/Iotatg). See [LICENSE](LICENSE).

<p align="center">
  <img src="https://readme-typing-svg.demolab.com?font=Fira+Code&weight=600&size=18&duration=2500&pause=1200&color=7C4DFF&center=true&vCenter=true&repeat=true&width=560&height=40&lines=Built+by+Iota+coder;github.com%2FIotatg%2Fgotdbot;Pure+Go+%C2%B7+Zero+CGO+%C2%B7+TDLib" alt="footer typing"/>
</p>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=gradient&customColorList=6,11,19,26&height=160&section=footer&text=Iota%20coder&fontSize=32&fontColor=ffffff&animation=twinkling&fontAlignY=68" width="100%" alt="Iota coder footer"/>
</p>
