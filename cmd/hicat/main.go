package main

import (
	"flag"
	"io"
	"os"
	"strings"

	"github.com/xyproto/files"
	"github.com/xyproto/syntax"
	"github.com/xyproto/textoutput"
)

func main() {
	flag.Parse()
	var (
		o    = textoutput.NewTextOutput(true, true)
		args = flag.Args()
	)

	if files.DataReadyOnStdin() { // Pri 1: read from stdin, if possible
		inputBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			o.ErrExit("Failed to read from stdin: " + err.Error())
		}
		if err := syntax.CatBytes(inputBytes, o); err != nil {
			o.ErrExit(err.Error())
		}
	} else if len(args) > 0 { // Pri 2: read from given filename(s), if possible
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
		if !anyFilesProcessed { // Pri 3: use the given arguments as the input data, as a last restort
			inputBytes := []byte(strings.Join(args, " "))
			if err := syntax.CatBytes(inputBytes, o); err != nil {
				o.ErrExit(err.Error())
			}
		}
	} else {
		// No data on stdin and no arguments
		o.ErrExit("Please provide input via stdin, given filenames or command line arguments.")
	}
}
