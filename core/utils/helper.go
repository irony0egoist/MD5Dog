package utils

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func ReadDictionary(path string, words chan string) {
	r, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(words)
}

func LineCounter(path string) (int, error) {

	r, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}
