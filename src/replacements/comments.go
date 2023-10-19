package replacements

import "strings"

// stripSingleLineComment will strip away trailing single-line comments
func StripSingleLineComment(line string) string {
	commentMarker := "//"
	if strings.Count(line, commentMarker) == 1 {
		p := strings.Index(line, commentMarker)
		return strings.TrimSpace(line[:p])
	}
	return line
}
