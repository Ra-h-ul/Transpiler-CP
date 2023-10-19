package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Ra-hu-l/Transpiler-CP/src/constants"
	"github.com/Ra-hu-l/Transpiler-CP/src/replacements"
	"github.com/Ra-hu-l/Transpiler-CP/src/utils"
)

func DeferCall(source string) string {
	trimmed := strings.TrimSpace(utils.After("defer", source))

	// This function handles three possibilities:
	// * defer f()
	// * defer func() { asdf }();
	// * "defer func() {" and then later "}()"

	if strings.HasPrefix(trimmed, "func() {") && strings.HasSuffix(trimmed, "}()") {
		// Anonymous function, on one line
		constants.DeferCounter++
		trimmed = strings.TrimSpace(utils.LeftBetween(trimmed, "func() {", "}()"))
		// TODO: let go2cpp() return pure source code + includes to place at the top, not just one large string
		return "// " + trimmed + "\nstd::shared_ptr<void> _defer" + strconv.Itoa(constants.DeferCounter) + "(nullptr, [](...) { " + go2cpp(trimmed) + "; });"
	} else if trimmed == "func() {" {
		// Anonymous function, on multiple lines
		constants.DeferCounter++
		constants.UnfinishedDeferFunction = true // output "});" later on, when "}()" is encountered in the Go code
		return "// " + trimmed + "\nstd::shared_ptr<void> _defer" + strconv.Itoa(constants.DeferCounter) + "(nullptr, [](...) { "
	} else {
		// Assume a regular function call
		return "// " + trimmed + "\nstd::shared_ptr<void> _defer" + strconv.Itoa(constants.DeferCounter) + "(nullptr, [](...) { " + trimmed + "; });"
	}
}

func IfSentence(source string) (output string) {
	output = source
	expression := strings.TrimSpace(utils.LeftBetweenRightmost(source, "if", "{"))
	return "if (" + expression + ") {"
}

func ElseIfSentence(source string) (output string) {
	output = source
	expression := strings.TrimSpace(utils.LeftBetweenRightmost(source, "} else if", "{"))
	return "} else if (" + expression + ") {"
}

func ForLoop(source string, encounteredHashMaps []string) string {
	expression := strings.TrimSpace(utils.LeftBetween(source, "for", "{"))
	if expression == "" {
		// endless loop
		return "for (;;) {"
	}
	// for range, with no comma
	if strings.Count(expression, ",") == 0 && strings.Contains(expression, "range") {
		fields := strings.Split(expression, " ")
		varName := fields[0]
		listName := fields[len(fields)-1]

		// for i := range l {
		// -->
		// for (auto i = 0; i < std::size(l); i++) {

		hashMapName := listName
		if utils.Has(encounteredHashMaps, hashMapName) {
			// looping over the key of a hash map, not over the index of a list
			return "for (const auto & [" + varName + ", " + varName + "__" + "] : " + hashMapName + ") {"
		} else if varName == "_" {
			return "for (const auto & [" + varName + "__" + ", " + varName + "___" + "] : " + hashMapName + ") {"
		} else {
			// looping over the index of a list
			return "for (std::size_t " + varName + " = 0; " + varName + " < std::size(" + listName + "); " + varName + "++) {"
		}
	}
	// for range, over index and element, or key and value
	if strings.Count(expression, ",") == 1 && strings.Contains(expression, "range") && strings.Contains(expression, ":=") {
		fields := strings.Split(expression, ":=")
		varnames := strings.Split(fields[0], ",")

		indexvar := varnames[0]
		elemvar := varnames[1]

		fields = strings.Split(expression, " ")
		listName := fields[len(fields)-1]
		hashMapName := listName

		if utils.Has(encounteredHashMaps, hashMapName) {
			if indexvar == "_" {
				// looping over the values of a hash map
				hashMapHashKey := hashMapName + constants.HashMapSuffix
				return "for (const auto & " + hashMapHashKey + " : " + hashMapName + ") {" + "\n" + "auto " + elemvar + " = " + hashMapHashKey + ".second"
			}
			// for k, v := range m
			keyvar := indexvar
			//hashMapHashKey := keyvar + hashMapSuffix + keysSuffix
			return "for (const auto & [" + keyvar + ", " + elemvar + "] : " + hashMapName + ") {"
			//return "for (auto " + hashMapHashKey + " : " + hashMapName + keysSuffix + ") {" + "\n" + "auto " + keyvar + " = " + hashMapHashKey + ".second;\nauto " + elemvar + " = " + hashMapName + ".at(" + hashMapHashKey + ".first)"
		}

		if indexvar == "_" {
			return "for (auto " + elemvar + " : " + listName + ") {"
		}
		return "for (std::size_t " + indexvar + " = 0; " + indexvar + " < std::size(" + listName + "); " + indexvar + "++) {" + "\n" + "auto " + elemvar + " = " + listName + "[" + indexvar + "]"
	}
	// not "for" + "range"
	if strings.Contains(expression, ":=") {
		if strings.HasPrefix(expression, "_,") && strings.Contains(expression, "range") {
			// For each, no index
			varname := utils.LeftBetween(expression, ",", ":")
			fields := strings.SplitN(expression, "range ", 2)
			listname := fields[1]
			// C++11 and later for each loop
			expression = "auto &" + varname + " : " + listname
		} else {
			expression = "auto " + strings.Replace(expression, ":=", "=", 1)
		}
	}
	return "for (" + expression + ") {"
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
	s := utils.LeftBetween(output, " ", ":")
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

// Return transformed line and the variable names
func VarDeclarations(source string) (string, []string) {

	// TODO: This is an ugly function. Refactor.

	if strings.Contains(source, "=") {
		parts := strings.SplitN(strings.TrimSpace(source), "=", 2)
		left := parts[0]
		right := strings.TrimSpace(parts[1])
		fields := strings.Split(strings.TrimSpace(left), " ")
		if fields[0] == "var" {
			fields = fields[1:]
		}
		if len(fields) == 2 {
			return replacements.TypeReplace(fields[1]) + " " + fields[0] + " = " + right, []string{fields[0]}
		} else if len(fields) > 2 {
			if strings.Contains(source, ",") {
				leftFields := strings.Fields(left)
				if leftFields[0] == "var" {
					leftFields = leftFields[1:]
				}
				rightFields := strings.Fields(right)
				if len(leftFields)-1 != len(rightFields) {
					panic("var declaration utils.Has mismatching number of variables and values: " + left + " VS " + right)
				}
				lastIndex := len(leftFields) - 1
				varType := leftFields[lastIndex]
				var sb strings.Builder
				var varNames []string

				for i, varName := range leftFields[:lastIndex] {
					if i > 0 {
						sb.WriteString(";")
					}
					if strings.HasSuffix(varName, ",") {
						varName = varName[:len(varName)-1]
					}
					varValue := rightFields[i]
					if strings.HasSuffix(varValue, ",") {
						varValue = varValue[:len(varValue)-1]
					}
					sb.WriteString(replacements.TypeReplace(varType) + " " + varName + " = " + varValue)
					varNames = append(varNames, varName)
				}
				return sb.String(), varNames
			}
			return replacements.TypeReplace(fields[1]) + " " + fields[0] + " " + strings.Join(fields[2:], " ") + " = " + right, []string{fields[0]}
		}
		leftFields := strings.Fields(left)
		if leftFields[0] == "var" {
			leftFields = leftFields[1:]
		}
		if len(leftFields) > 1 {
			panic("unsupported var declaration: " + source)
		}

		varName := leftFields[0]
		varType := right
		varValue := ""
		if strings.Contains(source, "=") {
			varType = "auto"
			lastIndex := len(leftFields) - 1
			if len(leftFields) > 1 && !strings.Contains(leftFields[lastIndex], ",") {
				varType = leftFields[lastIndex]
				//leftFields = leftFields[:lastIndex-1]
			}
			varValue = right
		}

		withBracket := false
		if strings.HasSuffix(right, "{") {
			varType = strings.TrimSpace(right[:len(right)-1])
			withBracket = true
		}

		varValue = strings.TrimPrefix(varValue, varType)

		s := replacements.TypeReplace(varType) + " " + varName + " = " + varValue
		if withBracket {
			if !strings.HasSuffix(s, "{") {
				s += "{"
			}
		} else {
			if !strings.Contains(source, "`") {
				// Only add a semicolon if it's not a multiline string and not an opening bracket
				s += ";"
			}
		}
		return s, []string{varName}
	}
	fields := strings.Fields(source)
	if fields[0] == "var" {
		fields = fields[1:]
	}
	if len(fields) == 2 {
		return replacements.TypeReplace(fields[1]) + " " + fields[0], []string{fields[0]}
	}
	if strings.Contains(source, ",") {
		// Comma separated variable names, with one common variable type,
		// and no value assignment
		lastIndex := len(fields) - 1
		varType := fields[lastIndex]
		var sb strings.Builder
		var varNames []string

		for i, varName := range fields[:lastIndex] {
			if i > 0 {
				sb.WriteString(";")
			}
			if strings.HasSuffix(varName, ",") {
				varName = varName[:len(varName)-1]
			}
			sb.WriteString(replacements.TypeReplace(varType) + " " + varName)
			varNames = append(varNames, varName)
		}

		return sb.String(), varNames
	}

	// Unrecognized
	panic("Unrecognized var declaration: " + source)

}

// TypeDeclaration returns a transformed string (from Go to C++),
// and a bool if a struct is opened (with {).
func TypeDeclaration(source string) (string, bool) {
	fields := strings.Split(strings.TrimSpace(source), " ")
	if fields[0] == "type" {
		fields = fields[1:]
	}
	left := strings.TrimSpace(fields[0])
	right := strings.TrimSpace(fields[1])
	words := strings.Split(left, " ")
	if len(fields) == 2 {
		// Type alias
		return "using " + left + " = " + replacements.TypeReplace(right), false
	} else if len(words) == 2 {
		// Type alias
		return "using " + words[1] + " " + words[0] + " = " + replacements.TypeReplace(right), false
	} else if strings.Contains(right, "struct") {
		// type Vec3 struct {
		// to
		// class Vec3 { public:
		// also the closing bracket must end with a semicolon
		return "class " + left + " { public:", true
	} else if len(words) == 1 {
		// Type alias
		return "using " + left + " = " + replacements.TypeReplace(right), false
	}
	// Unrecognized
	panic("Unrecognized type declaration: " + source)
}

func ConstDeclaration(source string) (output string) {
	output = source
	fields := strings.SplitN(source, "=", 2)
	if len(fields) == 0 {
		panic("no fields in const declaration")
	} else if len(fields) == 1 {
		// This happens if there is only a constant name, with no value assigned
		// Only simple iota incrementation is supported so far (not .. << ..)
		constants.IotaNumber++
		return "const auto " + strings.TrimSpace(fields[0]) + " = " + strconv.Itoa(constants.IotaNumber)
	}
	left := strings.TrimSpace(fields[0])
	right := strings.TrimSpace(fields[1])
	words := strings.Split(left, " ")
	if right == "iota" {
		constants.IotaNumber = 0
		right = strconv.Itoa(constants.IotaNumber)
	}
	if len(words) == 1 {
		// No type
		return "const auto " + left + " = " + right
	} else if len(words) == 2 {
		if words[0] == "const" {
			return "const auto " + words[1] + " = " + right
		}
		return "const " + replacements.TypeReplace(words[1]) + " " + words[0] + " = " + right
	}
	// Unrecognized
	panic("go2cpp: unrecognized const expression: " + source)
}

// HashElements transforms the contents of a map in Go to the contents of an unordered_map in C++
// keyType is the type of the key, in C++, for instance "std::string"
// if keyForBoth is true, a hash(key)->key map is created,
// if not, a hash(key)->value map is created.
// This will not work for multiline hash map initializations.
// TODO: Handle keys and values that look like this: "\": \"" (containing quotes, a colon and a space)
func HashElements(source, keyType string, keyForBoth bool) string {
	// Check if the given source line contains either a separating or a trailing comma
	if !strings.Contains(source, ",") {
		return source
	}
	// Check if there is only one pair
	if strings.Count(source, ": ") == 1 {
		pairElements := strings.SplitN(source, ": ", 2)
		if len(pairElements) != 2 {
			panic("This should be two elements, separated by a colon and a space " + source)
		}
		return "{ " + strings.TrimSpace(pairElements[0]) + ", " + strings.TrimSpace(pairElements[1]) + " }, "
	}
	// Multiple pairs
	pairs := strings.Split(source, ",")
	output := "{"
	first := true
	for _, pair := range pairs {
		if !first {
			output += ","
		} else {
			first = false
		}
		pairElements := strings.SplitN(pair, ": ", 2)
		//fmt.Println("HASH ELEMENTS", source)
		//fmt.Println("HASH ELEMENTS", pairs)
		if len(pairElements) != 2 {
			panic("This should be two elements, separated by a colon and a space: " + pair)
		}
		output += "{ " + strings.TrimSpace(pairElements[0]) + ", " + strings.TrimSpace(pairElements[1]) + " }"
	}

	return output + "}"
}

func go2cpp(source string) string {
	functionVarMap := map[string]string{} // variable names encountered in the function so far, and their corresponding smart names
	inMultilineString := false
	debugOutput := false
	lines := []string{}
	currentReturnType := ""
	currentFunctionName := ""
	inImport := false
	inVar := false
	inType := false
	inConst := false
	inHashMap := false
	hashKeyType := ""
	curlyCount := 0

	encounteredHashMaps := []string{}

	encounteredStructNames := []string{}
	inStruct := false
	// usePrettyPrint := false
	closingBracketNeedsASemicolon := false
	for _, line := range strings.Split(source, "\n") {

		if debugOutput {
			fmt.Fprintf(os.Stderr, "%s\n", line)
		}

		newLine := line

		trimmedLine := replacements.StripSingleLineComment(strings.TrimSpace(line))

		if strings.HasPrefix(trimmedLine, "//") {
			lines = append(lines, trimmedLine)
			continue
		}
		if strings.HasSuffix(trimmedLine, ";") {
			trimmedLine = trimmedLine[:len(trimmedLine)-1]
		}
		if len(trimmedLine) == 0 {
			lines = append(lines, newLine)
			continue
		}
		// Keep track of how deep we are into curly brackets
		curlyCount += (strings.Count(trimmedLine, "{") - strings.Count(trimmedLine, "}"))
		if inImport && strings.Contains(trimmedLine, ")") {
			inImport = false
			continue
		} else if inImport {
			continue
		} else if inVar && strings.Contains(trimmedLine, ")") {
			inVar = false
			continue
		} else if inType && strings.Contains(trimmedLine, ")") {
			inType = false
			continue
		} else if inConst && strings.Contains(trimmedLine, ")") {
			inConst = false
			continue
		} else if inHashMap && trimmedLine == "}" {
			inHashMap = false
			newLine = trimmedLine + ";"
		} else if inVar || (inStruct && trimmedLine != "}") {
			var varNames []string
			newLine, varNames = VarDeclarations(trimmedLine)
			if strings.HasSuffix(newLine, "{") {
				closingBracketNeedsASemicolon = true
			}

			if inStruct {
				// Gathering variable names from this struct
				encounteredStructNames = append(encounteredStructNames, varNames...)
			}
		} else if inType {
			prevInStruct := inStruct
			newLine, inStruct = TypeDeclaration(trimmedLine)
			if !prevInStruct && inStruct {
				// Entering struct, reset the slice that is used to gather variable names
				encounteredStructNames = []string{}
			}
		} else if inConst {
			newLine = ConstDeclaration(line)
		} else if inHashMap && !inMultilineString {
			newLine = HashElements(trimmedLine, hashKeyType, false)
		} else if strings.HasPrefix(trimmedLine, "func") {
			functionVarMap = map[string]string{}
			newLine, currentReturnType, currentFunctionName = utils.FunctionSignature(trimmedLine)
		} else if strings.HasPrefix(trimmedLine, "for") {
			newLine = ForLoop(line, encounteredHashMaps)
		} else if strings.HasPrefix(trimmedLine, "switch") {
			newLine = Switch(line)
		} else if strings.HasPrefix(trimmedLine, "case") {
			newLine = Case(line)
		} else if strings.HasPrefix(trimmedLine, "return") {
			if strings.HasPrefix(currentReturnType, constants.TupleType) {
				elems := strings.SplitN(newLine, "return ", 2)
				newLine = "return " + currentReturnType + "{" + elems[1] + "};"
				//} else {
				// Just use the standard tuple
			}
		} else if strings.HasPrefix(trimmedLine, "fmt.Print") || strings.HasPrefix(trimmedLine, "print") {
			// _ is if "pretty print" functionality may be needed, for non-literal strings and numbers
			// var pp bool
			newLine, _ = utils.PrintStatement(trimmedLine)
			// if pp {
			// 	usePrettyPrint = true
			// }
		} else if strings.Contains(trimmedLine, "=") && !strings.HasPrefix(trimmedLine, "var ") && !strings.HasPrefix(trimmedLine, "if ") && !strings.HasPrefix(trimmedLine, "const ") && !strings.HasPrefix(trimmedLine, "type ") {
			elem := strings.SplitN(trimmedLine, "=", 2)
			left := strings.TrimSpace(elem[0])
			declarationAssignment := false
			if strings.HasSuffix(left, ":") {
				declarationAssignment = true
				left = left[:len(left)-1]
			}
			right := strings.TrimSpace(elem[1])
			if strings.HasPrefix(right, "&") && strings.Contains(right, "{") && strings.Contains(right, "}") {
				right = "new " + right[1:]
			}
			if strings.Contains(left, ",") {
				if strings.Contains(left, "_") {
					varNames := strings.Split(left, ",")
					for _, name := range varNames {
						name = strings.TrimSpace(name)
						if value, found := functionVarMap[name]; found {
							// The key already exists, update the value
							if value == name {
								// Add a "0"
								functionVarMap[name] = value + "0"
							} else {
								// Increase the number in the current value by 1
								num := utils.TrailingNumber(value)
								num++
								functionVarMap[name] = name + strconv.Itoa(num)
							}
						} else {
							// The key does not exist, just add the name as it is
							functionVarMap[name] = name
						}
					}
					// The "varInFunction" map should now have been updated correctly, so use that

					useVarNames := []string{}
					for _, name := range varNames {
						name = strings.TrimSpace(name)
						useVarNames = append(useVarNames, functionVarMap[name])
					}
					newLine = "auto [" + strings.Join(useVarNames, ", ") + "] = " + right
					//fmt.Println("function var map:", functionVarMap)
					//panic(newLine)
				} else {
					newLine = "auto [" + left + "] = " + right
				}
			} else if declarationAssignment {
				if strings.HasPrefix(right, "[]") {
					if !strings.Contains(right, "{") {
						fmt.Fprintln(os.Stderr, "UNRECOGNIZED LINE: "+trimmedLine)
						//newLine = line

					}
					theType := replacements.TypeReplace(utils.LeftBetween(right, "]", "{"))
					fields := strings.SplitN(right, "{", 2)
					newLine = theType + " " + strings.TrimSpace(left) + "[] {" + fields[1]
				} else if strings.HasPrefix(right, "map[") {
					hashName := strings.TrimSpace(left)
					encounteredHashMaps = append(encounteredHashMaps, hashName)

					keyType := replacements.TypeReplace(utils.LeftBetween(right, "map[", "]"))
					valueType := replacements.TypeReplace(utils.LeftBetween(right, "]", "{"))

					closingBracket := strings.HasSuffix(strings.TrimSpace(right), "}")
					if !closingBracket {
						inHashMap = true
						hashKeyType = keyType
						newLine = "std::unordered_map<" + keyType + ", " + valueType + "> " + hashName + " {"
					} else {
						elements := utils.LeftBetween(right, "{", "}")
						newLine = "std::unordered_map<" + keyType + ", " + valueType + "> " + hashName + " " + HashElements(elements, keyType, false)
					}
				} else {
					varName := strings.TrimSpace(left)
					if value, found := functionVarMap[varName]; found {
						varName = value
					}

					for k, v := range functionVarMap {
						right = strings.Replace(right, k, v, -1)
					}

					newLine = "auto " + varName + " = " + strings.TrimSpace(right)
				}
			} else {
				newLine = left + " = " + right
			}
		} else if strings.HasPrefix(trimmedLine, "package ") {
			continue
		} else if strings.HasPrefix(trimmedLine, "import") {
			if strings.Contains(trimmedLine, "(") {
				inImport = true
			}
			if strings.Contains(trimmedLine, ")") {
				inImport = false
			}
			continue
		} else if strings.HasPrefix(trimmedLine, "defer ") {
			newLine = DeferCall(line)
		} else if strings.HasPrefix(trimmedLine, "if ") {
			newLine = IfSentence(line)
			// TODO: Short variable names utils.Has the potential to ruin if expressions this way, do a smarter replacement
			// TODO: Also do this for for loops, switches and other cases where this makes sense
			for k, v := range functionVarMap {
				newLine = strings.Replace(newLine, k, v, -1)
			}
		} else if strings.HasPrefix(trimmedLine, "} else if ") {
			newLine = ElseIfSentence(line)
		} else if trimmedLine == "var (" {
			inVar = true
			continue
		} else if trimmedLine == "type (" {
			inType = true
			continue
		} else if trimmedLine == "const (" {
			inConst = true
			continue
		} else if strings.HasPrefix(trimmedLine, "var ") {
			// Ignore variable name since it's not in a struct
			newLine, _ = VarDeclarations(line)
			if strings.HasSuffix(newLine, "{") {
				closingBracketNeedsASemicolon = true
			}
		} else if strings.HasPrefix(trimmedLine, "type ") {
			newLine, inStruct = TypeDeclaration(trimmedLine)
		} else if strings.HasPrefix(trimmedLine, "const ") {
			newLine = ConstDeclaration(trimmedLine)
		} else if trimmedLine == "fallthrough" {
			newLine = "goto " + LabelName() + "; // fallthrough"
			constants.SwitchLabel = LabelName()
			constants.LabelCounter++
		} else if constants.UnfinishedDeferFunction && trimmedLine == "}()" {
			constants.UnfinishedDeferFunction = false
			newLine = "});"
		} else if trimmedLine == "default:" {
			newLine = "} else { // default case"
			if constants.SwitchLabel != "" {
				newLine += "\n" + constants.SwitchLabel + ":"
				constants.SwitchLabel = ""
			}
		}

		if constants.CppHasStdFormat {
			// Special case for fmt.Sprintf -> std::format
			if strings.Contains(newLine, "fmt.Sprintf(") && strings.Contains(newLine, "%v") {
				newLine = strings.Replace(strings.Replace(newLine, "%v", "{}", -1), "fmt.Sprintf(", "std::format(", -1)
			}
		}

		if currentFunctionName == "main" && trimmedLine == "}" && curlyCount == 0 { // curlyCount utils.Has already been decreased for this line
			newLine = strings.Replace(trimmedLine, "}", "return 0;\n}", 1)
		}

		if strings.HasSuffix(trimmedLine, "}") {
			// If the struct is being closed, add a semicolon
			if inStruct {
				// Create a _str() method for this struct
				newLine = utils.CreateStrMethod(encounteredStructNames) + newLine + ";"

				inStruct = false
			} else if closingBracketNeedsASemicolon {
				newLine += ";"
				closingBracketNeedsASemicolon = false
			}
			newLine += "\n"
		}
		if (!strings.HasSuffix(newLine, ";") && !utils.Has(constants.Endings, utils.Lastchar(trimmedLine)) || strings.Contains(trimmedLine, "=")) && !strings.HasPrefix(trimmedLine, "//") && (!utils.Has(constants.Endings, utils.Lastchar(newLine)) && !strings.Contains(newLine, "//")) {
			if !inMultilineString {
				newLine += ";"
			}
		}

		// multiline strings
		for strings.Contains(newLine, "`") {
			if !inMultilineString {
				if strings.HasSuffix(newLine, ";") {
					newLine = newLine[:len(newLine)-1]
				}
				newLine = strings.Replace(newLine, "`", "R\"(", 1)
				inMultilineString = true
			} else {
				newLine = strings.Replace(newLine, "`", ")\"", 1)
				//if !strings.HasSuffix(newLine, ",") {
				newLine += ";"
				//}
				inMultilineString = false
			}
		}

		lines = append(lines, newLine)
	}
	output := strings.Join(lines, "\n")

	// The order matters
	output = replacements.WholeProgramReplace(output)

	// The order matters
	// output = AddFunctions(output, usePrettyPrint, len(encounteredStructNames) > 0)
	// output = AddIncludes(output)

	return output
}

func main() {

	debug := false
	compile := true
	clangFormat := true

	inputFilename := "./go_test_code/test2.txt"
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" {
			fmt.Println(constants.VersionString)
			return
		} else if os.Args[1] == "--help" {
			fmt.Println("supported arguments:")
			fmt.Println(" a .go file as the first argument")
			fmt.Println("supported options:")
			fmt.Println(" -o : Format with clang format")
			fmt.Println(" -O : Don't format with clang format")
			return
		}
		inputFilename = os.Args[1]
	}
	if len(os.Args) > 2 {
		if os.Args[2] == "-o" {
			clangFormat = true
		} else if os.Args[2] == "-O" {
			clangFormat = false
		} else if os.Args[2] != "-o" {
			log.Fatal("The second argument must be -o (format sources with clang-format) or -O (don't format sources with clang-format)")
		}
	}

	var sourceData []byte
	var err error
	if inputFilename != "" {
		sourceData, err = ioutil.ReadFile(inputFilename)
	} else {
		sourceData, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		fmt.Println(go2cpp(string(sourceData)))
		return
	}

	cppSource := "output.txt"
	if clangFormat {
		cmd := exec.Command("clang-format", "-style={BasedOnStyle: Webkit, ColumnLimit: 99}")
		cmd.Stdin = strings.NewReader(go2cpp(string(sourceData)))
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			// log.Println("clang-format is not available, the output will look ugly!")
			cppSource = go2cpp(string(sourceData))
		} else {
			cppSource = out.String()
		}
	} else {
		cppSource = go2cpp(string(sourceData))
	}

	if !compile {
		fmt.Println(cppSource)
		return
	}

	tempFile, err := ioutil.TempFile("", "go2cpp*")
	if err != nil {
		log.Fatal(err)
	}
	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName)

	// Compile the string in cppSource
	cpp := "g++"
	if cppenv := os.Getenv("CXX"); cppenv != "" {
		cpp = cppenv
	}
	cmd2 := exec.Command(cpp, "-x", "c++", "-std=c++2a", "-O2", "-pipe", "-fPIC", "-Wfatal-errors", "-fpermissive", "-s", "-o", tempFileName, "-")
	cmd2.Stdin = strings.NewReader(cppSource)
	var compiled bytes.Buffer
	var errors bytes.Buffer
	cmd2.Stdout = &compiled
	cmd2.Stderr = &errors
	err = cmd2.Run()
	if err != nil {
		//fmt.Println("Failed to compile this with g++:")
		fmt.Println("In")
		// fmt.Println(err)

		cppSource = `

template <typename T>
	void _format_output(std::ostream& out, const T& str) 
	{	
		out << str;
	}` + cppSource

		cppSource = `#include <bits/stdc++.h>` + cppSource
		cppSource = formatting(cppSource)
		fmt.Println(cppSource)
		outputFileName := "cpp_test_output/" + strings.Split(inputFilename, "/")[2] + ".cpp"
		writeFile(outputFileName, cppSource)

		// fmt.Println("Errors:")
		// fmt.Println(errors.String())
		//log.Fatal(err)
	}
	compiledBytes, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		log.Fatal(err)
	}
	outputFilename := ""
	if len(os.Args) > 3 {
		outputFilename = os.Args[3]
	}
	if outputFilename != "" {
		err = ioutil.WriteFile(outputFilename, compiledBytes, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(cppSource)
	}

}

func containsSubstring(str, substr string) bool {
	for i := 0; i < len(str)-len(substr)+1; i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func writeFile(filename, data string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func formatting(data string) string {
	// temporary
	ans := data
	ans = strings.Replace(data, "Math.Sqrt", "sqrt", -1)
	ans = strings.Replace(ans, "float64", "double", -1)
	ans = strings.Replace(ans, "+ =", "+=", -1)
	ans = strings.Replace(ans, "- =", "-=", -1)
	ans = strings.Replace(ans, "* =", "*=", -1)
	ans = strings.Replace(ans, "/ =", "/=", -1)

	return ans

}
