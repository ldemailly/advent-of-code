// Day 4 of advent of code 2025 - go version
// https://adventofcode.com/2025/day/4

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"time"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
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

func Remove(curGen int, grid [][]int) ([][]int, int) {
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
					v = -curGen
				}
			}
			nline = append(nline, v)
		}
		ngrid = append(ngrid, nline)
	}
	return ngrid, sum
}

func GridToImage(curGen int, grid [][]int) *image.RGBA {
	w := len(grid[0])
	h := len(grid)
	background := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			v := grid[y][x]
			var c color.RGBA
			switch v {
			case 0:
				c = background
			case 1:
				c = color.RGBA{R: 30, G: 170, B: 60, A: 255} // white
			default:
				kcol := tcolor.Oklchf(.2-float64(v)/70., 0.8, 0.08)
				ct, data := kcol.Decode()
				rgbc := tcolor.ToRGB(ct, data)
				c = color.RGBA{R: rgbc.R, G: rgbc.G, B: rgbc.B, A: 255}
			}
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func main() {
	ngrid := readInput()
	ap := ansipixels.NewAnsiPixels(0)
	ap.GetSize()
	_ = ap.ShowScaledImage(GridToImage(0, ngrid))
	fmt.Print("Initial...  ")
	time.Sleep(1 * time.Second)
	p1 := 0
	sum := 0
	var removed int
	curGen := 1
	for {
		ngrid, removed = Remove(curGen, ngrid)
		ap.StartSyncMode()
		ap.ClearScreen()
		_ = ap.ShowScaledImage(GridToImage(curGen, ngrid))
		ap.WriteAt(0, ap.H-4, "Generation %d removed: %d total so far %d   ", curGen, removed, sum)
		ap.EndSyncMode()
		sum += removed
		if p1 == 0 {
			p1 = removed
		}
		if removed == 0 {
			break
		}
		curGen++
		time.Sleep(time.Duration(removed) * time.Millisecond)
	}
	fmt.Println("\n1.Part 1:", p1)
	fmt.Println("2.Total removed:", sum)
}
