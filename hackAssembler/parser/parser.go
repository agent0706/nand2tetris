package parser

import (
	"regexp"
	"strings"
)

var (
	C_INSTRUCTION = "C_INSTRUCTION"
	A_INSTRUCTION = "A_INSTRUCTION"
	L_INSTRUCTION = "L_INSTRUCTION"
)

type Parser struct {
	currentInstructionType string
	symbol                 string
	dest                   string
	comp                   string
	jump                   string
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) InstructionType() string {
	return p.currentInstructionType
}

func (p *Parser) Symbol() string {
	return p.symbol
}

func (p *Parser) Dest() string {
	return p.dest
}

func (p *Parser) Comp() string {
	return p.comp
}

func (p *Parser) Jump() string {
	return p.jump
}

func (p *Parser) Parse(instruction string) {
	commentRegex := regexp.MustCompile("//.*$")
	formattedInstruction := commentRegex.ReplaceAllString(strings.Trim(instruction, " "), "")

	if len(formattedInstruction) == 0 {
		p.symbol = ""
		p.currentInstructionType = ""
		p.dest = ""
		p.comp = ""
		p.jump = ""
		return
	}

	LInstructionRegex := regexp.MustCompile(`^\((?P<symbol>.*)\)$`)
	AInstructionRegex := regexp.MustCompile("^@(?P<address>.*)$")
	CInstructionRegex := regexp.MustCompile("^((?P<dest>.*)=)?(?P<comp>.*?)?(;(?P<jump>.*))?$")

	if LInstructionRegex.MatchString(formattedInstruction) {
		p.parseLInstruction(formattedInstruction, LInstructionRegex)
	} else if AInstructionRegex.MatchString(formattedInstruction) {
		p.parseAInstruction(formattedInstruction, AInstructionRegex)
	} else if CInstructionRegex.MatchString(formattedInstruction) {
		p.parseCInstruction(formattedInstruction, CInstructionRegex)
	}
}

func (p *Parser) parseLInstruction(instruction string, matchRegex *regexp.Regexp) {
	matches := matchRegex.FindStringSubmatch(instruction)
	p.symbol = matches[1]
	p.dest = ""
	p.comp = ""
	p.jump = ""
	p.currentInstructionType = L_INSTRUCTION
}

func (p *Parser) parseAInstruction(instruction string, matchRegex *regexp.Regexp) {
	matches := matchRegex.FindStringSubmatch(instruction)
	p.symbol = matches[1]
	p.currentInstructionType = A_INSTRUCTION
	p.dest = ""
	p.comp = ""
	p.jump = ""
}

func (p *Parser) parseCInstruction(instruction string, matchRegex *regexp.Regexp) {
	matches := matchRegex.FindAllStringSubmatch(instruction, -1)
	p.dest = matches[0][2]
	p.comp = matches[0][3]
	p.jump = matches[0][5]
	p.currentInstructionType = C_INSTRUCTION
	p.symbol = ""
}
