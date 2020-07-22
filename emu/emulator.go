package emu

import "fmt"

type Emulator struct {
	cpu     *CPU
	memory  *Memory
	Display *Display
	// input   *Input
}

func (e *Emulator) Setup(romFilename string) {
	e.cpu = new(CPU)
	e.cpu.Setup()

	e.memory = new(Memory)
	e.memory.Setup()
	e.memory.LoadGame(romFilename)

	e.Display = new(Display)

	// input = new(Input)
}

func (e *Emulator) EmulateCycle() {
	// Fetch next instruction and advance the Program Counter
	var instruction uint16 = (uint16(e.memory.RAM[e.cpu.PC]) << 8) | uint16(e.memory.RAM[e.cpu.PC+1])
	e.cpu.PC += 2

	fmt.Printf("%X\n", instruction)

	// if instruction == 0x65EE {
	// 	fmt.Println("Breakpoint")
	// }

	switch instruction {
	// 00E0 - CLS
	// Clear the display.
	case 0x00E0:
		fmt.Println("--- 00E0")
		e.Display.Clear()

	// 00EE - RET
	// Return from a subroutine.
	// The interpreter sets the program counter to the address at the top of
	// the stack, then subtracts 1 from the stack pointer.
	case 0x00EE:
		fmt.Println("--- 00EE")
		e.cpu.PC = uint16(e.cpu.Stack[e.cpu.SP])
		e.cpu.SP--

	default:
		// panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		switch byte(instruction & 0xF000 >> 12) {

		// 0nnn - SYS addr
		// Jump to a machine code routine at nnn.
		// This instruction is only used on the old computers on which Chip-8 was
		// originally implemented. It is ignored by modern interpreters.
		case 0x0:
			fmt.Println("--- 0nnn")
			// Do nothing.
			// This instruction is ignored by modern interpreters

		// 1nnn - JP addr
		// Jump to location nnn.
		// The interpreter sets the program counter to nnn.
		case 0x1:
			fmt.Println("--- 1nnn")
			address := instruction & 0xFFF
			// fmt.Printf("jump %X\n", address)
			e.cpu.PC = address

		// 2nnn - CALL addr
		// Call subroutine at nnn.
		// The interpreter increments the stack pointer, then puts the current PC
		// on the top of the stack. The PC is then set to nnn.
		case 0x2:
			fmt.Println("--- 2nnn")
			address := instruction & 0xFFF
			e.cpu.SP++
			e.cpu.Stack[e.cpu.SP] = e.cpu.PC
			e.cpu.PC = address

		// 3xkk - SE Vx, byte
		// Skip next instruction if Vx = kk.
		// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
		case 0x3:
			fmt.Println("--- 3xkk")
			registerOffset := byte(instruction & 0xF00 >> 8)
			value := byte(instruction & 0xFF)
			if e.cpu.V[registerOffset] == value {
				e.cpu.PC += 2
			}

		// 4xkk - SNE Vx, byte
		case 0x4:
			fmt.Println("--- 4xkk")
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))

		// 5xy0 - SE Vx, Vy
		// Skip next instruction if Vx = Vy.
		// The interpreter compares register Vx to register Vy, and if they are equal,
		// increments the program counter by 2.
		case 0x5:
			fmt.Println("--- 5xy0")
			x := byte(instruction & 0xF00 >> 8)
			y := byte(instruction & 0xF0 >> 4)
			if x == y {
				e.cpu.PC += 2
			}

		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		// The interpreter puts the value kk into register Vx.
		case 0x6:
			fmt.Println("-- 6xkk")
			registerOffset := byte(instruction & 0xF00 >> 8)
			value := byte(instruction & 0xFF)
			e.cpu.V[registerOffset] = value

		// 7xkk - ADD Vx, byte
		// Set Vx = Vx + kk.
		// Adds the value kk to the value of register Vx, then stores the result in Vx.
		case 0x7:
			fmt.Println("--- 7xkk")
			registerOffset := byte(instruction & 0xF00 >> 8)
			value := byte(instruction & 0xFF)
			e.cpu.V[registerOffset] += value

		case 0x8:
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
		case 0x9:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))

		// Annn - LD I, addr
		// Set I = nnn.
		// The value of register I is set to nnn.
		case 0xA:
			fmt.Println("--- Annn")
			e.cpu.I = instruction & 0xFFF

		// Bnnn - JP V0, addr
		case 0xB:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// Cxkk - RND Vx, byte
		case 0xC:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		// Dxyn - DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		// The interpreter reads n bytes from memory, starting at the address stored in I.
		// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
		// Sprites are XORed onto the existing screen. If this causes any pixels to be erased,
		// VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it
		// is outside the coordinates of the display, it wraps around to the opposite side of the screen.
		case 0xD:
			fmt.Println("--- Dxyn")
			x := byte(instruction & 0xF00 >> 8)
			y := byte(instruction & 0xF0 >> 4)
			height := instruction & 0xF

			var i uint16 = 0
			for ; i < height; i++ {
				row := e.memory.RAM[e.cpu.I+i]

				if e.Display.DrawSprite(x, y, row) {
					e.cpu.V[0xF] = 1
				}
			}

		case 0xE:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			// Ex9E - SKP Vx
			// ExA1 - SKNP Vx
		case 0xF:
			switch instruction & 0xFF {

			// Fx07 - LD Vx, DT
			// Set Vx = delay timer value.
			// The value of DT is placed into Vx.
			case 0x07:
				fmt.Println("--- Fx07")
				offset := byte(instruction & 0xF00 >> 8)
				e.cpu.V[offset] = e.cpu.DelayTimer

			// Fx0A - LD Vx, K
			case 0x0A:
				panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			// Fx15 - LD DT, Vx
			case 0x15:
				panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			// Fx18 - LD ST, Vx
			case 0x18:
				panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			// Fx1E - ADD I, Vx
			case 0x1E:
				panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))

			// Fx29 - LD F, Vx
			// Set I = location of sprite for digit Vx.
			// The value of I is set to the location for the hexadecimal sprite corresponding
			// to the value of Vx. See section 2.4, Display, for more information on the
			// Chip-8 hexadecimal font.
			case 0x29:
				fmt.Println("--- Fx29")
				registerOffset := byte(instruction & 0xF00 >> 8)
				pixelFontAddress := RamFontStart + e.cpu.V[registerOffset]
				e.cpu.I = uint16(pixelFontAddress)

			// Fx33 - LD B, Vx
			// Store BCD representation of Vx in memory locations I, I+1, and I+2.
			// The interpreter takes the decimal value of Vx, and places the hundreds
			// digit in memory at location in I, the tens digit at location I+1,
			// and the ones digit at location I+2.
			case 0x33:
				fmt.Println("--- Fx33")
				registerOffset := byte(instruction & 0xF00 >> 8)
				decimalValue := e.cpu.V[registerOffset]
				e.memory.RAM[e.cpu.I] = decimalValue / 100  //hundreds
				e.memory.RAM[e.cpu.I+1] = decimalValue / 10 //tens
				e.memory.RAM[e.cpu.I+2] = decimalValue / 1  //ones

			// Fx55 - LD [I], Vx
			// Store registers V0 through Vx in memory starting at location I.
			// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
			case 0x55:
				fmt.Println("--- Fx55")
				maxRegisterOffset := byte(instruction & 0xF00 >> 8)
				var i byte = 0
				for ; i <= maxRegisterOffset; i++ {
					e.memory.RAM[e.cpu.I+uint16(i)] = e.cpu.V[i]
				}

			// Fx65 - LD Vx, [I]
			// Read registers V0 through Vx from memory starting at location I.
			// The interpreter reads values from memory starting at location I
			// into registers V0 through Vx.
			case 0x65:
				fmt.Println("--- Fx65")
				maxRegisterOffset := byte(instruction & 0xF00 >> 8)
				var i byte = 0
				for ; i <= maxRegisterOffset; i++ {
					e.cpu.V[i] = e.memory.RAM[e.cpu.I+uint16(i)]
				}

			default:
				panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
			}
		default:
			panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
		}
	}

	// Decode Opcode
	// Execute Opcode

	// Update timers

}
