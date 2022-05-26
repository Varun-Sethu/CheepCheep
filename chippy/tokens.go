package chippy

import (
	"fmt"
	"regexp"
)

// the full set of values a token type can take
type tokenType int

const (
	INSTRUCTION = iota
	VALUE
	COMMA
	NEWLINE
	COMMENT
)

// General token type is a union of all the above types
type Token struct {
	TokenType tokenType
	Value     []byte

	line   int
	column int
}

// SyntaxNode is a node in our "abstract syntax linked list"
type nodeType = uint16

// List of recognised nodes
const (
	Instruction           = 0
	ImmediateValue        = 1
	RegisterValue         = 2
	RegisterRelativeValue = 4
	PCRelativeValue       = 8
	Label                 = 16
	Addr                  = 32
)

type SyntaxNode struct {
	NodeType nodeType
	// Some nodes require 2 values to represent entirely (such as register relative)
	// in that case we store the register in value and the offset in argument
	Value    string
	Argument string

	// Children is a list of children, in the case of an instruction
	// the children form an argument list
	Children []SyntaxNode

	line int
	col  int
}

// cleanSyntaxNode takes a node and cleans the inner contents
// of the node
func cleanSyntaxNode(node SyntaxNode) SyntaxNode {
	if node.NodeType == ImmediateValue {
		node.Value = node.Value[1:]
	} else if node.NodeType == Addr {
		node.Value = node.Value[1 : len(node.Value)-1]
	}
	return node
}

// regular expressions for matching value types]
var labelExpr = `\.\w+`
var numericExpr = `#\d+`
var integerRegister = `r(?:[1-9]|10|11|12|13)`
var register = fmt.Sprintf(`(?:%s)|(?:sp|cmp|zero)`, integerRegister)

// For ease of parsing all these regular expressions return the matched value in the VALUE capturing group
var labelRegex = regexp.MustCompile(fmt.Sprintf(`^(?P<Value>%s)$`, labelExpr))
var addrRegex = regexp.MustCompile(`^\[(?P<Value>\d+)\]$`)
var immediateRegex = regexp.MustCompile(fmt.Sprintf(`^(?P<Value>%s)$`, numericExpr))
var registerRegex = regexp.MustCompile(fmt.Sprintf(`^\$(?P<Value>%s)$`, register))
var registerRelativeRegex = regexp.MustCompile(fmt.Sprintf(`^(?P<Argument>\d+)\+(?P<Value>\$%s)$`, register))
var pcRelativeRegex = regexp.MustCompile(`^#\((?P<Value>\d+)\)$`)

// Theres a few discrete values this could be
// we just verify what it is against the regular expressions at the bottom
var recognisedValueTypes = map[nodeType]*regexp.Regexp{
	Label:                 labelRegex,
	ImmediateValue:        immediateRegex,
	Addr:                  addrRegex,
	RegisterValue:         registerRegex,
	RegisterRelativeValue: registerRelativeRegex,
	PCRelativeValue:       pcRelativeRegex,
}
