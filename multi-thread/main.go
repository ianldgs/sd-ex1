package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	// "io/ioutil"
	"bufio"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func makeCharsMap() (charsMap map[string]int, charsOrdered []string) {
	charsMap = make(map[string]int)

	for i := 0; i < 10; i++ {
		// fmt.Println(i)
		charsMap[strconv.Itoa(i)] = 0
		charsOrdered = append(charsOrdered, strconv.Itoa(i))
	}

	bytes := make([]byte, 26)

	for i := range bytes {
		// fmt.Println(string('a' + byte(i)))
		charsMap[string('a'+byte(i))] = 0
		charsOrdered = append(charsOrdered, string('a'+byte(i)))
	}

	return
}

var (
	chars, ordered = makeCharsMap()
	charsMu        sync.Mutex
)

type CharCount struct {
	char  string
	count uint
}

func countChar(char string, ch chan CharCount) {
	var charCount = CharCount{
		char:  char,
		count: 0,
	}

	inFile, e := os.Open("./in.txt")

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	check(e)

	// start := time.Now()

	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())

		for _, c := range line {
			_char := string(c)

			if char != _char {
				continue
			}

			charCount.count++
		}

		// charCount.count += uint(strings.Count(line, char))
	}

	// fmt.Println("Processed char", char, "in", time.Since(start))

	ch <- charCount
}

//https://play.golang.org/p/coja1_w-fY

func main() {
	readStart := time.Now()

	ch := make(chan CharCount)

	for _, c := range ordered {
		go countChar(c, ch)

		// break
	}

	fmt.Println("Started all threads")

	for range ordered {
		// fmt.Println("waiting", i)

		<-ch

		// fmt.Println("Char:", charCount.char, "Ocurrencies:", charCount.count)

		// break
	}

	fmt.Println("Read took", time.Since(readStart))

	// return

	writeStart := time.Now()
	outFile, e := os.Create("./out.txt")

	check(e)

	defer outFile.Close()

	e = outFile.Truncate(0)
	check(e)

	_, e = outFile.Seek(0, 0)
	check(e)

	for _, char := range ordered {
		// fmt.Println("Char:", char, "Ocurrencies:", chars[char])
		for i := 0; i < chars[char]; i++ {
			fmt.Fprintf(outFile, char)

			// break
		}
	}

	outFile.Sync()

	fmt.Println("Write took", time.Since(writeStart))
}
