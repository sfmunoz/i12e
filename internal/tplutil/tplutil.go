package tplutil

import (
	"fmt"
	"strings"
	"text/template"
)

func add(a, b int) int {
	return a + b
}

func indent(spaces int, v fmt.Stringer) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v.String(), "\n", "\n"+pad)
}

func nindent(spaces int, v fmt.Stringer) string {
	return "\n" + indent(spaces, v)
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"add":     add,
		"indent":  indent,
		"nindent": nindent,
	}
}
