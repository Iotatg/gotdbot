package gotdbot

import "strings"

func NormalizeParseMode(parseMode string) string {
	switch strings.ToLower(strings.TrimSpace(parseMode)) {
	case "html":
		return ParseModeHTML
	case "markdown", "md":
		return ParseModeMarkdown
	case "markdownv2", "mdv2":
		return ParseModeMarkdownV2
	default:
		return ParseModeNone
	}
}

func (c *Client) HTML(text string) (*FormattedText, error) {
	return c.ParseText(text, ParseModeHTML)
}

func (c *Client) Markdown(text string) (*FormattedText, error) {
	return c.ParseText(text, ParseModeMarkdown)
}

func (c *Client) MarkdownV2(text string) (*FormattedText, error) {
	return c.ParseText(text, ParseModeMarkdownV2)
}

func (c *Client) Format(text, parseMode string) (*FormattedText, error) {
	return c.ParseText(text, NormalizeParseMode(parseMode))
}
