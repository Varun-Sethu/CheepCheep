package main

import (
	"chip8/emulator"
	"fmt"
	"os"
)


func main() {
	chip := emulator.NewChip()
	sourceRom := fmt.Sprintf("%s.chip", os.Args[1])
	chip.LoadROM(sourceRom)

	// forEVAAAA :)
	for {
		chip.PerformNextComputation()
	}
}