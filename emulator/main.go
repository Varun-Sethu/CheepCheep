package emulator

import (
	"fmt"
	"os"
)

/**
	Implementation Details:
		- The program to compute is located in memory from 0x0 onwards
		- Memory grows from there

		- The CPU can address 4k bytes of memory
		- There are 16 registers most are 8 bits except:
			- The stack pointer
			- The stack base pointer
			- Which are 16 bits long
			- The register codes are r14 and r15 respectively
 */




// defines the internal state of the Chip8
type Chipster struct {

	// Defines the basic entries in the emulator, memory related
	Memory     [4096]uint8
	Registers  [14]uint8
	StackRegisters [2]uint16

	// Control flow variables
	Pc uint16 // program counter


	// Of special importance in the VF flag are the first least significant two bits
	// LSB: if last operation was 0, LSB 2.0: if last operation resulted in a negative number
	Vf uint16 // flag register

}


// NewChip builds and returns a new chip
func NewChip() Chipster {
	return Chipster{
		Memory: [4096]uint8{0},
		Registers: [14]uint8{0},
		StackRegisters: [2]uint16{0},

		Pc: 0,
		Vf: 0,
	}
}





// LoadROM opens a file from the OS and reads it into memory
// note: it is assumed that the source file is simply a compiled file
func (c *Chipster) LoadROM(sourceFile string) {

	file, err := os.OpenFile(sourceFile, os.O_RDONLY, 0777)
	// TODO: fix this later
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read the file into a bytecode buffer
	fileSize, _ := file.Stat()
	buffer := make([]byte, fileSize.Size())
	file.Read(buffer)

	// finally copy the bytecode buffer into the chip
	copy(c.Memory[:], buffer)

}

// TODO: fix up addressing mode logic



// resolveParams retrieves the next set of parameters for an operation given the addressing mode, also returns the amount
// of consumed bytes.
func (c *Chipster) resolveParams(addrMode uint8, currentOffset uint16) (uint16, uint16) {
	// Note: this retrieves the VALUE "pointed" at by the addrMode

	switch addrMode {
		case immediate:
			return uint16(c.Memory[c.Pc + currentOffset]), 1
		case direct:
			var requestedAddr uint16 = uint16(c.Memory[c.Pc + currentOffset]) << 8 | uint16(c.Memory[c.Pc + 1 + currentOffset])
			return uint16(c.Memory[requestedAddr]), 2
		case registerDirect:
			return uint16(c.Registers[c.Memory[c.Pc + currentOffset]]), 1
		default:
			return 0, 0
	}
}





// PerformNextComputation reads the current instruction from memory and performs the dictated instruction
func (c *Chipster) PerformNextComputation() {

	// extract the current instruction, opcode and the addressing mode
	var currentInstruction uint8 = c.Memory[c.Pc]
	var opcode   uint8 = (currentInstruction & (0xf8)) >> 3
	var addrMode uint8 = currentInstruction & 0x7


	// represents how many bytes of memory the instruction we are current processing took up
	var consumedBytes uint16 = 1

	switch opcode {
		case PRINT:
			// Fetch the next byte from memory
			var targetRegister uint8 = c.Memory[c.Pc + consumedBytes]
			fmt.Printf("Outputted: %d\n", c.Registers[targetRegister])
			consumedBytes += 1
			break

		case LDR:
			var targetRegister uint8 = c.Memory[c.Pc + consumedBytes]
			consumedBytes += 1

			loadValue, operandBytes := c.resolveParams(addrMode, consumedBytes)
			c.Registers[targetRegister] = uint8(loadValue)
			consumedBytes += operandBytes
			break

		case ADD:
			var targetRegister uint8 = c.Memory[c.Pc + consumedBytes]
			consumedBytes += 1

			loadValue, operandBytes := c.resolveParams(addrMode, consumedBytes)
			c.Registers[targetRegister] += uint8(loadValue)
			consumedBytes += operandBytes
			break

		case CMP:
			// read the values from memory
			var targetRegister uint8 = c.Memory[c.Pc + consumedBytes]
			consumedBytes += 1
			loadValue, operandBytes := c.resolveParams(addrMode, consumedBytes)
			consumedBytes += operandBytes

			// compare the two values and based on the result of the comparison, set the corresponding flag register
			var comparison int8 = int8(c.Registers[targetRegister]) - int8(loadValue)
			c.Vf &= 0xfffc // unset the last two bits in the flag register
			switch {
				case comparison == 0:
					c.Vf |= 0x1 // set the last two bits to the appropriate value, in this case 01
					break
				case comparison < 0:
					c.Vf |= 0x2 // set last two bits to 10
					break
				case comparison > 0:
					// do nothing :)
					break
			}

			break

		case JMPL:
			// Read from the flag registers, since its JMPL the final two bits should read: 10
			var flagRegister uint16 = c.Vf & 0x3
			if flagRegister == 0x2 {
				// Perform the jump to the requested location
				// we shouldn't increment the bytes we consumed during a jump instruction
				var requestedAddr uint16 = uint16(c.Memory[c.Pc + consumedBytes]) << 8 | uint16(c.Memory[c.Pc + 1 + consumedBytes])
				c.Pc = requestedAddr
				return
			}

			break

		case HLT:
			//fmt.Printf("Unidenfitied opcode: %08b\n", opcode)
			return
	}


	// increment the program counter based on the amount of consumed bytes
	c.Pc += consumedBytes
}




















