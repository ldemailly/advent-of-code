package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var xmas = []string{"X", "M", "A", "S"}

func HasXMAS(matrix [][]string, i, j, deltaX, deltaY int) int {
	for step := 1; step <= 3; step++ {
		newI := i + step*deltaY
		newJ := j + step*deltaX
		if newI < 0 || newI >= len(matrix) || newJ < 0 || newJ >= len(matrix[i]) {
			return 0
		}
		if matrix[newI][newJ] != xmas[step] {
			return 0
		}
	}
	return 1
}

func HasMASCrossed(matrix [][]string, i, j int) bool {
	if matrix[i-1][j-1] != "M" && matrix[i-1][j-1] != "S" {
		return false
	}
	if matrix[i+1][j-1] != "M" && matrix[i+1][j-1] != "S" {
		return false
	}
	if matrix[i-1][j+1] != "M" && matrix[i-1][j+1] != "S" {
		return false
	}
	if matrix[i+1][j+1] != "M" && matrix[i+1][j+1] != "S" {
		return false
	}
	return matrix[i-1][j-1] != matrix[i+1][j+1] && matrix[i-1][j+1] != matrix[i+1][j-1]
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	dataStr := string(data)
	lines := strings.Split(strings.TrimSpace(dataStr), "\n")
	width := len(lines[0])
	height := len(lines)
	fmt.Printf("%d x %d\n", width, height)
	matrix := make([][]string, height)
	for i, line := range lines {
		matrix[i] = strings.Split(line, "")
		if len(matrix[i]) != width {
			log.Fatalf("At %d: Invalid matrix %d vs %d", i, len(matrix[i]), width)
		}
	}

	// Part 1
	directions := [][2]int{
		{1, 0},   // Right
		{-1, 0},  // Left
		{0, -1},  // Up
		{0, 1},   // Down
		{1, -1},  // Diagonal up-right
		{-1, -1}, // Diagonal up-left
		{1, 1},   // Diagonal down-right
		{-1, 1},  // Diagonal down-left
	}

	sum := 0
	for i := range matrix {
		for j := range matrix[i] {
			if matrix[i][j] != xmas[0] {
				continue
			}
			for _, dir := range directions {
				sum += HasXMAS(matrix, i, j, dir[0], dir[1])
			}
		}
	}
	fmt.Println(sum)

	// Part 2
	sum = 0
	for i := 1; i < len(matrix)-1; i++ {
		for j := 1; j < len(matrix[i])-1; j++ {
			if matrix[i][j] != "A" {
				continue
			}
			if HasMASCrossed(matrix, i, j) {
				sum++
			}
		}
	}
	fmt.Println(sum)
}
