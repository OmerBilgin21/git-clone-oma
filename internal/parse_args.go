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
	Key   string
	Value string
}

type Parsed struct {
	commands []Command
	flags    []Flag
}

type CLIArgs interface {
	Validate() error
	GetCommand(command *string) Command
	GetFlags() ([]Flag, error)
	GetFlag(key string) (Flag, error)
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
				Key:   keyPart,
				Value: valuePart,
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

	commitMissingMsgErrMsg := "commit command needs a --message='<your-message>' flag"
	revertMissingBackErrMsg := "revert command needs a --back=X flag"

	if parser.parsed.commands[0] == Commit && len(parser.parsed.flags) < 1 {
		return fmt.Errorf("%v\n", commitMissingMsgErrMsg)
	}

	if parser.parsed.commands[0] == Revert && len(parser.parsed.flags) >= 1 {
		found := false

		for _, flag := range parser.parsed.flags {
			if flag.Key == "message" {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("--message='<your-message>' was not found among your flags, %v\n", revertMissingBackErrMsg)
		}
	}

	if parser.parsed.commands[0] == Revert && len(parser.parsed.flags) < 1 {
		return fmt.Errorf("%v\n", revertMissingBackErrMsg)
	}

	if parser.parsed.commands[0] == Revert && len(parser.parsed.flags) >= 1 {
		found := false

		for _, flag := range parser.parsed.flags {
			if flag.Key == "back" {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("--back=X was not found among your flags, %v\n", revertMissingBackErrMsg)
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

func (parser *CLIArgsParser) GetFlags() ([]Flag, error) {
	err := parser.Validate()

	if err != nil {
		return nil, err
	}

	return parser.parsed.flags, nil
}

func (parser *CLIArgsParser) GetFlag(key string) (Flag, error) {
	for _, flag := range parser.parsed.flags {
		if flag.Key == key {
			return flag, nil
		}
	}

	return Flag{}, fmt.Errorf("flag was not found!")
}
