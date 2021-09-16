package main

import (
	"errors"
	"fmt"
	"log"
)

// character sprites used by chip8 programs
var sprites = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type cpu struct {
	mem [4096]uint8 // memory
	pc  uint16      // program counter
	v   [16]uint8   // generic registers
	//i     uint16         // special 16-bit 'index' register
	//dt    uint8          // delay timer
	//st    uint8          // sound timer
	sp    uint8      // stack pointer
	stack [16]uint16 // stack
	//keys  [16]uint8      // keyboard state
	disp    [64 * 32]uint8 // display
	noDebug bool           // print debug info
}

// set initial state, a prerequisite for all program execution
func (c *cpu) init(program []byte) {
	// load sprites into RAM
	copy(c.mem[0:], sprites)

	// load game into RAM
	copy(c.mem[0x0200:], []uint8(program))

	// set program counter
	c.pc = 0x0200

	// set stack pointer
	c.sp = 0x00
}

// fetch and execute a single opcode
func (c *cpu) cycle() bool {
	// fetch opcode
	opcode := c.fetch()
	// exec opcode
	ok, err := c.exec(opcode)
	if err != nil {
		log.Print(err)
	}
	return ok
}

// fetch the next opcode and advance the program counter
func (c *cpu) fetch() uint16 {
	// fetch opcode
	upper := uint16(c.mem[c.pc]) << 8
	lower := uint16(c.mem[c.pc+1])
	opcode := upper | lower

	// advance program counter
	c.pc += 2

	return opcode
}

// execute an opcode
func (c *cpu) exec(opcode uint16) (bool, error) {
	// decode
	family := opcode & 0xF000          // the highest 4 bits of the opcode
	nnn := opcode & 0x0FFF             // addr
	n := uint8(opcode & 0x000F)        // nibble
	x := uint8((opcode & 0x0F00) >> 8) // x operand
	y := uint8((opcode & 0x00F0) >> 4) // y operand
	kk := uint8(opcode & 0x00FF)       // byte

	// debug
	instruction := "" // generic name of instruction
	cPseudo := ""     // c pseudo code
	pc := c.pc - 2    // the address in memory whence the opcode was fetched

	// execute instruction
	switch family {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			instruction = "00E0"
			cPseudo = "clear()"
			for i := range c.disp {
				c.disp[i] = 0x00
			}
		case 0x00EE:
			instruction = "00EE"
			cPseudo = "return"
			c.sp = c.sp - 1
			c.pc = c.stack[c.sp]
			c.stack[c.sp] = 0x00
		default:
			msg := fmt.Sprintf("fatal error: unknown opcode 0x%X", opcode)
			return false, errors.New(msg)
		}
	case 0x1000:
		instruction = "1NNN"
		cPseudo = "jump"
		c.pc = nnn
	case 0x2000:
		instruction = "2NNN"
		cPseudo = "function call"
		c.stack[c.sp] = c.pc
		c.sp = c.sp + 1
		c.pc = nnn
	case 0x3000:
		instruction = "3XKK"
		cPseudo = "if v[x] == kk: continue"
		if c.v[x] == kk {
			c.pc = c.pc + 2
		}
	case 0x4000:
		instruction = "4XKK"
		cPseudo = "if v[x] != kk: continue"
		if c.v[x] != kk {
			c.pc = c.pc + 2
		}
	case 0x5000:
		switch n {
		case 0x0:
			instruction = "5XY0"
			cPseudo = "if v[x] == v[y]: continue"
			if c.v[x] == c.v[y] {
				c.pc = c.pc + 2
			}
		default:
			msg := fmt.Sprintf("fatal error: unknown opcode 0x%X", opcode)
			return false, errors.New(msg)
		}
	case 0x6000:
		instruction = "6XKK"
		cPseudo = "v[x] = kk"
		c.v[x] = kk
	case 0x7000:
		instruction = "7XKK"
		cPseudo = "v[x] = v[x] + kk"
		c.v[x] = c.v[x] + kk
	case 0x8000:
		switch n {
		case 0x0:
			instruction = "8XY0"
			cPseudo = "v[x] = v[y]"
			c.v[x] = c.v[y]
		case 0x1:
			instruction = "8XY1"
			cPseudo = "v[x] = v[x] | v[y]"
			c.v[x] = (c.v[x] | c.v[y])
		case 0x2:
			instruction = "8XY2"
			cPseudo = "v[x] = v[x] & v[y]"
			c.v[x] = (c.v[x] & c.v[y])
		case 0x3:
			instruction = "8XY3"
			cPseudo = "v[x] = v[x] ^ v[y]"
			c.v[x] = (c.v[x] ^ c.v[y])
		case 0x4:
			instruction = "8XY4"
			cPseudo = "if v[x] + v[y] > 0xFF: v[F] = 1 else: v[F] = 0; v[x] = v[x] + v[y]"
			if uint16(c.v[x])+uint16(c.v[y]) > 0xFF {
				c.v[0xF] = 0x01
			} else {
				c.v[0xF] = 0x00
			}
			c.v[x] = c.v[x] + c.v[y]
		case 0x5:
			instruction = "8XY5"
			cPseudo = "if v[x] > v[y]: v[F] = 1 else: v[F] = 0; v[x] = v[x] - v[y]"
			if c.v[x] > c.v[y] {
				c.v[0xF] = 0x01
			} else {
				c.v[0xF] = 0x00
			}
			c.v[x] = c.v[x] - c.v[y]
		case 0x6:
			instruction = "8XY6"
			cPseudo = "if v[x] & 0x01: v[F] = 1 else: v[F] = 0; v[x] = v[x] / 2"
			if c.v[x]&0x01 == 0x01 {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] = c.v[x] / 2
		case 0x7:
			instruction = "8XY7"
			cPseudo = "if v[y] > v[x]: v[F] = 1 else: v[F] = 0; v[x] = v[y] - v[x]"
			if c.v[y] > c.v[x] {
				c.v[0xF] = 0x01
			} else {
				c.v[0xF] = 0x00
			}
			c.v[x] = c.v[y] - c.v[x]
		case 0xE:
			instruction = "8XYE"
			cPseudo = "if v[x] >> 7 == 1: v[F] = 1 else: v[F] = 0; v[x] = v[x] * 2"
			if (c.v[x] >> 7) == 0x01 {
				c.v[0xF] = 0x01
			} else {
				c.v[0xF] = 0x00
			}
			c.v[x] = c.v[x] * 2
		default:
			msg := fmt.Sprintf("fatal error: unknown opcode 0x%X", opcode)
			return false, errors.New(msg)
		}
	case 0x9000:
		switch n {
		case 0x00:
			instruction = "9XY0"
			cPseudo = "if v[x] != v[y]: pc += 2"
			if c.v[x] != c.v[y] {
				c.pc = c.pc + 2
			}
		default:
			msg := fmt.Sprintf("fatal error: unknown opcode 0x%X", opcode)
			return false, errors.New(msg)
		}
	}

	if !c.noDebug {
		log.Printf(
			"opcode: 0x%X, instruction: %s, cPseudo: %s, memaddr: 0x%X",
			opcode,
			instruction,
			cPseudo,
			pc,
		)
	}

	return true, nil
}

func main() {
	c := cpu{}
	program := []byte{0x00, 0xE0, 0x80, 0x13}
	c.init(program)
	for c.cycle() {
	}
}
