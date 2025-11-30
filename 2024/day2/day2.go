// Day 2 of advent of code 2024
// https://adventofcode.com/2024/day/2

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	reports := readInput()
	safe := 0
	dsafe := 0
	for _, r := range reports {
		SafeCount(r, &safe, &dsafe)
	}
	fmt.Println(safe)
	fmt.Println(dsafe)
}

func SafeCount(r []int, safe, dsafe *int) {
	if Safe(r) {
		*safe++
		*dsafe++
		return
	}
	rr := make([]int, 0, len(r)-1)
	for i := range len(r) {
		rr = rr[:0]
		rr = append(rr, r[0:i]...)
		rr = append(rr, r[i+1:]...)
		// log.Printf("SafeCount: %v -> %v", r, rr)
		if Safe(rr) {
			*dsafe++
			return
		}
	}
}

func Safe(r []int) bool {
	cur := r[0]
	delta := r[1] - cur
	direction := 1
	if delta < 0 {
		direction = -1
	}
	for _, e := range r[1:] {
		delta = direction * (e - cur)
		if delta <= 0 {
			return false
		}
		if delta > 3 {
			return false
		}
		cur = e
	}
	log.Printf("Safe: %v", r)
	return true
}

func readInput() [][]int {
	reports := [][]int{}
	lines := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines++
		str := scanner.Text()
		r := []int{}
		for _, e := range strings.Split(str, " ") {
			n, err := strconv.Atoi(e)
			if err != nil {
				log.Fatalf("At %d Expected a number, got %q: %v", lines, e, err)
			}
			r = append(r, n)
		}
		if len(r) < 2 {
			log.Fatalf("At %d Expected at least 2 levels, %q -> %v", lines, str, r)
		}
		reports = append(reports, r)
	}
	return reports
}
