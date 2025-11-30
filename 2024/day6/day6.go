package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	north = iota
	east
	south
	west
	maxdir
)

type Guard struct {
	x, y      int
	direction int
}

func (g *Guard) TurnRight() {
	g.direction = (g.direction + 1) % maxdir
}

func (g *Guard) Move(backward bool) {
	delta := 1
	if backward {
		delta = -1
	}
	switch g.direction {
	case north:
		g.y -= delta
	case east:
		g.x += delta
	case south:
		g.y += delta
	case west:
		g.x -= delta
	}
}

const (
	GreenBackground = "\033[42m"
	// copied from my https://pkg.go.dev/fortio.org/terminal/ansipixels
	Reset          = "\033[0m"
	Green          = "\033[32m"
	FullPixel      = '█'
	BluePixel      = "\033[34m" + string(FullPixel) + Reset
	RedPixel       = "\033[31m" + string(FullPixel) + Reset
	ClearScreen    = "\033[2J"
	ClearEndOfLine = "\033[K"
)

func (g *Guard) Print() string {
	switch g.direction {
	case north:
		return GreenBackground + "⇑" + Reset
	case east:
		return GreenBackground + "⇒" + Reset
	case south:
		return GreenBackground + "⇓" + Reset
	case west:
		return GreenBackground + "⇐" + Reset
	default:
		panic("Invalid direction")
	}
}

const (
	empty = iota
	obstacle
	visited // + direction
)

// happens to be square.
type Map struct {
	m          [][]int
	l          int
	visited    int
	startX     int
	startY     int
	steps      int
	loop       bool
	newX, newY int
}

func (m *Map) Outside(g *Guard) bool {
	return g.x < 0 || g.y < 0 || g.x >= m.l || g.y >= m.l
}

func (m *Map) Next(g *Guard) bool {
	m.steps++
	if m.m[g.y][g.x] == empty {
		m.visited++
		m.m[g.y][g.x] = visited + g.direction
	} else if m.m[g.y][g.x] == visited+g.direction {
		// fmt.Println("Loop at", g.x, g.y, "after", m.steps, "steps")
		m.loop = true
		return false
	}
	g.Move(false)
	if m.Outside(g) {
		return false // done
	}
	if m.m[g.y][g.x] == obstacle {
		g.Move(true) // backward/undo
		g.TurnRight()
	}
	return true
}

func (m *Map) SetNew(x, y int) {
	m.newX = x
	m.newY = y
	m.m[y][x] = obstacle
}

func (m *Map) Print(g Guard) string {
	var b strings.Builder
	b.WriteString("\033[H")
	for y := range m.l {
		for x := range m.l {
			if g.x == x && g.y == y {
				b.WriteString(g.Print())
				continue
			}
			if x == m.startX && y == m.startY {
				b.WriteString(Green + "↑" + Reset)
				continue
			}
			switch m.m[y][x] {
			case empty:
				b.WriteRune(FullPixel)
			case obstacle:
				if x == m.newX && y == m.newY {
					b.WriteString(BluePixel)
				} else {
					b.WriteString(RedPixel)
				}
			case visited + north:
				b.WriteString("↑")
			case visited + east:
				b.WriteString("→")
			case visited + south:
				b.WriteString("↓")
			case visited + west:
				b.WriteString("←")
			default:
				log.Fatalf("Invalid state %q at %d,%d", m.m[y][x], x, y)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// copy without visited.
func (m *Map) CleanClone() Map {
	var m2 Map
	m2.l = m.l
	m2.startX = m.startX
	m2.startY = m.startY
	m2.steps = 0
	m2.visited = 0
	m2.loop = false
	m2.m = make([][]int, 0, m.l)
	for y := range m.l {
		row := make([]int, 0, m.l)
		for x := range m.l {
			v := m.m[y][x]
			if v != obstacle {
				v = empty
			}
			row = append(row, v)
		}
		m2.m = append(m2.m, row)
	}
	return m2
}

var (
	doPrint = true
	slow    time.Duration
)

func (m *Map) MayPrint(g Guard, slowMult int, msg ...interface{}) {
	if doPrint {
		fmt.Print(m.Print(g))
		if len(msg) > 0 {
			fmt.Println(msg...)
		}
		time.Sleep(time.Duration(slowMult) * slow)
	}
}

func main() {
	flag.BoolVar(&doPrint, "print", false, "Print the map")
	flag.DurationVar(&slow, "delay", 0, "Slow down the output by this amount")
	flag.Parse()
	s := bufio.NewScanner(os.Stdin)
	var m Map
	var guard Guard
	y := 0
	for s.Scan() {
		line := s.Text()
		elems := strings.Split(line, "")
		var row []int
		for x, e := range elems {
			switch e {
			case "#":
				row = append(row, obstacle)
			case ".":
				row = append(row, empty)
			case "^":
				row = append(row, empty)
				guard.x = x
				guard.y = y
			default:
				log.Fatalf("Invalid char %q", e)
			}
		}
		m.m = append(m.m, row)
		y++
	}
	m.l = len(m.m)
	if len(m.m[0]) != m.l {
		log.Fatalf("Invalid matrix %d vs %d", len(m.m[0]), m.l)
	}
	initialGuard := guard
	m.startX = guard.x
	m.startY = guard.y
	if doPrint {
		// clear screen
		fmt.Print(ClearScreen)
	}
	m.MayPrint(guard, 5, "Initial, guard at", guard.x, guard.y)
	for m.Next(&guard) {
		m.MayPrint(guard, 1, "Guard at", guard.x, guard.y, guard.Print(), ClearEndOfLine)
	}
	part1 := m.visited
	m.MayPrint(guard, 5, "Done after", m.steps, "steps:", part1, ClearEndOfLine)
	loops := 0
	for y := range m.l {
		for x := range m.l {
			if x == m.startX && y == m.startY {
				continue
			}
			if m.m[y][x] >= visited {
				m2 := m.CleanClone()
				m2.SetNew(x, y)
				guard = initialGuard
				for m2.Next(&guard) {
				}
				if m2.loop {
					loops++
					m2.MayPrint(guard, 5, "Loop", loops, "adding obstacle at", x, y)
				}
			}
		}
	}
	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", loops)
}
