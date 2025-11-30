package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var ansi256Gray = []int{
	237, // 58
	239, // 78
	241, // 98
	243, // 118
	245, // 138
	247, // 158
	249, // 178
	251, // 198
	253, // 218
	15,  // 255 white
}

// copied from  https://pkg.go.dev/fortio.org/terminal fps -fire mode
//
// --------------------------- 0   1   2   3    4    5    6    7    8    9.
var ansiColorGradiant = []int{52, 88, 124, 166, 202, 208, 214, 220, 226, 228}

const (
	ClearScreen = "\033[2J"
)

func StartSyncMode() {
	fmt.Print("\033[?2026h")
}

// End sync (and flush).
func EndSyncMode() {
	fmt.Print("\033[?2026l")
}

func HeightToPixel(h int, color bool) string {
	if h < 0 || h >= len(ansi256Gray) {
		log.Fatalf("Invalid height %d", h)
	}
	fg := 0
	if h < 5 {
		fg = 15
	}
	c := ansi256Gray[h]
	if color {
		c = ansiColorGradiant[h]
	}
	// Copied/inspired from my https://pkg.go.dev/fortio.org/terminal/ansipixels
	// with numbers for verification.
	return fmt.Sprintf("\033[48;5;%dm\033[38;5;%dm %d\033[0m", c, fg, h)
	// just colors/gray: (no number)
	// return fmt.Sprintf("\033[48;5;%dm\033[38;5;%dm  \033[0m", c, fg)
}

type Point struct {
	height int
	path   bool
}

// happens to be square.
type Map struct {
	points [][]Point
	l      int
	part1  bool
}

func (m *Map) Print() string {
	var b strings.Builder
	b.WriteString("\033[H")
	for y := range m.l {
		line := m.points[y]
		for x := range m.l {
			b.WriteString(HeightToPixel(line[x].height, line[x].path))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *Map) FindPaths() {
	sum := 0
	// find all the starts
	fmt.Print(ClearScreen)
	sleep := 10 * time.Millisecond
	if m.part1 {
		sleep *= 3
	}
	for y := range m.l {
		for x := range m.l {
			if m.points[y][x].height == 0 {
				p := m.CheckPath(x, y, 0)
				sum += p
				if p > 0 {
					StartSyncMode()
					fmt.Print(m.Print())
					fmt.Printf("Found %d path(s) at %d %d\n", p, x, y)
					EndSyncMode()
					time.Sleep(sleep)
					if m.part1 {
						m.ClearPaths()
					}
				}
			}
		}
	}
	part := "2"
	if m.part1 {
		part = "1"
	}
	fmt.Printf("### Part %s found %d paths\n", part, sum)
	time.Sleep(700 * time.Millisecond)
}

func (m *Map) CheckPath(x, y, cur int) int {
	// fmt.Printf("Checking %d %d %d\n", x, y, cur)
	if x < 0 || x >= m.l || y < 0 || y >= m.l {
		return 0
	}
	h := m.points[y][x].height
	if h != cur {
		return 0
	}
	if h == 9 {
		if m.part1 && m.points[y][x].path {
			return 0 // part1 only counts how many 9s are reachable, part2 counts all paths.
		}
		m.points[y][x].path = true
		return 1
	}
	cur++
	// check all not just first match
	p := m.CheckPath(x+1, y, cur) // check right
	p += m.CheckPath(x, y-1, cur) // check up
	p += m.CheckPath(x-1, y, cur) // check left
	p += m.CheckPath(x, y+1, cur) // check down
	if p > 0 {
		m.points[y][x].path = true
		return p
	}
	return 0
}

func (m *Map) ClearPaths() {
	for y := range m.l {
		for x := range m.l {
			m.points[y][x].path = false
		}
	}
}

func main() {
	//	flag.BoolVar(&doPrint, "print", false, "Print the map")
	//	flag.Parse()
	s := bufio.NewScanner(os.Stdin)
	var m Map
	y := 0
	for s.Scan() {
		line := s.Text()
		elems := strings.Split(line, "")
		var row []Point
		for _, e := range elems {
			n, err := strconv.Atoi(e)
			if err != nil {
				log.Fatalf("Invalid height %q: %v", e, err)
			}
			row = append(row, Point{height: n})
		}
		m.points = append(m.points, row)
		y++
	}
	m.l = len(m.points)
	if len(m.points[0]) != m.l {
		log.Fatalf("Invalid matrix %d vs %d", len(m.points[0]), m.l)
	}
	m.part1 = true
	m.FindPaths()
	m.part1 = false
	m.FindPaths()
}
