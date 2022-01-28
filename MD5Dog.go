package main

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/irony0egoist/MD5Dog/core/semaphore"
	"github.com/irony0egoist/MD5Dog/core/utils"
	"log"
	"os"
	"strings"
	"sync"
)

type HashType int

const (
	md5Salt HashType = 10
	hmacMd5 HashType = 20
)

func init() {
	flag.Usage = func() {
		h := []string{
			"MD5Dog - Use Golang to Crack Hash",
			"",
			"Options:",
			"  -H, --hash <data>		Hash value",
			"  -s, --salt <data>		Salt value",
			// "  -S, --secret <data>		Secret value",
			"  -p, --path <data>		Path to passwords dictionary",
			"  -t, --type <data>		Hash type",
			"  -c, --concurrency <data>	Concurrency count AND file buffer size",
			"",
			"Usage:",
			"*md5(password.salt):",
			"	Example: MD5Dog.exe -t 10 -c 1000 -s abc -H 1b7ff998949c08bfa0d399d41aa0cbdf -p dic\\pass.txt",
			"*hmac.md5(password) ;secret=md5(password):",
			"	Example: MD5Dog.exe -t 20 -c 1000 -H 92d858ce796c86d55090ce1f1bb7be9a -p dic\\pass.txt",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func main() {

	var hash string
	flag.StringVar(&hash, "H", "", "")
	flag.StringVar(&hash, "hash", "", "")

	var salt string
	flag.StringVar(&salt, "s", "", "")
	flag.StringVar(&salt, "salt", "", "")

	var path string
	flag.StringVar(&path, "p", "", "")
	flag.StringVar(&path, "path", "", "")

	var hashType int
	flag.IntVar(&hashType, "t", 0, "")
	flag.IntVar(&hashType, "type", 0, "")

	var concurrency int
	flag.IntVar(&concurrency, "c", 1000, "")
	flag.IntVar(&concurrency, "concurrency", 1000, "")

	flag.Parse()

	sem := semaphore.New(concurrency)
	words := make(chan string, concurrency)

	go utils.ReadDictionary(path, words)

	passwordsCount, err := utils.LineCounter(path)
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
			bruteHash(words, salt, hash, HashType(hashType))
		}()
	}
	wg.Wait()
	log.Printf("md5Dog finished!")
}

func bruteHash(words chan string, salt, hash string, hashType HashType) {
	wordTrim := strings.TrimSpace(<-words)
	if hashType == md5Salt {
		rst := fmt.Sprintf("%x", md5.Sum([]byte(wordTrim+salt)))
		if rst == strings.ToLower(hash) {
			log.Printf("successful: %s\n", wordTrim)
		}
	} else if hashType == hmacMd5 {
		secret := fmt.Sprintf("%x", md5.Sum([]byte(wordTrim)))
		h := hmac.New(md5.New, []byte(secret))
		h.Write([]byte(wordTrim))
		rst := hex.EncodeToString(h.Sum(nil))
		if rst == strings.ToLower(hash) {
			log.Printf("successful: %s\n", wordTrim)
		}
	} else {
		log.Fatal("hash type not support")
	}

}
