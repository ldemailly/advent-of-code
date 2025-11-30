package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type File struct {
	id   int
	size int
}

const EMPTY = -1

type FS struct {
	s      []File // fileids
	size   int
	lastID int
	flat   []int
}

func (fs *FS) NewEmpty(lg int) {
	fs.size += lg
	fs.s = append(fs.s, File{id: EMPTY, size: lg})
}

func (fs *FS) NewFile(lg int) {
	fs.size += lg
	fs.s = append(fs.s, File{id: fs.lastID, size: lg})
	fs.lastID++
}

func (f File) ID2Str() string {
	if f.id == EMPTY {
		return "."
	}
	if f.id > 9 {
		return "*"
	}
	return string(byte(f.id + '0'))
}

func (f File) String() string {
	return strings.Repeat(f.ID2Str(), f.size)
}

func (fs FS) String() string {
	var buf strings.Builder
	for _, f := range fs.s {
		buf.WriteString(f.String())
	}
	return buf.String()
}

func CharToLen(c byte) int {
	return int(c - '0')
}

func (fs *FS) Flatten() {
	fs.flat = make([]int, fs.size)
	i := 0
	for _, f := range fs.s {
		for range f.size {
			fs.flat[i] = f.id
			i++
		}
	}
}

func (fs *FS) Defrag() {
	i := 0
	j := fs.size - 1
	for i <= j {
		if fs.flat[j] == EMPTY {
			j--
			continue
		}
		if fs.flat[i] != EMPTY {
			i++
			continue
		}
		fs.flat[i], fs.flat[j] = fs.flat[j], fs.flat[i]
		i++
		j--
	}
}

func (fs *FS) Defrag2() {
	j := len(fs.s) - 1
	for j >= 0 {
		if fs.s[j].id == EMPTY {
			j--
			continue
		}
		i := 0
		swapped := false
		for i < j {
			if fs.s[i].id != EMPTY {
				i++
				continue
			}
			if fs.s[i].size >= fs.s[j].size {
				newHole := fs.s[i].size - fs.s[j].size
				if newHole > 0 {
					fs.s[i] = fs.s[j]
					fs.s[j].id = EMPTY
					// insert remaining hole
					newFS := make([]File, 0, len(fs.s)+1)
					newFS = append(newFS, fs.s[:i+1]...)
					newFS = append(newFS, File{id: EMPTY, size: newHole})
					newFS = append(newFS, fs.s[i+1:]...)
					fs.s = newFS
					// j is already decremented by making fs bigger/inserting the leftover hole
				} else {
					fs.s[i], fs.s[j] = fs.s[j], fs.s[i]
					j--
				}
				swapped = true
				break
			}
			i++
		}
		if !swapped {
			j--
		}
	}
}

func (fs *FS) FlatString() string {
	var buf strings.Builder
	for _, id := range fs.flat {
		if id == EMPTY {
			buf.WriteByte('.')
		} else {
			if id > 9 {
				buf.WriteRune('*')
			} else {
				buf.WriteByte(byte(id + '0'))
			}
		}
	}
	return buf.String()
}

func (fs *FS) Checksum() int {
	sum := 0
	for i, id := range fs.flat {
		mul := id
		if id == EMPTY {
			mul = 0
		}
		sum += i * mul
	}
	return sum
}

var debug = false

func DebugPrintf(format string, a ...interface{}) {
	if debug {
		fmt.Printf(format, a...)
	}
}

func main() {
	flag.BoolVar(&debug, "debug", false, "debug mode: only for uppercase input")
	flag.Parse()
	fs := FS{}
	m, _ := io.ReadAll(os.Stdin)
	for i, c := range m {
		if c == '\n' {
			break
		}
		if i%2 == 0 {
			fs.NewFile(CharToLen(c))
		} else {
			fs.NewEmpty(CharToLen(c))
		}
	}
	DebugPrintf("FS is %d blocks (%s)\n", fs.size, fs)
	fs.Flatten()
	DebugPrintf("Flattened1 is %s (%v)\n", fs.FlatString(), fs.flat)
	fs.Defrag()
	DebugPrintf("Defragged1 is %s (%v)\n", fs.FlatString(), fs.flat)
	fmt.Println("Checksum Part1:", fs.Checksum())
	// Part1
	fs.Defrag2()
	fs.Flatten()
	DebugPrintf("Defragged2 is %s (%v)\n", fs.FlatString(), fs.flat)
	fmt.Println("Checksum Part2:", fs.Checksum())
}
