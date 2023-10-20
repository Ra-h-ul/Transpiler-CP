package main

import (
	"fmt"
	"strings"
)

type Tokenx struct {
	Type  string
	Value string
}

type Node struct {
	Type     string
	Value    string
	Children []Node
}

func NewNode(Type string, Value string) *Node {
	return &Node{Type: Type, Value: Value}
}

func AddChild(node *Node, child *Node) {
	node.Children = append(node.Children, *child)
}

func ParseAST(tokens []Token) *Node {
	root := NewNode("Root", "")

	for _, token := range tokens {
		switch token.Type {
		case "package":
			AddChild(root, NewNode("Idenpackagetifier", token.Value))
		case "StringLiteral":
			AddChild(root, NewNode("StringLiteral", token.Value))
		case "NumberLiteral":
			AddChild(root, NewNode("NumberLiteral", token.Value))
		case "func":
			AddChild(root, NewNode("func", token.Value))
		case "main":
			AddChild(root, NewNode("main", token.Value))

		case "OpenParen":
			AddChild(root, NewNode("OpenParen", token.Value))
		case "CloseParen":
			AddChild(root, NewNode("CloseParen", token.Value))
		case "OpenBrace":
			AddChild(root, NewNode("OpenBrace", token.Value))
		case "CloseBrace":
			AddChild(root, NewNode("CloseBrace", token.Value))
		case "Plus":
			AddChild(root, NewNode("Plus", token.Value))
		case "Minus":
			AddChild(root, NewNode("Minus", token.Value))
		case "Asterisk":
			AddChild(root, NewNode("Asterisk", token.Value))
		case "Slash":
			AddChild(root, NewNode("Slash", token.Value))
		case "Equals":
			AddChild(root, NewNode("Equals", token.Value))
		case "blank space":
			AddChild(root, NewNode("blank space", token.Value))
		case "double quotes":
			AddChild(root, NewNode("double quotes", token.Value))
		case "print function":
			AddChild(root, NewNode("print function", token.Value))

		default:
			AddChild(root, NewNode("identifier", token.Value))
		}

	}

	return root
}

func PrintAST(node *Node, indent int) {
	fmt.Printf("%s%s: %s\n", strings.Repeat(" ", indent), node.Type, node.Value)

	for _, child := range node.Children {
		PrintAST(&child, indent+1)
	}
}

func Ast(T []Token) {
	tokens := T

	root := ParseAST(tokens)

	// Print the AST tree.
	PrintAST(root, 0)
}
