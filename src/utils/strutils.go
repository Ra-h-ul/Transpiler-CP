package utils

import (
	"strings"
)

func Lastchar(line string) string {
	if len(line) > 0 {
		return string(line[len(line)-1])
	}
	return ""
}

func Has(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

func SplitAtAndTrim(s string, poss []int) []string {
	l := make([]string, len(poss)+1)
	startpos := 0
	for i, pos := range poss {
		l[i] = strings.TrimSpace(s[startpos:pos])
		startpos = pos + 1
	}
	l[len(poss)] = strings.TrimSpace(s[startpos:])
	return l
}

func CreateStrMethod(varNames []string) string {
	var sb strings.Builder
	sb.WriteString("std::string _str() {\n")
	sb.WriteString("  std::stringstream ss;\n")
	sb.WriteString("  ss << \"{\";\n")
	for i, varName := range varNames {
		if i > 0 {
			sb.WriteString("  ss << \" \";\n")
		}
		sb.WriteString("  _format_output(ss, ")
		sb.WriteString(varName)
		sb.WriteString(");\n")
	}
	sb.WriteString("  ss << \"}\";")
	sb.WriteString("  return ss.str();\n")
	sb.WriteString("}\n")
	return sb.String()
}

func After(keyword, line string) string {
	pos := strings.Index(line, keyword)
	if pos == -1 {
		return line
	}
	return line[pos+len(keyword):]
}
