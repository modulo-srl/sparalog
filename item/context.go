package item

import (
	"fmt"

	"github.com/modulo-srl/sparalog"
)

// Context implements sparalog.Context
type Context struct {
	tags map[string]string
	data map[string]interface{}

	prefixFmt  string
	prefixTags []string
}

func (c *Context) SetPrefix(format string, tags []string) {
	c.prefixFmt = format
	c.prefixTags = tags
}

// AssignContext clone a context to current.
func (c *Context) AssignContext(fromCtx sparalog.Context, clone bool) {
	tags := fromCtx.Tags()
	data := fromCtx.Data()
	prefixFmt, prefixTags := fromCtx.Prefix()

	c.prefixFmt = prefixFmt

	if !clone {
		c.tags = tags
		c.data = data
		c.prefixTags = prefixTags
		return
	}

	if tags != nil {
		c.tags = make(map[string]string)
		for k, v := range tags {
			c.tags[k] = v
		}
	}

	if data != nil {
		c.data = make(map[string]interface{})
		for k, v := range data {
			c.data[k] = v
		}
	}

	if prefixTags != nil {
		c.prefixTags = make([]string, len(prefixTags))
		copy(c.prefixTags, prefixTags)
	}
}

func (c *Context) SetTag(name, value string) {
	if c.tags == nil {
		c.tags = make(map[string]string)
	}

	c.tags[name] = value
}

func (c *Context) Tags() map[string]string {
	return c.tags
}

func (c *Context) SetData(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}

	c.data[key] = value
}

func (c *Context) Data() map[string]interface{} {
	return c.data
}

func (c *Context) Prefix() (string, []string) {
	return c.prefixFmt, c.prefixTags
}

func (c *Context) RenderPrefix() string {
	if c.prefixFmt == "" {
		return ""
	}

	if len(c.prefixTags) == 0 || len(c.tags) == 0 {
		return c.prefixFmt
	}

	args := make([]interface{}, 0, len(c.prefixTags))
	for _, name := range c.prefixTags {
		args = append(args, c.tags[name])
	}

	return fmt.Sprintf(c.prefixFmt, args...)
}
