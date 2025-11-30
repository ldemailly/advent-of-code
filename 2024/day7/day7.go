package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Val struct {
	val    int
	len    int
	powmul int
}

type Eq struct {
	result   Val
	operands []Val
}

func IntPow10(n int) int {
	res := 1
	for range n {
		res *= 10
	}
	return res
}

func Str2Val(s string) Val {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid operand %q: %v", s, err)
	}
	return Val{val: n, len: len(s), powmul: IntPow10(len(s))}
}

func main() {
	var eqs []Eq
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var eq Eq
		line := s.Text()
		lr := strings.Split(line, ": ")
		if len(lr) != 2 {
			log.Fatalf("Invalid line %q", line)
		}
		eq.result = Str2Val(lr[0])
		elems := strings.Split(lr[1], " ")
		sumLen := 0
		for _, e := range elems {
			n := Str2Val(e)
			sumLen += n.len
			eq.operands = append(eq.operands, n)
		}
		if eq.result.len > sumLen {
			// sadly this doesn't happen/doesn't remove any case.
			fmt.Println("Skipping impossible", eq)
			continue
		}
		eqs = append(eqs, eq)
	}
	var sum int64
	for _, eq := range eqs {
		if evalPart1(eq) {
			// fmt.Println("Found", eq)
			sum += int64(eq.result.val)
		}
	}
	fmt.Println("Part1:", sum)
	sum = 0
	for _, eq := range eqs {
		if evalPart2(eq) {
			// fmt.Println("Found", eq)
			sum += int64(eq.result.val)
		}
	}
	fmt.Println("Part2:", sum)
}

func evalPart1(eq Eq) bool {
	return TryPlus1(eq.operands[1:], eq.operands[0].val, eq.result.val) ||
		TryTimes1(eq.operands[1:], eq.operands[0].val, eq.result.val)
}

func TryPlus1(operands []Val, acc, target int) bool {
	// fmt.Printf("TryPlus(%v, %d, %d)\n", operands, acc, target)
	acc += operands[0].val
	if len(operands) == 1 {
		return acc == target
	}
	if acc > target {
		return false
	}
	return TryPlus1(operands[1:], acc, target) || TryTimes1(operands[1:], acc, target)
}

func TryTimes1(operands []Val, acc, target int) bool {
	// fmt.Printf("TryPlus(%v, %d, %d)\n", operands, acc, target)
	acc *= operands[0].val
	if len(operands) == 1 {
		return acc == target
	}
	if acc > target {
		return false
	}
	return TryPlus1(operands[1:], acc, target) || TryTimes1(operands[1:], acc, target)
}

// --- could change part1 too but... it's late already.

func evalPart2(eq Eq) bool {
	return TryConcat(eq.operands[1:], eq.operands[0].val, eq.result.val) ||
		TryPlus2(eq.operands[1:], eq.operands[0].val, eq.result.val) ||
		TryTimes2(eq.operands[1:], eq.operands[0].val, eq.result.val)
}

func TryPlus2(operands []Val, acc, target int) bool {
	// fmt.Printf("TryPlus(%v, %d, %d)\n", operands, acc, target)
	acc += operands[0].val
	if len(operands) == 1 {
		return acc == target
	}
	if acc > target {
		return false
	}
	return TryPlus2(operands[1:], acc, target) || TryTimes2(operands[1:], acc, target) || TryConcat(operands[1:], acc, target)
}

func TryTimes2(operands []Val, acc, target int) bool {
	// fmt.Printf("TryTimes(%v, %d, %d)\n", operands, acc, target)
	acc *= operands[0].val
	if len(operands) == 1 {
		return acc == target
	}
	if acc > target {
		return false
	}
	return TryPlus2(operands[1:], acc, target) || TryTimes2(operands[1:], acc, target) || TryConcat(operands[1:], acc, target)
}

func TryConcat(operands []Val, acc, target int) bool {
	// fmt.Printf("TryConcat(%v, %d, %d)", operands, acc, target)
	acc = acc*operands[0].powmul + operands[0].val
	// fmt.Printf(" -> %d\n", acc)
	if len(operands) == 1 {
		return acc == target
	}
	if acc > target {
		return false
	}
	return TryPlus2(operands[1:], acc, target) || TryTimes2(operands[1:], acc, target) || TryConcat(operands[1:], acc, target)
}

// initial float math version.
func IntConcat1(a, b int) int {
	numDigits := int(math.Floor(math.Log10(float64(b))) + 1)
	return a*int(math.Pow10(numDigits)) + b
}

// Pure integer version (also not used now).
func IntConcat2(a, b int) int {
	big := 1
	for b >= big {
		big *= 10
	}
	return (a * big) + b
}
