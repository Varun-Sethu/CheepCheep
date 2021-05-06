package emulator

import (
	"fmt"
	"os"
	"os/exec"
)

// font data for the emulator
var fontData = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}


// defines the internal state of the Chip8
type Chip8 struct {

	// Defines the basic entries in the emulator, memory related
	Memory     [4096]uint8
	Display    [64][32]bool
	Registers  [16]uint8
	Stack	   [16]uint16

	// Control flow variables
	Pc uint16 // program counter
	Sp uint16 // stack pointer
	Ir uint16 // index register
	Vf uint16 // flag register


	// General timers
	DelayTimer uint8
	SoundTimer uint8
}





// NewChip allocates a new chip and writes the font data to memory
func NewChip() Chip8 {
	chip := Chip8{
		Memory:     [4096]uint8{},
		Display:    [64][32]bool{},
		Registers:  [16]uint8{},
		Stack: 		[16]uint16{},
		Pc:         0x200,
		Sp:         0,
		Ir:         0,
		DelayTimer: 0,
		SoundTimer: 0,
	}

	// Write the font data to memory
	copy(chip.Memory[0x50 - 1:], fontData[:])
	return chip
}


// Just prints the entire contents of RAM to the screen
func (chip* Chip8) DumpMemory() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 64; j++ {
			fmt.Printf("%3d ", chip.Memory[64*i + j])
		}
		fmt.Print("\n")
	}
	fmt.Print("\n\n")
}



func (chip *Chip8) LoadROM(romName string) {
	// TODO: implement
	rom, err := os.OpenFile(romName, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}

	defer rom.Close()

	// open the binary file and read
	romStat, _ := rom.Stat()
	buffer     := make([]byte, romStat.Size())
	_, err = rom.Read(buffer)
	if err != nil {
		panic(err)
	}

	// load the rom into memory
	for i := 0; i < len(buffer); i++ {
		chip.Memory[i + 0x200] = buffer[i]
	}
}




// Fetches the current command pointed at by the PC
func (chip* Chip8) Fetch() uint16 {

	command := (uint16(chip.Memory[chip.Pc]) << 8) | uint16(chip.Memory[chip.Pc + 1])
	chip.Pc += 2

	return command
}



// Tick executes the command currently pointed at by the PC
func (chip *Chip8) Tick(command uint16) {


	// determine the opcode and arguments for the current instruction
	var opcode 		 	  	 = uint8((command & 0xF000) >> 8)
	var targetRegister 		 = uint8((command & 0xF00) >> 8)
	var arguments  		 	 = uint8(command & 0xFF)


	// Chip8 instructions are broken up as follows:
	/*
		[a a a a] [b b b b] [c c c c] [d d d d]
		[a a a a] contains the opcode
		[b b b b] contains a lookup for the registers
		[c c c c] contains another lookup for the registers
		[d d d d] some 4 bit number

		We denote and instruction as follows
		[opcode]KNN
	 */

	// Implementations of behaviour for each opcode
	switch opcode {

	// If we receive a "system" opcode, just check if it is a request to leave a subroutine or clear the screen
	case SYS:
		switch uint8(command & 0xff)  {
			case CLR:
				// Iterate over the display and set everything to o
				for y := 0; y < len(chip.Display); y++ {
					for x := 0; x < len(chip.Display[0]); x++ {
						chip.Display[y][x] = false
					}
				}
				break
			case RET:
				// Pop the current entry in the stack off and change the PC to the stack entry
				chip.Sp -= 1
				chip.Pc = chip.Stack[chip.Sp]
				chip.Pc += 2

				break
		}

	// Control flow implementations
	case CALL:
		chip.Sp += 1
		chip.Stack[chip.Sp] = chip.Pc
		// Compute the address of the requested method
		chip.Pc = (uint16(targetRegister) << 8) | uint16(arguments)

		break

	case JMP:
		// Determine what index in memory to point at
		var memIndex = (uint16(targetRegister) << 8) | uint16(arguments)
		chip.Pc = memIndex

		break
	case SE:
		// Retrieve the value at Vk
		data := chip.Registers[targetRegister]
		if data == arguments {
			chip.Pc += 2
		}
		break
	case SNE:
		data := chip.Registers[targetRegister]
		if data != targetRegister {
			chip.Pc += 2
		}
		break
	case SER:
		datax := chip.Registers[targetRegister]
		datay := chip.Registers[(arguments & 0xf0) >> 4]
		if datax == datay {
			chip.Pc += 2
		}
		break
	case SNEK:
		// Skips the next value if the data at register Vx != kk
		data := chip.Registers[targetRegister]
		if data != arguments {
			chip.Pc += 2
		}
		break





	// Arithmetic and logical operations
	case ALU:
		regx := targetRegister
		regy := (arguments & 0xf0) >> 4

		// Switch on all the possible operations
		switch arguments & 0xf {
			case LDR:
				chip.Registers[regx] = chip.Registers[regy]
				break
			case OR:
				chip.Registers[regx] |= chip.Registers[regy]
				break
			case AND:
				chip.Registers[regx] &= chip.Registers[regy]
				break
			case XOR:
				chip.Registers[regx] ^= chip.Registers[regy]
				break
			case ADD:
				result := chip.Registers[regx] + chip.Registers[regy]
				if result > 255 {
					chip.Vf = 1
				}
				chip.Registers[regx] = uint8(result & 0xff)
				break
			case SUB:
				if chip.Registers[regy] > chip.Registers[regx] {
					chip.Vf = 1
					chip.Registers[regx] = chip.Registers[regy] - chip.Registers[regx]
				} else {
					chip.Registers[regx] -= chip.Registers[regy]
				}
				break
			case SHR:
				if (chip.Registers[regx] & 0x1) == uint8(1) {
					chip.Vf = 1
				}
				chip.Registers[regx] >>= 1
				break
			case SUBN:
				if chip.Registers[regy] < chip.Registers[regx] {
					chip.Vf = 1
					chip.Registers[regx] = chip.Registers[regx] - chip.Registers[regy]
				} else {
					chip.Registers[regx] = chip.Registers[regy] - chip.Registers[regx]
				}
				break
			case SHL:
				if (chip.Registers[regx] & 0x80) == uint8(1) {
					chip.Vf = 1
				}
				chip.Registers[regx] <<= 1
				break

		}

		break











	case LD:
		// load NN into register K
		chip.Registers[targetRegister] = arguments
		break

	case ADDK:
		// Add the value of the argument to the specified target register
		chip.Registers[targetRegister] += arguments
		break

	case LDI:
		var value = (uint16(targetRegister) << 8) | uint16(arguments)
		chip.Ir = value
		break




	case DISP:
		// Write the data in the specified register locations to the display array
		var requestedSprite = chip.Ir
		var vx = targetRegister
		var vy = (arguments & 0xf0) >> 4
		var n  = uint16(arguments & 0xf)

		x := uint16(chip.Registers[vx] % 64)
		y := uint16(chip.Registers[vy] % 32)

		// Modify the display in accordance to what was requested
		for offset := uint16(0); offset < n; offset++ {
			// If we are going to print over the screen just stop trying

			// Determine what to "modify"
			var spriteRow = chip.Memory[requestedSprite + offset]

			// iterate over each bit in the byte and output everything to the correct pixel
			for i := uint8(0); i < uint8(8); i++ {
				// Determine if we should draw this row or not
				if x + uint16(i) > 64 {
					break
				}

				var v bool = (((0x80 >> i) & spriteRow) >> (7 - i)) != 0 // Magic :O
				if v {
					// TODO: Flick indicated register
					chip.Display[x + uint16(i)][y] = true
				}
			}
			y += 1
		}

		// Flush standard output
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()



		// Display the screen of the chip8
		for x := 0; x < len(chip.Display[0]); x++ {
			for y := 0; y < len(chip.Display); y++ {
				if chip.Display[y][x] {
					fmt.Print("█	█")
				} else {
					fmt.Print("  ")
				}
			}
			fmt.Print("\n")
		}


		break
	}


}

























