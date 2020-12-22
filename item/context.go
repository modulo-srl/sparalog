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

// SetPrefix change the prefix for the context.
// Tags contains the list of tags names used as params for format.
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

// SetTag changes a specific tag.
func (c *Context) SetTag(name, value string) {
	if c.tags == nil {
		c.tags = make(map[string]string)
	}

	c.tags[name] = value
}

// Tags returns all the tags.
func (c *Context) Tags() map[string]string {
	return c.tags
}

// SetData changes a specific data field.
func (c *Context) SetData(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}

	c.data[key] = value
}

// Data returns all the data fields.
func (c *Context) Data() map[string]interface{} {
	return c.data
}

// Prefix returns the prefix format and tags names.
func (c *Context) Prefix() (format string, tags []string) {
	return c.prefixFmt, c.prefixTags
}

// RenderPrefix returns a prefix as rendered string.
func (c *Context) RenderPrefix() string {
	if c.prefixFmt == "" {
		return ""
	}

	if len(c.prefixTags) == 0 {
		return c.prefixFmt
	}

	args := make([]interface{}, 0, len(c.prefixTags))
	for _, name := range c.prefixTags {
		args = append(args, c.tags[name])
	}

	return fmt.Sprintf(c.prefixFmt, args...)
}
