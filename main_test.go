package main

import (
	"testing"
)

func mockAllOnDisplay() [32][64]uint8 {
	disp := [32][64]uint8{}
	for i := 0;  i < 32; i++ {
		for j := 0; j < 64; j++ {
			disp[i][j] = 1
		}
	}
	return disp
}

func TestExec(t *testing.T) {
	cases := []struct {
		desc     string
		opcode   uint16
		cpu      cpu
		expected cpu
	}{
		{
			"00E0",
			0x00E0,
			cpu{disp: mockAllOnDisplay()},
			cpu{},
		},
		{
			"00EE",
			0x00EE,
			cpu{
				pc:    0x0222,
				sp:    0x02,
				stack: [16]uint16{0x0444, 0x0333, 0x0000},
			},
			cpu{
				pc:    0x0333,
				sp:    0x01,
				stack: [16]uint16{0x0444, 0x0000},
			},
		},
		{
			"1NNN",
			0x1333,
			cpu{pc: 0x0222},
			cpu{pc: 0x0333},
		},
		{
			"2NNN",
			0x2333,
			cpu{
				pc:    0x0222,
				sp:    0x00,
				stack: [16]uint16{},
			},
			cpu{
				pc:    0x0333,
				sp:    0x01,
				stack: [16]uint16{0x0222},
			},
		},
		{
			"3XKK",
			0x32FF,
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0x00, 0xFF},
			},
			cpu{
				pc: 0x0224,
				v:  [16]uint8{0x00, 0x00, 0xFF},
			},
		},
		{
			"3XKK",
			0x32FF,
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0x00, 0xEE},
			},
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0x00, 0xEE},
			},
		},
		{
			"4XKK",
			0x42FF,
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0x00, 0xEE},
			},
			cpu{
				pc: 0x0224,
				v:  [16]uint8{0x00, 0x00, 0xEE},
			},
		},
		{
			"4XKK",
			0x42FF,
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0x00, 0xFF},
			},
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0x00, 0xFF},
			},
		},
		{
			"5XY0",
			0x5120,
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0xFF, 0xFF},
			},
			cpu{
				pc: 0x0224,
				v:  [16]uint8{0x00, 0xFF, 0xFF},
			},
		},
		{
			"5XY0",
			0x5120,
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0xEE, 0xFF},
			},
			cpu{
				pc: 0x0222,
				v:  [16]uint8{0x00, 0xEE, 0xFF},
			},
		},
		{
			"6XKK",
			0x60AB,
			cpu{v: [16]uint8{}},
			cpu{v: [16]uint8{0xAB}},
		},
		{
			"7XKK",
			0x7012,
			cpu{v: [16]uint8{0x35}},
			cpu{v: [16]uint8{0x47}},
		},
		{
			"8XY0",
			0x8120,
			cpu{v: [16]uint8{0x00, 0x35, 0x47}},
			cpu{v: [16]uint8{0x00, 0x47, 0x47}},
		},
		{
			"8XY1",
			0x8121,
			cpu{v: [16]uint8{0x00, 0x01, 0x02}},
			cpu{v: [16]uint8{0x00, 0x03, 0x02}},
		},
		{
			"8XY2",
			0x8122,
			cpu{v: [16]uint8{0x00, 0x03, 0x02}},
			cpu{v: [16]uint8{0x00, 0x02, 0x02}},
		},
		{
			"8XY3",
			0x8123,
			cpu{v: [16]uint8{0x00, 0x03, 0x02}},
			cpu{v: [16]uint8{0x00, 0x01, 0x02}},
		},
		{
			"8XY4",
			0x8014,
			cpu{
				v: [16]uint8{
					0x01, 0x02, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x03, 0x02, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			"8XY4",
			0x8014,
			cpu{
				v: [16]uint8{
					0xFF, 0x02, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x01, 0x02, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x01,
				},
			},
		},
		{
			"8XY5",
			0x8015,
			cpu{
				v: [16]uint8{
					0xFF, 0x02, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0xFD, 0x02, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x01,
				},
			},
		},
		{
			"8XY5",
			0x8015,
			cpu{
				v: [16]uint8{
					0x02, 0xFF, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x03, 0xFF, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			"8XY6",
			0x8DE6,
			cpu{
				v: [16]uint8{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x0F, 0x00, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x07, 0x00, 0x01,
				},
			},
		},
		{
			"8XY7",
			0x8DE7,
			cpu{
				v: [16]uint8{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x01, 0x0A, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x09, 0x0A, 0x01,
				},
			},
		},
		{
			"8XY7",
			0x8DE7,
			cpu{
				v: [16]uint8{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x0A, 0x01, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0xF7, 0x01, 0x00,
				},
			},
		},
		{
			"8XYE",
			0x800E,
			cpu{
				v: [16]uint8{
					0x80, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x01,
				},
			},
		},
		{
			"8XYE",
			0x801E,
			cpu{
				v: [16]uint8{
					0x30, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				v: [16]uint8{
					0x60, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			"9XY0",
			0x9010,
			cpu{
				pc: 0x222,
				v: [16]uint8{
					0x30, 0x30, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				pc: 0x222,
				v: [16]uint8{
					0x30, 0x30, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			"9XY0",
			0x9010,
			cpu{
				pc: 0x222,
				v: [16]uint8{
					0x30, 0x50, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				pc: 0x224,
				v: [16]uint8{
					0x30, 0x50, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			"ANNN",
			0xA123,
			cpu{},
			cpu{i: 0x0123},
		},
		{
			"BNNN",
			0xB123,
			cpu{
				pc: 0x0333,
				v: [16]uint8{
					0x30, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				pc: 0x0153,
				v: [16]uint8{
					0x30, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		// TODO: test CNNN (random mask)
		{
			"DXYN",
			0xD005,
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0,
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x0A,
				disp: [32][64]uint8{},
			},
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0,
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x0A,
				disp: [32][64]uint8{
					{1, 1, 1, 1},
					{0, 0, 0, 1},
					{1, 1, 1, 1},
					{1, 0, 0, 0},
					{1, 1, 1, 1},
				},
				v: [16]uint8{
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
				},
			},
		},
		{
			"DXYN",
			0xD005,
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x00,
				disp: [32][64]uint8{
					{1, 1, 1, 1},
					{0, 1, 1, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
				},
			},
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x00,
				disp: [32][64]uint8{
					{0, 0, 0, 0},
					{1, 1, 1, 1},
					{1, 0, 0, 1},
					{1, 0, 0, 1},
					{1, 1, 1, 1},
				},
				v: [16]uint8{
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x1,
				},
			},
		},
		{
			"DXYN",
			0xD015,
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x00,
				v: [16]uint8{
					0x3E, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x00,
				disp: [32][64]uint8{
					{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
					{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
					{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
					{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
					{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
				},
				v: [16]uint8{
					0x3E, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0,
				},
			},
		},
		{
			"DXYN",
			0xD015,
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x00,
				disp: [32][64]uint8{},
				v: [16]uint8{
					0x00, 0x1E, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0,
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0x00,
				disp: [32][64]uint8{
					{1, 0, 0, 1},
					{1, 0, 0, 1},
					{1, 1, 1, 1},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{0, 0, 0, 0},
					{1, 1, 1, 1},
					{1, 0, 0, 1},
				},
				v: [16]uint8{
					0x00, 0x1E, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			"EX9E",
			0xE09E,
			cpu{
				pc: 0x222,
				v: [16]uint8{
					2, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				keys: [16]uint8{
					0, 0, 1, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
			},
			cpu{
				pc: 0x224,
				v: [16]uint8{
					2, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				keys: [16]uint8{
					0, 0, 1, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
			},
		},
		{
			"EXA1",
			0xE0A1,
			cpu{
				pc: 0x222,
				v: [16]uint8{
					2, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				keys: [16]uint8{
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
			},
			cpu{
				pc: 0x224,
				v: [16]uint8{
					2, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				keys: [16]uint8{
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
			},
		},
		{
			"FX07",
			0xF107,
			cpu{
				dt: 9,
			},
			cpu{
				dt: 9,
				v: [16]uint8{
					0, 9, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
			},
		},
		// TODO: test FX0A (halt until keypress)
		{
			"FX15",
			0xF015,
			cpu{
				v: [16]uint8{
					1, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				dt: 0,
			},
			cpu{
				v: [16]uint8{
					1, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				dt: 1,
			},
		},
		{
			"FX18",
			0xF018,
			cpu{
				v: [16]uint8{
					3, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				st: 0,
			},
			cpu{
				v: [16]uint8{
					3, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				st: 3,
			},
		},
		{
			"FX1E",
			0xF01E,
			cpu{
				v: [16]uint8{
					3, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				i: 1,
			},
			cpu{
				v: [16]uint8{
					3, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
					0, 0, 0, 0,
				},
				i: 4,
			},
		},
		{
			"FX29",
			0xF129,
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0,
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 0,
				v: [16]uint8{
					0x00, 0x09, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
			cpu{
				mem: [4096]uint8{
					0xF0, 0x90, 0x90, 0x90, 0xF0,
					0x20, 0x60, 0x20, 0x20, 0x70,
					0xF0, 0x10, 0xF0, 0x80, 0xF0,
					0xF0, 0x10, 0xF0, 0x10, 0xF0,
					0x90, 0x90, 0xF0, 0x10, 0x10,
					0xF0, 0x80, 0xF0, 0x10, 0xF0,
					0xF0, 0x80, 0xF0, 0x90, 0xF0,
					0xF0, 0x10, 0x20, 0x40, 0x40,
					0xF0, 0x90, 0xF0, 0x90, 0xF0,
					0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
					0xF0, 0x90, 0xF0, 0x90, 0x90,
					0xE0, 0x90, 0xE0, 0x90, 0xE0,
					0xF0, 0x80, 0x80, 0x80, 0xF0,
					0xE0, 0x90, 0x90, 0x90, 0xE0,
					0xF0, 0x80, 0xF0, 0x80, 0xF0,
					0xF0, 0x80, 0xF0, 0x80, 0x80,
				},
				i: 45,
				v: [16]uint8{
					0x00, 0x09, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			"",
			0x0000,
			cpu{},
			cpu{},
		},
	}

	for _, tc := range cases {
		tc.cpu.exec(tc.opcode)
		// mem
		for i := range tc.cpu.mem {
			if tc.cpu.mem[i] != tc.expected.mem[i] {
				t.Fatalf(
					"fatal memory error for %s: expected 0x%X, got 0x%X at mem index 0x%X",
					tc.desc,
					tc.expected.mem[i],
					tc.cpu.mem[i],
					i,
				)
			}
		}
		// pc
		if tc.cpu.pc != tc.expected.pc {
			t.Fatalf(
				"fatal program counter error for %s: expected 0x%X, got 0x%X",
				tc.desc,
				tc.expected.pc,
				tc.cpu.pc,
			)
		}
		// v
		for i := range tc.cpu.v {
			if tc.cpu.v[i] != tc.expected.v[i] {
				t.Fatalf(
					"fatal register error for %s: expected 0x%X, got 0x%X at register index 0x%X",
					tc.desc,
					tc.expected.v[i],
					tc.cpu.v[i],
					i,
				)
			}
		}
		// i
		if tc.cpu.i != tc.expected.i {
			t.Fatalf(
				"fatal i-reg error for %s: expected 0x%X, got 0x%X",
				tc.desc,
				tc.expected.i,
				tc.cpu.i,
			)
		}
		// dt
		if tc.cpu.dt != tc.expected.dt {
			t.Fatalf(
				"fatal delay timer error for %s: expected 0x%X, got 0x%X",
				tc.desc,
				tc.expected.dt,
				tc.cpu.dt,
			)
		}
		// st
		if tc.cpu.st != tc.expected.st {
			t.Fatalf(
				"fatal sound timer error for %s: expected 0x%X, got 0x%X",
				tc.desc,
				tc.expected.st,
				tc.cpu.st,
			)
		}
		// sp
		if tc.cpu.sp != tc.expected.sp {
			t.Fatalf(
				"fatal stack pointer error for %s: expected 0x%X, got 0x%X",
				tc.desc,
				tc.expected.sp,
				tc.cpu.sp,
			)
		}
		// stack
		for i := range tc.cpu.stack {
			if tc.cpu.stack[i] != tc.expected.stack[i] {
				t.Fatalf(
					"fatal stack error for %s: expected 0x%X, got 0x%X",
					tc.desc,
					tc.expected.stack[i],
					tc.cpu.stack[i],
				)
			}
		}
		// keys
		for i := range tc.cpu.keys {
			if tc.cpu.keys[i] != tc.expected.keys[i] {
				t.Fatalf(
					"fatal display error for %s: expected %d, got %d at key %X",
					tc.desc,
					tc.cpu.keys[i],
					tc.expected.keys[i],
					i,
				)
			}
		}
		// disp
		for i := 0; i < 32; i++ {
			for j := 0; j < 64; j++ {
				if tc.cpu.disp[i][j] != tc.expected.disp[i][j] {
					t.Fatalf(
						"fatal display error for %s: expected %d, got %d at pixel (%d,%d)",
						tc.desc,
						tc.expected.disp[i][j],
						tc.cpu.disp[i][j],
						j,
						i,
					)
				}
			}
		}
	}
}
