package cmd

import "os"

type Context struct {
	flags []*Flag
}

func newContex(flags []*Flag) *Context {
	return &Context{flags: flags}
}

// Short will find a flag with short name `name`.
// If found, return `flag.Value()`, or return (nil, false)
func (c *Context) Short(name string) (v any, ok bool) {
	for i := range c.flags {
		if c.flags[i].Short == name {
			return c.flags[i].Value()
		}
	}
	return nil, false
}

// Long will find a flag with long name `name`.
// If found, return `flag.Value()`, or return (nil, false)
func (c *Context) Long(name string) (v any, ok bool) {
	for i := range c.flags {
		if c.flags[i].Short == name {
			return c.flags[i].Value()
		}
	}
	return nil, false
}

// Executable return executable program path.
// It will panic if an error occur.
func (c *Context) Executable() string {
	p, e := os.Executable()
	if e != nil {
		panic(e)
	}
	return p
}

// Working return current working directory path.
// It will panic if an error occur.
func (c *Context) Working() string {
	p, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	return p
}
