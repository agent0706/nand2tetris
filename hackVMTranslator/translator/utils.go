package translator

import (
	"fmt"
	"log"
	"nandtotetris/vmtranslator/parser"
)

func initializeRegister(registerName string, value int) string {
	instructions := fmt.Sprintf(`
// initiliazing %s register
@%d
D=A
@%s
M=D
	`, registerName, value, registerName)
	return instructions
}

func decrementStackPointer() string {
	return `@SP
M=M-1`
}

func getAndStoreSegmentAddress(segment string, index string, registerToStoreIn string) string {
	return fmt.Sprintf(
		`@%s
D=M
@%s
D=D+A
@%s
M=D
`,
		segment,
		index,
		registerToStoreIn,
	)
}

var (
	segmentMap map[string]string = map[string]string{
		"local":    "LCL",
		"argument": "ARG",
		"this":     "THIS",
		"that":     "THAT",
		"temp":     "TEMP",
		"constant": "constant",
		"static":   "static",
		"pointer":  "pointer",
	}
)

func doubleOperandOperationInstrucitons(instruction *parser.ParsedInstruction, operation string) string {
	return fmt.Sprintf(`
//%s
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
%s
@SP
M=M+1
`,
		instruction.OriginalInstruction(),
		operation,
	)
}

func singleOperandOperationInstrucitons(instruction *parser.ParsedInstruction, operation string) string {
	return fmt.Sprintf(`
//%s
@SP
M=M-1
A=M
%s
@SP
M=M+1
`,
		instruction.OriginalInstruction(),
		operation,
	)
}

func comparisonOperation(instruction *parser.ParsedInstruction, operation string, count int, jumpCommand string) string {
	return fmt.Sprintf(`
// %s
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@%s_%d
D;%s
@SP
A=M-1
M=0
@%s_%d_END
0;JMP
(%s_%d)
@SP
A=M-1
M=-1
(%s_%d_END)
	`,
		instruction.OriginalInstruction(),
		operation,
		count,
		jumpCommand,
		operation,
		count,
		operation,
		count,
		operation,
		count,
	)
}

func getRegisterNameFromSegmentName(instruction *parser.ParsedInstruction) string {
	segment, ok := segmentMap[instruction.Arg1()]
	if !ok {
		log.Fatalln("unknown segment name ", instruction.Arg1())
	}

	return segment
}

func getIndexvalue(instruction *parser.ParsedInstruction) string {
	index := instruction.Arg2()
	if index == "true" {
		return "0"
	} else if index == "false" {
		return "-1"
	} else {
		return index
	}
}
