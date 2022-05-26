package chippy

import (
	"bufio"
	"bytes"
	"fmt"
)

// Parse takes a bufio reader (stream) and turns it into a list to tokens
func Parse(stream bufio.Reader) []SyntaxNode {
	tokens := cleanTokens(
		tokeniseFileStream(stream))

	// single pass to clean nodes
	return transform(tokens)
}

// printSyntaxNodes is for debugging purposes, it just prints the syntax nodes in a
// nice format to stdout
func PrintSyntaxNodes(syntaxNodes []SyntaxNode) {
	for _, node := range syntaxNodes {
		fmt.Printf("Instruction: %s\n Children: \n", node.Value)
		for _, child := range node.Children {
			fmt.Printf("	Node Type: %d, Node Value: %s, Node Argument: %s\n", child.NodeType, child.Value, child.Argument)
		}
	}
}

// transform takes a stream of tokes and transforms them into a set
// of "syntax nodes" this part is mostly the identification of register addressing modes
// as well as lables, note that this function assumes the token stream is CLEAN
func transform(tokens []Token) []SyntaxNode {
	nodes := []SyntaxNode{}
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// we want to consume the current token and any subsequent tokens
		// if this is an instruction token we consume the number of arguments
		// the instruction expects
		if token.TokenType == INSTRUCTION {
			// feels like converting to a string is a bad idea :L
			toConsume := OPCODES[string(token.Value)][NUMARGS]
			if i+int(toConsume+1) > len(tokens) {
				panic(fmt.Sprintf(`Error - Unexpected EOF on line %d.`, token.line+1))
			}

			nodes = append(nodes, cleanSyntaxNode(SyntaxNode{
				NodeType: Instruction,
				Value:    string(token.Value),
				Children: consumeChildren(tokens[i+1:], int(toConsume)),
				line:     token.line, col: token.column,
			}))

			i += int(toConsume)
		} else {
			// after cleaning we only reach here if and onlf if this is a value token
			// either this is a label which we admit or we throw an error
			matches, matched := matchNamedGroups(labelRegex, token.Value)
			if !matched {
				panic(fmt.Sprintf(`Error - Unexpected identifier "%s" on line %d column %d.`, token.Value, token.line, token.column))
			}

			nodes = append(nodes, cleanSyntaxNode(SyntaxNode{
				NodeType: Label,
				Value:    string(matches["Value"]),
				line:     token.line, col: token.column,
			}))
		}
	}

	return nodes
}

// consumeChildren consumes the children of a INSTRUCTION node
func consumeChildren(tokens []Token, toConsume int) []SyntaxNode {
	children := []SyntaxNode{}
	for j := 0; j < toConsume; j++ {
		child := tokens[j]
		if child.TokenType == INSTRUCTION {
			panic(fmt.Sprintf(`Error - Unexpected instruction "%s" on line %d cloumn %d.`, child.Value, child.line, child.column))
		}

		childNode, isValid := createValueNode(child)
		if !isValid {
			panic(fmt.Sprintf(`Error - Invalid identifier "%s" on line %d column %d.`, child.Value, child.line, child.column))
		}
		children = append(children, cleanSyntaxNode(childNode))
	}

	return children
}

// createValueNode takes a token and returns a SyntaxNode of the corresponding type
// it also returns a boolean indicating of the value is recognised
func createValueNode(token Token) (SyntaxNode, bool) {
	for nodeType, r := range recognisedValueTypes {
		// try and match and extract the substrings
		matches, matched := matchNamedGroups(r, token.Value)
		if !matched {
			continue
		}

		// construct a new syntax node now :D
		// if theres an argument node return that too
		node := SyntaxNode{
			NodeType: nodeType,
			Value:    string(matches["Value"]),
		}
		if val, exists := matches["Argument"]; exists {
			node.Argument = string(val)
		}
		return node, true
	}
	return SyntaxNode{}, false
}

// cleanTokens takes an array of tokens and cleans them by removing
// empty tokes, comments and newlines
func cleanTokens(tokens []Token) []Token {
	withinComment := false
	return mapOnto(tokens, func(token Token) (Token, bool) {

		shouldAdmit := false
		if token.TokenType == COMMENT {
			withinComment = true
		} else if token.TokenType == NEWLINE {
			withinComment = false
		} else {
			isMeaningful := token.Value != nil &&
				len(bytes.TrimSpace(token.Value)) != 0 &&
				!withinComment

			// note that we discard comma tokens, they dont have any
			// real meaning and only serve as a convenient separator
			shouldAdmit = isMeaningful
		}

		return token, shouldAdmit
	})
}

// mapOnto just maps a function onto a stream of tokesn
func mapOnto[T any](tokens []Token, functor func(Token) (T, bool)) []T {
	mappedValues := []T{}

	for _, token := range tokens {
		computedValue, admit := functor(token)
		if admit {
			mappedValues = append(mappedValues, computedValue)
		}
	}
	return mappedValues
}
