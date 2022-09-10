package gocmd

import (
	"errors"
	"strings"
)

type FlagType byte

const (
	FlagTypeBool  FlagType = iota // a flag receiving none value
	FlagTypeValue                 // a flag receiving a single value
	FlagTypeMulti                 // a flag receiving multi values
)

type Flag struct {
	Long  string // long name
	Short string // short name
	Usage string
	Type  FlagType
	value any
}

// IsSet check whether the value for the flag is set. If the flag value's
// type is not matched to flag's type, (e.x., FlagTypeBool should match a
// bool value), it also return false.
func (f *Flag) IsSet() bool {
	if f.value == nil {
		return false
	}
	switch f.value.(type) {
	case bool:
		return f.Type == FlagTypeBool
	case string:
		return f.Type == FlagTypeValue
	case []string:
		return f.Type == FlagTypeMulti
	}
	return false
}

// Value return flag's value. If not set, return (nil, false).
// If set, it ensures the returned `v` is in proper type:
//
//	f.Type == FlagTypeBool   ->  v is bool
//	f.Type == FlagTypeValue  ->  v is string
//	f.Type == FlagTypeMulti  ->  v is []string
func (f *Flag) Value() (v any, ok bool) {
	if f.IsSet() {
		return f.value, true
	}
	return nil, false
}

func isFlag(arg string) (yes bool, isLongFlag bool, name string) {
	if strings.HasPrefix(arg, "--") {
		return true, true, arg[2:]
	} else if strings.HasPrefix(arg, "-") {
		return true, false, arg[1:]
	}
	return false, false, ""
}

func matchFlag(name string, isLong bool, candidates []*Flag) *Flag {
	for _, f := range candidates {
		if (isLong && f.Long == name) || (!isLong && f.Short == name) {
			return f
		}
	}
	return nil
}

// ParseFlags will regard args that start with "-"/"--" as flag arg.
// The arg "--xx" and arg "-y" will be regarded as a long flag name "xx" and a
// short flag name "y".
//
// For a flag arg, if found a same name flag in `candidates`, the arg is matched,
// otherwise, the arg is returned in `unknown`.
//
// For a matched flag, if its Type is `FlagTypeBool`, it will set the flag's value to
// true; if Type is `FlagTypeValue`, it will read the next non-flag arg as its value;
// if Type is `FlagTypeMulti`, it will read all following args before the next flag arg
// as its value. If Type is `FlagTypeValue` or `FlagTypeMulti`, but no next proper
// args can be read as value, it will return ""no value given..." error. Examples:
//
//	f is a flag of FlagTypeBool type with short name "f":
//	  '-f foo' -> value of f is set to true, "foo" is not read
//	f is a flag of FlagTypeValue type with short name "f":
//	  '-f value1 value2' -> value of f is set to "value1", "value2" is not read
//	f is a flag of FlagTypeMulti type with short name "f":
//	  '-f value1 value2 --foo' -> value of f is set to ["value1", "value2"], "--foo" is not read
//
// It will parse `args` to proper flags in `candidates` one by one until arg silce ends or a
// non-flag and non-flag-value arg is read. Remaining args will be returned as `remaining`.
func ParseFlags(args []string, candidates []*Flag) (remaining, unknown []string, err error) {
	remaining = args
	for len(remaining) > 0 {
		var flag *Flag
		arg := remaining[0]
		if yes, isLong, name := isFlag(arg); !yes {
			break
		} else {
			remaining = remaining[1:]
			if flag = matchFlag(name, isLong, candidates); flag == nil {
				unknown = append(unknown, arg)
				continue
			}
		}

		// read value(s) for the flag
		switch flag.Type {
		case FlagTypeBool:
			flag.value = true
		case FlagTypeValue:
			// if next arg exists and is not a flag arg
			if len(remaining) > 0 && !strings.HasPrefix(remaining[0], "-") {
				flag.value = remaining[0]
				remaining = remaining[1:]
			} else {
				return remaining, unknown, errors.New("no value given for flag: " + arg)
			}
		case FlagTypeMulti:
			// read all non-flag args as values
			var v []string
			for len(remaining) > 0 && !strings.HasPrefix(remaining[0], "-") {
				v = append(v, remaining[0])
				remaining = remaining[1:]
			}
			if len(v) < 1 {
				return remaining, unknown, errors.New("no value given for flag: " + arg)
			}
			flag.value = v
		}
	}
	return remaining, unknown, nil
}
