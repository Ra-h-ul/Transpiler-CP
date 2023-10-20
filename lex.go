package main

import (
	"fmt"
	_ "fmt"
	"regexp"
)

type Token struct {
	Type  string
	Value string
}

type Lexer struct {
	input  string
	pos    int
	tokens []Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) NextToken() *Token {
	if l.pos >= len(l.input) {
		return nil
	}

	// Match regular expressions to identify tokens

	re := regexp.MustCompile(`([[:space:]]+)|([a-zA-Z]+)|([0-9]+)|([+\-*/%=])|(\(|\))|(\{|\})`)
	if match := re.FindStringIndex(l.input[l.pos:]); match != nil {
		l.pos += match[1]
		tokenType := ""
		tokenValue := l.input[l.pos-match[1] : l.pos]

		switch tokenValue {
		case "package main":
			tokenType = "package"
		case "func":
			tokenType = "Func"
		case "main":
			tokenType = "main"
		case "(":
			tokenType = "OpenParen"
		case ")":
			tokenType = "CloseParen"
		case "{":
			tokenType = "OpenBrace"
		case "}":
			tokenType = "CloseBrace"
		case "+":
			tokenType = "Plus"
		case "-":
			tokenType = "Minus"
		case "*":
			tokenType = "Asterisk"
		case "/":
			tokenType = "Slash"
		case "=":
			tokenType = "Equals"
		case " ":
			tokenType = "blank space"
		case "\"":
			tokenType = "double quotes"
		case "fmt.Println":
			tokenType = "print function"
		default:
			tokenType = "identifier"
		}

		return &Token{Type: tokenType, Value: tokenValue}
	}

	// If no regular expression matches, return an error token

	return &Token{Type: "Error", Value: string(l.input[l.pos])}
}

func LEX(input string) {

	var res []Token

	//	i:=0
	l := NewLexer(input)
	t := true
	for token := l.NextToken(); token != nil; token = l.NextToken() {

		if token.Value == "(" {
			t = false
		} else if token.Value == ")" {
			t = true
		}

		if token.Value != " " && t == true {
			//fmt.Printf("%s: %s\n", token.Type, token.Value)

		} else if t == false {
			//fmt.Printf("%s: %s\n", token.Type, token.Value)

		}

		if token.Value == " " || token.Value == "\n" || token.Value == "\t" {
			continue
		}

		//		fmt.Printf("%s \n", token.Value)

		res = append(res, Token{
			Type:  token.Type,
			Value: token.Value,
		})
	}
	fmt.Println(res)

	println("__________")

	Ast(res)

}
