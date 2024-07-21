package translator

import (
	"fmt"
	"log"
	"nand2tetris/assembler/parser"
	"strconv"
)

type Translator struct {
	symbolTable map[string]int
	parser      *parser.Parser
}

var dmap = map[string]string{
	"null": "000",
	"M":    "001",
	"D":    "010",
	"MD":   "011",
	"A":    "100",
	"AM":   "101",
	"AD":   "110",
	"AMD":  "111",
}

var jmap = map[string]string{
	"null": "000",
	"JGT":  "001",
	"JEQ":  "010",
	"JGE":  "011",
	"JLT":  "100",
	"JNE":  "101",
	"JLE":  "110",
	"JMP":  "111",
}

var cmap = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",
	"M":   "1110000",
	"!M":  "1110001",
	"-M":  "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
}

func NewTranslator(symbolTable map[string]int, parser *parser.Parser) Translator {
	return Translator{
		symbolTable: symbolTable,
		parser:      parser,
	}
}

func (t *Translator) Translate() string {
	if t.parser.InstructionType() == parser.A_INSTRUCTION {
		symbol := t.parser.Symbol()
		symbolVal, ok := t.symbolTable[symbol]
		var symbolInt int
		if ok {
			symbolInt = symbolVal
		} else {
			var err error
			symbolInt, err = strconv.Atoi(symbol)
			if err != nil {
				log.Fatalln("error converting string")
			}
		}

		symbolBinary := fmt.Sprintf("%015b", symbolInt)
		symbolBinary = "0" + symbolBinary
		return symbolBinary
	} else if t.parser.InstructionType() == parser.C_INSTRUCTION {
		dest := t.parser.Dest()
		if len(dest) == 0 {
			dest = "null"
		}
		comp := t.parser.Comp()
		jump := t.parser.Jump()
		if len(jump) == 0 {
			jump = "null"
		}

		result := "111"

		result += cmap[comp] + dmap[dest] + jmap[jump]
		return result
	}
	return ""
}
