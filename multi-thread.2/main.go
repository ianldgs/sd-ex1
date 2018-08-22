package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	// "io/ioutil"
	"bufio"
	"os"
	"strconv"
)

//https://stackoverflow.com/questions/18267460/how-to-use-a-goroutine-pool

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

	start := time.Now()

	for {
		line, more := <-in

		if !more {
			fmt.Println("Processed char", char, "in", time.Since(start))
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
	ins := make(map[string]chan string)

	for _, char := range ordered {
		ins[char] = make(chan string)

		go countChar(char, ins[char], out)
	}

	inFile, e := os.Open("./in.txt")

	check(e)

	defer inFile.Close()

	scanner := bufio.NewReader(inFile)

	fmt.Println("Started all routines")

	for {
		charr, _, err := scanner.ReadRune()

		char := strings.ToLower(string(charr))

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		if _, ok := ins[char]; ok {
			ins[char] <- char
		}
	}

	fmt.Println("Queued all work", time.Since(readStart))

	return

	for {
		_, more := <-out

		if !more {
			break
		}
	}

	for range ordered {
		charCount := <-out

		fmt.Println("Char:", charCount.char, "Ocurrencies:", charCount.count)
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
