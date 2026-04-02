package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var std = func() *log.Logger {
	spew.Config = spew.ConfigState{
		Indent:                  "\t",
		MaxDepth:                30,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		SortKeys:                true,
		SpewKeys:                true,
	}

	return log.New(os.Stderr, "", 0)
}()

func Logger(args ...any) {
	if len(args) == 0 {
		return
	}

	level := "INFO"
	color := Green

	var err error
	var msg string
	var dumpArgs []any

	for _, a := range args {
		switch v := a.(type) {
		case string:
			if msg == "" {
				msg = v
			} else {
				dumpArgs = append(dumpArgs, v)
			}
		case error:
			if v != nil && err == nil {
				err = v
				level = "ERROR"
				color = Red
				if msg != "" {
					msg = fmt.Errorf("%s:\n%w", msg, err).Error()
				} else {
					msg = err.Error()
				}
			}
		default:
			dumpArgs = append(dumpArgs, v)
		}
	}

	ts := time.Now().UTC().Format("02-01-2006 at 15:04")

	// color only the prefix
	prefix := fmt.Sprintf("%s%s [%s] %s", color, ts, level, Reset)
	std.SetPrefix(prefix)

	if msg == "" && len(dumpArgs) == 0 {
		return
	}

	if len(dumpArgs) == 0 {
		std.Print(msg)

		return
	}

	if msg == "" {
		std.Print(spew.Sdump(dumpArgs...))
		return
	}

	std.Printf("\n%s\n%s\n", msg, spew.Sdump(dumpArgs...))
}

func LogAndExit(args ...any) {
	Logger(args...)
	os.Exit(1)
}
