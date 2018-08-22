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
	count int
}

func countChar(char string, in chan string, out chan CharCount) {
	var charCount = CharCount{
		char:  char,
		count: 0,
	}

	// start := time.Now()

	for {
		line, more := <-in

		if !more {
			// fmt.Println("Processed char", char, "in", time.Since(start))
			out <- charCount
			return
		}

		charCount.count += strings.Count(line, char)
	}
}

//https://play.golang.org/p/coja1_w-fY

func main() {
	readStart := time.Now()

	out := make(chan CharCount)
	ins := make([]chan string, 0)
	for range ordered {
		in := make(chan string)

		ins = append(ins, in)
	}

	inFile, e := os.Open("./in.txt")

	check(e)

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	for i, c := range ordered {
		go countChar(c, ins[i], out)

		// break
	}

	fmt.Println("Started all routines")

	for scanner.Scan() {
		for _, in := range ins {
			in <- strings.ToLower(scanner.Text())
		}
	}

	fmt.Println("Queued all work")

	for _, in := range ins {
		close(in)
	}

	for range ordered {
		charCount := <-out

		fmt.Println("Char:", charCount.char, "Ocurrencies:", charCount.count)

		// break
	}

	close(out)

	fmt.Println("Read took", time.Since(readStart))

	return

	writeStart := time.Now()
	outFile, e := os.Create("./out.txt")

	check(e)

	defer outFile.Close()

	e = outFile.Truncate(0)
	check(e)

	_, e = outFile.Seek(0, 0)
	check(e)

	for _, char := range ordered {
		fmt.Println("Char:", char, "Ocurrencies:", chars[char])
		for i := 0; i < chars[char]; i++ {
			fmt.Fprintf(outFile, char)

			// break
		}
	}

	outFile.Sync()

	fmt.Println("Write took", time.Since(writeStart))
}
