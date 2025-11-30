package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"fortio.org/sets"
)

type Coord struct {
	x, y int
}

type Region struct {
	id        byte
	coords    sets.Set[Coord]
	perimeter int
	corners   int
}

type Map struct {
	plots   [][]byte
	seen    sets.Set[Coord]
	inner   sets.Set[Coord]
	regions []*Region
}

func (m *Map) Regions() {
	for y, row := range m.plots {
		for x, c := range row {
			co := Coord{x, y}
			if m.seen.Has(co) {
				continue
			}
			r := &Region{id: c, coords: sets.New[Coord]()}
			m.regions = append(m.regions, r)
			m.ExpandRegion(x, y, c, r)
		}
	}
}

var (
	west  = Coord{-1, 0}
	south = Coord{0, 1}
	north = Coord{0, -1}
	east  = Coord{1, 0}
)

// West, South, North, East // optimized for fill order.
var directions = []Coord{west, south, north, east}

func (m *Map) IsEdge(c byte, x, y int, direction Coord) bool {
	xx := x + direction.x
	yy := y + direction.y
	return m.Outside(xx, yy) || m.plots[yy][xx] != c
}

/* Rotate by 45 degrees clockwise ie N->NE; NE->E; E->SE... */
func (c Coord) RotateRight() Coord {
	if c.x == 0 { // N/S to NE/SW.
		c.x = -c.y
		return c
	}
	if c.y == 0 { // E/W to SE/NW.
		c.y = c.x
		return c
	}
	if c.x == -c.y { // NE -> E
		c.y = 0
		return c
	}
	c.x = 0
	return c
}

/*
Handle (for each 4 directions, this example is north)

	3
	13

	ie. 3 and 3 + 1 in middle by checking the Diagonal.
*/
func (m *Map) IsSpecial(c byte, x, y int, d Coord) bool {
	diag := d.RotateRight()
	xd := x + diag.x
	yd := y + diag.y
	if m.Outside(xd, yd) {
		return false
	}
	if m.plots[yd][xd] == c {
		return false
	}
	// We already check that direction d isn't an edge, so only the 90deg right needs checking.
	next := diag.RotateRight()
	xdn := x + next.x
	ydn := y + next.y
	return m.plots[ydn][xdn] == c
}

func (m *Map) Edges(x, y int) (int, int) {
	edges := 0
	c := m.plots[y][x]
	var e [4]bool
	corners := 0
	for i, d := range directions {
		if m.IsEdge(c, x, y, d) {
			e[i] = true
			edges++
		} else if m.IsSpecial(c, x, y, d) {
			// fmt.Printf("Special %c at %d,%d for %v (%v)\n", r.id, x, y, d, d.RotateRight())
			corners++
		}
	}
	switch edges {
	case 4:
		corners += 4
	case 3:
		corners += 2
	case 2:
		// only a corner if the edges are not opposite
		if !((e[0] && e[3]) || (e[1] && e[2])) {
			corners++
		}
	}
	// fmt.Printf("Region  %c at %d,%d edges %d  (%v) corners %d\n", r.id, x, y, edges, e, corners)
	return edges, corners
}

func (m *Map) Outside(x, y int) bool {
	return x < 0 || y < 0 || x >= len(m.plots) || y >= len(m.plots[0])
}

func (m *Map) ExpandRegion(x, y int, c byte, r *Region) {
	if m.Outside(x, y) {
		return
	}
	if m.plots[y][x] != c {
		return
	}
	co := Coord{x, y}
	if m.seen.Has(co) {
		return
	}
	r.coords.Add(co)
	m.seen.Add(co)
	for _, d := range directions {
		m.ExpandRegion(x+d.x, y+d.y, c, r)
	}
	perim, corners := m.Edges(x, y)
	r.perimeter += perim
	r.corners += corners
	if perim == 0 {
		m.inner.Add(co)
	}
}

func Color(c, what byte) string {
	idx := int(c) - int('A')
	if idx < 0 || idx >= 26 {
		log.Fatalf("Invalid color %c %d", c, c)
	}
	color := 18 + idx*216/26
	return fmt.Sprintf("\033[48;5;%dm%c\033[0m", color, what)
}

func (m *Map) String() string {
	sum1 := 0
	sum2 := 0
	var b strings.Builder
	for y, row := range m.plots {
		for x, c := range row {
			what := c
			if m.inner.Has(Coord{x, y}) {
				what = ' '
			}
			b.WriteString(Color(c, what))
		}
		b.WriteByte('\n')
	}
	for _, r := range m.regions {
		s := len(r.coords)
		p := r.perimeter
		cost1 := s * p
		sum1 += cost1
		e := r.corners
		cost2 := s * e
		sum2 += cost2
		b.WriteString(fmt.Sprintf("Region %c: s %d p %d = %d - c %d = %d\n", r.id, s, p, cost1, e, cost2))
	}
	b.WriteString(fmt.Sprintf("Total with perimeter : %d\n", sum1))
	b.WriteString(fmt.Sprintf("Total with corners   : %d\n", sum2))
	return b.String()
}

func main() {
	// read input list
	data, _ := io.ReadAll(os.Stdin)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	// Part1
	m := Map{plots: make([][]byte, len(lines)), seen: sets.New[Coord](), inner: sets.New[Coord]()}
	for i, line := range lines {
		m.plots[i] = []byte(line)
	}
	m.Regions()
	fmt.Println(m.String())
}
