package emu

import (
	"fmt"
	"math/rand"
)

type Emulator struct {
	cpu     *CPU
	memory  *Memory
	Display *Display
	Input   *Input

	waitingForInputRegisterOffset byte
}

func (e *Emulator) Setup(romFilename string) {
	e.cpu = new(CPU)
	e.cpu.Setup()

	e.memory = new(Memory)
	e.memory.Setup()
	e.memory.LoadGame(romFilename)

	e.Display = new(Display)

	e.Input = new(Input)
}

func (e *Emulator) CatchInput(keyIndex byte) {
	// fmt.Println("~ ~ ~ ~ ", e.waitingForInputRegisterOffset)
	e.cpu.V[e.waitingForInputRegisterOffset] = keyIndex
	e.Input.WaitingForInput = false
}

func (e *Emulator) EmulateCycle() {
	if e.Input.WaitingForInput {
		return
	}

	// Fetch instruction at program counter
	var instruction uint16 = (uint16(e.memory.RAM[e.cpu.PC]) << 8) | uint16(e.memory.RAM[e.cpu.PC+1])

	// Advance program counter
	e.cpu.PC += 2

	// Execute instruction
	switch instruction {
	case 0x00E0:
		e.op00E0(instruction)
	case 0x00EE:
		e.op00EE(instruction)
	default:
		switch byte(instruction & 0xF000 >> 12) {
		case 0x0:
			e.op0nnn(instruction)
		case 0x1:
			e.op1nnn(instruction)
		case 0x2:
			e.op2nnn(instruction)
		case 0x3:
			e.op3xkk(instruction)
		case 0x4:
			e.op4xkk(instruction)
		case 0x5:
			e.op5xy0(instruction)
		case 0x6:
			e.op6xkk(instruction)
		case 0x7:
			e.op7xkk(instruction)
		case 0x8:
			switch instruction & 0xF {
			case 0x0:
				e.op8xy0(instruction)
			case 0x1:
				e.op8xy1(instruction)
			case 0x2:
				e.op8xy2(instruction)
			case 0x3:
				e.op8xy3(instruction)
			case 0x4:
				e.op8xy4(instruction)
			case 0x5:
				e.op8xy5(instruction)
			case 0x6:
				e.op8xy6(instruction)
			case 0x7:
				e.op8xy7(instruction)
			case 0xE:
				e.op8xyE(instruction)
			default:
				panicInstructionNotImplemented(instruction)
			}
		case 0x9:
			e.op9xy0(instruction)
		case 0xA:
			e.opAnnn(instruction)
		case 0xB:
			e.opBnnn(instruction)
		case 0xC:
			e.opCxkk(instruction)
		case 0xD:
			e.opDxyn(instruction)
		case 0xE:
			switch instruction & 0xFF {
			case 0x9E:
				e.opEx9E(instruction)
			case 0xA1:
				e.opExA1(instruction)
			default:
				panicInstructionNotImplemented(instruction)
			}
		case 0xF:
			switch instruction & 0xFF {
			case 0x07:
				e.opFx07(instruction)
			case 0x0A:
				e.opFx0A(instruction)
			case 0x15:
				e.opFx15(instruction)
			case 0x18:
				e.opFx18(instruction)
			case 0x1E:
				e.opFx1E(instruction)
			case 0x29:
				e.opFx29(instruction)
			case 0x33:
				e.opFx33(instruction)
			case 0x55:
				e.opFx55(instruction)
			case 0x65:
				e.opFx65(instruction)
			default:
				panicInstructionNotImplemented(instruction)
			}
		default:
			panicInstructionNotImplemented(instruction)
		}
	}

	// Update timers
	if e.cpu.DelayTimer > 0 {
		e.cpu.DelayTimer--
	}
	if e.cpu.SoundTimer > 0 {
		e.cpu.SoundTimer--
	}
}

// 00E0 - CLS
// Clear the display.
func (e *Emulator) op00E0(instruction uint16) {
	fmt.Println("--- 00E0")
	e.Display.Clear()
}

// 00EE - RET
// Return from a subroutine.
// The interpreter sets the program counter to the address at the top of
// the stack, then subtracts 1 from the stack pointer.
func (e *Emulator) op00EE(instruction uint16) {
	fmt.Println("--- 00EE")
	e.cpu.PC = uint16(e.cpu.Stack[e.cpu.SP])
	e.cpu.SP--
}

// 0nnn - SYS addr
// Jump to a machine code routine at nnn.
// This instruction is only used on the old computers on which Chip-8 was
// originally implemented. It is ignored by modern interpreters.
func (e *Emulator) op0nnn(instruction uint16) {
	fmt.Println("--- 0nnn")
	// Do nothing.
	// This instruction is ignored by modern interpreters
}

// 1nnn - JP addr
// Jump to location nnn.
// The interpreter sets the program counter to nnn.
func (e *Emulator) op1nnn(instruction uint16) {
	fmt.Println("--- 1nnn")
	address := instruction & 0xFFF
	// fmt.Printf("jump %X\n", address)
	e.cpu.PC = address
}

// 2nnn - CALL addr
// Call subroutine at nnn.
// The interpreter increments the stack pointer, then puts the current PC
// on the top of the stack. The PC is then set to nnn.
func (e *Emulator) op2nnn(instruction uint16) {
	fmt.Println("--- 2nnn")
	address := instruction & 0xFFF
	e.cpu.SP++
	e.cpu.Stack[e.cpu.SP] = e.cpu.PC
	e.cpu.PC = address
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
func (e *Emulator) op3xkk(instruction uint16) {
	fmt.Println("--- 3xkk")
	registerOffset := byte(instruction & 0xF00 >> 8)
	value := byte(instruction & 0xFF)
	if e.cpu.V[registerOffset] == value {
		e.cpu.PC += 2
	}
}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
// The interpreter compares register Vx to kk, and if they are not equal,
// increments the program counter by 2.
func (e *Emulator) op4xkk(instruction uint16) {
	fmt.Println("--- 4xkk")
	registerOffset := byte(instruction & 0xF00 >> 8)
	value := byte(instruction & 0xFF)
	if e.cpu.V[registerOffset] != value {
		e.cpu.PC += 2
	}
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
// The interpreter compares register Vx to register Vy, and if they are equal,
// increments the program counter by 2.
func (e *Emulator) op5xy0(instruction uint16) {
	fmt.Println("--- 5xy0")
	x := byte(instruction & 0xF00 >> 8)
	y := byte(instruction & 0xF0 >> 4)
	if x == y {
		e.cpu.PC += 2
	}
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
// The interpreter puts the value kk into register Vx.
func (e *Emulator) op6xkk(instruction uint16) {
	fmt.Println("--- 6xkk")
	registerOffset := byte(instruction & 0xF00 >> 8)
	value := byte(instruction & 0xFF)
	e.cpu.V[registerOffset] = value
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
// Adds the value kk to the value of register Vx, then stores the result in Vx.
func (e *Emulator) op7xkk(instruction uint16) {
	fmt.Println("--- 7xkk")
	registerOffset := byte(instruction & 0xF00 >> 8)
	value := byte(instruction & 0xFF)
	e.cpu.V[registerOffset] += value
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
// Stores the value of register Vy in register Vx.
func (e *Emulator) op8xy0(instruction uint16) {
	fmt.Println("--- 7xkk")
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)
	e.cpu.V[registerOffsetX] = e.cpu.V[registerOffsetY]
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
// Performs a bitwise OR on the values of Vx and Vy, then stores the
// result in Vx. A bitwise OR compares the corrseponding bits from two
// values, and if either bit is 1, then the same bit in the result is
// also 1. Otherwise, it is 0.
func (e *Emulator) op8xy1(instruction uint16) {
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)

	e.cpu.V[registerOffsetX] |= e.cpu.V[registerOffsetY]
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
// Performs a bitwise AND on the values of Vx and Vy, then stores the result
// in Vx. A bitwise AND compares the corrseponding bits from two values,
// and if both bits are 1, then the same bit in the result is also 1.
// Otherwise, it is 0.
func (e *Emulator) op8xy2(instruction uint16) {
	fmt.Println("--- 8xy2")
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)
	e.cpu.V[registerOffsetX] &= e.cpu.V[registerOffsetY]
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
// Performs a bitwise exclusive OR on the values of Vx and Vy, then
// stores the result in Vx. An exclusive OR compares the corrseponding
// bits from two values, and if the bits are not both the same, then the
// corresponding bit in the result is set to 1. Otherwise, it is 0.
func (e *Emulator) op8xy3(instruction uint16) {
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)

	e.cpu.V[registerOffsetX] ^= e.cpu.V[registerOffsetY]
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
// The values of Vx and Vy are added together. If the result is greater
// than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest
// 8 bits of the result are kept, and stored in Vx.
func (e *Emulator) op8xy4(instruction uint16) {
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)

	sum := uint16(e.cpu.V[registerOffsetX]) + uint16(e.cpu.V[registerOffsetY])
	var overflowStatus byte
	if sum > 0xFFFF {
		overflowStatus = 1
	}
	e.cpu.V[0xF] = overflowStatus
	e.cpu.V[registerOffsetX] = byte(sum)
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted
// from Vx, and the results stored in Vx.
func (e *Emulator) op8xy5(instruction uint16) {
	fmt.Println("--- 8xy5")
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)

	var noBorrow byte
	if e.cpu.V[registerOffsetX] > e.cpu.V[registerOffsetY] {
		noBorrow = 1
	}
	e.cpu.V[0xF] = noBorrow
	e.cpu.V[registerOffsetX] -= e.cpu.V[registerOffsetY]
}

// 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
func (e *Emulator) op8xy6(instruction uint16) {
	fmt.Println("--- 8xy6")
	registerIndex := byte(instruction & 0xF00 >> 8)

	var lsbIsOne byte
	if (e.cpu.V[registerIndex] & 0xF) == 1 {
		lsbIsOne = 1
	}
	e.cpu.V[0xF] = lsbIsOne
	e.cpu.V[registerIndex] >>= 1
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted
// from Vy, and the results stored in Vx.
func (e *Emulator) op8xy7(instruction uint16) {
	fmt.Println("--- 8xy7")
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)

	var noBorrow byte
	if e.cpu.V[registerOffsetY] > e.cpu.V[registerOffsetX] {
		noBorrow = 1
	}
	e.cpu.V[0xF] = noBorrow

	e.cpu.V[registerOffsetX] = e.cpu.V[registerOffsetY] - e.cpu.V[registerOffsetX]
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
// If the most-significant bit of Vx is 1, then VF is set to 1,
// otherwise to 0. Then Vx is multiplied by 2.
func (e *Emulator) op8xyE(instruction uint16) {
	registerOffsetX := byte(instruction & 0xF00 >> 8)

	var msbIsOne byte
	if (e.cpu.V[registerOffsetX] >> 7) == 1 {
		msbIsOne = 1
	}
	e.cpu.V[0xF] = msbIsOne

	e.cpu.V[registerOffsetX] <<= 1
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
// The values of Vx and Vy are compared, and if they are not equal,
// the program counter is increased by 2.
func (e *Emulator) op9xy0(instruction uint16) {
	fmt.Println("--- 9xy0")
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)
	if e.cpu.V[registerOffsetX] != e.cpu.V[registerOffsetY] {
		e.cpu.PC += 2
	}
}

// Annn - LD I, addr
// Set I = nnn.
// The value of register I is set to nnn.
func (e *Emulator) opAnnn(instruction uint16) {
	fmt.Println("--- Annn")
	e.cpu.I = instruction & 0xFFF
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
// The program counter is set to nnn plus the value of V0.
func (e *Emulator) opBnnn(instruction uint16) {
	fmt.Println("--- Annn")
	address := instruction & 0xFFF

	e.cpu.PC = address + uint16(e.cpu.V[0])
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
// The interpreter generates a random number from 0 to 255, which is then
// ANDed with the value kk. The results are stored in Vx. See instruction
// 8xy2 for more information on AND.
func (e *Emulator) opCxkk(instruction uint16) {
	fmt.Println("--- Cxkk")
	registerOffset := byte(instruction & 0xF00 >> 8)
	andValue := byte(instruction & 0xFF)

	randomValue := byte(rand.Uint32() % 255)
	e.cpu.V[registerOffset] = randomValue & andValue
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
// The interpreter reads n bytes from memory, starting at the address stored in I.
// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
// Sprites are XORed onto the existing screen. If this causes any pixels to be erased,
// VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it
// is outside the coordinates of the display, it wraps around to the opposite side of the screen.
func (e *Emulator) opDxyn(instruction uint16) {
	fmt.Println("--- Dxyn")
	registerOffsetX := byte(instruction & 0xF00 >> 8)
	registerOffsetY := byte(instruction & 0xF0 >> 4)
	numRows := byte(instruction & 0xF)

	x := e.cpu.V[registerOffsetX]
	y := e.cpu.V[registerOffsetY]

	e.cpu.V[0xF] = 0
	var i byte = 0
	for ; i < numRows; i++ {
		row := e.memory.RAM[e.cpu.I+uint16(i)]

		if e.Display.DrawSprite(x, y+i, row) {
			e.cpu.V[0xF] = 1
		}
	}
}

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
// Checks the keyboard, and if the key corresponding to the value
// of Vx is currently in the down position, PC is increased by 2.
func (e *Emulator) opEx9E(instruction uint16) {
	fmt.Println("--- Ex9E")
	keyIndex := byte(instruction & 0xF00 >> 8)
	// fmt.Println("==================")
	// fmt.Println(".......................... check press ", keyIndex)
	// fmt.Println("==================")
	if e.Input.IsPressed(keyIndex) {
		// fmt.Println(".............. pressed ", keyIndex)
		// fmt.Println(".......................... ", keyIndex)
		e.cpu.PC += 2
	}
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
// Checks the keyboard, and if the key corresponding to the value of
// Vx is currently in the up position, PC is increased by 2.
func (e *Emulator) opExA1(instruction uint16) {
	fmt.Println("--- ExA1")
	keyIndex := byte(instruction & 0xF00 >> 8)
	// fmt.Println("...check NOT press ", keyIndex)
	// fmt.Println("==================")
	// fmt.Println(".......................... check NOT press ", keyIndex)
	if !e.Input.IsPressed(keyIndex) {
		// fmt.Println("================== ", keyIndex)
		e.cpu.PC += 2
	}
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
// The value of DT is placed into Vx.
func (e *Emulator) opFx07(instruction uint16) {
	fmt.Println("--- Fx07")
	offset := byte(instruction & 0xF00 >> 8)
	e.cpu.V[offset] = e.cpu.DelayTimer
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (e *Emulator) opFx0A(instruction uint16) {
	registerOffsetX := byte(instruction & 0xF00 >> 8)

	e.Input.WaitingForInput = true
	e.waitingForInputRegisterOffset = registerOffsetX
	fmt.Println("$$$$$$$$$$$$$$")
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
// DT is set equal to the value of Vx.
func (e *Emulator) opFx15(instruction uint16) {
	fmt.Println("--- Fx15")
	registerOffset := byte(instruction & 0xF00 >> 8)
	e.cpu.DelayTimer = e.cpu.V[registerOffset]
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
// ST is set equal to the value of Vx.
func (e *Emulator) opFx18(instruction uint16) {
	fmt.Println("--- Fx18")
	registerOffset := byte(instruction & 0xF00 >> 8)
	e.cpu.SoundTimer = e.cpu.V[registerOffset]
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
// The values of I and Vx are added, and the results are stored in I.
func (e *Emulator) opFx1E(instruction uint16) {
	fmt.Println("--- Fx15")
	registerOffset := byte(instruction & 0xF00 >> 8)
	e.cpu.I += uint16(e.cpu.V[registerOffset])
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
// The value of I is set to the location for the hexadecimal sprite corresponding
// to the value of Vx. See section 2.4, Display, for more information on the
// Chip-8 hexadecimal font.
func (e *Emulator) opFx29(instruction uint16) {
	fmt.Println("--- Fx29")
	registerOffset := byte(instruction & 0xF00 >> 8)
	pixelFontAddress := RamFontStart + e.cpu.V[registerOffset]
	e.cpu.I = uint16(pixelFontAddress)
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
// The interpreter takes the decimal value of Vx, and places the hundreds
// digit in memory at location in I, the tens digit at location I+1,
// and the ones digit at location I+2.
func (e *Emulator) opFx33(instruction uint16) {
	fmt.Println("--- Fx33")
	registerOffset := byte(instruction & 0xF00 >> 8)
	decimalValue := e.cpu.V[registerOffset]

	e.memory.RAM[e.cpu.I] = decimalValue / 100 //hundreds
	decimalValue -= e.memory.RAM[e.cpu.I] * 100

	e.memory.RAM[e.cpu.I+1] = decimalValue / 10 //tens
	decimalValue -= e.memory.RAM[e.cpu.I+1] * 10

	e.memory.RAM[e.cpu.I+2] = decimalValue / 1 //ones
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
// The interpreter copies the values of registers V0 through Vx into
// memory, starting at the address in I.
func (e *Emulator) opFx55(instruction uint16) {
	fmt.Println("--- Fx55")
	maxRegisterOffset := byte(instruction & 0xF00 >> 8)
	var i byte = 0
	for ; i <= maxRegisterOffset; i++ {
		e.memory.RAM[e.cpu.I+uint16(i)] = e.cpu.V[i]
	}
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
// The interpreter reads values from memory starting at location I
// into registers V0 through Vx.
func (e *Emulator) opFx65(instruction uint16) {
	fmt.Println("--- Fx65")
	maxRegisterOffset := byte(instruction & 0xF00 >> 8)
	var i byte = 0
	for ; i <= maxRegisterOffset; i++ {
		fmt.Println("maxRegisterOffset: ", maxRegisterOffset)
		e.cpu.V[i] = e.memory.RAM[e.cpu.I+uint16(i)]
	}
}

func panicInstructionNotImplemented(instruction uint16) {
	panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
}
