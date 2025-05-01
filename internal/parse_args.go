package internal

import (
	"fmt"
	"slices"
	"strings"
)

type Command string

const (
	Commit Command = "commit"
	Init   Command = "init"
	Diff   Command = "diff"
	Revert Command = "revert"
)

var CommandArr []Command = []Command{Commit, Init, Diff, Revert}

type Flag struct {
	key   string
	value string
}

type Parsed struct {
	commands []Command
	flags    []Flag
}

type CLIArgs interface {
	Validate() error
	GetCommand(command *string) Command
	GetFlags(flags *[]Flag) error
}

func parse(args []string) Parsed {
	var foundFlags []Flag
	var commands []Command

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			temp := strings.Split(arg, "=")
			keyPart := temp[0]
			valuePart := temp[1]
			keyPart = strings.TrimPrefix(keyPart, "--")

			foundFlags = append(foundFlags, Flag{
				key:   keyPart,
				value: valuePart,
			})
		} else if slices.Contains(CommandArr, Command(arg)) {
			commands = append(commands, Command(arg))
		}
	}

	return Parsed{
		commands: commands,
		flags:    foundFlags,
	}
}

type CLIArgsParser struct {
	args   []string
	parsed Parsed
}

func NewCLIArgsParser(args []string) *CLIArgsParser {
	parsed := parse(args)

	return &CLIArgsParser{
		args:   args,
		parsed: parsed,
	}
}

func (parser *CLIArgsParser) Validate() error {
	if len(parser.parsed.commands) != 1 {
		return fmt.Errorf("unexpected amount of commands\n")
	}

	if parser.parsed.commands[0] == Revert && len(parser.parsed.flags) < 1 {
		found := false
		for _, flag := range parser.parsed.flags {
			if flag.key == "back" {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("revert command needs a --back=X flag\n")

		}
	}

	return nil
}

func (parser *CLIArgsParser) GetCommand(command *Command) error {
	err := parser.Validate()

	if err != nil {
		return err
	}

	*command = parser.parsed.commands[0]

	return nil
}

func (parser *CLIArgsParser) GetFlags(flags *[]Flag) error {
	err := parser.Validate()

	if err != nil {
		return err
	}

	*flags = parser.parsed.flags
	return nil
}
