package ui

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/lthiagol/maitre/internal/core"
)

// Everything the console renders came from somewhere an agent could write: a title, a
// command, a note. A newline in a title would corrupt the screen, and an escape
// sequence in one would inject colour. So the text is sanitised once, on the way in.
//
// The file on disk keeps what was written — the sanitising is a rendering concern, not a
// claim about the data.

var ansiPattern = regexp.MustCompile("\x1b\\[[0-9;?]*[ -/]*[@-~]")

// bidiControls are the Unicode direction overrides. They reorder what is displayed
// without changing the bytes — the "Trojan Source" trick.
const bidiControls = "\u202a\u202b\u202c\u202d\u202e\u2066\u2067\u2068\u2069\u200e\u200f"

// sanitize replaces the characters that would corrupt the terminal.
//
// Control characters become a space; bidi overrides are dropped; ANSI escapes are
// removed.
func sanitize(s string) string {
	if s == "" {
		return s
	}
	s = ansiPattern.ReplaceAllString(s, "")
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\n' || r == '\r' || r == '\t':
			return ' '
		case r < 0x20 || r == 0x7f:
			return ' '
		case strings.ContainsRune(bidiControls, r):
			return -1
		}
		return r
	}, s)
}

// sanitizePlan sanitises every string in the loaded plan.
//
// The walk is reflective so that a field added to the model later is covered without
// anyone remembering to add it here.
func sanitizePlan(recipes []core.Recipe, specials []core.Special, notes []core.PrepNote, journal core.Journal) {
	for i := range recipes {
		sanitizeValue(reflect.ValueOf(&recipes[i]).Elem())
	}
	for i := range specials {
		sanitizeValue(reflect.ValueOf(&specials[i]).Elem())
	}
	for i := range notes {
		sanitizeValue(reflect.ValueOf(&notes[i]).Elem())
	}
	for i := range journal.Events {
		sanitizeValue(reflect.ValueOf(&journal.Events[i]).Elem())
	}
}

// sanitizeValue walks a value and sanitises every string it finds.
func sanitizeValue(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.String:
		if v.CanSet() {
			v.SetString(sanitize(v.String()))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			sanitizeValue(v.Field(i))
		}
	case reflect.Pointer:
		if !v.IsNil() {
			sanitizeValue(v.Elem())
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			sanitizeValue(v.Index(i))
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			sanitizeValue(v.Index(i))
		}
	}
}
