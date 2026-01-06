package main

import (
	"testing"
)

func TestReplaceEmojiShortcuts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple smiley",
			input:    "Hello :)",
			expected: "Hello :slightly_smiling_face:",
		},
		{
			name:     "Laughing face",
			input:    "That's funny =))",
			expected: "That's funny :laughing:",
		},
		{
			name:     "Multiple shortcuts",
			input:    "Hi :) How are you? =))",
			expected: "Hi :slightly_smiling_face: How are you? :laughing:",
		},
		{
			name:     "Heart emoji",
			input:    "I love this <3",
			expected: "I love this :heart:",
		},
		{
			name:     "Shortcut at start",
			input:    ":) Good morning!",
			expected: ":slightly_smiling_face: Good morning!",
		},
		{
			name:     "No shortcuts",
			input:    "Just a normal message",
			expected: "Just a normal message",
		},
		{
			name:     "Empty message",
			input:    "",
			expected: "",
		},
		{
			name:     "Wink emoji",
			input:    "Just kidding ;)",
			expected: "Just kidding :wink:",
		},
		{
			name:     "Sad face",
			input:    "That's sad :(",
			expected: "That's sad :slightly_frowning_face:",
		},
		{
			name:     "XD laughing",
			input:    "So funny XD",
			expected: "So funny :joy:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceEmojiShortcuts(tt.input)
			if result != tt.expected {
				t.Errorf("ReplaceEmojiShortcuts(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessMessageForEmoji(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal message with shortcut",
			input:    "Hello :)",
			expected: "Hello :slightly_smiling_face:",
		},
		{
			name:     "Code block should be preserved",
			input:    "```\n:) this is code\n```",
			expected: "```\n:) this is code\n```",
		},
		{
			name:     "Inline code should be preserved",
			input:    "Use `:)` for smiley",
			expected: "Use `:)` for smiley",
		},
		{
			name:     "Mixed content",
			input:    "Hello :) check `code :)` and :D",
			expected: "Hello :slightly_smiling_face: check `code :)` and :grinning:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessMessageForEmoji(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessMessageForEmoji(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsCodeBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Fenced code block",
			input:    "```\ncode here\n```",
			expected: true,
		},
		{
			name:     "Inline code",
			input:    "`code`",
			expected: true,
		},
		{
			name:     "Normal text",
			input:    "Just text",
			expected: false,
		},
		{
			name:     "Mixed with code",
			input:    "Some `code` here",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCodeBlock(tt.input)
			if result != tt.expected {
				t.Errorf("isCodeBlock(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
