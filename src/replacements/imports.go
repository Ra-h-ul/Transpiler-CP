package replacements

import (
	"strings"

	"github.com/Ra-hu-l/Transpiler-CP/src/constants"
)

func AddIncludes(source string) (output string) {
	output = source

	includeString := ""
	for k, v := range constants.IncludeMap {
		if strings.Contains(output, k) {
			newInclude := "" + v + ""
			if !strings.Contains(includeString, newInclude) {
				includeString += newInclude
			}

		}
	}
	if constants.CppHasStdFormat {
		//"std::format":                      "format",
		k := "std::format"
		v := "format"
		if strings.Contains(output, k) {
			newInclude := "#include <" + v + ">\n"
			if !strings.Contains(includeString, newInclude) {
				includeString += newInclude
			}
		}
	}
	return includeString + "\n" + output
}
