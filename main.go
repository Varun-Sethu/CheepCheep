package main

import (
	"chip8/emulator"
)



// yes.
func main() {
	chip := emulator.NewChip()
	chip.LoadROM("IBMLogo.ch8")
	//chip.DumpMemory()

	// infinite loooop :)
	for {
		command := chip.Fetch()
		chip.Tick(command)
	}
}
// 0b11100000