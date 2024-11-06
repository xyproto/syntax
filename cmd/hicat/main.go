package main

import (
	"flag"
	"io"
	"os"
	"strings"

	"github.com/xyproto/files"
	"github.com/xyproto/mode"
	"github.com/xyproto/syntax"
	"github.com/xyproto/textoutput"
)

func main() {
	flag.Parse()
	var (
		o          = textoutput.NewTextOutput(true, true)
		args       = flag.Args()
		lenargs    = len(args)
		inputBytes []byte
	)
	if lenargs != 1 {
		o.ErrExit("Must specify exactly 1 filename argument.")
	}
	if files.DataReadyOnStdin() || (lenargs > 0 && args[lenargs-1] == "-") {
		if data, err := io.ReadAll(os.Stdin); err == nil { // success
			inputBytes = data
		}
	} else { // Set msg to the given arguments (if any)
		inputBytes = []byte(strings.Join(args, " "))
		// If the first argument is a filename, then use that one for the input
		if lenargs > 0 && files.Exists(args[0]) {
			data, err := os.ReadFile(args[0])
			if err == nil { // success
				inputBytes = data
			}
		}
	}
	mode := mode.SimpleDetectBytes(inputBytes)
	text, err := syntax.AsText(inputBytes, mode)
	if err != nil {
		o.ErrExit(err.Error())
	}
	textWithTags := string(text)
	o.OutputTags(textWithTags)
}
