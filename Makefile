ASSEMBLER_NAME = chippy.out
EMULATOR_NAME = emulator.out
ROM_DIR = ./ROMs
ROM_OUT_DIR = ./binaries

.PHONY: emulator
emulator:
	go build -o ${EMULATOR_NAME} emulator.go

assembler:
	go build -o ${ASSEMBLER_NAME} chippy.go

all: assembler emulator

roms: assembler
	$(foreach file, $(wildcard $(ROM_DIR)/*), ./${ASSEMBLER_NAME} ${file} binaries/$(basename $(notdir $(file))) false;)	

.PHONY: clean
clean:
	-rm -f binaries/*.chip
	-rm ${EMULATOR_NAME}
	-rm ${ASSEMBLER_NAME}
