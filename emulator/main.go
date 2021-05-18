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

// computeOperand determines what values should be inputted into the operation, note its only called by functions
// that support multiple addressing modes, additionally; depending on the function, the amount of desired "bytes" in the
// operand is provided: eg. if the operand is expected to be an address (like jump commands) then it requests 2 bytes
// along with the value we just computed we return how many bytes were "used" to get that value from the code segment
func (c *Chipster) computeOperand(addrMode uint8, requestedBytes uint16) (uint16, uint16) {

	// First compute the physical location in memory specified by the addressing mode, if we are accessing a register
	// just flick its flag
	var memoryLocation uint16 = c.Pc
	var usedBytes uint16 = requestedBytes
	var isRegisterAccess bool = false

	switch addrMode {
		case immediate:
			break
		case direct:
			memoryLocation = uint16(c.Memory[c.Pc]) << 8 | uint16(c.Memory[c.Pc + 1])
			usedBytes = 2
			break
		case indirect:
			baseLocation := uint16(c.Memory[c.Pc]) << 8 | uint16(c.Memory[c.Pc + 1])
			memoryLocation = uint16(c.Memory[baseLocation]) << 8 | uint16(c.Memory[baseLocation + 1])
			usedBytes = 2
			break
		case registerDirect:
			memoryLocation = uint16(c.Memory[c.Pc])
			isRegisterAccess = true
			break
	}

	// Now resolve and compute the operand data
	if isRegisterAccess {
		return uint16(c.Registers[memoryLocation]), usedBytes
	} else {
		// 2 possible situations: if the requested bytes was a single value or if the requested byte was 2 bytes
		// 2 bytes implies we read the next two values, otherwise we read a single value
		var highByte uint8 = 0
		var lowByte uint8  = c.Memory[c.Pc]
		if requestedBytes == 2 {
			highByte = c.Memory[c.Pc]
			lowByte = c.Memory[c.Pc + 1]
		}

		return uint16(highByte) << 8 | uint16(lowByte), usedBytes
	}
}





// PerformNextComputation reads the current instruction from memory and performs the dictated instruction
func (c *Chipster) PerformNextComputation() {

	// extract the current instruction, opcode and the addressing mode
	var currentInstruction uint8 = c.Memory[c.Pc]
	var opcode   uint8 = (currentInstruction & (0xf8)) >> 3
	var addrMode uint8 = currentInstruction & 0x7
	c.Pc += 1

	switch opcode {
		case PRINT:
			// Fetch the next byte from memory
			var targetRegister uint8 = c.Memory[c.Pc]
			c.Pc += 1
			fmt.Printf("Outputted: %d\n", c.Registers[targetRegister])
			break

		case LDR:
			// fetch the target register
			var targetRegister uint8 = c.Memory[c.Pc]
			c.Pc += 1

			// fetch the operand and increment the program counter
			loadValue, usedBytes := c.computeOperand(addrMode, 1)
			c.Registers[targetRegister] = uint8(loadValue)
			c.Pc += usedBytes

			break

		case ADD:
			// fetch the target register
			var targetRegister uint8 = c.Memory[c.Pc]
			c.Pc += 1

			// fetch the operand and increment the program counter
			additiveValue, usedBytes := c.computeOperand(addrMode, 1)
			c.Registers[targetRegister] += uint8(additiveValue)
			c.Pc += usedBytes

			break

		case CMP:
			// fetch the target register
			var targetRegister uint8 = c.Memory[c.Pc]
			c.Pc += 1

			// fetch the operand and increment the program counter
			valueToCompare, usedBytes := c.computeOperand(addrMode, 1)
			c.Pc += usedBytes

			// compare the two values and based on the result of the comparison, set the corresponding flag register
			var comparison int8 = int8(c.Registers[targetRegister]) - int8(valueToCompare)
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
			jumpDestination, consumedBytes := c.computeOperand(addrMode, 2)
			c.Pc += consumedBytes

			if flagRegister == 0x2 {
				// Perform the jump to the requested location
				// we shouldn't increment the bytes we consumed during a jump instruction
				c.Pc = jumpDestination
				return
			}

			break

		case HLT:
			//fmt.Printf("Unidentified opcode: %08b\n", opcode)
			c.Pc -= 1
			return
	}
}




















