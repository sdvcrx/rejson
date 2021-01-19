package rejson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	tests := []struct {
		name string
		args string
		want tag
	}{
		{
			name: "ignore",
			args: "-",
			want: tag{
				Type:  tagTypeIgnore,
				Value: "",
			},
		}, {
			name: "empty",
			args: "",
			want: tag{
				Type:  tagTypeEmpty,
				Value: "",
			},
		}, {
			name: "path",
			args: "first_name",
			want: tag{
				Type:  tagTypePath,
				Value: "first_name",
			},
		}, {
			name: "func",
			args: "func:FormatTime",
			want: tag{
				Type:  tagTypeFunc,
				Value: "FormatTime",
			},
		}, {
			name: "complex",
			args: "xxx.@comp",
			want: tag{
				Type:  tagTypePath,
				Value: "xxx.@comp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTag(tt.args)
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.Value, got.Value)
		})
	}
}
