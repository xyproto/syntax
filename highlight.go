// Package syntax provides syntax highlighting for source code using the same
// approach as Orbiton. It tokenizes input via text/scanner, classifies tokens
// into kinds, and wraps them in color tags that can be converted to ANSI escape
// codes by the vt package.
//
// Theme selection is supported via the O_THEME (or THEME) environment variable.
package syntax

import (
	"bytes"
	"io"
	"text/scanner"
	"unicode"

	"github.com/sourcegraph/annotate"
	"github.com/xyproto/mode"
)

// Kind represents a syntax highlighting kind (class) assigned to tokens.
type Kind uint8

// Supported highlighting kinds.
const (
	Whitespace Kind = iota
	AndOr
	AngleBracket
	AssemblyEnd
	Class
	Comment
	Decimal
	Dollar
	Literal
	Keyword
	Mut
	Plaintext
	Private
	Protected
	Public
	Punctuation
	Self
	Star
	Static
	String
	Tag
	TextAttrName
	TextAttrValue
	TextTag
	Type
)

// TextConfig maps token kinds to color tag names used by vt.TextOutput.DarkTags.
type TextConfig struct {
	AndOr         string
	AngleBracket  string
	AssemblyEnd   string
	Class         string
	Comment       string
	Decimal       string
	Dollar        string
	Keyword       string
	Literal       string
	Mut           string
	Plaintext     string
	Private       string
	Protected     string
	Public        string
	Punctuation   string
	Self          string
	Star          string
	Static        string
	String        string
	Tag           string
	TextAttrName  string
	TextAttrValue string
	TextTag       string
	Type          string
	Whitespace    string
}

// GetClass returns the color tag name for a given token kind.
func (c TextConfig) GetClass(kind Kind) string {
	switch kind {
	case String:
		return c.String
	case Keyword:
		return c.Keyword
	case Comment:
		return c.Comment
	case Type:
		return c.Type
	case Literal:
		return c.Literal
	case Punctuation:
		return c.Punctuation
	case Plaintext:
		return c.Plaintext
	case Tag:
		return c.Tag
	case TextTag:
		return c.TextTag
	case TextAttrName:
		return c.TextAttrName
	case TextAttrValue:
		return c.TextAttrValue
	case Decimal:
		return c.Decimal
	case AndOr:
		return c.AndOr
	case AngleBracket:
		return c.AngleBracket
	case Dollar:
		return c.Dollar
	case Star:
		return c.Star
	case Static:
		return c.Static
	case Self:
		return c.Self
	case Class:
		return c.Class
	case Public:
		return c.Public
	case Private:
		return c.Private
	case Protected:
		return c.Protected
	case AssemblyEnd:
		return c.AssemblyEnd
	case Mut:
		return c.Mut
	}
	return ""
}

// Option is a function that modifies a TextConfig.
type Option func(*TextConfig)

// Printer renders highlighted output.
type Printer interface {
	Print(w io.Writer, kind Kind, tokText string) error
}

// TextPrinter wraps TextConfig to implement Printer, emitting color tags.
type TextPrinter TextConfig

// Print writes token text wrapped in <color>...<off> tags based on its kind.
func (p TextPrinter) Print(w io.Writer, kind Kind, tokText string) error {
	class := TextConfig(p).GetClass(kind)
	if class != "" {
		if _, err := io.WriteString(w, "<"+class+">"); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, tokText); err != nil {
		return err
	}
	if class != "" {
		if _, err := io.WriteString(w, "<off>"); err != nil {
			return err
		}
	}
	return nil
}

// Annotator produces syntax highlighting annotations.
type Annotator interface {
	Annotate(start int, kind Kind, tokText string) (*annotate.Annotation, error)
}

// TextAnnotator wraps TextConfig to implement Annotator.
type TextAnnotator TextConfig

// Annotate returns an annotation for the given token.
func (a TextAnnotator) Annotate(start int, kind Kind, tokText string) (*annotate.Annotation, error) {
	class := TextConfig(a).GetClass(kind)
	if class != "" {
		left := []byte("<" + class + ">")
		return &annotate.Annotation{
			Start: start, End: start + len(tokText),
			Left: left, Right: []byte("<off>"),
		}, nil
	}
	return nil, nil
}

// DefaultTextConfig provides color names matching Orbiton's default theme.
var DefaultTextConfig = TextConfig{
	AndOr:         "red",
	AngleBracket:  "red",
	AssemblyEnd:   "lightyellow",
	Class:         "white",
	Comment:       "darkgray",
	Decimal:       "red",
	Dollar:        "white",
	Keyword:       "red",
	Literal:       "white",
	Mut:           "magenta",
	Plaintext:     "white",
	Private:       "red",
	Protected:     "red",
	Public:        "red",
	Punctuation:   "red",
	Self:          "magenta",
	Star:          "white",
	Static:        "lightyellow",
	String:        "lightwhite",
	Tag:           "white",
	TextAttrName:  "white",
	TextAttrValue: "white",
	TextTag:       "white",
	Type:          "white",
	Whitespace:    "",
}

// Print scans tokens from s and writes highlighted output using Printer p.
func Print(s *scanner.Scanner, w io.Writer, p Printer, m mode.Mode) error {
	switch m {
	case mode.C3:
		s.IsIdentRune = func(ch rune, i int) bool {
			return ch == '$' || ch == '@' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
		}
	case mode.Clojure, mode.Lisp:
		s.IsIdentRune = func(ch rune, i int) bool {
			return ch == '*' || ch == '-' || ch == '+' || ch == '/' || ch == '?' || ch == '!' || ch == '.' || ch == ':' || ch == '&' || ch == '<' || ch == '>' || ch == '=' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
		}
	case mode.Shell, mode.Make, mode.Just:
		s.IsIdentRune = func(ch rune, i int) bool {
			return ch == '-' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
		}
	case mode.Swift:
		s.IsIdentRune = func(ch rune, i int) bool {
			return ch == '#' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
		}
	case mode.Vibe67:
		s.IsIdentRune = func(ch rune, i int) bool {
			return ch == '&' || ch == '<' || ch == '>' || ch == '^' || ch == '|' || ch == '~' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
		}
	}
	inComment := false
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tokText := s.TokenText()
		if err := p.Print(w, tokenKind(tok, tokText, &inComment, m), tokText); err != nil {
			return err
		}
	}
	return nil
}

// Annotate tokenizes src and returns annotations for mode m.
func Annotate(src []byte, a Annotator, m mode.Mode) (annotate.Annotations, error) {
	var (
		anns      annotate.Annotations
		s         = NewScanner(src)
		read      = 0
		inComment = false
		tok       = s.Scan()
	)
	for tok != scanner.EOF {
		tokText := s.TokenText()
		ann, err := a.Annotate(read, tokenKind(tok, tokText, &inComment, m), tokText)
		if err != nil {
			return nil, err
		}
		read += len(tokText)
		if ann != nil {
			anns = append(anns, ann)
		}
		tok = s.Scan()
	}
	return anns, nil
}

// AsText returns src highlighted for mode m, applying options to TextConfig.
func AsText(src []byte, m mode.Mode, options ...Option) ([]byte, error) {
	cfg := DefaultTextConfig
	for _, opt := range options {
		opt(&cfg)
	}
	var buf bytes.Buffer
	if err := Print(NewScanner(src), &buf, TextPrinter(cfg), m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NewScanner returns a scanner.Scanner configured for syntax highlighting.
func NewScanner(src []byte) *scanner.Scanner {
	return NewScannerReader(bytes.NewReader(src))
}

// NewScannerReader returns a scanner.Scanner configured for syntax highlighting from r.
func NewScannerReader(r io.Reader) *scanner.Scanner {
	var s scanner.Scanner
	s.Init(r)
	s.Error = func(*scanner.Scanner, string) {}
	s.Whitespace = 0
	s.Mode ^= scanner.SkipComments
	return &s
}
