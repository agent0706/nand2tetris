package parser

import (
	"regexp"
	"strings"
)

var (
	C_ARITHMETIC = "C_ARITHMETIC"
	C_PUSH       = "C_PUSH"
	C_POP        = "C_POP"
	C_LABEL      = "C_LABEL"
	C_GOTO       = "C_GOTO"
	C_IF         = "C_IF"
	C_FUNCTION   = "C_FUNCTION"
	C_RETURN     = "C_RETURN"
	C_CALL       = "C_CALL"
)

type ParsedInstruction struct {
	command     string
	arg1        string
	arg2        string
	commandType string
	instruction string
	fileName    string
}

func getCommandType(command string) string {
	switch strings.ToLower(command) {
	case "add":
		fallthrough
	case "sub":
		fallthrough
	case "neg":
		fallthrough
	case "eq":
		fallthrough
	case "gt":
		fallthrough
	case "lt":
		fallthrough
	case "and":
		fallthrough
	case "or":
		fallthrough
	case "not":
		return C_ARITHMETIC
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	default:
		return ""
	}
}

func ParseInstruction(instruction string, fileName string) *ParsedInstruction {
	// remove comments
	commentRegex := regexp.MustCompile(`//.*`)
	temp := commentRegex.ReplaceAllString(instruction, "")

	temp = strings.Trim(temp, " ")

	if len(temp) == 0 {
		return &ParsedInstruction{}
	}

	splittedInstruction := strings.Split(instruction, " ")

	var arg1 string
	var arg2 string
	if len(splittedInstruction) == 2 {
		arg1 = splittedInstruction[1]
	} else if len(splittedInstruction) == 3 {
		arg1 = splittedInstruction[1]
		arg2 = splittedInstruction[2]
	}

	var command string = splittedInstruction[0]

	commandType := getCommandType(command)

	return &ParsedInstruction{
		command:     command,
		arg1:        arg1,
		arg2:        arg2,
		commandType: commandType,
		instruction: instruction,
		fileName:    fileName,
	}
}

func (i *ParsedInstruction) Command() string {
	return i.command
}

func (i *ParsedInstruction) Arg1() string {
	return i.arg1
}

func (i *ParsedInstruction) Arg2() string {
	return i.arg2
}

func (i *ParsedInstruction) OriginalInstruction() string {
	return i.instruction
}

func (i *ParsedInstruction) FileName() string {
	return i.fileName
}
