package gormqs_test

import (
	"testing"

	"github.com/foxie.io/gormqs"
)

func TestSafeTextForSql(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "hello world"},
		{"hello*world", "hello%world"},
		{"hello*world%", "hello%world"},
		{"**hello_world**", "%%hello_world%%"},
		{"**hello_world;", "%%hello_world"},
		{"hello_world/+;value=2", "hello_world"},
	}

	for _, test := range tests {
		result := gormqs.SafeTextForSql(test.input)
		if result != test.expected {
			t.Errorf("SafeTextForSql('%s') got '%s', want '%s'", test.input, result, test.expected)
		}
	}
}
