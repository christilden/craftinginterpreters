package main

import (
	"fmt"
)


type VM struct {
	chunk *Chunk
	stack []Value
}

type InterpretResult uint8

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

var vm VM

func resetStack() {
	vm.stack = make([]Value, 8)
}

func initVM() {
	resetStack()
}

func interpret(source *string) InterpretResult {
	compile(source)
	return INTERPRET_OK
}

func run() InterpretResult {
	l := vm.CodeLength()
	for i := 0; i < l; i++ {
		traceExecution(i)
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

func (vm *VM) CodeLength() int {
	return len(vm.chunk.code)
}

func readByte(offset int) uint8 {
	return vm.chunk.code[offset]
}

func readConstant(offset int) Value {
	return vm.chunk.constants[readByte(offset)]
}

func push(value Value) {
	vm.stack = append(vm.stack, value)
}

func pop() Value {
	key := len(vm.stack) - 1
	value := vm.stack[key]
	vm.stack[key] = 0 //write zero value to avoid memory leak
	vm.stack = vm.stack[:key]
	return value
}

func traceExecution(offset int) {
	fmt.Printf("          ")
	for _, value := range vm.stack {
		fmt.Printf("[ %g ]", value)
	}
	fmt.Printf("\n")
	disassembleInstruction(vm.chunk, offset)
}
