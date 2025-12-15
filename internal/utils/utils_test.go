package utils

import (
	"testing"
)

func TestIsProfane(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected bool
	}{
		{
			name:     "profane word: kerfuffle",
			word:     "kerfuffle",
			expected: true,
		},
		{
			name:     "profane word: sharbert",
			word:     "sharbert",
			expected: true,
		},
		{
			name:     "profane word: fornax",
			word:     "fornax",
			expected: true,
		},
		{
			name:     "profane word uppercase: KERFUFFLE",
			word:     "KERFUFFLE",
			expected: true,
		},
		{
			name:     "profane word mixed case: ShArBeRt",
			word:     "ShArBeRt",
			expected: true,
		},
		{
			name:     "non-profane word: hello",
			word:     "hello",
			expected: false,
		},
		{
			name:     "non-profane word: world",
			word:     "world",
			expected: false,
		},
		{
			name:     "empty string",
			word:     "",
			expected: false,
		},
		{
			name:     "similar but not profane: kerfluffles",
			word:     "kerfluffles",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsProfane(tt.word)
			if result != tt.expected {
				t.Errorf("IsProfane(%q) = %v, want %v", tt.word, result, tt.expected)
			}
		})
	}
}

func TestCleanMessageProfane(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "single profane word",
			message:  "kerfuffle",
			expected: "****",
		},
		{
			name:     "multiple profane words",
			message:  "kerfuffle sharbert fornax",
			expected: "**** **** ****",
		},
		{
			name:     "mixed profane and clean words",
			message:  "hello kerfuffle world",
			expected: "hello **** world",
		},
		{
			name:     "profane words with uppercase",
			message:  "This is KERFUFFLE and sharbert",
			expected: "This is **** and ****",
		},
		{
			name:     "no profane words",
			message:  "hello world this is clean",
			expected: "hello world this is clean",
		},
		{
			name:     "empty string",
			message:  "",
			expected: "",
		},
		{
			name:     "single clean word",
			message:  "hello",
			expected: "hello",
		},
		{
			name:     "profane word at start",
			message:  "kerfuffle at the start",
			expected: "**** at the start",
		},
		{
			name:     "profane word at end",
			message:  "at the end is fornax",
			expected: "at the end is ****",
		},
		{
			name:     "consecutive profane words",
			message:  "kerfuffle sharbert fornax together",
			expected: "**** **** **** together",
		},
		{
			name:     "repeated profane words",
			message:  "kerfuffle kerfuffle kerfuffle",
			expected: "**** **** ****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanMessageProfane(tt.message)
			if result != tt.expected {
				t.Errorf("CleanMessageProfane(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}
