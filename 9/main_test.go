package main

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		// обычные кейсы
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"", "", false},

		// ошибки
		{"45", "", true},
		{"3abc", "", true},
		{`ab\`, "", true},

		// экранирование
		{`qwe\4\5`, "qwe45", false},
		{`qwe\45`, "qwe44444", false},
		{`\\5`, `\\\\\`, false}, // экранированный слеш, повторяем 5 раз
	}

	for _, tt := range tests {
		got, err := Unpack(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("unpack(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.expected {
			t.Errorf("unpack(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
