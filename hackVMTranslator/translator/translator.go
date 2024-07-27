package translator

import (
	"fmt"
	"nandtotetris/vmtranslator/parser"
	"strings"
)

type Translator struct {
	eqInstructionsCount int
	ltInstructionsCount int
	gtInstructionsCount int
}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t *Translator) GetInitializationInstructions() string {
	return t.initialInstructions()
}

func (t *Translator) Translate(parsedInstruction *parser.ParsedInstruction) string {
	command := parsedInstruction.Command()

	switch strings.ToLower(command) {
	case "pop":
		return t.translatePOPInstruction(parsedInstruction)
	case "push":
		return t.translatePUSHInstruction(parsedInstruction)
	case "add":
		return t.translateADDInstruction(parsedInstruction)
	case "sub":
		return t.translateSUBInstruction(parsedInstruction)
	case "neg":
		return t.translateNEGInstruction(parsedInstruction)
	case "and":
		return t.translateANDInstruction(parsedInstruction)
	case "or":
		return t.translateORInstruction(parsedInstruction)
	case "not":
		return t.translateNOTInstruction(parsedInstruction)
	case "eq":
		t.eqInstructionsCount += 1
		return t.translateEQInstruction(parsedInstruction)
	case "gt":
		t.gtInstructionsCount += 1
		return t.translateGTInstruction(parsedInstruction)
	case "lt":
		t.ltInstructionsCount += 1
		return t.translateLTInstruction(parsedInstruction)
	default:
		// log.Fatalln("unkonwn vm command")
		return ""
	}
}

func (t *Translator) GetTerminalInstrucitons() string {
	return `
(END)
@END
0;JMP`
}

func (t *Translator) initialInstructions() string {
	// initialize virtual registers
	var assemblyInstructions string
	// initilaize SP register
	assemblyInstructions += initializeRegister("SP", 0)
	// initialize LCL register
	assemblyInstructions += initializeRegister("LCL", 2050)
	// initialize ARG register
	assemblyInstructions += initializeRegister("ARG", 3050)
	// initialize ARG register
	assemblyInstructions += initializeRegister("THIS", 4050)
	// initialize ARG register
	assemblyInstructions += initializeRegister("THAT", 5050)
	return assemblyInstructions
}

func (t *Translator) translatePOPInstruction(instruction *parser.ParsedInstruction) string {
	segment := getRegisterNameFromSegmentName(instruction)
	index := getIndexvalue(instruction)
	originalVMInstruction := instruction.OriginalInstruction()

	if segment == "static" {
		return fmt.Sprintf(
			`
// %s
%s
A=M
D=M
@%s.%s
M=D
`,
			originalVMInstruction,
			decrementStackPointer(),
			instruction.FileName(),
			index,
		)
	}

	if segment == "pointer" {
		if index == "0" {
			return fmt.Sprintf(`
// %s
%s
A=M
D=M
@THIS
M=D
			`, originalVMInstruction, decrementStackPointer())
		}
		if index == "1" {
			return fmt.Sprintf(`
// %s
%s
A=M
D=M
@THAT
M=D
			`, originalVMInstruction, decrementStackPointer())
		}
	}

	if segment == "TEMP" {
		return fmt.Sprintf(`
// %s
@5
D=A
@%s
D=D+A
@R13
M=D
%s
A=M
D=M
@R13
A=M
M=D
		`,
			originalVMInstruction,
			index,
			decrementStackPointer(),
		)
	}

	return fmt.Sprintf(`
// %s
%s
%s
A=M
D=M
@R13
A=M
M=D
`,
		originalVMInstruction,
		getAndStoreSegmentAddress(segment, index, "R13"),
		decrementStackPointer(),
	)
}

func (t *Translator) translatePUSHInstruction(instruction *parser.ParsedInstruction) string {
	segment := getRegisterNameFromSegmentName(instruction)
	index := getIndexvalue(instruction)
	originalVMInstruction := instruction.OriginalInstruction()

	if segment == "constant" {
		return fmt.Sprintf(
			`
// %s
@%s
D=A
@SP
A=M
M=D
@SP
M=M+1
`,
			originalVMInstruction,
			index,
		)
	}

	if segment == "pointer" {
		if index == "0" {
			return fmt.Sprintf(`
// %s
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
			`, originalVMInstruction)
		}
		if index == "1" {
			return fmt.Sprintf(`
// %s
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1
			`, originalVMInstruction)
		}
	}

	if segment == "static" {
		return fmt.Sprintf(
			`
// %s
@%s.%s
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
			originalVMInstruction,
			instruction.FileName(),
			index,
		)
	}

	if segment == "TEMP" {
		return fmt.Sprintf(
			`
// %s
@5
D=A
@%s
D=D+A
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1
		`,
			originalVMInstruction,
			index,
		)
	}

	return fmt.Sprintf(
		`
// %s
@%s
D=M
@%s
D=D+A
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1
`,
		originalVMInstruction,
		segment,
		index,
	)

}

func (t *Translator) translateADDInstruction(instruction *parser.ParsedInstruction) string {
	return doubleOperandOperationInstrucitons(instruction, "M=M+D")
}

func (t *Translator) translateSUBInstruction(instruction *parser.ParsedInstruction) string {
	return doubleOperandOperationInstrucitons(instruction, "M=M-D")
}

func (t *Translator) translateNEGInstruction(instruction *parser.ParsedInstruction) string {
	return singleOperandOperationInstrucitons(instruction, "M=-M")
}

func (t *Translator) translateEQInstruction(instruction *parser.ParsedInstruction) string {
	return comparisonOperation(instruction, "EQ", t.eqInstructionsCount, "JEQ")
}

func (t *Translator) translateGTInstruction(instruction *parser.ParsedInstruction) string {
	return comparisonOperation(instruction, "GT", t.gtInstructionsCount, "JGT")
}

func (t *Translator) translateLTInstruction(instruction *parser.ParsedInstruction) string {
	return comparisonOperation(instruction, "LT", t.ltInstructionsCount, "JLT")
}

func (t *Translator) translateANDInstruction(instruction *parser.ParsedInstruction) string {
	return doubleOperandOperationInstrucitons(instruction, "M=M&D")
}

func (t *Translator) translateORInstruction(instruction *parser.ParsedInstruction) string {
	return doubleOperandOperationInstrucitons(instruction, "M=M|D")
}

func (t *Translator) translateNOTInstruction(instruction *parser.ParsedInstruction) string {
	return singleOperandOperationInstrucitons(instruction, "M=!M")
}
