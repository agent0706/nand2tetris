package writer

import (
	"log"
	"os"
)

type Writer struct {
	file     *os.File
	fileName string
}

func NewWriter(fileName string) *Writer {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error open file to write output")
	}
	return &Writer{
		file:     file,
		fileName: fileName,
	}
}

func (w *Writer) WriteLines(lines []string) {
	for _, line := range lines {
		w.WriteLine(line)
	}
}

func (w *Writer) WriteLine(line string) {
	_, err := w.file.WriteString(line)

	if err != nil {
		log.Fatalf("Error writing line to output file %s\n", w.fileName)
	}
}

func (w *Writer) Close() {
	w.file.Close()
}
