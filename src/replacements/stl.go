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
