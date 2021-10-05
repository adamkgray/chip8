# CHIP8

A Chip-8 emulator written in Go

## About

I made this project to learn the basics of emulators. I wanted to execute some opcodes, manage a stack, learn how to handle input, and output to anywhere. The final product is not perfect, but in terms of my learning goals I am satisfied. I have enough knowledge to go forward and try my hand at a GameBoy emulator - which of course was th real goal all along.

## What I learned

1. Making the virtual machine is the easy part. There are only 36 opcodes, and they are 99% functional. Writing unit tests, therefor, is easy. The stack is just like the stack you learn a CS101, and the memory is just an array. If you think you did it right - you probably did!

2. The hard part is everything else. How do you make sound? How to you make a screen? How to you read keydown and keyup events? I was determined to do this in pure go, as I didn't want to splelunk into a cave of low-level system depencies. Unfortunately, I was unable to find a pure go library that would handle all of this for me. Luckily, when researching other implementations of Chip-8 (of which there are many), the go-to library is `sdl`. So while this c/c++ package may have archaic looking interfaces, the benefit of using it is that everyone else already is. So when I go on to make a GameBoy emulator, I will pick `sdl`. I have no desire to learn it completely now, after all, actually playing the Chip-8 isn't all that fun. But when I need those features, I know where to look.

3. It is possible to do things you don't know a lot about! With enough *determination* you can do anything!