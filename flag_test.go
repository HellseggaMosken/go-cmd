package cmd_test

import (
	"reflect"
	"strings"
	"testing"

	cmd "github.com/HellseggaMosken/go-cmd"
)

func TestParseFlagNoParse(t *testing.T) {
	flags := []*cmd.Flag{}
	args := []string{"a"}
	r, u, e := cmd.ParseFlags(args, flags)
	if e != nil {
		t.Error(e)
	}
	if !reflect.DeepEqual(args, r) {
		t.Error("remaining args should be same to orginal args, r:", r)
	}
	if len(u) > 0 {
		t.Error("should not have unknown flag args, u:", u)
	}
}

func TestParseFlagLongAndShort(t *testing.T) {
	flags := []*cmd.Flag{
		{
			Long:  "long",
			Short: "l",
			Type:  cmd.FlagTypeBool,
		},
		{
			Long:  "short",
			Short: "s",
			Type:  cmd.FlagTypeBool,
		},
	}
	args := []string{"--long", "-s"}
	_, _, _ = cmd.ParseFlags(args, flags)
	if !flags[0].IsSet() {
		t.Error("long flag is not set")
	}
	if !flags[1].IsSet() {
		t.Error("short flag is not set")
	}
}

func TestParseFlagUnknownAndRemaining(t *testing.T) {
	flags := []*cmd.Flag{}
	args := []string{"-b", "--c", "-", "--", "foo"}
	r, u, e := cmd.ParseFlags(args, flags)
	if e != nil {
		t.Error(e)
	}
	if !reflect.DeepEqual([]string{"foo"}, r) {
		t.Error("remaining args wrong.", r)
	}
	if !reflect.DeepEqual(u, []string{"-b", "--c", "-", "--"}) {
		t.Error("wrong unknown flag arg:", u[0])
	}
}

func TestParseFlagBool(t *testing.T) {
	flags := []*cmd.Flag{
		{
			Long:  "bool",
			Short: "b",
			Type:  cmd.FlagTypeBool,
		},
	}
	args := []string{"-b"}
	_, _, e := cmd.ParseFlags(args, flags)
	if e != nil {
		t.Error(e)
	}
	if v, ok := flags[0].Value(); !ok {
		t.Error("flag is not set")
	} else if _, ok := v.(bool); !ok {
		t.Error("flag type is wrong")
	}
}

func TestParseFlagValue(t *testing.T) {
	flags := []*cmd.Flag{
		{
			Long:  "value",
			Short: "v",
			Type:  cmd.FlagTypeValue,
		},
	}
	_, _, e := cmd.ParseFlags([]string{"-v", "value1"}, flags)
	if e != nil {
		t.Error(e)
	}
	if v, ok := flags[0].Value(); !ok {
		t.Error("flag is not set")
	} else if v, ok := v.(string); !ok {
		t.Error("flag type is wrong")
	} else if v != "value1" {
		t.Error("flag value is wrong, v:", v)
	}
	_, _, e = cmd.ParseFlags([]string{"-v"}, flags)
	if e == nil {
		t.Error("should return error")
	}
	_, _, e = cmd.ParseFlags([]string{"-v", "-v2", "foo"}, flags)
	if e == nil {
		t.Error("should return error")
	}
}

func TestParseFlagMulti(t *testing.T) {
	flags := []*cmd.Flag{
		{
			Long:  "multi",
			Short: "m",
			Type:  cmd.FlagTypeMulti,
		},
	}
	_, _, e := cmd.ParseFlags([]string{"-m", "value1", "value2"}, flags)
	if e != nil {
		t.Error(e)
	}
	if v, ok := flags[0].Value(); !ok {
		t.Error("flag is not set")
	} else if v, ok := v.([]string); !ok {
		t.Error("flag type is wrong")
	} else if !reflect.DeepEqual(v, []string{"value1", "value2"}) {
		t.Error("flag value is wrong, v:", v)
	}
	_, _, e = cmd.ParseFlags([]string{"-m"}, flags)
	if e == nil {
		t.Error("should return error")
	}
	_, _, e = cmd.ParseFlags([]string{"-m", "--m2", "foo"}, flags)
	if e == nil {
		t.Error("should return error")
	}
}

func TestParseFlagFull(t *testing.T) {
	flags := []*cmd.Flag{
		{
			Long:  "bool1",
			Short: "b1",
			Type:  cmd.FlagTypeBool,
		},
		{
			Long:  "bool2",
			Short: "b2",
			Type:  cmd.FlagTypeBool,
		},
		{
			Long:  "bool3",
			Short: "b3",
			Type:  cmd.FlagTypeBool,
		},
		{
			Long:  "value1",
			Short: "v1",
			Type:  cmd.FlagTypeValue,
		},
		{
			Long:  "value2",
			Short: "v2",
			Type:  cmd.FlagTypeValue,
		},
		{
			Long:  "value3",
			Short: "v3",
			Type:  cmd.FlagTypeValue,
		},
		{
			Long:  "multi1",
			Short: "m1",
			Type:  cmd.FlagTypeMulti,
		},
		{
			Long:  "multi2",
			Short: "m2",
			Type:  cmd.FlagTypeMulti,
		},
		{
			Long:  "multi3",
			Short: "m3",
			Type:  cmd.FlagTypeMulti,
		},
	}
	raw := "-b1 -v1 foo1 -u1 --bool2 -m1 foo2 foo3 foo4 --value2 foo5 --multi2 foo6 --u2 r1 r2"
	args := strings.Split(raw, " ")
	r, u, e := cmd.ParseFlags(args, flags)
	if e != nil {
		t.Error(e)
	}
	if !reflect.DeepEqual(r, []string{"r1", "r2"}) {
		t.Error("remaining args wrong, r:", r)
	}
	if !reflect.DeepEqual(u, []string{"-u1", "--u2"}) {
		t.Error("unknown flags wrong, r:", u)
	}
	if v, _ := flags[0].Value(); !reflect.DeepEqual(v, true) {
		t.Error("wrong value:", v)
	}
	if v, _ := flags[1].Value(); !reflect.DeepEqual(v, true) {
		t.Error("wrong value:", v)
	}
	if flags[2].IsSet() {
		t.Error("value should not be set:", flags[2])
	}
	if v, _ := flags[3].Value(); !reflect.DeepEqual(v, "foo1") {
		t.Error("wrong value:", v)
	}
	if v, _ := flags[4].Value(); !reflect.DeepEqual(v, "foo5") {
		t.Error("wrong value:", v)
	}
	if flags[5].IsSet() {
		t.Error("value should not be set:", flags[5])
	}
	if v, _ := flags[6].Value(); !reflect.DeepEqual(v, []string{"foo2", "foo3", "foo4"}) {
		t.Error("wrong value:", v)
	}
	if v, _ := flags[7].Value(); !reflect.DeepEqual(v, []string{"foo6"}) {
		t.Error("wrong value:", v)
	}
	if flags[8].IsSet() {
		t.Error("value should not be set:", flags[8])
	}
}
