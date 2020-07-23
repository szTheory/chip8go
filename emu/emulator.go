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

	nnn := instruction & 0xFFF
	n := byte(instruction & 0xF)
	x := byte(instruction & 0xF00 >> 8)
	y := byte(instruction & 0xF0 >> 4)
	kk := byte(instruction & 0xFF)

	// Execute instruction
	switch instruction {
	case 0x00E0:
		e.op00E0()
	case 0x00EE:
		e.op00EE()
	default:
		switch byte(instruction & 0xF000 >> 12) {
		case 0x0:
			e.op0nnn(nnn)
		case 0x1:
			e.op1nnn(nnn)
		case 0x2:
			e.op2nnn(nnn)
		case 0x3:
			e.op3xkk(x, kk)
		case 0x4:
			e.op4xkk(x, kk)
		case 0x5:
			e.op5xy0(x, y)
		case 0x6:
			e.op6xkk(x, kk)
		case 0x7:
			e.op7xkk(x, kk)
		case 0x8:
			switch instruction & 0xF {
			case 0x0:
				e.op8xy0(x, y)
			case 0x1:
				e.op8xy1(x, y)
			case 0x2:
				e.op8xy2(x, y)
			case 0x3:
				e.op8xy3(x, y)
			case 0x4:
				e.op8xy4(x, y)
			case 0x5:
				e.op8xy5(x, y)
			case 0x6:
				e.op8xy6(x, y)
			case 0x7:
				e.op8xy7(x, y)
			case 0xE:
				e.op8xyE(x, y)
			default:
				panicInstructionNotImplemented(instruction)
			}
		case 0x9:
			e.op9xy0(x, y)
		case 0xA:
			e.opAnnn(nnn)
		case 0xB:
			e.opBnnn(nnn)
		case 0xC:
			e.opCxkk(x, kk)
		case 0xD:
			e.opDxyn(x, y, n)
		case 0xE:
			switch instruction & 0xFF {
			case 0x9E:
				e.opEx9E(x)
			case 0xA1:
				e.opExA1(x)
			default:
				panicInstructionNotImplemented(instruction)
			}
		case 0xF:
			switch instruction & 0xFF {
			case 0x07:
				e.opFx07(x)
			case 0x0A:
				e.opFx0A(x)
			case 0x15:
				e.opFx15(x)
			case 0x18:
				e.opFx18(x)
			case 0x1E:
				e.opFx1E(x)
			case 0x29:
				e.opFx29(x)
			case 0x33:
				e.opFx33(x)
			case 0x55:
				e.opFx55(x)
			case 0x65:
				e.opFx65(x)
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
func (e *Emulator) op00E0() {
	fmt.Println("--- 00E0")
	e.Display.Clear()
}

// 00EE - RET
// Return from a subroutine.
// The interpreter sets the program counter to the address at the top of
// the stack, then subtracts 1 from the stack pointer.
func (e *Emulator) op00EE() {
	fmt.Println("--- 00EE")
	e.cpu.PC = uint16(e.cpu.Stack[e.cpu.SP])
	e.cpu.SP--
}

// 0nnn - SYS addr
// Jump to a machine code routine at nnn.
// This instruction is only used on the old computers on which Chip-8 was
// originally implemented. It is ignored by modern interpreters.
func (e *Emulator) op0nnn(addr uint16) {
	fmt.Println("--- 0nnn")
	// Do nothing.
	// This instruction is ignored by modern interpreters
}

// 1nnn - JP addr
// Jump to location nnn.
// The interpreter sets the program counter to nnn.
func (e *Emulator) op1nnn(addr uint16) {
	fmt.Println("--- 1nnn")
	e.cpu.PC = addr
}

// 2nnn - CALL addr
// Call subroutine at nnn.
// The interpreter increments the stack pointer, then puts the current PC
// on the top of the stack. The PC is then set to nnn.
func (e *Emulator) op2nnn(addr uint16) {
	fmt.Println("--- 2nnn")
	e.cpu.SP++
	e.cpu.Stack[e.cpu.SP] = e.cpu.PC
	e.cpu.PC = addr
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
func (e *Emulator) op3xkk(x, kk byte) {
	fmt.Println("--- 3xkk")
	if e.cpu.V[x] == kk {
		e.cpu.PC += 2
	}
}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
// The interpreter compares register Vx to kk, and if they are not equal,
// increments the program counter by 2.
func (e *Emulator) op4xkk(x, kk byte) {
	fmt.Println("--- 4xkk")
	if e.cpu.V[x] != kk {
		e.cpu.PC += 2
	}
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
// The interpreter compares register Vx to register Vy, and if they are equal,
// increments the program counter by 2.
func (e *Emulator) op5xy0(x, y byte) {
	fmt.Println("--- 5xy0")
	if e.cpu.V[x] == e.cpu.V[y] {
		e.cpu.PC += 2
	}
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
// The interpreter puts the value kk into register Vx.
func (e *Emulator) op6xkk(x, kk byte) {
	fmt.Println("--- 6xkk")
	e.cpu.V[x] = kk
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
// Adds the value kk to the value of register Vx, then stores the result in Vx.
func (e *Emulator) op7xkk(x, kk byte) {
	fmt.Println("--- 7xkk")
	e.cpu.V[x] += kk
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
// Stores the value of register Vy in register Vx.
func (e *Emulator) op8xy0(x, y byte) {
	fmt.Println("--- 7xkk")
	e.cpu.V[x] = e.cpu.V[y]
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
// Performs a bitwise OR on the values of Vx and Vy, then stores the
// result in Vx. A bitwise OR compares the corrseponding bits from two
// values, and if either bit is 1, then the same bit in the result is
// also 1. Otherwise, it is 0.
func (e *Emulator) op8xy1(x, y byte) {
	e.cpu.V[x] |= e.cpu.V[y]
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
// Performs a bitwise AND on the values of Vx and Vy, then stores the result
// in Vx. A bitwise AND compares the corrseponding bits from two values,
// and if both bits are 1, then the same bit in the result is also 1.
// Otherwise, it is 0.
func (e *Emulator) op8xy2(x, y byte) {
	fmt.Println("--- 8xy2")
	e.cpu.V[x] &= e.cpu.V[y]
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
// Performs a bitwise exclusive OR on the values of Vx and Vy, then
// stores the result in Vx. An exclusive OR compares the corrseponding
// bits from two values, and if the bits are not both the same, then the
// corresponding bit in the result is set to 1. Otherwise, it is 0.
func (e *Emulator) op8xy3(x, y byte) {
	e.cpu.V[x] ^= e.cpu.V[y]
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
// The values of Vx and Vy are added together. If the result is greater
// than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest
// 8 bits of the result are kept, and stored in Vx.
func (e *Emulator) op8xy4(x, y byte) {
	sum := uint16(e.cpu.V[x]) + uint16(e.cpu.V[y])

	var overflowStatus byte
	if sum > 0xFFFF {
		overflowStatus = 1
	}
	e.cpu.V[0xF] = overflowStatus

	e.cpu.V[x] = byte(sum)
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted
// from Vx, and the results stored in Vx.
func (e *Emulator) op8xy5(x, y byte) {
	fmt.Println("--- 8xy5")
	var noBorrow byte
	if e.cpu.V[x] > e.cpu.V[y] {
		noBorrow = 1
	}
	e.cpu.V[0xF] = noBorrow
	e.cpu.V[x] -= e.cpu.V[y]
}

// 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
func (e *Emulator) op8xy6(x, y byte) {
	fmt.Println("--- 8xy6")

	var lsbIsOne byte
	if (e.cpu.V[x] & 0xF) == 1 {
		lsbIsOne = 1
	}
	e.cpu.V[0xF] = lsbIsOne
	e.cpu.V[x] >>= 1
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted
// from Vy, and the results stored in Vx.
func (e *Emulator) op8xy7(x, y byte) {
	fmt.Println("--- 8xy7")
	var noBorrow byte
	if e.cpu.V[y] > e.cpu.V[x] {
		noBorrow = 1
	}
	e.cpu.V[0xF] = noBorrow

	e.cpu.V[x] = e.cpu.V[y] - e.cpu.V[x]
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
// If the most-significant bit of Vx is 1, then VF is set to 1,
// otherwise to 0. Then Vx is multiplied by 2.
func (e *Emulator) op8xyE(x, y byte) {
	var msbIsOne byte
	if (e.cpu.V[x] >> 7) == 1 {
		msbIsOne = 1
	}
	e.cpu.V[0xF] = msbIsOne

	e.cpu.V[x] <<= 1
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
// The values of Vx and Vy are compared, and if they are not equal,
// the program counter is increased by 2.
func (e *Emulator) op9xy0(x, y byte) {
	fmt.Println("--- 9xy0")
	if e.cpu.V[x] != e.cpu.V[y] {
		e.cpu.PC += 2
	}
}

// Annn - LD I, addr
// Set I = nnn.
// The value of register I is set to nnn.
func (e *Emulator) opAnnn(addr uint16) {
	fmt.Println("--- Annn")
	e.cpu.I = addr
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
// The program counter is set to nnn plus the value of V0.
func (e *Emulator) opBnnn(addr uint16) {
	fmt.Println("--- Annn")
	e.cpu.PC = addr + uint16(e.cpu.V[0])
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
// The interpreter generates a random number from 0 to 255, which is then
// ANDed with the value kk. The results are stored in Vx. See instruction
// 8xy2 for more information on AND.
func (e *Emulator) opCxkk(x, kk byte) {
	fmt.Println("--- Cxkk")
	randomValue := byte(rand.Uint32() % 255)

	e.cpu.V[x] = randomValue & kk
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
// The interpreter reads n bytes from memory, starting at the address stored in I.
// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
// Sprites are XORed onto the existing screen. If this causes any pixels to be erased,
// VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it
// is outside the coordinates of the display, it wraps around to the opposite side of the screen.
func (e *Emulator) opDxyn(x, y, n byte) {
	fmt.Println("--- Dxyn")
	xVal := e.cpu.V[x]
	yVal := e.cpu.V[y]

	e.cpu.V[0xF] = 0
	var i byte = 0
	for ; i < n; i++ {
		row := e.memory.RAM[e.cpu.I+uint16(i)]

		if e.Display.DrawSprite(xVal, yVal+i, row) {
			e.cpu.V[0xF] = 1
		}
	}
}

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
// Checks the keyboard, and if the key corresponding to the value
// of Vx is currently in the down position, PC is increased by 2.
func (e *Emulator) opEx9E(x byte) {
	fmt.Println("--- Ex9E")
	// fmt.Println("==================")
	// fmt.Println(".......................... check press ", keyIndex)
	// fmt.Println("==================")
	if e.Input.IsPressed(x) {
		// fmt.Println(".............. pressed ", keyIndex)
		// fmt.Println(".......................... ", keyIndex)
		e.cpu.PC += 2
	}
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
// Checks the keyboard, and if the key corresponding to the value of
// Vx is currently in the up position, PC is increased by 2.
func (e *Emulator) opExA1(x byte) {
	fmt.Println("--- ExA1")
	// fmt.Println("...check NOT press ", keyIndex)
	// fmt.Println("==================")
	// fmt.Println(".......................... check NOT press ", keyIndex)
	if !e.Input.IsPressed(x) {
		// fmt.Println("================== ", keyIndex)
		e.cpu.PC += 2
	}
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
// The value of DT is placed into Vx.
func (e *Emulator) opFx07(x byte) {
	fmt.Println("--- Fx07")
	e.cpu.V[x] = e.cpu.DelayTimer
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (e *Emulator) opFx0A(x byte) {
	e.Input.WaitingForInput = true
	e.waitingForInputRegisterOffset = x
	fmt.Println("$$$$$$$$$$$$$$")
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
// DT is set equal to the value of Vx.
func (e *Emulator) opFx15(x byte) {
	fmt.Println("--- Fx15")
	e.cpu.DelayTimer = e.cpu.V[x]
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
// ST is set equal to the value of Vx.
func (e *Emulator) opFx18(x byte) {
	fmt.Println("--- Fx18")
	e.cpu.SoundTimer = e.cpu.V[x]
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
// The values of I and Vx are added, and the results are stored in I.
func (e *Emulator) opFx1E(x byte) {
	fmt.Println("--- Fx15")
	e.cpu.I += uint16(e.cpu.V[x])
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
// The value of I is set to the location for the hexadecimal sprite corresponding
// to the value of Vx. See section 2.4, Display, for more information on the
// Chip-8 hexadecimal font.
func (e *Emulator) opFx29(x byte) {
	fmt.Println("--- Fx29")
	e.cpu.I = uint16(RamFontStart) + uint16(e.cpu.V[x])
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
// The interpreter takes the decimal value of Vx, and places the hundreds
// digit in memory at location in I, the tens digit at location I+1,
// and the ones digit at location I+2.
func (e *Emulator) opFx33(x byte) {
	fmt.Println("--- Fx33")
	decimalValue := e.cpu.V[x]

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
func (e *Emulator) opFx55(x byte) {
	fmt.Println("--- Fx55")
	var i byte = 0
	for ; i <= x; i++ {
		e.memory.RAM[e.cpu.I+uint16(i)] = e.cpu.V[i]
	}
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
// The interpreter reads values from memory starting at location I
// into registers V0 through Vx.
func (e *Emulator) opFx65(x byte) {
	fmt.Println("--- Fx65")
	var i byte = 0
	for ; i <= x; i++ {
		fmt.Println("maxRegisterOffset: ", x)
		e.cpu.V[i] = e.memory.RAM[e.cpu.I+uint16(i)]
	}
}

func panicInstructionNotImplemented(instruction uint16) {
	panic("Instruction not implemented 0x" + fmt.Sprintf("%X", instruction))
}
