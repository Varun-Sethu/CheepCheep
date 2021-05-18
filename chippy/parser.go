package chippy

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

/**
parser.go
	- The purpose of this file is to parse the input stream into a series of tokens
	- The file also defines methods for validating syntax and determining the addressing mode of a parameter
	- If the syntax is invalid parsing will fail and an error will be thrown
*/

// TODO: Add support for labelling arbitrary memory locations

// Represents a token within the assembly file
type token struct {
	tokenType int

	tokenData   string
	tokenParams []param
	lineNumber  int
}

const LABEL int = 0
const OPERATION int = 1

// Since we are also dealing with command parameters, we need to be able to distinguish between them
// param defines a representation of a parameter within the file
type param struct {
	rawData string

	// The actual data for this specific parameter
	addressingMode uint16
	addressData    []string
}




// ParseFile takes in a file, its name and spits out a series of tokens IFF it has valid syntax, if not
// ParseFile returns false and a nil slice
func ParseFile(fileName string, fileData string) ([]token, error) {
	fileTokens := tokenizeSource(fileData)
	validFile, syntaxErr := validateSyntaxAndResolveTokens(fileTokens)

	if !validFile {
		return nil, fmt.Errorf("Parse Error:\n%s:%s\n",
			fileName, syntaxErr.Error())
	}

	return fileTokens, nil
}



// tokenizeSource just tokenizes the source file under the assumption that everything within the file is valid
// after a file has been tokenized we verify its correctness within "validateSyntax"
func tokenizeSource(fileData string) []token {
	// Iterate line by line over the source file, each non-empty line represents a token
	scanner := bufio.NewScanner(strings.NewReader(fileData))
	currentLine := 0
	// The set of tokens we are computing
	var generatedTokens []token
	// Regex to remove comments from a line
	removeComments := regexp.MustCompile(`//.*`)

	// Iterate over each line within the scanner
	for scanner.Scan() {
		// Scan in and split up the line by spaces, also increment the line we are looking at
		rawLine := removeComments.ReplaceAllString(scanner.Text(), "")
		line := strings.Fields(rawLine)
		currentLine += 1
		// the token we are trying to compute
		var tokenForLine = token{}

		// for each line process it properly
		switch {
		// The line we are looking at is completely empty
		case len(line) == 0:
			continue
		// The line we are looking at is a label
		case line[0][0] == '.':
			tokenForLine.tokenType = LABEL
			tokenForLine.tokenData = line[0][1:]
			break
		// The line just contains an operation
		default:
			tokenForLine.tokenType = OPERATION
			tokenForLine.tokenData = line[0]
			// All arguments to a command are split by spaces so they will
			// all trail arr[1:], we can use this to compute the set of parameters
			for _, rawParam := range line[1:] {
				tokenForLine.tokenParams = append(tokenForLine.tokenParams, param{
					rawData: rawParam,
				})
			}
			break
		}

		// append the token for this line to the token array
		tokenForLine.lineNumber = currentLine
		generatedTokens = append(generatedTokens, tokenForLine)
	}

	return generatedTokens
}



// validateSyntaxAndResolveTokens takes the tokenized assembly file and determines if it fits the required syntax conventions
// for the .chippy format, additionally it attempts to determine the addressing mode for all the parameters and
// validate the naming for each label
func validateSyntaxAndResolveTokens(tokens []token) (bool, error) {
	// Syntax validation for an assembly program is rather simple, since each token in a program needs to be valid
	// we simply just have to iterate over each token and validate

	// verification for a validVariableLabel
	validLabelName := regexp.MustCompile(`^([a-zA-Z]+)$`)

	for _, token := range tokens {
		// Only operations need to be validated, labels can be left as is
		switch token.tokenType {

		case OPERATION:
			// Determine if the token data exists in the opcode table

			// Before we do anything, just make the current string uppercase

			if opcodeData, ok := OPCODES[strings.ToUpper(token.tokenData)]; ok {
				if len(token.tokenParams) != int(opcodeData[OPTARG]) {
					return false, fmt.Errorf("%d: operation \"%s\" requires %d argument(s), %d were given",
						token.lineNumber, token.tokenData, opcodeData[OPTARG], len(token.tokenParams))
				}

				// It looks like the operation code is valid, we just need to verify that the parameters are of the
				// right addressing mode
				for i, paramEntity := range token.tokenParams {
					isValidAddress, paramEntity := determineAddressingMode(paramEntity)
					token.tokenParams[i] = paramEntity

					if !isValidAddress {
						return false, fmt.Errorf("%d: the adressing mode of \"%s\" is unclear", token.lineNumber, paramEntity.rawData)
					// if we are expecting addresses, just ensure they are of the correct type
					} else {
						// Just assert this parameter is in the valid addressing mode, check mnemonics.go for more
						// information, paramEntity.addressingMode would have been computed in the determineAddressingMode call
						if opcodeData[2 + i] & paramEntity.addressingMode == 0 {
							return false, fmt.Errorf("%d: the arguments for \"%s\" are invalid",
								token.lineNumber, token.tokenData)
						}
					}
				}
			} else {
				return false, fmt.Errorf("%d: unidentified mneumonic \"%s\"", token.lineNumber, token.tokenData)
			}
			break

		case LABEL:
			if !validLabelName.MatchString(token.tokenData) {
				return false, fmt.Errorf("%d: invalid label: \"%s\"", token.lineNumber, token.tokenData)
			}
			break

		}
	}

	return true, nil
}





// determineAddressingMode computes the appropriate addressing mode for for a parameter, it additionally enters
// this computed information into the passed parameter
func determineAddressingMode(parameter param) (bool, param) {
	// This function has a rather messy implementation unfortunately :(
	rawParam := parameter.rawData

	if paramData := ImmediateR.FindStringSubmatch(rawParam); len(paramData) != 0 {
		parameter.addressData = paramData[1:]
		parameter.addressingMode = IMMEDIATE
	} else if paramData := DirectR.FindStringSubmatch(rawParam); len(paramData) != 0 {
		parameter.addressData = paramData[1:]
		parameter.addressingMode = DIRECT
	} else if paramData := IndirectR.FindStringSubmatch(rawParam); len(paramData) != 0 {
		parameter.addressData = paramData[1:]
		parameter.addressingMode = INDIRECT
	} else if paramData := RegDirectR.FindStringSubmatch(rawParam); len(paramData) != 0 {
		parameter.addressData = paramData[1:]
		parameter.addressingMode = REGDIRECT
	} else if paramData := RegIndirectR.FindStringSubmatch(rawParam); len(paramData) != 0 {
		parameter.addressData = paramData[1:]
		parameter.addressingMode = REGINDIRECT
	} else if paramData := IndexedR.FindStringSubmatch(rawParam); len(paramData) != 0 {
		parameter.addressData = paramData[1:]
		parameter.addressingMode = INDEXED
	} else if paramData := IndexScaledR.FindStringSubmatch(rawParam); len(paramData) != 0 {
		parameter.addressData = paramData[1:]
		parameter.addressingMode = INDEXSCALED
	} else {
		return false, parameter
	}

	return true, parameter
}
