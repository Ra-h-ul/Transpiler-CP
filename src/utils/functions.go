package utils

import (
	"strings"

	"github.com/Ra-hu-l/Transpiler-CP/src/replacements"
)

// Name and type is used to keep a variable name and a variable type
type NameAndType struct {
	name string
	typ  string
}

// FunctionArguments transforms the arguments given to a function
func FunctionArguments(source string) string {
	namesAndTypes := make([]NameAndType, 0)
	// First find all names and all types
	currentType := ""
	currentName := ""
	split := strings.Split(source, ",")
	for i := len(split) - 1; i >= 0; i-- {
		nameAndMaybeType := strings.TrimSpace(split[i])
		if strings.Contains(nameAndMaybeType, " ") {
			nameAndType := strings.Split(nameAndMaybeType, " ")
			currentType = replacements.TypeReplace(strings.Join(nameAndType[1:], " "))
			currentName = nameAndType[0]
		} else {
			currentName = nameAndMaybeType
		}
		namesAndTypes = append(namesAndTypes, NameAndType{currentName, currentType})
		//fmt.Println("NAME: " + currentName + ", TYPE: " + currentType)
	}
	cppSignature := ""
	for i := len(namesAndTypes) - 1; i >= 0; i-- {
		//fmt.Println(namesAndTypes[i])
		cppSignature += namesAndTypes[i].typ + " " + namesAndTypes[i].name
		if i > 0 {
			cppSignature += ", "
		}
	}
	return strings.TrimSpace(cppSignature)
}

// FunctionRetvals transforms the return values from a function
func FunctionRetvals(source string) (output string) {
	if len(strings.TrimSpace(source)) == 0 {
		return source
	}
	output = source
	if strings.Contains(output, "(") {
		s := GreedyBetween(output, "(", ")")
		retvals := FunctionArguments(s)
		if strings.Contains(retvals, ",") {
			output = "(" + retvals + ")"
		} else {
			output = retvals
		}
	}
	return strings.TrimSpace(output)
}
