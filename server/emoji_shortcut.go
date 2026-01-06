package main

import (
	"regexp"
	"strings"
)

// EmojiMapping defines the mapping from shortcut to standard Mattermost emoji
type EmojiMapping struct {
	Shortcut string // The shortcut pattern to match (e.g., "=))", ":))")
	Emoji    string // The Mattermost emoji to replace with (e.g., "smile")
}

// EmojiShortcuts contains all the emoji shortcut mappings
// Add or modify mappings here as needed
var EmojiShortcuts = []EmojiMapping{
	// Smileys
	{"Shortcut": ":))", "Emoji": ":yh_laughing:"},
	{"Shortcut": "=))", "Emoji": ":yh_rolling:"},
	{"Shortcut": ";))", "Emoji": ":yh_giggle:"},
	{"Shortcut": ";)", "Emoji": ":yh_winking:"},
	{"Shortcut": ";;)", "Emoji": ":yh_eyelashes:"},
	{"Shortcut": ":|", "Emoji": ":yh_straight_face:"},
	{"Shortcut": "\/:)", "Emoji": ":yh_raised_eyebrow:"},
	{"Shortcut": ":)", "Emoji": ":yh_smile:"},
	{"Shortcut": ":-?", "Emoji": ":yh_thinking:"},
	{"Shortcut": ":p", "Emoji": ":yh_sticking_tongue_out:"},
	{"Shortcut": ":D", "Emoji": ":yh_grin:"},
	{"Shortcut": ":((", "Emoji": ":yh_crying:"},
	{"Shortcut": ":(", "Emoji": ":yh_sad:"},
	{"Shortcut": ":-$", "Emoji": ":yh_shh:"},
	{"Shortcut": ":\">", "Emoji": ":yh_blushing:"},
	{"Shortcut": ":-s", "Emoji": ":yh_worried:"},
	{"Shortcut": ":o", "Emoji": ":yh_surprised:"},
	{"Shortcut": "\\m\/", "Emoji": ":yh_rocking:"},
	{"Shortcut": "=P~", "Emoji": ":yh_drooling:"},
	{"Shortcut": ":-j", "Emoji": ":yh_oh_go_on:"},
	{"Shortcut": "=D>", "Emoji": ":yh_applause:"},
	{"Shortcut": ":->", "Emoji": ":yh_smug:"},
	{"Shortcut": ":-w", "Emoji": ":yh_waiting:"},
	{"Shortcut": ":-x", "Emoji": ":yh_love_struck:"},
	{"Shortcut": ":-??", "Emoji": ":yh_dont_know:"},
	{"Shortcut": ":-\"", "Emoji": ":yh_whistling:"},
	{"Shortcut": "):", "Emoji": ":yh_sad:"},
}

// compiledPatterns holds the precompiled regex patterns for better performance
type compiledPattern struct {
	regex *regexp.Regexp
	emoji string
}

var compiledPatterns []compiledPattern

// init precompiles all regex patterns for emoji shortcuts
func init() {
	compilePatterns()
}

// compilePatterns creates regex patterns for all emoji shortcuts
func compilePatterns() {
	compiledPatterns = make([]compiledPattern, 0, len(EmojiShortcuts))

	for _, mapping := range EmojiShortcuts {
		// Escape special regex characters in the shortcut
		escaped := regexp.QuoteMeta(mapping.Shortcut)

		// Build pattern that matches the shortcut when:
		// - At the start of string or after whitespace/punctuation
		// - At the end of string or before whitespace/punctuation
		// - Not inside a code block or existing emoji
		pattern := `(?:^|[\s\(\[\{])` + escaped + `(?:$|[\s\)\]\}\.,!?])`

		regex, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		compiledPatterns = append(compiledPatterns, compiledPattern{
			regex: regex,
			emoji: mapping.Emoji,
		})
	}
}

// ReplaceEmojiShortcuts replaces all emoji shortcuts in the given message with their corresponding emoji codes
func ReplaceEmojiShortcuts(message string) string {
	if message == "" {
		return message
	}

	// Skip if message is inside a code block
	if isCodeBlock(message) {
		return message
	}

	result := message

	// Process each mapping
	for i, mapping := range EmojiShortcuts {
		if i >= len(compiledPatterns) {
			break
		}

		pattern := compiledPatterns[i]
		replacement := ":" + pattern.emoji + ":"

		// Replace all occurrences while preserving surrounding whitespace
		result = replacePreservingContext(result, mapping.Shortcut, replacement)
	}

	return result
}

// replacePreservingContext replaces shortcut with emoji while preserving surrounding context
func replacePreservingContext(text, shortcut, replacement string) string {
	escaped := regexp.QuoteMeta(shortcut)

	// Pattern to match shortcut with word boundaries or at string edges
	// This ensures we don't replace shortcuts inside URLs, code, or other words
	patterns := []string{
		`(^|\s)` + escaped + `($|\s)`,           // Surrounded by whitespace or at edges
		`(^|\s)` + escaped + `([.,!?\)\]\}])`,   // Followed by punctuation
		`([\(\[\{])` + escaped + `($|\s)`,       // Preceded by opening bracket
		`([\(\[\{])` + escaped + `([.,!?\)\]\}])`, // Between brackets/punctuation
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return strings.Replace(match, shortcut, replacement, 1)
		})
	}

	return text
}

// isCodeBlock checks if the entire message appears to be a code block
func isCodeBlock(message string) bool {
	trimmed := strings.TrimSpace(message)

	// Check for fenced code blocks
	if strings.HasPrefix(trimmed, "```") && strings.HasSuffix(trimmed, "```") {
		return true
	}

	// Check for inline code that spans the whole message
	if strings.HasPrefix(trimmed, "`") && strings.HasSuffix(trimmed, "`") && strings.Count(trimmed, "`") == 2 {
		return true
	}

	return false
}

// ProcessMessageForEmoji processes a message and replaces emoji shortcuts
// while avoiding code blocks and already-formatted emojis
func ProcessMessageForEmoji(message string) string {
	if message == "" {
		return message
	}

	// Split by code blocks and process only non-code parts
	result := processWithCodeBlocks(message)

	return result
}

// processWithCodeBlocks handles the message while preserving code blocks
func processWithCodeBlocks(message string) string {
	// Match fenced code blocks (```)
	fencedCodeRegex := regexp.MustCompile("(?s)```.*?```")

	// Match inline code (`)
	inlineCodeRegex := regexp.MustCompile("`[^`]+`")

	// Replace code blocks with placeholders
	var codeBlocks []string
	placeholder := "\x00CODE_BLOCK_%d\x00"

	// Extract fenced code blocks
	message = fencedCodeRegex.ReplaceAllStringFunc(message, func(match string) string {
		idx := len(codeBlocks)
		codeBlocks = append(codeBlocks, match)
		return strings.Replace(placeholder, "%d", string(rune('0'+idx)), 1)
	})

	// Extract inline code
	message = inlineCodeRegex.ReplaceAllStringFunc(message, func(match string) string {
		idx := len(codeBlocks)
		codeBlocks = append(codeBlocks, match)
		return strings.Replace(placeholder, "%d", string(rune('0'+idx)), 1)
	})

	// Process emoji shortcuts in the remaining text
	message = ReplaceEmojiShortcuts(message)

	// Restore code blocks
	for i, block := range codeBlocks {
		p := strings.Replace(placeholder, "%d", string(rune('0'+i)), 1)
		message = strings.Replace(message, p, block, 1)
	}

	return message
}
