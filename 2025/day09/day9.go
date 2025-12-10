package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// Day 9 of advent of code 2025 - go version
// https://adventofcode.com/2025/day/9

func readInput() [][2]int {
	res := [][2]int{}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		res = append(res, [2]int{x, y})
	}
	return res
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func surf(a, b [2]int) int {
	dx := abs(a[0]-b[0]) + 1
	dy := abs(a[1]-b[1]) + 1
	return dx * dy
}

func sortFunc3(a, b [3]int) int {
	// reverse sort by first number (higher first)
	if a[0] < b[0] {
		return 1
	} else if a[0] > b[0] {
		return -1
	}
	// then by point1
	if a[1] < b[1] {
		return -1
	} else if a[1] > b[1] {
		return 1
	}
	return 0
}

func main() {
	fmt.Println("Part 1:")
	points := readInput()
	fmt.Println("input:", points)

	npts := len(points)
	ndist := (npts * (npts - 1)) / 2
	fmt.Println("Total points:", npts, "=> total rectangles:", ndist)
	surfaces := make([][3]int, ndist) // [distance, point1, point2]
	count := 0
	for i := range npts - 1 {
		p1 := points[i]
		for j := i + 1; j < npts; j++ {
			p2 := points[j]
			s := surf(p1, p2)
			surfaces[count] = [3]int{s, i, j}
			count++
		}
		fmt.Println("Point", i, "done.")
	}
	fmt.Println("Total surfaces:", len(surfaces))
	// sort surfaces
	slices.SortFunc(surfaces, sortFunc3)

	fmt.Println("Sorted surfaces:", surfaces)

	fmt.Println("Part 1: Result:", surfaces[0][0])

	// segments
	prev := points[0]
	var horizontal [][3]int
	var vertical [][3]int
	for i := 1; i < npts; i++ {
		x0 := prev[0]
		y0 := prev[1]
		curr := points[i]
		x1 := curr[0]
		y1 := curr[1]
		if x0 == x1 {
			// vertical
			height := abs(y1 - y0)
			vertical = append(vertical, [3]int{height, i - 1, i})
		} else {
			// horizontal
			width := abs(x1 - x0)
			horizontal = append(horizontal, [3]int{width, i - 1, i})
		}
		prev = curr
	}
	// bigger segments first, most likely to intersect
	slices.SortFunc(horizontal, sortFunc3)
	slices.SortFunc(vertical, sortFunc3)
	fmt.Println("Horizontal segments:", horizontal)
	fmt.Println("Vertical segments:", vertical)
	for _, s := range surfaces {
		idx1 := s[1]
		idx2 := s[2]
		corner1 := points[idx1]
		corner2 := points[idx2]
		bad := false
		for _, hseg := range horizontal {
			if Intersects(points[hseg[1]], points[hseg[2]], corner1, corner2) {
				fmt.Println("Surface", s, corner1, corner2, "intersects horizontal segment", hseg)
				bad = true
				break
			}
		}
		if bad {
			continue
		}
		for _, vseg := range vertical {
			if Intersects(points[vseg[1]], points[vseg[2]], corner1, corner2) {
				fmt.Println("Surface", s, corner1, corner2, "intersects vertical segment", vseg)
				bad = true
				break
			}
		}
		if bad {
			continue
		}
		fmt.Println("Part 2: Result:", s, corner1, corner2)
		topLeft, topRight, bottomLeft, bottomRight := getCorners(corner1, corner2)
		fmt.Printf(`<polygon points="%.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f" style="fill:blue;stroke:black;stroke-width=1;opacity:0.8"/>`+"\n",
			float64(topLeft[0])/100., float64(topLeft[1])/100.,
			float64(topRight[0])/100., float64(topRight[1])/100.,
			float64(bottomRight[0])/100., float64(bottomRight[1])/100.,
			float64(bottomLeft[0])/100., float64(bottomLeft[1])/100.,
		)
		return
	}
}

// Intersects returns true if segment is either fully or partially inside the rectangle
// defined by corner1 and corner2 but not just touching an edge or a corner.
func Intersects(segStart, segEnd, corner1, corner2 [2]int) bool {
	x, y, w, h := normalizeRectangle(corner1, corner2)
	start, end := normalizeSegment(segStart, segEnd)
	switch {
	case start[0] >= x+w:
		// segment entirely right of rectangle
		return false
	case end[0] <= x:
		// segment entirely left of rectangle
		return false
	case start[1] >= y+h:
		// segment entirely below rectangle
		return false
	case end[1] <= y:
		// segment entirely above rectangle
		return false
	}
	return true
}

func orderedSegment(a, b int) [2]int {
	if a < b {
		return [2]int{a, b}
	}
	return [2]int{b, a}
}

// Given 2 points of a segment returns the first,second such as first is leftmost or topmost
func normalizeSegment(c1, c2 [2]int) (first, second [2]int) {
	if c1[0] < c2[0] {
		return c1, c2
	} else if c1[0] > c2[0] {
		return c2, c1
	} else {
		// same x, vertical segment: order by y
		if c1[1] < c2[1] {
			return c1, c2
		} else {
			return c2, c1
		}
	}
}

// Given two corners, return coords of top left and width and height
func normalizeRectangle(c1, c2 [2]int) (x, y, w, h int) {
	x0 := min(c1[0], c2[0])
	y0 := min(c1[1], c2[1])
	x1 := max(c1[0], c2[0])
	y1 := max(c1[1], c2[1])
	return x0, y0, x1 - x0, y1 - y0
}

// Given two corners, return topleft, topright, bottomleft, bottomright corners
func getCorners(c1, c2 [2]int) ([2]int, [2]int, [2]int, [2]int) {
	x0 := min(c1[0], c2[0])
	y0 := min(c1[1], c2[1])
	x1 := max(c1[0], c2[0])
	y1 := max(c1[1], c2[1])
	return [2]int{x0, y0}, [2]int{x1, y0}, [2]int{x0, y1}, [2]int{x1, y1}
}

// Check if horizontal segmen hseg crosses line p1-p2
func CrossH(points [][2]int, hseg [3]int, p1, p2 [2]int) bool {
	fmt.Println("CrossH: hseg", hseg, "p1", p1, "p2", p2)
	// horizontal segment
	h1 := points[hseg[1]]
	h2 := points[hseg[2]]
	yh := h1[1]
	if h2[1] != yh {
		panic("Horizontal segment not horizontal")
	}
	xh1 := min(h1[0], h2[0])
	xh2 := max(h1[0], h2[0])
	// check p1,p2 are vertical
	if p1[0] != p2[0] {
		panic("CrossH: p1,p2 not vertical")
	}
	xp := p1[0]
	ymin := min(p1[1], p2[1])
	ymax := max(p1[1], p2[1])
	fmt.Println("xh:", xh1, xh2, "yh:", yh, "xp:", xp, "ymin:", ymin, "ymax:", ymax)
	// x of p1,p2 outside or just touching horizontal segment: not crossing
	if xp <= xh1 || xp >= xh2 {
		fmt.Println("xp outside hseg", xp, xh1, xh2)
		return false
	}
	if yh <= ymin || yh >= ymax {
		fmt.Println("yh outside px", yh, ymin, ymax)
		return false
	}
	// crossing
	return true
}

// Check if vertical segment vseg crosses line p1-p2
func CrossV(points [][2]int, vseg [3]int, p1, p2 [2]int) bool {
	fmt.Println("CrossV: vseg", vseg, "p1", p1, "p2", p2)
	// vertical segment
	v1 := points[vseg[1]]
	v2 := points[vseg[2]]
	xv := v1[0]
	if v2[0] != xv {
		panic("Vertical segment not vertical")
	}
	yv1 := min(v1[1], v2[1])
	yv2 := max(v1[1], v2[1])
	// check p1,p2 are horizontal
	if p1[1] != p2[1] {
		panic("CrossV: p1,p2 not horizontal")
	}
	yp := p1[1]
	xmin := min(p1[0], p2[0])
	xmax := max(p1[0], p2[0])
	fmt.Println("yv:", yv1, yv2, "xv:", xv, "yp:", yp, "xmin:", xmin, "xmax:", xmax)
	// y of p1,p2 outside or just touching vertical segment: not crossing
	if yp <= yv1 || yp >= yv2 {
		fmt.Println("yp outside vseg", yp, yv1, yv2)
		return false
	}
	if xv <= xmin || xv >= xmax {
		fmt.Println("xv outside px", xv, xmin, xmax)
		return false
	}
	// crossing
	return true
}
