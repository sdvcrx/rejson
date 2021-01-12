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

func splitTag(s string) tag {
	idx := strings.Index(s, ":")
	if idx == -1 {
		return tag{Value: s}
	}

	return tag{
		Type:  s[:idx],
		Value: s[idx+1:],
	}
}

func parseTag(t string) (tag, error) {
	v := strings.TrimSpace(t)

	switch {
	case v == tagTypeEmpty:
		return tag{Type: tagTypeEmpty}, nil
	case v == tagTypeIgnore:
		return tag{Type: tagTypeIgnore}, nil
	default:
		tg := splitTag(v)
		if tg.Type == "" {
			tg.Type = tagTypePath
		}
		return tg, nil
	}
}
