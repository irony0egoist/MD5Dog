package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"log"
	"md5/core/semaphore"
	"md5/core/utils"
	"os"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) != 4 {
		log.Printf("Crack md5(password.salt)\nUsage: %s Salt Hash PasswordPath", os.Args[0])
		return
	}
	salt, sourceHash, passwordPath := os.Args[1], os.Args[2], os.Args[3]

	sem := semaphore.New(1000)
	words := make(chan string, 1000)

	go readPasswords(passwordPath, words)

	file, err := os.Open(passwordPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	passwordsCount, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("load passwords: %d lines", passwordsCount+1)

	var wg sync.WaitGroup
	for i := 0; i <= passwordsCount; i++ {
		wg.Add(1)
		sem.Acquire()
		go func() {
			defer wg.Done()
			defer sem.Release()
			brutePasswords(words, salt, sourceHash)
		}()
	}
	wg.Wait()
	log.Printf("md5Dog finished!")
}

func readPasswords(passwordPath string, words chan string) {
	file, err := os.Open(passwordPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(words)
}

func brutePasswords(words chan string, salt, sourceHash string) {
	wordTrim := strings.TrimSpace(<-words)
	md5Data := fmt.Sprintf("%x", md5.Sum([]byte(wordTrim+salt)))
	if md5Data == strings.ToLower(sourceHash) {
		log.Printf("successful: %s\n", wordTrim)
	}

}
