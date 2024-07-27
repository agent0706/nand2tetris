package reader

import (
	"bufio"
	"log"
	"os"
)

type Reader struct {
	hasMoreLines bool
	scanner      *bufio.Scanner
	fd           *os.File
}

func NewReader(filePath string) *Reader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Error opening file")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return &Reader{
		scanner:      scanner,
		hasMoreLines: true,
		fd:           file,
	}
}

func (r *Reader) Next() (string, bool) {
	r.hasMoreLines = r.scanner.Scan()
	if r.hasMoreLines {
		result := r.scanner.Text()
		return result, r.hasMoreLines
	}
	r.fd.Close()
	return "", false
}
