package main

import (
	"bytes"
	"fmt"
	"strings"
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

var chars, ordered = makeCharsMap()

func main() {
	readStart := time.Now()
	inFile, e := os.Open("./in.txt")

	check(e)

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	for scanner.Scan() {
		// fmt.Println(scanner.Text())

		for _, c := range scanner.Text() {
			char := strings.ToLower(string(c))

			if _, ok := chars[char]; !ok {
				continue
			}

			chars[char]++
		}

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
		fmt.Println("Char:", char, "Ocurrencies:", chars[char])

		var text bytes.Buffer

		for i := 0; i < chars[char]; i++ {
			text.WriteString(char)
		}

		text.WriteTo(outFile)
	}

	outFile.Sync()

	fmt.Println("Write took", time.Since(writeStart))
}
