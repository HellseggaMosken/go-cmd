package main

import (
	"fmt"

	gocmd "github.com/HellseggaMosken/go-cmd"
)

func main() {
	// create a command
	cmd := gocmd.New(
		"example-app",                  // command's name
		"An example command line app.", // command's usage or description
	).Flag( // chaining call
		// add a bool type flag; user should give flag as "-a" or "--aflag"
		gocmd.FlagTypeBool,
		"a",                  // flag's short name; user can give flag as "-a"
		"aflag",              // flag's long name; user can give flag as "--aflag"
		"This is bool flag.", // flag's description
	).Flag(
		// add a value type flag; user should give flag as "-b 'a value'" or "--bflag 'a value'"
		gocmd.FlagTypeValue,
		"b",
		"bflag",
		"This is value flag.",
	).Flag(
		// add a multi-value type flag; user should give flag as "-c value1 value2 ..." or "--cflag value1 value2 ..."
		gocmd.FlagTypeMulti,
		"c",
		"cflag",
		"This is multi-value flag.",
	).FlagService(
		// flag service is similar as flag, but it will run its service when the flag is set
		gocmd.FlagTypeMulti,
		"s",
		"start",
		"Start this service. You can give a value as your start arg."+
			"The usage may be vary long, but the package will wrap lines properly "+
			"when outputing help message.",
		// this is the flag's service that will be called when the flag is set
		func(ctx *gocmd.Context, value any) error {
			// you can safely convert value's type to []string for "Multi" type flag;
			// also, if flag's type is "Value", you can safely convert value to string,
			// and if flag's type is "Bool", you can safely convert value to bool
			arg := value.([]string)
			fmt.Println(arg)

			// you can get a defined flag's value;
			// the following "b" flag is "bflag" flag which is defined above
			if b, isSet := ctx.Short("b"); isSet {
				// just like converting value above, you can safely convert flag here
				v := b.(string)
				fmt.Println("b is set, its value is: %v", v)
			} else {
				fmt.Println("b is not set")
			}

			// you can get app's executable path
			fmt.Println(ctx.Executable())

			// you can get app's working dir path
			fmt.Println(ctx.Working())

			return nil
		},
	).Service(
		// you can regard service as an anonymous flag service.
		// if not provide a service, the package will use a "print help"
		// function as its service.
		func(ctx *gocmd.Context) error {
			fmt.Println("This is a service for the example app.")
			return nil
		},
	).Sub(
		// you can add sub commands. By combining small commands, you
		// can build a complex command app easily.
		subCommand,
	)

	// RunWithArgs will use "os.Args[1:]" as its args
	if err := cmd.RunWithArgs(); err != nil {
		fmt.Println(err)
	}

	// or, you can manually give your args:
	// if err := cmd.Run(os.Args[1:]); err != nil {
	// 	fmt.Println(err)
	// }
}

// a sub command
var subCommand = gocmd.New(
	"sub",
	"A sub command.",
)
