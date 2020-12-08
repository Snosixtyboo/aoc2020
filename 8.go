package main

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type InstructionKind uint

const (
	Nop InstructionKind = 0
	Acc InstructionKind = 1
	Jmp InstructionKind = 2
)

type Instruction struct {
	kind   InstructionKind
	values []int
}

func (instruction Instruction) String() string {
	str := ""
	switch instruction.kind {
	case Nop:
		str += "nop"
	case Acc:
		str += "acc"
	case Jmp:
		str += "jmp"
	default:
		str += "NUL"
	}
	for _, v := range instruction.values {
		str += " " + strconv.Itoa(v)
	}
	return str
}

func (instruction *Instruction) execute() error {
	switch instruction.kind {
	case Nop:
	case Acc:
		Accumulator += instruction.values[0]
	case Jmp:
		InstructionPointer += (instruction.values[0] - 1)
	default:
		return errors.New("Unsupported instruction: " + instruction.String())
	}
	InstructionPointer++
	return nil
}

type ParserError int

const (
	WrongNumArgs  ParserError = 0
	FailedParsing ParserError = 1
)

func ParsingErrorString(pError ParserError, args []string) string {
	origCommand := strings.Join(args, " ")
	switch pError {
	case WrongNumArgs:
		return args[0] + ": wrong number args \"" + origCommand + "\""
	case FailedParsing:
		return args[0] + ": failed to parse \"" + origCommand + "\""
	default:
		return args[0] + ": Unknown error"
	}
}

func parseInstruction(args []string) (Instruction, error) {
	var instruction Instruction

	if len(args) != 2 {
		return instruction, errors.New(ParsingErrorString(WrongNumArgs, args))
	}
	val, err := strconv.Atoi(args[1])
	if err != nil {
		return instruction, errors.New(ParsingErrorString(FailedParsing, args))
	}

	switch args[0] {
	case "nop":
		instruction = Instruction{Nop, []int{val}}
	case "acc":
		instruction = Instruction{Acc, []int{val}}
	case "jmp":
		instruction = Instruction{Jmp, []int{val}}
	default:
		return instruction, errors.New("Unsupported instruction!")
		break
	}
	return instruction, nil
}

type LineMetric struct {
	visited int
}

var FullProgram []Instruction
var ProgramLength int
var InstructionPointer int
var Accumulator int

var Version1 = false

func main() {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	lines := strings.Split(string(bytes), "\n")

	for l, line := range lines {
		parts := strings.Split(line, " ")
		instruction, err := parseInstruction(parts)
		if err != nil {
			log.Fatalf("Line %v: %v", l, err.Error())
			return
		}
		FullProgram = append(FullProgram, instruction)
	}
	ProgramLength = len(FullProgram)

	lineMetrics := make([]LineMetric, ProgramLength)
	visited := list.New()
	changeable := list.New()
	initialRun := true
	manipulation := false

	for InstructionPointer = 0; true; {
		if InstructionPointer == ProgramLength { // Success
			break
		}

		if InstructionPointer < 0 || InstructionPointer > ProgramLength {
			log.Fatal("Invalid instruction pointer!")
			return
		}

		lineMetrics[InstructionPointer].visited++

		if lineMetrics[InstructionPointer].visited > 1 {
			initialRun = false
			log.Printf("%d > run : %v --> already visited \n", InstructionPointer, FullProgram[InstructionPointer])
			if Version1 {
				break
			} else {
				toChange := changeable.Back().Value.(int)
				changeable.Remove(changeable.Back())

				for true {
					undo := visited.Back().Value.(int)
					if FullProgram[undo].kind == Acc {
						undoInstruction := FullProgram[undo]
						undoInstruction.values[0] *= -1
						undoInstruction.execute()
					}
					visited.Remove(visited.Back())
					lineMetrics[undo].visited--

					if undo == toChange {
						break
					}
				}

				log.Printf("Manipulating line %d", toChange)
				manipulation = true
				InstructionPointer = toChange
			}
		}

		instruction := FullProgram[InstructionPointer]
		if (instruction.kind == Nop || instruction.kind == Jmp) && initialRun {
			changeable.PushBack(InstructionPointer)
		}
		visited.PushBack(InstructionPointer)

		if manipulation {
			log.Printf("Manipulating %d...\n", InstructionPointer)
			if instruction.kind == Nop {
				instruction.kind = Jmp
			} else {
				instruction.kind = Nop
			}
			manipulation = false
		}
		log.Printf("%v > run : %v\n", InstructionPointer, instruction)
		instruction.execute()
	}
	fmt.Printf("Instruction Pointer: %v of %d\n", InstructionPointer, ProgramLength)
	fmt.Printf("Accumulator: %v\n", Accumulator)
}
