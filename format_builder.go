package cmd

import (
	"strings"

	ptxt "github.com/jedib0t/go-pretty/v6/text"
)

type formatBuilder struct {
	level   int
	maxLen  int // max length for a row
	builder *strings.Builder
}

func (fb *formatBuilder) String() string {
	return fb.builder.String()
}

// add a indent
func (fb *formatBuilder) nextLevel() formatBuilder {
	return formatBuilder{
		level:   fb.level + 1,
		maxLen:  fb.maxLen,
		builder: fb.builder,
	}
}

// remove a indent
func (fb *formatBuilder) lastLevel() formatBuilder {
	l := fb.level
	if l > 0 {
		l -= 1
	}
	return formatBuilder{
		level:   l,
		maxLen:  fb.maxLen,
		builder: fb.builder,
	}
}

// receive multi strings and add new line to the end of each string.
// If no args given, it will add a new line.
func (fb *formatBuilder) out(strs ...string) {
	if len(strs) < 1 {
		fb.builder.WriteRune('\n')
		return
	}

	maxLen := fb.maxLen - 2*fb.level // each level has a more indent (two spaces)

	for _, s := range strs {
		s = ptxt.WrapSoft(s, maxLen)
		ss := strings.Split(s, "\n")
		for _, s := range ss {
			fb.builder.WriteString(strings.Repeat("  ", fb.level)) // add indent
			fb.builder.WriteString(s)
			fb.builder.WriteRune('\n')
		}
	}
}

// example:
//
//	 outWithLeading("xxx:", 6, "yyyy yyyyyyyy") // with level:0 maxlen:10
//	 =>
//	"xxx:  yyyy
//	       yyyyyyyy
//	"
func (fb *formatBuilder) outWithLeading(leading string, leftLen int, right string) {
	maxLen := fb.maxLen - 2*fb.level // each level has a more indent (two spaces)

	rights := strings.Split(ptxt.WrapSoft(right, maxLen-leftLen), "\n")
	fb.builder.WriteString(strings.Repeat("  ", fb.level)) // add indent
	fb.builder.WriteString(leading)
	fb.builder.WriteString(strings.Repeat(" ", leftLen-len(leading)))
	fb.builder.WriteString(rights[0])
	fb.builder.WriteRune('\n')

	for _, r := range rights[1:] {
		// add indent and empty spaces for left
		fb.builder.WriteString(strings.Repeat(" ", 2*fb.level+leftLen))

		fb.builder.WriteString(r)
		fb.builder.WriteRune('\n')
	}
}
