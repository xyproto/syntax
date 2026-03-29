package main

import (
	"flag"
	"io"
	"os"
	"strings"

	"github.com/xyproto/files"
	"github.com/xyproto/syntax"
	"github.com/xyproto/vt"
)

func main() {
	flag.Parse()
	var (
		o    = vt.New()
		args = flag.Args()
	)

	// Apply the O_THEME setting.
	syntax.SetDefaultTextConfigFromEnv()

	if files.DataReadyOnStdin() {
		inputBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			o.ErrExit("Failed to read from stdin: " + err.Error())
		}
		if err := syntax.CatBytes(inputBytes, o); err != nil {
			o.ErrExit(err.Error())
		}
	} else if len(args) > 0 {
		var anyFilesProcessed bool
		for _, arg := range args {
			if files.Exists(arg) {
				inputBytes, err := os.ReadFile(arg)
				if err != nil {
					o.ErrExit("Failed to read file '" + arg + "': " + err.Error())
				}
				if err := syntax.CatBytes(inputBytes, o); err != nil {
					o.ErrExit(err.Error())
				}
				anyFilesProcessed = true
			}
		}
		if !anyFilesProcessed {
			inputBytes := []byte(strings.Join(args, " "))
			if err := syntax.CatBytes(inputBytes, o); err != nil {
				o.ErrExit(err.Error())
			}
		}
	} else {
		o.ErrExit("Please provide input via stdin, given filenames or command line arguments.")
	}
}
