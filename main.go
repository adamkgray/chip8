package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"io/ioutil"
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
	mem 	[4096]uint8		// memory
	pc  	uint16			// programme counter
	v 		[16]uint8		// general registers
	i 		uint16			// special 'i' register
	dt 		uint8			// delay timer
	st 		uint8			// sound timer
	sp 		uint8			// stack pointer
	stack 	[16]uint16		// stack
	keys 	[16]uint8		// keyboard
	disp 	[32][64]uint8 	// display
	noDebug bool			// debug switch
}

// set initial state, prerequisite for all program execution
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

// fetch and execute single opcode
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

// fetch next opcode and advance program counter
func (c *cpu) fetch() uint16 {
	// fetch opcode
	upper := uint16(c.mem[c.pc]) << 8
	lower := uint16(c.mem[c.pc+1])
	opcode := upper | lower

	// advance program counter
	c.pc += 2

	return opcode
}

// execute opcode
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
			for i := 0; i < 32; i ++ {
				for j := 0; j < 64; j ++ {
					c.disp[i][j] = 0x00
				}
			}
		case 0x00EE:
			instruction = "00EE"
			cPseudo = "return"
			c.sp -= 1
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
		c.sp += 1
		c.pc = nnn
	case 0x3000:
		instruction = "3XKK"
		cPseudo = "if v[x] == kk: continue"
		if c.v[x] == kk {
			c.pc += 2
		}
	case 0x4000:
		instruction = "4XKK"
		cPseudo = "if v[x] != kk: continue"
		if c.v[x] != kk {
			c.pc += 2
		}
	case 0x5000:
		switch n {
		case 0x0:
			instruction = "5XY0"
			cPseudo = "if v[x] == v[y]: continue"
			if c.v[x] == c.v[y] {
				c.pc += 2
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
			c.v[x] += c.v[y]
		case 0x5:
			instruction = "8XY5"
			cPseudo = "if v[x] > v[y]: v[F] = 1 else: v[F] = 0; v[x] = v[x] - v[y]"
			if c.v[x] > c.v[y] {
				c.v[0xF] = 0x01
			} else {
				c.v[0xF] = 0x00
			}
			c.v[x] -= c.v[y]
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
			cPseudo = "if v[x] != v[y]: pc = pc + 2"
			if c.v[x] != c.v[y] {
				c.pc += 2
			}
		default:
			msg := fmt.Sprintf("fatal error: unknown opcode 0x%X", opcode)
			return false, errors.New(msg)
		}
	case 0xA000:
		instruction = "ANNN"
		cPseudo = "i = nnn"
		c.i = nnn
	case 0xB000:
		instruction = "BNNN"
		cPseudo = "pc = v[0] + nnn"
		c.pc = uint16(c.v[0x0]) + nnn
	case 0xC000: // TODO: unit test
		instruction = "CNNN"
		cPseudo = "v[x] = rand-byte & kk"
		c.v[x] = uint8(rand.Uint32()) & kk
	case 0xD000:
		instruction = "DXYN"
		cPseudo = "/* write n-rows of sprite to disp */"

		// assume no pixels will be erased
		c.v[0xF] = 0x00

		// iterate through sprite rows
		var rows uint8
		for rows = 0; rows < n; rows++ {
			// iterate through bits of sprite
			var cols uint8
			for cols = 0; cols < 8; cols++ {
				// handle x wrap
				dispX := c.v[x] + cols
				if dispX >= 64 {
					dispX -= 64
				}

				// handle y wrap
				dispY := c.v[y] + rows
				if dispY >= 32 {
					dispY -= 32
				}

				// was the pixel on?
				pixelWasOn := c.disp[dispY][dispX] > 0

				// write to display
				// how?
				// get the sprite row from memory
				// bit shift it to the left for the correct pixel
				// mask it with 0x80 to get only the leftmost bit
				// shift that bit all the way back to the right to get a 1 or 0
				pixel := ((uint8(c.mem[c.i + uint16(rows)])<<cols) & 0x80)>>0x07
				c.disp[dispY][dispX] = c.disp[dispY][dispX] ^ pixel

				// is the pixel now off?
				pixelNowOff := c.disp[dispY][dispX] == 0

				// flag VF if any pixels were erased

				if pixelWasOn && pixelNowOff {
					c.v[0xF] = 0x01
				}

			}
		}
		fmt.Print("\033[H\033[2J")
		for row := range c.disp {
			for col := range c.disp[row] {
				if c.disp[row][col] == 1 {
					fmt.Printf("██")
				} else {
					fmt.Printf("  ")
				}
			}
			fmt.Print("\n")
		}
	case 0xE000:
		switch kk {
		case 0x9E:
			instruction = "EX9E"
			cPseudo = "if keys[v[x]] == PRESSED: pc += 2"
			if c.keys[c.v[x]] == 1 {
				c.pc += 2
			}
		case 0xA1:
			instruction = "EX9E"
			cPseudo = "if keys[v[x]] == !PRESSED: pc += 2"
			if c.keys[c.v[x]] == 0 {
				c.pc += 2
			}
		default:
			msg := fmt.Sprintf("fatal error: unknown opcode 0x%X", opcode)
			return false, errors.New(msg)
		}
	case 0xF000:
		switch kk {
		case 0x07:
			instruction = "FX07"
			cPseudo = "v[x] = dt"
			c.v[x] = c.dt
		case 0x0A:
			instruction = "FX0A"
			cPseudo = "v[x] = getkey()"
			c.v[x] = c.getkey()
		case 0x15:
			instruction = "FX15"
			cPseudo = "dt = v[x]"
			c.dt = c.v[x]
		case 0x18:
			instruction = "FX18"
			cPseudo = "st = v[x]"
			c.st = c.v[x]
		case 0x1E:
			instruction = "FX1E"
			cPseudo = "i += v[x]"
			c.i += uint16(c.v[x])
		case 0x29:
			instruction = "FX29"
			cPseudo = "i = &SPRITE(v[x])"
			c.i = uint16(5 * c.v[x])
		case 0x33:
			instruction = "FX33"
			cPseudo = "mem[i], mem[i+1], mem[i+2] = BCD(v[x])"
			c.mem[c.i] 		= c.v[x] / 100
			c.mem[c.i + 1]	= (c.v[x] % 100) / 10
			c.mem[c.i + 2]	= ((c.v[x] % 100) % 10) / 1
		case 0x55:
			instruction = "FX55"
			cPseudo = "mem[i:i+x] = v[0:x]"
			var j uint8
			for j = 0; j <= x; j++ {
				c.mem[c.i + uint16(j)] = c.v[j]
			}
		case 0x65:
			instruction = "FX65"
			cPseudo = "v[0:x] = mem[i:i+x]"
			var j uint8
			for j = 0; j <= x; j++ {
				c.v[j] = c.mem[c.i + uint16(j)]
			}
		}
	}

	if c.noDebug {
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

// halt until any key is pressed, return key value
func (c *cpu) getkey() uint8 {
	// TODO: get actual key press
	return 0
}

func main() {
	c := &cpu{}
	program, _ := ioutil.ReadFile("Maze (alt) [David Winter, 199x].ch8")
	c.init(program)
	for c.cycle() {
	}
}
