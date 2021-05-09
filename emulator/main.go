package emulator



// NOTE: this file used to contain code for a small chip8 emulator i built, I am yet to actually implement an emulator for my ISA




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




























