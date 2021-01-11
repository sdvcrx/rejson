package rejson

import (
	"strings"
)

const (
	tagName = "jsonp"
)

const (
	tagTypeEmpty  = ""
	tagTypeIgnore = "-"
	tagTypeFunc   = "func"
	tagTypePath   = "path"
)

type tag struct {
	Type  string
	Value string
}

func parseTag(t string) (tag, error) {
	// var tags []tag

	v := strings.TrimSpace(t)

	switch {
	case v == tagTypeEmpty:
		return tag{Type: tagTypeEmpty}, nil
	case v == tagTypeIgnore:
		return tag{Type: tagTypeIgnore}, nil
	default:
		t := strings.SplitN(strings.TrimSpace(t), ":", 2)
		if len(t) == 2 {
			// "func:FuncName"
			return tag{
				Type:  t[0],
				Value: t[1],
			}, nil
		}

		// "path"
		return tag{
			Type:  tagTypePath,
			Value: t[0],
		}, nil
	}
}
