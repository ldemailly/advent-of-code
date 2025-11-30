package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	dataStr := string(data)
	fmt.Println(Mul(dataStr))

	// Part 2
	sum := 0
	for {
		idx1 := strings.Index(dataStr, "don't()")
		sum += Mul(dataStr[:idx1])
		if idx1 == -1 {
			break
		}
		idx2 := strings.Index(dataStr[idx1:], "do()")
		if idx2 == -1 {
			break
		}
		dataStr = dataStr[idx1+idx2+4:]
	}
	fmt.Println(sum)
}

func Mul(data string) int {
	re := regexp.MustCompile(`mul\(([0-9]+),([0-9]+)\)`)
	res := re.FindAllStringSubmatch(data, -1)
	sum := 0
	for _, r := range res {
		a, b := r[1], r[2]
		x, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(b)
		if err != nil {
			panic(err)
		}
		sum += x * y
	}
	return sum
}
