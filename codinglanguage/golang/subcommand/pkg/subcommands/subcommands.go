package subcommands

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	commands map[string]SubCommand
	mu       sync.Mutex
)

type SubCommand struct {
	flagSet *flag.FlagSet
	run     func()
}

func RegisterSubCommand(name string, flagSet *flag.FlagSet, run func()) {
	mu.Lock()
	defer mu.Unlock()
	if commands == nil {
		commands = map[string]SubCommand{
			name: {flagSet: flagSet, run: run},
		}
		return
	}

	if _, ok := commands[name]; ok {
		fmt.Println("name already registered.")
		os.Exit(1)
	}

	commands[name] = SubCommand{flagSet: flagSet, run: run}
}

func Usage() {
	subcommands := []string{}
	for cmd := range commands {
		subcommands = append(subcommands, cmd)
	}
	fmt.Println("expected subcommands: ", strings.Join(subcommands, ", "))
	os.Exit(1)
}

func Run(cmd string) {
	subCmd, ok := commands[cmd]
	if !ok {
		Usage()
		return
	}

	subCmd.flagSet.Parse(os.Args[2:])
	fmt.Println("tail: ", subCmd.flagSet.Args())

	subCmd.run()
}
