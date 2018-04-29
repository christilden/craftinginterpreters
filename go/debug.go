package main

import (
	"fmt"
)


func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, chunk *Chunk, offset int) int {
	constant := chunk.code[offset + 1]
	fmt.Printf("%-16s %4d '", name, constant)
	fmt.Printf("%g", chunk.constants[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func disassembleChunk(chunk *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	i := 0
	for (i < len(chunk.code)) {
		i = disassembleInstruction(chunk, i)
	}
}

func disassembleInstruction(chunk *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if (offset > 0 && chunk.lines[offset] == chunk.lines[offset - 1]) {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", chunk.lines[offset])
	}

	instruction := chunk.code[offset]
	switch (instruction) {
		case OP_CONSTANT:
			return constantInstruction("OP_CONSTANT", chunk, offset)
		case OP_ADD:
			return simpleInstruction("OP_ADD", offset)
		case OP_SUBTRACT:
			return simpleInstruction("OP_SUBTRACT", offset)
		case OP_MULTIPLY:
			return simpleInstruction("OP_MULTIPLY", offset)
		case OP_DIVIDE:
			return simpleInstruction("OP_DIVIDE", offset)
		case OP_NEGATE:
			return simpleInstruction("OP_NEGATE", offset)
		case OP_RETURN:
			return simpleInstruction("OP_RETURN", offset)
		default:
			fmt.Printf("Unknown opcode %d\n", instruction)
			return offset + 1
	}
}
