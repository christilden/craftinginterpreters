package main

const (
	OP_CONSTANT uint8 = iota
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_RETURN
)

type Chunk struct {
	code []uint8
	lines []int
	constants []float64
}

func initChunk(chunk *Chunk) {
	chunk = new(Chunk)
}

func writeChunk(chunk *Chunk, byte uint8, line int) {
	chunk.code = append(chunk.code, byte)
	chunk.lines = append(chunk.lines, line)
}

func addConstant(chunk *Chunk, value float64) uint8 {
	chunk.constants = append(chunk.constants, value)
	return uint8(len(chunk.constants) - 1)
}
