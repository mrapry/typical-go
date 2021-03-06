package typast

import (
	"fmt"
	"strings"
)

type (
	// Signature for code generation by annotation
	Signature struct {
		TagName string
		Help    string
	}
)

var _ fmt.Stringer = (*Signature)(nil)

func (s Signature) String() string {
	var out strings.Builder
	fmt.Fprintln(&out, "Autogenerated by Typical-Go. DO NOT EDIT.")
	if s.TagName != "" {
		fmt.Fprintf(&out, "\nTagName:\n\t%s\n", s.TagName)
	}
	if s.Help != "" {
		fmt.Fprintf(&out, "\nHelp:\n\t%s\n", s.Help)
	}
	return out.String()
}
