package emu

import "fmt"

func Run(romFilename string) {
	cpu := new(CPU)
	mem := new(Memory)
	display := new(Display)
	input := new(Input)

	cpu.Setup()
	mem.Setup()
	mem.LoadGame(romFilename)

	loop(cpu, mem, display, input)
}

func loop(cpu *CPU, mem *Memory, display *Display, input *Input) {
	for {
		emulateCycle(cpu, mem)
		display.Draw(mem)
		input.Update()
	}
}

func emulateCycle(cpu *CPU, mem *Memory) {
	// fmt.Println("Emu Cycle")

	// Fetch next instruction and advance the Program Counter
	var instruction uint16 = (uint16(mem.RAM[cpu.PC]) << 8) | uint16(mem.RAM[cpu.PC+1])
	// cpu.PC += 2

	fmt.Printf("%X\n", instruction)

	// opcode := mem.RAM[cpu.PC]<<8 | mem.RAM[cpu.PC+1]

	switch instruction {
	// 00E0 - CLS
	case 0x00E0:

	// 00EE - RET
	case 0x00EE:
	default:
		// panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		switch instruction & 0xF000 {
		default:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))

		// 0nnn - SYS addr
		case 0x0000:
		// 1nnn - JP addr
		case 0x1000:
		// 2nnn - CALL addr
		case 0x2000:
		// 3xkk - SE Vx, byte
		case 0x3000:
		// 4xkk - SNE Vx, byte
		case 0x4000:
		// 5xy0 - SE Vx, Vy
		case 0x5000:
		// 6xkk - LD Vx, byte
		case 0x6000:
		// 7xkk - ADD Vx, byte
		case 0x7000:
		case 0x8000:
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
		// Annn - LD I, addr
		case 0xA000:
		// Bnnn - JP V0, addr
		case 0xB000:
		// Cxkk - RND Vx, byte
		case 0xC000:
		// Dxyn - DRW Vx, Vy, nibble
		case 0xD000:
		case 0xE000:
			// Ex9E - SKP Vx
			// ExA1 - SKNP Vx
		case 0xF000:
			// Fx07 - LD Vx, DT
			// Fx0A - LD Vx, K
			// Fx15 - LD DT, Vx
			// Fx18 - LD ST, Vx
			// Fx1E - ADD I, Vx
			// Fx29 - LD F, Vx
			// Fx33 - LD B, Vx
			// Fx55 - LD [I], Vx
			// Fx65 - LD Vx, [I]
		}
	}

	// Decode Opcode
	// Execute Opcode

	// Update timers

}
