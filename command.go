package gotdbot

import (
	"strings"
	"unicode"
)

type ParsedCommand struct {
	Prefix  string
	Name    string
	Mention string
	Args    []string
}

func (m *Message) Command() string {
	parsed := ParseCommand(m.GetText(), "/")
	if parsed == nil {
		return ""
	}
	return parsed.Name
}

func (m *Message) CommandMention() string {
	parsed := ParseCommand(m.GetText(), "/")
	if parsed == nil {
		return ""
	}
	return parsed.Mention
}

func (m *Message) CommandArgsQuoted() []string {
	parsed := ParseCommand(m.GetText(), "/")
	if parsed == nil {
		return nil
	}
	return parsed.Args
}

func ParseCommand(text string, prefixes ...string) *ParsedCommand {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if len(prefixes) == 0 {
		prefixes = []string{"/"}
	}

	matchedPrefix := ""
	found := false
	hasEmpty := false
	for _, p := range prefixes {
		if p == "" {
			hasEmpty = true
			continue
		}
		if strings.HasPrefix(text, p) {
			matchedPrefix = p
			found = true
			break
		}
	}
	if !found && hasEmpty {
		found = true
		matchedPrefix = ""
	}
	if !found {
		return nil
	}

	rest := strings.TrimSpace(text[len(matchedPrefix):])
	if rest == "" {
		return nil
	}

	cmdPart, argText, _ := strings.Cut(rest, " ")
	cmdPart = strings.TrimSpace(cmdPart)
	if cmdPart == "" {
		return nil
	}

	name := cmdPart
	mention := ""
	if i := strings.Index(cmdPart, "@"); i != -1 {
		name = cmdPart[:i]
		mention = cmdPart[i+1:]
	}
	if name == "" {
		return nil
	}

	return &ParsedCommand{
		Prefix:  matchedPrefix,
		Name:    name,
		Mention: mention,
		Args:    splitCommandArgs(strings.TrimSpace(argText)),
	}
}

func splitCommandArgs(s string) []string {
	if s == "" {
		return nil
	}
	var args []string
	var buf strings.Builder
	var quote rune
	escaped := false
	for _, r := range s {
		switch {
		case escaped:
			buf.WriteRune(r)
			escaped = false
		case r == '\\' && quote != 0:
			escaped = true
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				buf.WriteRune(r)
			}
		case r == '"' || r == '\'':
			quote = r
		case unicode.IsSpace(r):
			if buf.Len() > 0 {
				args = append(args, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		args = append(args, buf.String())
	}
	return args
}

func IsCommandText(text string, prefixes ...string) bool {
	return ParseCommand(text, prefixes...) != nil
}
