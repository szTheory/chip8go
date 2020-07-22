package emu

import "fmt"

type Emulator struct {
	cpu     *CPU
	memory  *Memory
	display *Display
	// input   *Input
}

func (e *Emulator) Setup(romFilename string) {
	e.cpu = new(CPU)
	e.cpu.Setup()

	e.memory = new(Memory)
	e.memory.Setup()
	e.memory.LoadGame(romFilename)

	e.display = new(Display)

	// input = new(Input)
}

func (e *Emulator) EmulateCycle() {
	// Fetch next instruction and advance the Program Counter
	var instruction uint16 = (uint16(e.memory.RAM[e.cpu.PC]) << 8) | uint16(e.memory.RAM[e.cpu.PC+1])
	e.cpu.PC += 2

	fmt.Printf("%X\n", instruction)

	switch instruction {
	// 00E0 - CLS
	case 0x00E0:
		panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
	// 00EE - RET
	case 0x00EE:
		panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
	default:
		// panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		switch instruction & 0xF000 {
		// 0nnn - SYS addr
		case 0x0000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// 1nnn - JP addr
		case 0x1000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// 2nnn - CALL addr
		case 0x2000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// 3xkk - SE Vx, byte
		case 0x3000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// 4xkk - SNE Vx, byte
		case 0x4000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// 5xy0 - SE Vx, Vy
		case 0x5000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))

		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		// The interpreter puts the value kk into register Vx.
		case 0x6000:
			register := instruction & 0x0F00 >> 8
			value := uint8(instruction & 0x00FF)
			e.cpu.V[register] = value

		// 7xkk - ADD Vx, byte
		case 0x7000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		case 0x8000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			// 8xy0 - LD Vx, Vy
			// 8xy1 - OR Vx, Vy
			// 8xy2 - AND Vx, Vy
			// 8xy3 - XOR Vx, Vy
			// 8xy4 - ADD Vx, Vy
			// 8xy5 - SUB Vx, Vy
			// 8xy6 - SHR Vx {, Vy}
			// 8xy7 - SUBN Vx, Vy
			// 8xyE - SHL Vx {, Vy}
		// 9xy0 - SNE Vx, Vy
		case 0x9000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))

		// Annn - LD I, addr
		// Set I = nnn.
		// The value of register I is set to nnn.
		case 0xA000:
			e.cpu.I = instruction & 0x0FFF

		// Bnnn - JP V0, addr
		case 0xB000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// Cxkk - RND Vx, byte
		case 0xC000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// Dxyn - DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		// The interpreter reads n bytes from memory, starting at the address stored in I.
		// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
		// Sprites are XORed onto the existing screen. If this causes any pixels to be erased,
		// VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it
		// is outside the coordinates of the display, it wraps around to the opposite side of the screen.
		case 0xD000:
			x := byte(instruction & 0x0F00 >> 8)
			y := byte(instruction & 0x00F0 >> 4)
			height := instruction & 0x000F

			var i uint16 = 0
			for ; i < height; i++ {
				row := e.memory.RAM[e.cpu.I+i]

				if e.display.DrawSprite(x, y, row) {
					e.cpu.V[0xF] = 1
				}
			}

		case 0xE000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			// Ex9E - SKP Vx
			// ExA1 - SKNP Vx
		case 0xF000:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			// Fx07 - LD Vx, DT
			// Fx0A - LD Vx, K
			// Fx15 - LD DT, Vx
			// Fx18 - LD ST, Vx
			// Fx1E - ADD I, Vx
			// Fx29 - LD F, Vx
			// Fx33 - LD B, Vx
			// Fx55 - LD [I], Vx
			// Fx65 - LD Vx, [I]
		default:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		}
	}

	// Decode Opcode
	// Execute Opcode

	// Update timers

}
