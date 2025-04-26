package internal

import (
	"errors"
	"os"
	"slices"
	"strings"
)

type Command string

const (
	Commit Command = "commit"
	Init   Command = "init"
	Diff   Command = "diff"
)

var CommandArr []Command = []Command{Commit, Init, Diff}

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

func parse() Parsed {
	asd := os.Args[1:]
	var foundFlags []Flag
	var args []Command

	for _, arg := range asd {
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
			args = append(args, Command(arg))
		}
	}

	return Parsed{
		commands: args,
		flags:    foundFlags,
	}
}

type CLIArgsParser struct {
	args   []string
	parsed Parsed
}

func NewCLIArgsParser(args []string) *CLIArgsParser {
	parsed := parse()

	return &CLIArgsParser{
		args:   args,
		parsed: parsed,
	}
}

func (self *CLIArgsParser) Validate() error {
	if len(self.parsed.commands) > 1 {
		return errors.New("can not process more than one argument at a time")
	}

	return nil
}

func (self *CLIArgsParser) GetCommand(command *Command) error {
	err := self.Validate()

	if err != nil {
		return err
	}

	*command = self.parsed.commands[0]

	return nil
}

func (self *CLIArgsParser) GetFlags(flags *[]Flag) error {
	err := self.Validate()

	if err != nil {
		return err
	}

	*flags = self.parsed.flags
	return nil
}
