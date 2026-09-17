package gotdbot

import "strings"

type PluginHandler struct {
	Handler Handler
	Group   int
}

type Plugin struct {
	Name     string
	Handlers []PluginHandler
}

type PluginLoadOpts struct {
	Include []string
	Exclude []string
}

func NewPlugin(name string) *Plugin {
	return &Plugin{Name: name}
}

func (p *Plugin) Add(handler Handler, group int) *Plugin {
	if p == nil {
		return p
	}
	p.Handlers = append(p.Handlers, PluginHandler{Handler: handler, Group: group})
	return p
}

func (p *Plugin) OnMessage(handler func(client *Client, message *Message) error, filter func(msg *Message) bool, group int) *Plugin {
	return p.Add(NewMessageHandler(handler, filter), group)
}

func (p *Plugin) OnCommand(command string, handler func(client *Client, message *Message) error, group int) *Plugin {
	return p.Add(NewCommandHandler(command, handler), group)
}

func (c *Client) AddPlugin(p *Plugin) {
	if p == nil {
		return
	}
	for _, h := range p.Handlers {
		if h.Handler == nil {
			continue
		}
		c.AddHandlerGroup(h.Handler, h.Group)
	}
	if c.Logger != nil {
		c.Logger.Info("Loaded plugin", "name", p.Name, "handlers", len(p.Handlers))
	}
}

func (c *Client) LoadPlugins(plugins []*Plugin, opts *PluginLoadOpts) {
	for _, p := range plugins {
		if p == nil {
			continue
		}
		if !pluginAllowed(p.Name, opts) {
			if c.Logger != nil {
				c.Logger.Info("Skipped plugin", "name", p.Name)
			}
			continue
		}
		c.AddPlugin(p)
	}
}

func pluginAllowed(name string, opts *PluginLoadOpts) bool {
	if opts == nil {
		return true
	}
	if len(opts.Include) > 0 && !containsFold(opts.Include, name) {
		return false
	}
	if len(opts.Exclude) > 0 && containsFold(opts.Exclude, name) {
		return false
	}
	return true
}

func containsFold(items []string, want string) bool {
	for _, item := range items {
		if strings.EqualFold(item, want) {
			return true
		}
	}
	return false
}
