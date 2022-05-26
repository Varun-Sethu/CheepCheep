package chippy

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Compilation is a rather easy process, we simply take our list of syntax nodes
// and evaluate each instruction + its operands individually, however before doing this
// we need to construct a relocation table for deadling with labels, the function
// returns a buffered reader that computes the bytecode on the fly
func Compile(nodes []SyntaxNode) *bufio.Reader {
	// validate and compute a relocation table for labels
	relocationTable := computeRelocationTable(nodes)
	validateInstructionOperands(nodes, relocationTable)

	// Finally return a reader for our compiled bytecodes
	return bufio.NewReader(
		&CompiledOps{
			nodes:             nodes,
			relocationTable:   relocationTable,
			nodePosition:      0,
			instructionBuffer: [4]byte{},
			bufferSize:        0,
		},
	)
}

// CompiledOps is an implementation of io.Reader
// implementing io.Reader allows us to write directly to a file
// without having to create a buffer in memory and then copying that over to a file
type CompiledOps struct {
	nodes             []SyntaxNode
	relocationTable   map[string]uint16
	nodePosition      int
	instructionBuffer [4]byte
	bufferSize        int
}

func (c *CompiledOps) Read(p []byte) (int, error) {
	// first empty the instruction buffer, there should be at most
	// 4 values in the instruction buffer
	pn := 0
	for i := 4 - c.bufferSize; i < 4; i++ {
		p[pn] = c.instructionBuffer[i]
		pn++
	}

	// now translate the current instruction
	for pn < len(p) && c.nodePosition < len(c.nodes) {
		if c.nodes[c.nodePosition].NodeType == Label {
			c.nodePosition++
			continue
		}

		// translate to the buffer and write it out
		c.bufferSize = 4
		translateToBuffer(c.nodes[c.nodePosition], c.relocationTable, &c.instructionBuffer)
		// write the buffer out
		for i := 0; i < 4 && pn < len(p); i++ {
			p[pn] = c.instructionBuffer[i]
			c.bufferSize--
			pn++
		}
		c.nodePosition++
	}

	if pn == 0 {
		return 0, io.EOF
	}
	return pn, nil
}

// translateToBuffer writes an instruction to a 4 byte buffer
// write out the opcode
// we need to identify the addressing mode as well
// finally we need to pack the operands into the final 3 bytes, the general rules are:
// 	- register values take 4 bits
// 	- normal values take up 16 bits
//	- register relative values take 20 bits
//  - PC relative values take 16 bits
//	- addresses take 16 bits
// eg: add $r1, $r2, 3 would encode to:
//	- [0001][0010][0000 0000][0000 0011]
// eg: ldr $r1, 3+$r3 would encode to
//	- [0001][0011][0000 0000][0000 0011]
func translateToBuffer(node SyntaxNode, relocationTable map[string]uint16, output *[4]byte) {
	translationTable := OPCODES[node.Value]

	// pack in the instruction first
	var instruction uint32 = (uint32(translationTable[0]) << 26) | (uint32(resolveAddressingMode(node, translationTable)) << 24)
	var consumedBits uint32 = 8

	// now we pack the individual operands into the instruction, this is mostly just a bunch of
	// case work
	for _, token := range node.Children {
		if consumedBits == 32 {
			panic("Compilation Failure - Internal error: cannot pack instructions into 4 bytes.")
		}

		if token.NodeType == RegisterValue {
			consumeRegister(token.Value, &instruction, &consumedBits)

		} else if token.NodeType == ImmediateValue || token.NodeType == Addr || token.NodeType == Instruction || token.NodeType == PCRelativeValue {
			consumeTwoBytes(token.Value, &instruction, &consumedBits)

		} else if token.NodeType == RegisterRelativeValue {
			consumeRegister(token.Value, &instruction, &consumedBits)
			consumeTwoBytes(token.Argument, &instruction, &consumedBits)

		} else if token.NodeType == Label {
			// translate the label and write it out, i hate that im doing this
			var translation uint16 = relocationTable[token.Value]
			consumeTwoBytes(strconv.Itoa(int(translation)), &instruction, &consumedBits)
		}
	}

	// finally write out the instructio to the buffer
	for i := 3; i >= 0; i-- {
		(*output)[i] = byte((instruction >> (8 * uint32(i))) & 255)
	}
}

// consume function adds a register value to the instruction so far
func consumeRegister(value string, instruction *uint32, consumedBits *uint32) {
	regValue := REGISTERS[value]
	*instruction |= (uint32(regValue) << (32 - *consumedBits - 4))
	*consumedBits += 4
}

// consume two bytes consumes a two byte value and adds it to the
// instruction
func consumeTwoBytes(value string, instruction *uint32, consumedBits *uint32) {
	val, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		panic(err)
	}
	*instruction |= (uint32(val)) << (32 - *consumedBits - 16)
	*consumedBits += 16
}

// resolveAddressingMode resolves the addressing mode for an instruction
// given its opcode table entry
func resolveAddressingMode(node SyntaxNode, opEntry []uint16) uint8 {
	// for each operation exactly 1 parameter supports different ways of being called
	// we just need to identify that parameter for this opEntry, the complication
	// is that IMMEDIATE values can be of several differnt types
	supportsManyArgs := 2
	for supportsManyArgs < len(opEntry) && (opEntry[supportsManyArgs] == RegisterValue ||
		opEntry[supportsManyArgs] == RegisterRelativeValue ||
		opEntry[supportsManyArgs] == PCRelativeValue ||
		isImmediate(opEntry[supportsManyArgs])) {
		supportsManyArgs++
	}

	// we have to return the appropriate bitcode for this addrmode
	argIndex := max(0, supportsManyArgs-2)
	// there is an edge case here with an operation has a single argument
	if len(node.Children) == 1 {
		argIndex = 0
	}

	return ADDRMODES[node.Children[argIndex].NodeType]
}

// is immediate addr mode tells us if a integer of accepted types is all
// immediate
func isImmediate(nodeType uint16) bool {
	return nodeType&RegisterValue == 0 &&
		nodeType&RegisterRelativeValue == 0 &&
		nodeType&PCRelativeValue == 0
}

// validateInstructionOperands iterates over every instruction and validates
// if its operands are valid
func validateInstructionOperands(nodes []SyntaxNode, relocationTable map[string]uint16) {
	for _, node := range nodes {
		if node.NodeType == Label {
			continue
		}

		// iterate over the children for an instruction node
		opData := OPCODES[node.Value]

		for i, arg := range node.Children {
			argType := arg.NodeType
			if argType == Label {
				if _, ok := relocationTable[arg.Value]; !ok {
					panic(fmt.Sprintf(`Compilation Error - Undefined label "%s" on line: %d.`,
						arg.Value, node.line))
				}
			}

			if opData[2+i]&argType == 0 {
				panic(fmt.Sprintf(`Compilation Error - Invalid argument type of "%s" for instruction "%s" on line: %d.`,
					arg.Value, node.Value, node.line))
			}
		}
	}
}

// note on addresses: addresses are all 16 bit unsigned integers
func computeRelocationTable(nodes []SyntaxNode) map[string]uint16 {
	var currentAddr uint16 = 0
	var relocationTable = make(map[string]uint16)

	for _, node := range nodes {
		if node.NodeType == Label {
			// resolve this label by first checking if its been relocated yet
			if _, ok := relocationTable[node.Value]; ok {
				panic(fmt.Sprintf(`Compilation Error - Duplicate definition of label: "%s"`, node.Value))
			} else {
				relocationTable[node.Value] = currentAddr + 4
			}
		}

		currentAddr += 4
	}
	return relocationTable
}
