package cli

import (
	"fmt"
	"strings"
)

// flagSet is a small argument parser.
//
// The standard library's flag package stops at the first positional argument, which
// breaks the natural `maitre-cli recipe approve RCP07 --by alice` ordering. This parser
// accepts flags before, between and after positionals.
type flagSet struct {
	name  string
	bools map[string]*bool
	strs  map[string]*string
	ints  map[string]*int
	lists map[string]*[]string
	known map[string]bool
	args  []string
}

// newFlagSet returns the flag set for one verb.
//
// Every verb carries the three identity flags — who is acting, in which session, in
// which role — so a command can be stamped without an environment variable. They may
// also be given before the command: `maitre-cli --actor runner recipe claim RCP01`.
func newFlagSet(name string) *flagSet {
	fs := &flagSet{
		name:  name,
		bools: map[string]*bool{},
		strs:  map[string]*string{},
		ints:  map[string]*int{},
		lists: map[string]*[]string{},
		known: map[string]bool{},
	}
	fs.String("actor", "", "who is acting")
	fs.String("session", "", "the session this belongs to")
	fs.String("role", "", "the role this belongs to")
	return fs
}

// String declares a string flag.
func (f *flagSet) String(name, def, _ string) *string {
	v := def
	f.strs[name] = &v
	f.known[name] = true
	return &v
}

// Int declares an integer flag.
func (f *flagSet) Int(name string, def int, _ string) *int {
	v := def
	f.ints[name] = &v
	f.known[name] = true
	return &v
}

// Bool declares a boolean flag.
func (f *flagSet) Bool(name string, def bool, _ string) *bool {
	v := def
	f.bools[name] = &v
	f.known[name] = true
	return &v
}

// StringList declares a repeatable string flag.
func (f *flagSet) StringList(name string) *[]string {
	v := []string{}
	f.lists[name] = &v
	f.known[name] = true
	return &v
}

// StringValue returns a string flag's value, after parsing.
func (f *flagSet) StringValue(name string) string {
	if v, ok := f.strs[name]; ok {
		return *v
	}
	return ""
}

// parse reads the arguments.
func (f *flagSet) parse(args []string) error {
	positionalOnly := false
	for i := 0; i < len(args); i++ {
		a := args[i]

		if positionalOnly || !strings.HasPrefix(a, "-") || a == "-" {
			f.args = append(f.args, a)
			continue
		}
		if a == "--" {
			positionalOnly = true
			continue
		}

		name := strings.TrimLeft(a, "-")
		value := ""
		hasValue := false
		if eq := strings.Index(name, "="); eq >= 0 {
			name, value, hasValue = name[:eq], name[eq+1:], true
		}

		switch {
		case f.bools[name] != nil:
			if hasValue {
				switch strings.ToLower(value) {
				case "true", "1", "yes":
					*f.bools[name] = true
				case "false", "0", "no":
					*f.bools[name] = false
				default:
					return fmt.Errorf("--%s: expected a boolean, got %q", name, value)
				}
			} else {
				*f.bools[name] = true
			}

		case f.strs[name] != nil:
			if !hasValue {
				if i+1 >= len(args) {
					return fmt.Errorf("--%s needs a value", name)
				}
				i++
				value = args[i]
			}
			*f.strs[name] = value

		case f.ints[name] != nil:
			if !hasValue {
				if i+1 >= len(args) {
					return fmt.Errorf("--%s needs a value", name)
				}
				i++
				value = args[i]
			}
			n, err := atoi(value)
			if err != nil {
				return fmt.Errorf("--%s: %w", name, err)
			}
			*f.ints[name] = n

		case f.lists[name] != nil:
			if !hasValue {
				if i+1 >= len(args) {
					return fmt.Errorf("--%s needs a value", name)
				}
				i++
				value = args[i]
			}
			*f.lists[name] = append(*f.lists[name], value)

		default:
			return fmt.Errorf("unknown flag --%s", name)
		}
	}
	return nil
}

// nArg is the number of positional arguments.
func (f *flagSet) nArg() int { return len(f.args) }

// arg returns the i-th positional argument.
func (f *flagSet) arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}
