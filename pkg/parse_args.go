package pkg

import ()

type Commands string

const (
	Commit Commands = "commit"
	Init   Commands = "init"
	Diff   Commands = "diff"
)

type CLIArgs interface {
	GetCommand() Commands
	GetPrefixAndValue(prefix string) (string, string)
	Validate() bool
	GetNamedArgCount() int
}

type CLIArgsParser struct {
	args []string
}

func NewCLIArgsParser(args []string) *CLIArgsParser {
	return &CLIArgsParser{
		args: args,
	}
}

func (self *CLIArgsParser) GetCommand() string {
	// TODO: implement
	return "commit"
}

func (self *CLIArgsParser) GetPrefixAndValue(prefix string) (string, string) {
	// TODO: implement
	return "a", "c"
}

func (self *CLIArgsParser) Validate() bool {
	// TODO: implement
	return true
}

func (self *CLIArgsParser) GetNamedArgCount() int {
	return 1
}
