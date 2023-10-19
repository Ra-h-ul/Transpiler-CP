package replacements

import "strings"

func WholeProgramReplace(source string) (output string) {
	output = source

	// TODO: Add these in a smarter way, with more supported types
	replacements := map[string]string{
		" string ":         " " + TypeReplace("string") + " ",
		"(string ":         "(" + TypeReplace("string") + " ",
		"return string":    "return std::to_string",
		"make([]string, ":  "std::vector<" + TypeReplace("string") + "> (",
		"make([]int, ":     "std::vector<" + TypeReplace("int") + "> (",
		"make([]uint, ":    "std::vector<" + TypeReplace("uint") + "> (",
		"make([]float64, ": "std::vector<" + TypeReplace("double") + "> (",
		"make([]float32, ": "std::vector<" + TypeReplace("float") + "> (",
		"-> []string":      "-> std::vector<" + TypeReplace("string") + ">",
		"-> []int":         "-> std::vector<" + TypeReplace("int") + ">",
		"-> []uint":        "-> std::vector<" + TypeReplace("uint") + ">",
		"-> []float64":     "-> std::vector<" + TypeReplace("double") + ">",
		"-> []float32":     "-> std::vector<" + TypeReplace("float") + ">",
		"= nil)":           "= std::nullopt)",
	}
	for k, v := range replacements {
		output = strings.Replace(output, k, v, -1)
	}
	return output
}

func TypeReplace(source string) string {
	// TODO: uintptr, complex64 and complex128
	trimmed := strings.TrimSpace(source)
	// For pointer types, move the star
	if strings.HasPrefix(trimmed, "*") {
		trimmed = trimmed[1:] + "*"
	}
	switch trimmed {
	case "string":
		return "std::string"
	case "float64":
		return "double"
	case "float32":
		return "float"
	case "uint64":
		return "std::uint64_t"
	case "uint32":
		return "std::uint32_t"
	case "uint16":
		return "std::uint16_t"
	case "uint8":
		return "std::uint8_t"
	case "int64":
		return "std::int64_t"
	case "int32":
		return "std::int32_t"
	case "int16":
		return "std::int16_t"
	case "int8":
		return "std::int8_t"
	case "byte":
		return "std::uint8_t"
	case "rune":
		return "std::int32_t"
	case "uint":
		return "unsigned int"
	default:
		if strings.HasPrefix(trimmed, "[]") {
			innerType := trimmed[2:]
			return "std::vector<" + TypeReplace(innerType) + ">"
		}
		return trimmed
	}
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
