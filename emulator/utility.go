package emulator

import (
	"fmt"
	"os"
)

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



