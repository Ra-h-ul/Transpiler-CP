package utils

import (
	"strings"

	"github.com/Ra-hu-l/Transpiler-CP/src/constants"
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

func FunctionSignature(source string) (output, returntype, name string) {
	if len(strings.TrimSpace(source)) == 0 {
		return source, "", ""
	}
	output = source
	args := FunctionArguments(LeftBetween(output, "(", ")"))
	// Has return values in a parenthesis
	var rets string
	if strings.Contains(output, ") (") {
		// There is a parenthesis with return types in the function signature
		rets = FunctionRetvals(Between(output, ")", "{", false, true))
	} else {
		// There is not a parenthesis with return types in the function signature
		rets = FunctionRetvals(Between(output, ")", "{", true, true))
	}
	if strings.Contains(rets, ",") {
		// Multiple return
		rets = constants.TupleType + "<" + CPPTypes(rets) + ">"
	}
	name = LeftBetween(output, "func ", "(")
	if name == "main" {
		rets = "int"
	}
	if len(strings.TrimSpace(rets)) == 0 {
		rets = "void"
	}
	output = "auto " + name + "(" + args + ") -> " + rets + " {"
	return strings.TrimSpace(output), rets, name
}

// Split arguments. Handles quoting 1 level deep.
func SplitArgs(s string) []string {
	inQuote := false
	inSingleQuote := false
	inPar := false
	inCurly := false
	var args []string
	word := ""
	for _, letter := range s {
		switch letter {
		case '"':
			inQuote = !inQuote
		case '\'':
			inSingleQuote = !inSingleQuote
		}
		if letter == '(' && !inQuote && !inSingleQuote && !inPar && !inCurly {
			inPar = true
		}
		if letter == ')' && !inQuote && !inSingleQuote {
			inPar = false
		}
		if letter == '{' && !inQuote && !inSingleQuote && !inPar && !inCurly {
			inCurly = true
		}
		if letter == '}' && !inQuote && !inSingleQuote {
			inCurly = false
		}
		if letter == ',' && !inQuote && !inSingleQuote && !inPar && !inCurly {
			args = append(args, strings.TrimSpace(word))
			word = ""
		} else {
			word += string(letter)
		}
	}
	args = append(args, strings.TrimSpace(word))
	return args
}

// CPPTypes picks out the types given a list of C++ arguments with name and type
func CPPTypes(args string) string {
	words := strings.Split(LeftBetween(args, "(", ")"), ",")
	var atypes []string
	for _, word := range words {
		elems := strings.Split(strings.TrimSpace(word), " ")
		t := replacements.TypeReplace(elems[0])
		atypes = append(atypes, t)
	}
	return strings.Join(atypes, ", ")
}
