// Day 4 of advent of code 2025 - go version
// https://adventofcode.com/2025/day/4

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func readInput() [][]int {
	var res [][]int
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		var row []int
		for _, vs := range strings.Split(line, "") {
			v := 0
			if vs == "@" {
				v = 1 // 1 for initial live cell
			}
			row = append(row, v)
		}
		res = append(res, row)
	}
	return res
}

func countNeighbors(grid [][]int, x, y int) int {
	sum := 0
	w := len(grid[0])
	h := len(grid)
	for j := y - 1; j < y+2; j++ {
		for i := x - 1; i < x+2; i++ {
			if i == x && j == y { //center
				continue
			}
			if i < 0 || j < 0 || j >= h || i >= w {
				continue
			}
			if grid[j][i] == 1 {
				sum++
			}
		}
	}
	return sum
}

const (
	FULL     = "█"
	TOP      = "▀"
	BOTTOM   = "▄"
	GREEN    = "\x1b[32m"
	GREEN_BG = "\x1b[42m"
	RED      = "\x1b[31m"
	RED_BG   = "\x1b[41m"
	RESET    = "\x1b[0m"

	HOME       = "\x1b[H"
	CLEAR      = "\x1b[2J"
	SYNC_START = "\x1b[?2026h"
	SYNC_END   = "\x1b[?2026l"
)

func color(state int) string {
	if state == 0 {
		return RESET
	}
	if state == 1 {
		return GREEN
	}
	return RED
}

func printGrid(grid [][]int) {
	w := len(grid[0])
	h := len(grid)
	for y := 0; y < h; y += 2 {
		line := ""
		for x := range w {
			top := grid[y][x]
			bottom := grid[y+1][x]
			if top == bottom {
				if top == 0 {
					line += color(top) + " "
				} else {
					line += color(top) + FULL
				}
			} else if bottom == 0 {
				line += color(top) + TOP
			} else if top == 0 {
				line += color(bottom) + BOTTOM
			} else {
				if top == 2 {
					line += RED_BG
				} else {
					line += GREEN_BG
				}
				line += color(bottom) + BOTTOM + RESET
			}
		}
		fmt.Println(line)
	}
}

func CountGridNeighbors(grid [][]int) {
	w := len(grid[0])
	h := len(grid)
	for y := range h {
		for x := range w {
			n := countNeighbors(grid, x, y)
			fmt.Printf("%d", n)
		}
		fmt.Println()
	}
}

func Remove(grid [][]int) ([][]int, int) {
	w := len(grid[0])
	h := len(grid)
	sum := 0
	var ngrid [][]int
	for y := range h {
		var nline []int
		for x := range w {
			v := grid[y][x]
			if v == 1 {
				n := countNeighbors(grid, x, y)
				if n < 4 {
					sum++
					v = 2
				}
			}
			nline = append(nline, v)
		}
		ngrid = append(ngrid, nline)
	}
	return ngrid, sum
}

func main() {
	fmt.Print(HOME + CLEAR)
	input := readInput()
	ngrid := input
	printGrid(ngrid)
	fmt.Print("Initial...  ")
	time.Sleep(1 * time.Second)
	p1 := 0
	sum := 0
	var removed int
	for {
		ngrid, removed = Remove(ngrid)
		fmt.Print(HOME + SYNC_START)
		// fmt.Println()
		printGrid(ngrid)
		fmt.Print(RESET+"Next generation removed:", removed, "  total so far ", sum, "  ")
		fmt.Print(SYNC_END)
		sum += removed
		if p1 == 0 {
			p1 = removed
		}
		if removed == 0 {
			break
		}
		time.Sleep(time.Duration(removed) * time.Millisecond)
	}
	fmt.Println("\n1.Part 1:", p1)
	fmt.Println("2.Total removed:", sum)
}
