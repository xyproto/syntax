package main

import (
	"flag"
	//"fmt"
	"io/ioutil"

	"github.com/xyproto/syntax"
	"github.com/xyproto/textoutput"
)

// highlight cat

func main() {
	flag.Parse()

	o := textoutput.NewTextOutput(true, true)

	if flag.NArg() != 1 {
		o.ErrExit("Must specify exactly 1 filename argument.")
	}

	input, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		o.ErrExit(err.Error())
	}

	text, err := syntax.AsText(input)
	if err != nil {
		o.ErrExit(err.Error())
	}

	textWithTags := string(text)
	//fmt.Println(textWithTags)
	o.OutputTags(textWithTags)
}
