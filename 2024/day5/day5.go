package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"fortio.org/sets"
)

type Rule struct {
	Before, After int
}

type Data struct {
	Row []int // raw row
	Set sets.Set[int]
}

func (d Data) Middle() int {
	return d.Row[len(d.Row)/2]
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	var rules []Rule
	// fmt.Println("Rules:")
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break // 2nd section
		}
		var r Rule
		_, err := fmt.Sscanf(line, "%d|%d", &r.Before, &r.After)
		if err != nil {
			panic(err)
		}
		rules = append(rules, r)
	}
	// fmt.Println(rules)
	// fmt.Println("Data:")
	var data []Data
	for s.Scan() {
		line := s.Text()
		elems := strings.Split(line, ",")
		var row Data
		for _, e := range elems {
			v, err := strconv.Atoi(e)
			if err != nil {
				panic(err)
			}
			row.Row = append(row.Row, v)
		}
		row.Set = sets.FromSlice(row.Row)
		data = append(data, row)
	}
	// fmt.Println(data)

	// Part 1 + Part 2 prep
	sum1 := 0
	var badRows []Data
	for i, row := range data {
		rowOk := true
		for j, r := range rules {
			if !(row.Set.Has(r.Before) && row.Set.Has(r.After)) {
				// need both to be there for the rule to matter
				continue
			}
			i1 := slices.Index(row.Row, r.Before)
			i2 := slices.Index(row.Row, r.After)
			if i1 == -1 || i2 == -1 {
				log.Fatalf("BUG Row %d: Rule %d: %d or %d not found in %v despite check", i+1, j+1, r.Before, r.After, row)
			}
			if i1 > i2 {
				fmt.Printf("Row %d eliminated by rule %d: %v: %d not before %d\n", i+1, j+1, row.Row, r.Before, r.After)
				rowOk = false
				break
			}
		}
		if rowOk {
			fmt.Printf("Row %d is ok, adding %d: %v\n", i+1, row.Middle(), row.Row)
			sum1 += row.Middle()
		} else {
			badRows = append(badRows, row)
		}
	}
	fmt.Println(sum1)
	sum2 := 0
	for i, row := range badRows {
		for j := 0; j < len(rules); j++ {
			// log.Printf("Checking row %d with rule %d", i+1, j+1)
			didFix := false
			r := rules[j]
			if !row.Set.Has(r.Before) || !row.Set.Has(r.After) {
				continue
			}
			for {
				i1 := slices.Index(row.Row, r.Before)
				i2 := slices.Index(row.Row, r.After)
				if i1 == -1 || i2 == -1 {
					log.Fatalf("BUG Row %d: Rule %d: %d or %d not found in %v", i+1, j+1, r.Before, r.After, row)
				}
				if i1 < i2 {
					break
				}
				// move the after to the left
				fmt.Printf("Row %d fixing for rule %d (%d before %d): before %v -> ", i+1, j+1, r.Before, r.After, row.Row)
				row.Row[i1], row.Row[i1-1] = row.Row[i1-1], row.Row[i1]
				fmt.Printf("now %v\n", row.Row)
				didFix = true
			}
			if didFix {
				// recheck from the start
				fmt.Printf("Rechecking %d from 0 as rule %d made changes\n", i+1, j+1)
				j = -1 // will be incremented to 0
			}
		}
		sum2 += row.Middle()
	}
	fmt.Println(sum2)
}
