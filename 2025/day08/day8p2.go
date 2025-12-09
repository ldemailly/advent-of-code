package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// Day 8 of advent of code 2025 - go version
// https://adventofcode.com/2025/day/8

func readInput() [][3]int {
	res := [][3]int{}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		z, _ := strconv.Atoi(parts[2])
		res = append(res, [3]int{x, y, z})
	}
	return res
}

func dist(a, b [3]int) int {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]
	return dx*dx + dy*dy + dz*dz
}

func main() {
	// fmt.Println("Part 1:")
	points := readInput()
	// fmt.Println("input:", points)

	npts := len(points)
	ndist := (npts * (npts - 1)) / 2
	// fmt.Println("Total points:", npts, "=> total distances:", ndist)
	distances := make([][3]int, ndist) // [distance, point1, point2]
	count := 0
	for i := range npts - 1 {
		p1 := points[i]
		for j := i + 1; j < npts; j++ {
			p2 := points[j]
			d := dist(p1, p2)
			distances[count] = [3]int{d, i, j}
			count++
		}
		// fmt.Println("Point", i, "done.")
	}
	// fmt.Println("Total distances:", len(distances))
	// sort distances
	slices.SortFunc(distances, func(a, b [3]int) int {
		if a[0] < b[0] {
			return -1
		} else if a[0] > b[0] {
			return 1
		}
		// then by point1
		if a[1] < b[1] {
			return -1
		} else if a[1] > b[1] {
			return 1
		}
		return 0
	})

	// Build circuits

	pt2circuit := make(map[int]int) // point index to circuit id
	circuits := make(map[int][]int) // circuit id to list of point indices

	merge := func(c1, c2 int) {
		if c1 == c2 {
			return
		}
		// fmt.Println("Merging circuits", c1, "and", c2)
		for _, pt := range circuits[c2] {
			pt2circuit[pt] = c1
		}
		circuits[c1] = append(circuits[c1], circuits[c2]...)
		delete(circuits, c2)
	}

	cnum := 0

	for _, dval := range distances {
		// d := dval[0]
		i := dval[1]
		j := dval[2]
		if _, found := pt2circuit[i]; !found {
			// fmt.Println("Creating circuit", cnum, "for point", i)
			pt2circuit[i] = cnum
			circuits[cnum] = []int{i}
			cnum++
		}
		cid := pt2circuit[i]
		if _, found := pt2circuit[j]; found {
			merge(cid, pt2circuit[j])
		} else {
			circuits[cid] = append(circuits[cid], j)
		}
		pt2circuit[j] = cid
		// fmt.Println("Circuits", circuits)
		if len(circuits[cid]) == npts {
			fmt.Println("Part 2: Result:", i, j, points[i], points[j], points[i][0]*points[j][0])
			break
		}
	}
}
