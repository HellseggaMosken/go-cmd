package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type (
	Service func(ctx *Context) error
	// `value` is the flag's value
	FlagService func(ctx *Context, value any) error
)

type Command struct {
	name    string
	usage   string
	service Service
	subs    []*Command // subsidiary commands
	flags   []struct {
		*Flag
		FlagService
	}
}

// New return a command whose default service is printing help messages to os.Stdout,
// and it will also add a "--help/-h" flag to this command for the same behavior.
func New(name string, usage string) *Command {
	c := &Command{name: name, usage: usage}
	c.service = func(ctx *Context) error {
		fmt.Print(c.Help())
		return nil
	}
	c.FlagService(FlagTypeBool,
		"h", "help",
		"Print help message for command '"+name+"'.",
		func(_ *Context, _ any) error {
			fmt.Print(c.Help())
			return nil
		})
	return c
}

// NewWithoutHelp return a command whose default service is nil
func NewWithoutHelp(name string, usage string) *Command {
	return &Command{name: name, usage: usage}
}

// Run will run the command. It may be helpful ro use `Run(os.Args[1:])`.
// Run will behave as the following order:
//
// 1, Use `cmd.ParseFlags` function to parse all avaiable flags from `args`, and
// return error if any.
//
// 2, If there are any unknown flag args it will return an error.
//
// 3, If there are any parsed flags with a defined service it will (randomly) select
// one to run and return its running error.
//
// 4, If there are any sub commands in the remaining args, it will match a command
// from defined subsidiary commands. If not matched, return an error, or it will
// run the command with remaining args and return its running error.
//
// 5, If the command itself has a defined service, it will run the service and return
// its running error, or it will return "no defined operation..." error.
func (c *Command) Run(args []string) error {
	var flags []*Flag
	for _, f := range c.flags {
		flags = append(flags, f.Flag)
	}

	remaining, unknown, err := ParseFlags(args, flags)
	if err != nil {
		return err
	}
	if len(unknown) > 0 {
		return fmt.Errorf("unknown flag(s): %v", unknown)
	}

	for _, f := range c.flags {
		if f.IsSet() && f.FlagService != nil {
			return f.FlagService(newContex(flags), f.value)
		}
	}

	if len(remaining) > 0 {
		// run the first matched sub command
		for _, sc := range c.subs {
			if sc.name == remaining[0] {
				return sc.Run(remaining[1:])
			}
		}
		return errors.New("unknown command: " + remaining[0])
	}

	if c.service != nil {
		return c.service(newContex(flags))
	}

	return fmt.Errorf("no defined operation for '%v'", c.name)
}

// RunWithArgs is a shortcut for `(cmd.Command).Run(os.Args[1:])`
func (c *Command) RunWithArgs() error {
	// ignore the first os arg (program name)
	return c.Run(os.Args[1:])
}

// Flag will add a flag to the command
func (c *Command) Flag(t FlagType, short string, long string, description string) *Command {
	return c.FlagService(t, short, long, description, nil)
}

// FlagService will add a flag to the command. It ensures that `onSet`
// won't run if the flag's value is not set, which means the `value` param
// of `onSet` func is always valid and can be safely converted to proper
// type:
//
//	FlagTypeBool   ->  v is bool
//	FlagTypeValue  ->  v is string
//	FlagTypeMulti  ->  v is []string
func (c *Command) FlagService(t FlagType, short string, long string, usage string, onSet FlagService) *Command {
	f := &Flag{
		Long:  long,
		Short: short,
		Usage: usage,
		Type:  t,
	}
	c.flags = append(c.flags, struct {
		*Flag
		FlagService
	}{f, onSet})
	return c
}

// Sub will add a subsidiary command to the command
func (c *Command) Sub(command *Command) *Command {
	c.subs = append(c.subs, command)
	return c
}

// Service will set the default service for the command
func (c *Command) Service(s Service) *Command {
	c.service = s
	return c
}

// Help will generate help string for the command, including info messages
// for the command, as well as all defined flags and subsidiary commands.
func (c *Command) Help() string {
	fb := formatBuilder{
		level:   0,
		maxLen:  75,
		builder: &strings.Builder{},
	}
	c.help(fb)
	return fb.String()
}

func (c *Command) help(fb formatBuilder) {
	fb.out(c.name)
	fb = fb.nextLevel()
	fb.out(c.usage)
	fb.out()

	var leftCol []string
	var rightCol []string
	leftLen := 0
	for _, f := range c.flags {
		left := "-" + f.Short + "/--" + f.Long
		if f.Type == FlagTypeValue {
			left = left + " <arg>"
		} else if f.Type == FlagTypeMulti {
			left = left + " <arg ...>"
		}
		leftCol = append(leftCol, left)
		rightCol = append(rightCol, f.Usage)
		if leftLen < len(left) {
			leftLen = len(left)
		}
	}

	leftLen += 2 // this will add two spaces between left and right

	for i := range leftCol {
		fb.outWithLeading(leftCol[i], leftLen, rightCol[i])
	}

	for _, sc := range c.subs {
		fb.out()
		sc.help(fb)
	}
}
