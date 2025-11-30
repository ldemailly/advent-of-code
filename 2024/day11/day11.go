package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Stone struct {
	n    int
	nStr string
}

func evenDigits(n Stone) bool {
	return len(n.nStr)%2 == 0
}

func leftRight(n Stone) []Stone {
	left := n.nStr[:len(n.nStr)/2]
	right := n.nStr[len(n.nStr)/2:]
	return []Stone{StoneFromStr(left), StoneFromStr(right)}
}

var One = Stone{1, "1"}

func StoneFromNum(n int) Stone {
	return Stone{n, strconv.Itoa(n)}
}

func StoneFromStr(s string) Stone {
	n, _ := strconv.Atoi(s)
	return Stone{n, strings.TrimLeft(s, "0")}
}

type CacheKey struct {
	stone Stone
	depth int
}

var cache = make(map[CacheKey]int)

func ApplyRules(v Stone, depth int) int {
	// fmt.Println(v, depth)
	if depth == 0 {
		return 1
	}
	cacheKey := CacheKey{v, depth}
	if r, ok := cache[cacheKey]; ok {
		return r
	}
	r := 0
	switch {
	case v.n == 0:
		r = ApplyRules(One, depth-1)
	case evenDigits(v):
		lr := leftRight(v)
		r = ApplyRules(lr[0], depth-1) + ApplyRules(lr[1], depth-1)
	default:
		r = ApplyRules(StoneFromNum(2024*v.n), depth-1)
	}
	cache[cacheKey] = r
	return r
}

func main() {
	// read input list
	data, _ := io.ReadAll(os.Stdin)
	inp := strings.TrimSpace(string(data))
	inpList := strings.Split(inp, " ")
	// Split into list of ints.
	stones := make([]Stone, 0, len(inpList))
	for _, n := range inpList {
		stones = append(stones, StoneFromStr(n))
	}
	// Part1
	sum := 0
	for v := range stones {
		sum += ApplyRules(stones[v], 25)
	}
	fmt.Println("Part 1", sum)
	// Part2
	sum = 0
	for v := range stones {
		sum += ApplyRules(stones[v], 75)
	}
	fmt.Println("Part 2", sum)
}
