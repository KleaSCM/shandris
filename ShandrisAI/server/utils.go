package server

import "regexp"

// Function to remove ANSI escape codes
func CleanANSI(input string) string {
	ansiRegex := regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	return ansiRegex.ReplaceAllString(input, "")
}
