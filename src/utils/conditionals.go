package utils

import (
	"strconv"
	"strings"

	"github.com/Ra-hu-l/Transpiler-CP/src/constants"
)

func IfSentence(source string) (output string) {
	output = source
	expression := strings.TrimSpace(LeftBetweenRightmost(source, "if", "{"))
	return "if (" + expression + ") {"
}

func ElseIfSentence(source string) (output string) {
	output = source
	expression := strings.TrimSpace(LeftBetweenRightmost(source, "} else if", "{"))
	return "} else if (" + expression + ") {"
}

func SwitchExpressionVariable() string {
	return constants.SwitchPrefix + strconv.Itoa(constants.SwitchExpressionCounter)
}

func LabelName() string {
	return constants.LabelPrefix + strconv.Itoa(constants.LabelCounter)
}

func Switch(source string) (output string) {
	output = strings.TrimSpace(source)[len("switch "):]
	if strings.HasSuffix(output, "{") {
		output = strings.TrimSpace(output[:len(output)-1])
	}
	constants.SwitchExpressionCounter++
	constants.FirstCase = true
	return "auto&& " + SwitchExpressionVariable() + " = " + output + "; // switch on " + output
}

func Case(source string) (output string) {
	output = source
	s := LeftBetween(output, " ", ":")
	if constants.FirstCase {
		constants.FirstCase = false
		output = "if ("
	} else {
		output = "} else if ("
	}
	output += SwitchExpressionVariable() + " == " + s + ") { // case " + s
	if constants.SwitchLabel != "" {
		output += "\n" + constants.SwitchLabel + ":"
		constants.SwitchLabel = ""
	}
	return output
}
