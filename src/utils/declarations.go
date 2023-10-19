package utils

import (
	"strconv"
	"strings"

	"github.com/Ra-hu-l/Transpiler-CP/src/constants"
	"github.com/Ra-hu-l/Transpiler-CP/src/replacements"
)

// Return transformed line and the variable names
func VarDeclarations(source string) (string, []string) {

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
					panic("var declaration has mismatching number of variables and values: " + left + " VS " + right)
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
