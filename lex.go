package main

import (
	"unicode"
)

type Tokenkind uint

const (
	syntaxToken     Tokenkind = iota // (  )
	integerToken                     //1 2 3 12
	identifierToken                  //+ -

)

type Token struct {
	value    string
	kind     Tokenkind
	location int
}

/*
remove white space till a character is found and returns
the index of the character
*/
func removespace(source []rune, cursor int) int {
	for cursor < len(source) {
		if unicode.IsSpace(source[cursor]) {
			cursor++
			continue
		}
		break
	}
	return cursor
}

/*
	search for SyntaxTokens , if token is found return next locations

and and the token with its location else return locations and nil
*/
func lexSyntaxToken(source []rune, cursor int) (int, *Token) {
	if source[cursor] == '(' || source[cursor] == ')' {
		return cursor + 1, &Token{
			value:    string([]rune{source[cursor]}),
			kind:     syntaxToken,
			location: cursor,
		}
	}
	return cursor, nil
}

/* append all the integer tokens and if no integer token is found
we will return the origin cursor index */

func lexIntegerToken(source []rune, cursor int) (int, *Token) {
	originalCursor := cursor
	var value []rune
	for cursor < len(source) {
		r := source[cursor]
		if r >= '0' && r <= '9' {
			value = append(value, r)
			cursor++
			continue
		}
		break
	}
	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &Token{
		value:    string(value),
		kind:     integerToken,
		location: originalCursor,
	}
}

/*
 */
func lexIdentifierToken(source []rune, cursor int) (int, *Token) {
	originalCursor := cursor
	var value []rune
	for cursor < len(source) {

		if !unicode.IsSpace(source[cursor]) {
			value = append(value, source[cursor])
			cursor++
			continue
		}
		break
	}
	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &Token{
		value:    string(value),
		kind:     identifierToken,
		location: originalCursor,
	}
}

/*
Takes a slice/arrays of rune type and returns
slice  of items(Tokens) which will be used for
parsing
*/
func Lex(input string) []Token {
	source := []rune(input)
	var Tokens []Token
	var t *Token
	cursor := 0

	for cursor < len(source) {
		// remove white space
		cursor = removespace(source, cursor)
		if cursor == len(source) {
			break
		}
		// check for syntax tokens : ( )
		cursor, t = lexSyntaxToken(source, cursor)
		if t != nil {
			Tokens = append(Tokens, *t)
			continue
		}

		// check for integer tokens : 1 2 31
		cursor, t = lexIntegerToken(source, cursor)
		if t != nil {
			Tokens = append(Tokens, *t)
			continue
		}

		// check for identifier tokens : + -
		cursor, t = lexIdentifierToken(source, cursor)
		if t != nil {
			Tokens = append(Tokens, *t)
			continue
		}

		// exceptions
		println("error occured")
		println(source[cursor])

	}
	return Tokens
}
