package interpret

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

var ChipPixels = make([]byte, 64*32)

func loadFile(fileName string) ([]byte, error) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

func Run(fileName string) error {
	input, err := loadFile(fileName)
	if err != nil {
		return err
	}

	var memory [4096]byte
	for i, b := range input {
		memory[0x200+i] = b // programs should be loaded at 0x200
	}

	var pc uint16 = 0x200

	var registers [16]byte // 16 - 8 bit registers (V)
	var iRegister [2]byte  // single - 16 bit register

	for {
		instr := memory[pc : pc+2]

		pc += 2

		switch {
		case bytes.Equal(instr, []byte{0x00, 0xe0}): // CLS - clear screen
		case instr[0] >= 0x10 && instr[0] <= 0x1f: // 1NNN jump TODO: bitwise this
			pc = uint16(instr[0]-0x10)<<8 + uint16(instr[1])
		case instr[0] >= 0x60 && instr[0] <= 0x6f: // 6XNN set register VX
			regNum := instr[0] - 0x60
			registers[regNum] = instr[1]
		case instr[0] >= 0x70 && instr[0] <= 0x7f: // 7XNN add to register VX
			regNum := instr[0] - 0x70
			registers[regNum] = registers[regNum] + instr[1]
		case instr[0] >= 0xa0 && instr[0] <= 0xaf: // aNNN set register I
			iRegister[0] = instr[0] - 0xa0
			iRegister[1] = instr[1]
		case instr[0] >= 0xd0 && instr[0] <= 0xdf: // dXYN draw
			x := instr[0] & 0x0f
			y := (instr[1] & 0xf0) >> 4
			spriteByteCount := instr[1] & 0x0f

			xStart := int(registers[x] % 64)
			yStart := int(registers[y] % 32)

			registers[0xf] = 0

			spriteAddr := toUint16(iRegister[:])
			for i, row := range memory[0x200+spriteAddr : 0x200+spriteAddr+uint16(spriteByteCount)] { // each row
				for j := 0; j < 8; j++ { // each pixel
					xOffset := xStart + j
					yOffset := (int(yStart) + i) * 64
					idx := xOffset + yOffset
					pixVal := ChipPixels[idx]

					maskedVal := row & 1 << (7 - j)
					if maskedVal != 0 { // bit is set so we need to flip pixel
						if pixVal > 0 {
							ChipPixels[idx] = 0
							registers[0xf] = 1
						} else {
							ChipPixels[idx] = 255
						}
					}
				}
			}

			for i := 0; i < 64; i++ {
				for j := 0; j < 32; j++ {
					if ChipPixels[i+j*64] == 255 {
						fmt.Printf("1")
					} else {
						fmt.Printf("0")
					}
				}
				fmt.Println()
			}
		default:
			fmt.Println("UNKNOWN INSTRUCTION")
		}
	}

	return nil
}

func toUint16(hexSlice []byte) uint16 {
	if len(hexSlice) != 2 {
		panic("toUint16 requires a slice of len 2")
	}

	return uint16(hexSlice[0]<<8) + uint16(hexSlice[1])
}
