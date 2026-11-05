package syntax

import (
	"maps"
	"strings"
	"testing"

	"github.com/xyproto/mode"
)

func TestShellKeywords(t *testing.T) {
	defaultKeywords := maps.Clone(Keywords)
	defer func() { Keywords = defaultKeywords }()
	tests := []struct {
		m       mode.Mode
		word    string
		keyword bool
	}{
		{mode.Shell, "set", false},
		{mode.Shell, "done", true},
		{mode.Fish, "set", true},
		{mode.Fish, "function", true},
		{mode.Fish, "end", true},
		{mode.Oil, "proc", true},
		{mode.Oil, "setvar", true},
		{mode.Oil, "var", true},
	}
	for _, tc := range tests {
		Keywords = maps.Clone(defaultKeywords)
		AdjustKeywords(tc.m)
		highlighted, err := AsText([]byte(tc.word+" x"), tc.m)
		if err != nil {
			t.Fatal(err)
		}
		wordTag, _, _ := strings.Cut(string(highlighted), tc.word)
		_, rest, _ := strings.Cut(string(highlighted), " ")
		plainTag, _, _ := strings.Cut(rest, "x")
		if got := wordTag != plainTag; got != tc.keyword {
			t.Errorf("expected %q to be highlighted as a keyword in %s mode: %v, got %q", tc.word, tc.m, tc.keyword, highlighted)
		}
	}
}
