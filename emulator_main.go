package main

import "chip8/emulator"


func main() {
	chip := emulator.NewChip()
	chip.LoadROM("print_ten.chip")

	// forEVAAAA :)
	for {
		chip.PerformNextComputation()
	}
}