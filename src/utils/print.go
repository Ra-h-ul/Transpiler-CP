package utils

import "strings"

// Will return the transformed string, and a bool if pretty printing may be needed
func PrintStatement(source string) (string, bool) {

	// Pick out and trim all arguments given to the print functon
	args := SplitArgs(GreedyBetween(strings.TrimSpace(source), "(", ")"))

	// Identify the print function
	if !strings.Contains(source, "(") {
		// Not a function call
		return source, false
	}

	fname := strings.TrimSpace(source[:strings.Index(source, "(")])
	//fmt.Println("FNAME", fname)

	// Check if the function call ends with "ln" (println, fmt.Println)
	addNewline := strings.HasSuffix(fname, "ln")
	//fmt.Println("NEWLINE", addNewline)

	// Check if the function call starts with "print" (as opposed to "Print")
	lowercasePrint := strings.HasPrefix(fname, "print")
	//fmt.Println("LOWERCASE PRINT", lowercasePrint)

	// Check if all the arguments are literal strings
	allLiteralStrings := true
	for _, arg := range args {
		if !strings.HasPrefix(arg, "\"") {
			allLiteralStrings = false
		}
	}

	// Check if all the arguments are literal numbers
	allLiteralNumbers := true
	for _, arg := range args {
		if !isNum(arg) {
			allLiteralNumbers = false
		}
	}

	mayNeedPrettyPrint := !allLiteralStrings || !allLiteralNumbers

	// --- enough information gathered, it's time to build the output code ---

	if strings.HasSuffix(fname, "rintf") {
		output := source
		// TODO: Also support fmt.Fprintf, and format %v values differently.
		//       Converting to an iostream expression is one possibility.
		output = strings.Replace(output, "fmt.Printf", "printf", 1)
		output = strings.Replace(output, "fmt.Fprintf", "fprintf", 1)
		output = strings.Replace(output, "fmt.Sprintf", "sprintf", 1)
		if strings.Contains(output, "%v") {
			// TODO: Upgrade this in the future
			output = strings.Replace(output, "%v", "%s", -1)
			//panic("support for %v is not implemented yet")
		}
		return output, mayNeedPrettyPrint
	}

	outputName := "std::cout"
	if lowercasePrint {
		// print and println outputs to stderr
		outputName = "std::cerr"
	}
	//fmt.Println("OUTPUT NAME", outputName)

	// Useful values
	pipe := " << "
	blank := "\" \""
	nl := "std::endl"

	// Silence pipeNewline if the print function does not end with "ln"
	pipeNewline := pipe + nl
	if !addNewline {
		pipeNewline = ""
	}

	// No arguments given?
	if len(args) == 0 {
		// Just output a newline
		if addNewline {
			return outputName + pipeNewline, false
		}
	}

	// Only one argument given?
	if len(args) == 1 {
		if strings.TrimSpace(args[0]) == "" {
			// Just output a newline
			if addNewline {
				return outputName + pipeNewline, false
			}
		}
		if allLiteralStrings || allLiteralNumbers {
			return outputName + pipe + args[0] + pipeNewline, false
		}
		output := "_format_output(" + outputName + ", " + args[0] + ")"
		if addNewline {
			output += ";\n" + outputName + pipeNewline
		}
		return output, true
	}

	// Several arguments given
	//fmt.Println("SEVERAL ARGUMENTS", args)

	// Almost everything should start with "pipe" and almost nothing should end with "pipe"
	output := outputName
	lastIndex := len(args) - 1
	for i, arg := range args {
		//fmt.Println("ARGUMENT", i, arg)
		if strings.HasPrefix(arg, "\"") {
			// Literal string
			output += pipe + arg
		} else if isNum(arg) {
			// Literal number
			output += pipe + arg
		} else {
			if i == 0 {
				output = ""
			} else {
				output += ";\n"
			}
			output += "_format_output(" + outputName + ", " + arg + ");\n" + outputName
		}
		if i < lastIndex {
			output += pipe + blank
		} else {
			output += pipeNewline
		}
	}

	//fmt.Println("GENERATED OUTPUT", output)

	return output, mayNeedPrettyPrint
}
