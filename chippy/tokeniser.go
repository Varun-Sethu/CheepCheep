package chippy

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

// tokeniseFileStream reads a file stream and converts it to an array
// of tokens, it takes a stream as allocating a buffer in memory is obviously
// too memory intensive
func tokeniseFileStream(stream bufio.Reader) []Token {
	// buffer for storing the current token we are constructing
	tokenBuffer := [100]byte{}
	tokenSize := 0
	tokens := []Token{}

	// for book-keeping the line and column count
	columnCount := 0
	lineCount := 0

	for c, err := stream.ReadByte(); err == nil; c, err = stream.ReadByte() {
		// we trigger the evaluation of the troken buffer if we arrive at a separator
		if c == ' ' || c == '\n' || c == ',' {
			tokenBuffer[tokenSize] = 0

			tokens = append(tokens, resolveToken(tokenBuffer, tokenSize, lineCount, columnCount))
			tokens = append(tokens, terminalToToken(c, lineCount, columnCount))

			clearBuffer(&tokenBuffer)
			tokenSize = 0

			// increment the line count and reset the column count
			if c == '\n' {
				lineCount++
				columnCount = 0
			}
		} else {
			// if we're not at a terminating charachter adding this charachter
			// may mark the start of a comment
			if tokenSize != 0 && tokenBuffer[tokenSize-1] == '/' && c == '/' {
				tokenBuffer[tokenSize-1] = ' '

				tokens = append(tokens, resolveToken(tokenBuffer, tokenSize, lineCount, columnCount))
				tokens = append(tokens, Token{
					TokenType: COMMENT,
					line:      lineCount,
					column:    columnCount,
				})

				tokenSize = 0
			} else {
				// otherwise just append this charachter the buffer
				tokenBuffer[tokenSize] = c
				tokenSize++
			}
		}
		columnCount++
	}

	// resolve the final token and return
	return append(tokens, resolveToken(tokenBuffer, tokenSize, lineCount, columnCount))
}

// converts a single terminal character to a token
func terminalToToken(c byte, line, col int) Token {
	switch c {
	case ',':
		return Token{
			TokenType: COMMA,
			line:      line,
			column:    col,
		}
	case '\n':
		return Token{
			TokenType: NEWLINE,
			line:      line,
			column:    col,
		}
	default:
		return Token{
			TokenType: VALUE,
			Value:     nil,
			line:      line,
			column:    col,
		}
	}
}

// resolveToken converts the contents of a token buffer
// into a singular token
func resolveToken(buffer [100]byte, bufferSize, line, col int) Token {
	if commentRegex.Match(buffer[:]) {
		return Token{
			TokenType: COMMENT,
		}

	} else {
		// Either the buffer is a value type or an instruction type
		// this can be resolved rather easilly by searching in the instrucition table
		// note that the conversion to a string results in an allocation on the heap which isnt
		// ideal, TODO: change this i guess
		value := make([]byte, bufferSize)
		copy(value, buffer[:])
		value = bytes.TrimSpace(value)
		cleanedOpcode := strings.ToUpper(string(value))

		if _, ok := OPCODES[cleanedOpcode]; ok {
			return Token{
				TokenType: INSTRUCTION,
				Value:     []byte(cleanedOpcode),
				line:      line,
				column:    col,
			}
		}

		return Token{
			TokenType: VALUE,
			Value:     value,
			line:      line,
			column:    col,
		}
	}
}

var commentRegex = regexp.MustCompile(`^\\\\\s*$`)
