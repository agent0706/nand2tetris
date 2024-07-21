package main

import (
	"flag"
	"log"
	"nand2tetris/assembler/assembler"
)

func getFileName() string {
	flag.Parse()
	fileName := flag.Arg(0)
	if len(fileName) == 0 {
		log.Fatalln("please provide path of the assembly file")
	}
	return fileName
}

func main() {
	//read hack file name from command line

	filePath := getFileName()

	hackAssembler := assembler.Assembler{
		SymbolTable: make(map[string]int),
		FilePath:    filePath,
	}
	hackAssembler.FirstPass()
	hackAssembler.SecondPass()
}
