package main

import (
	"fmt"
)


type VM struct {
	chunk *Chunk
	stack []float64
}

type InterpretResult float64

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

var vm VM

func resetStack() {
	vm.stack = make([]float64, 8)
}

func initVM() {
	resetStack()
}

func interpret(chunk *Chunk) InterpretResult {
	vm.chunk = chunk
	return run()
}

func run() InterpretResult {
	l := len(vm.chunk.code)
	for i := 0; i < l; i++ {
		fmt.Printf("          ")
		for _, value := range vm.stack {
			fmt.Printf("[ ")
			fmt.Printf("%g", value)
			fmt.Printf(" ]")
		}
		fmt.Printf("\n")
		instruction := readByte(i)
		switch(instruction) {
			case OP_CONSTANT: {
				i++
				constant := readConstant(i)
				push(constant)
				break
			}
			case OP_ADD: {
				b := pop()
				a := pop()
				push(a + b)
				break
			}
			case OP_SUBTRACT: {
				b := pop()
				a := pop()
				push(a - b)
				break
			}
			case OP_MULTIPLY: {
				b := pop()
				a := pop()
				push(a * b)
				break
			}
			case OP_DIVIDE: {
				b := pop()
				a := pop()
				push(a / b)
				break
			}
			case OP_NEGATE: {
				push(-pop())
				break
			}
			case OP_RETURN: {
				fmt.Printf("%g", pop())
				fmt.Printf("\n")
				return INTERPRET_OK
			}
			default: {
				fmt.Printf("Unknown opcode %d\n", instruction)
				return INTERPRET_RUNTIME_ERROR
			}
		}
	}

	return INTERPRET_RUNTIME_ERROR
}

func readByte(i int) uint8 {
	return vm.chunk.code[i]
}

func readConstant(i int) float64 {
	return vm.chunk.constants[readByte(i)]
}

func push(value float64) {
	vm.stack = append(vm.stack, value)
}

func pop() float64 {
	value := vm.stack[len(vm.stack)-1]
	vm.stack[len(vm.stack)-1] = 0 //write zero value to avoid memory leak
	vm.stack = vm.stack[:len(vm.stack)-1]
	return value
}
