// Day 10 part 2 of advent of code 2025
// https://adventofcode.com/2025/day/10
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Machine struct {
	buttons  [][]int
	joltages []int
}

func readInput(scanner *bufio.Scanner) []Machine {
	res := []Machine{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		// ignore lights (parts[0])
		btns := [][]int{}
		for _, b := range parts[1 : len(parts)-1] {
			btn := []int{}
			trimmed := strings.Trim(b, "()")
			for _, i := range strings.Split(trimmed, ",") {
				val, _ := strconv.Atoi(i)
				btn = append(btn, val)
			}
			btns = append(btns, btn)
		}
		// joltages
		j := []int{}
		trimmed := strings.Trim(parts[len(parts)-1], "{}")
		for _, i := range strings.Split(trimmed, ",") {
			val, _ := strconv.Atoi(i)
			j = append(j, val)
		}
		fmt.Println("Machine buttons", btns, "joltages", j)
		res = append(res, Machine{buttons: btns, joltages: j})
	}
	return res
}

func press(state, button []int) []int {
	// Make a copy of state
	newState := make([]int, len(state))
	copy(newState, state)
	for _, i := range button {
		newState[i] += 1
	}
	return newState
}

func tooHigh(state, target []int) bool {
	for i := 0; i < len(state); i++ {
		if state[i] > target[i] {
			return true
		}
	}
	return false
}

type QueueItem struct {
	state   []int
	presses int
	nextIdx int
}

func solve(state, target []int, buttons [][]int) int {
	queue := []QueueItem{{state: state, presses: 0, nextIdx: 0}}
	numButtons := len(buttons)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:] // pop front
		// Found solution - return immediately
		if slices.Equal(current.state, target) {
			fmt.Println("Found solution with presses:", current.presses)
			return current.presses
		}
		// Only try buttons from nextIdx onwards to avoid exploring symmetrical permutations
		for bidx := current.nextIdx; bidx < numButtons; bidx++ {
			b := buttons[bidx]
			newState := press(current.state, b)
			if tooHigh(newState, target) {
				continue
			}
			queue = append(queue, QueueItem{state: newState, presses: current.presses + 1, nextIdx: bidx})
		}
	}
	return 0
}

func main() {
	fmt.Println("Part 2:")
	scanner := bufio.NewScanner(os.Stdin)
	machines := readInput(scanner)

	for _, m := range machines {
		fmt.Println("Machine:", m)
	}

	sum := 0
	for i, m := range machines {
		btns := m.buttons
		j := m.joltages
		fmt.Printf("Solving machine %d with joltages %v and buttons %v\n", i+1, j, btns)
		numState := len(j)
		state := make([]int, numState)
		res := solve(state, j, btns)
		fmt.Printf("Machine %d result: %d\n", i+1, res)
		sum += res
	}

	fmt.Println("Part 2: Result:", sum)
}
