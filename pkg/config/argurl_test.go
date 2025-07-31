package config_test

import (
	"reflect"
	"testing"

	"github.com/egoodhall/go-links/pkg/config"
)

func TestArgUrl(t *testing.T) {
	tests := []struct {
		name string
		url  config.ArgUrl
		want []string
	}{
		{
			name: "simple",
			url:  config.ArgUrl("https://example.com/$arg1/"),
			want: []string{"arg1"},
		},
		{
			name: "multiple",
			url:  config.ArgUrl("https://example.com/$arg1/$arg2?x=$arg3"),
			want: []string{"arg1", "arg2", "arg3"},
		},
		{
			name: "colons",
			url:  config.ArgUrl("https://example.com/:arg1"),
			want: []string{"arg1"},
		},
		{
			name: "repeats",
			url:  config.ArgUrl("https://example.com/$arg1/$arg1"),
			want: []string{"arg1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.url.Args()
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Args() = %v, want %v", got, test.want)
			}
		})
	}
}
