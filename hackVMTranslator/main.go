package main

import (
	"flag"
	"log"
	"nandtotetris/vmtranslator/parser"
	"nandtotetris/vmtranslator/reader"
	"nandtotetris/vmtranslator/translator"
	"nandtotetris/vmtranslator/writer"
	"os"
	"path"
	"strings"
	"unicode"
)

func main() {
	filePath := getFilePath()

	validationError := validateFilePath(filePath)
	if validationError != nil {
		log.Fatalln(validationError)
	}

	fileName := getFileName(filePath)
	outputFileName := getOutputFilePath(filePath, fileName)

	instReader := reader.NewReader(filePath)
	outputWriter := writer.NewWriter(outputFileName)
	instTranslator := translator.NewTranslator()

	// asmInst := translator.GetInitializationInstructions()
	// outputWriter.WriteLine(asmInst)

	// read an instruction -> parse -> translate -> write
	for {
		inst, hasMoreLines := instReader.Next()
		if !hasMoreLines {
			break
		}
		parsedInst := parser.ParseInstruction(inst, fileName)
		translatedInstructions := instTranslator.Translate(parsedInst)
		outputWriter.WriteLine(translatedInstructions)
	}

	asmInst := instTranslator.GetTerminalInstrucitons()
	outputWriter.WriteLine(asmInst)

	outputWriter.Close()
}

func getFilePath() string {
	flag.Parse()
	fileName := flag.Arg(0)
	if len(fileName) == 0 {
		log.Fatalln("please provide path of the assembly file")
	}
	return fileName
}

func getFileName(inputFilePath string) string {
	extension := path.Ext(inputFilePath)
	file := path.Base(inputFilePath)

	fileName := strings.Replace(file, extension, "", -1)
	return fileName
}

func getOutputFilePath(inputFilePath string, fileName string) string {
	dir := path.Dir(inputFilePath)
	return path.Join(dir, fileName+".asm")
}

func validateFilePath(filePath string) *string {
	_, err := os.Stat(filePath)
	var errMsg string
	if err != nil {
		errMsg = "Provided file path does not exist. Please provide a valid path"
		return &errMsg
	}

	ext := path.Ext(filePath)
	if ext != ".vm" {
		errMsg = "Provided file is not a vm file"
		return &errMsg
	}

	fileName := path.Base(filePath)
	fileName = strings.Replace(fileName, ".vm", "", -1)

	if !(unicode.IsUpper(rune(fileName[0]))) {
		errMsg = "Filename should start with capital letter"
		return &errMsg
	}

	// filename should start with captial letter

	return nil
}
