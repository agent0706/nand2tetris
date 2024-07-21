package assembler

import (
	"log"
	"nand2tetris/assembler/parser"
	"nand2tetris/assembler/reader"
	"nand2tetris/assembler/translator"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	predefinedLabels map[string]int = map[string]int{
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	MAX_MEMORY_ADDRESS                 = 16384
	validJumpLabels    map[string]bool = map[string]bool{
		"JMP": true,
		"JGT": true,
		"JEQ": true,
		"JGE": true,
		"JLT": true,
		"JNE": true,
		"JLE": true,
	}
)

type Assembler struct {
	SymbolTable map[string]int
	FilePath    string
}

func (assembler *Assembler) FirstPass() {
	instructionReader := reader.NewReader(assembler.FilePath)
	instructionParser := parser.NewParser()

	var currentInstructionNumber int = 0
	//var variablesCount int = 0
	var variables []string
	var lineNumber int = 0

	for {
		instruction, hasMoreLines := instructionReader.Next()
		if !hasMoreLines {
			break
		}
		instructionParser.Parse(instruction)

		lineNumber += 1

		if instructionParser.InstructionType() == parser.A_INSTRUCTION || instructionParser.InstructionType() == parser.L_INSTRUCTION {
			currentSymbol := instructionParser.Symbol()
			validLabelRegex := regexp.MustCompile(`^[a-z|A-Z|_|\.|$|;][a-z|A-Z|_|\.|$|;|\d]*`)
			fullNumberRegex := regexp.MustCompile(`^\d*$`)

			if validLabelRegex.MatchString(currentSymbol) {
				label := validLabelRegex.FindString(currentSymbol)
				if val, ok := predefinedLabels[label]; ok {
					assembler.SymbolTable[label] = val
				} else {
					if instructionParser.InstructionType() == parser.L_INSTRUCTION {
						assembler.SymbolTable[label] = currentInstructionNumber
						// variablesCount -= 1
					} else {
						if _, ok := assembler.SymbolTable[label]; !ok {
							assembler.SymbolTable[label] = -1
							// variablesCount += 1
							variables = append(variables, label)
						}
					}
				}

			} else if !fullNumberRegex.MatchString(currentSymbol) {
				log.Fatalf("Error at Line:%d  illegal symbol name %s. should not start with number\n", lineNumber, instructionParser.Symbol())
			}
		} else {
			if len(instructionParser.Jump()) > 0 {
				if _, ok := validJumpLabels[instructionParser.Jump()]; !ok {
					log.Fatalf("Error ar Line:%d invalid jump label\n", currentInstructionNumber)
				}
			}
		}

		if instructionParser.InstructionType() == parser.A_INSTRUCTION || instructionParser.InstructionType() == parser.C_INSTRUCTION {
			currentInstructionNumber += 1
		}
	}

	variableMemoryAddress := 16
	for _, variable := range variables {

		if assembler.SymbolTable[variable] == -1 {
			assembler.SymbolTable[variable] = variableMemoryAddress
			variableMemoryAddress += 1
		}
	}
}

func getOutputFilePath(filePath string) string {
	extension := path.Ext(filePath)
	fileName := path.Base(filePath)
	return strings.Replace(fileName, extension, "", -1) + ".hack"
}

func (assembler *Assembler) SecondPass() {
	instructionReader := reader.NewReader(assembler.FilePath)
	instructionParser := parser.NewParser()
	translator := translator.NewTranslator(assembler.SymbolTable, instructionParser)

	outPutFilePath := getOutputFilePath(assembler.FilePath)
	file, err := os.OpenFile(outPutFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("error opening the output file")
	}
	defer file.Close()

	for {
		instruction, hasMoreLines := instructionReader.Next()
		if !hasMoreLines {
			file.Close()
			return
		}

		instructionParser.Parse(instruction)
		translatedInstruction := translator.Translate()

		if len(translatedInstruction) > 0 {
			_, err := file.WriteString(translatedInstruction + "\n")
			if err != nil {
				log.Fatalln("error writing translated instruction to output file")
			}
		}
	}
}
