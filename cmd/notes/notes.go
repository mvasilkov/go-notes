package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
	tty "github.com/nsf/termbox-go"
	filepathx "github.com/yargevad/filepathx"
)

// Constants
const (
	Prompt = "> "
	Cursor = "* "
	Blank  = "  "
)

// Map function
func Map(a []string, cb func(string) string) []string {
	b := make([]string, len(a))
	for i, str := range a {
		b[i] = cb(str)
	}
	return b
}

func printString(x int, y int, str string) int {
	for _, c := range str {
		tty.SetCell(x, y, c, tty.ColorDefault, tty.ColorDefault)
		x += runewidth.RuneWidth(c)
	}
	return x
}

func printRunes(x int, y int, str []rune) int {
	for _, c := range str {
		tty.SetCell(x, y, c, tty.ColorDefault, tty.ColorDefault)
		x += runewidth.RuneWidth(c)
	}
	return x
}

// Input struct
type Input struct {
	text []rune
}

// Append (Input)
func (input *Input) Append(c rune) {
	input.text = append(input.text, c)
}

// Pop (Input)
func (input *Input) Pop() rune {
	if len(input.text) == 0 {
		return 0
	}
	c := input.text[len(input.text)-1]
	input.text = input.text[:len(input.text)-1]
	return c
}

// Paint (Input)
func (input *Input) Paint(x int, y int) {
	printString(x, y, Prompt)
	pos := printRunes(x+2, y, input.text)
	tty.SetCursor(pos, y)
}

var input Input
var selection int

func paint(notes []string, height int) {
	tty.Clear(tty.ColorDefault, tty.ColorDefault)
	defer tty.Flush()

	input.Paint(0, 0)

	for i, name := range notes {
		p := Blank
		if i == selection {
			p = Cursor
		}
		printString(0, i+1, p)
		printString(2, i+1, name)

		if i > height-2 {
			break
		}
	}
}

// LoadNotes function
func LoadNotes(dir string) []string {
	notes, err := filepathx.Glob(filepath.Join(dir, "**", "*.go"))
	if err != nil {
		panic(err)
	}
	notes = Map(notes, func(a string) string {
		return strings.Replace(a, dir+string(os.PathSeparator), "", 1)
	})
	return notes
}

func isMatching(a string) bool {
	for _, c := range input.text {
		pos := strings.IndexRune(a, c)
		if pos == -1 {
			return false
		}
		a = a[pos+1:]
	}
	return true
}

// FilterNotes function
func FilterNotes(notes []string) []string {
	var results []string
	for _, name := range notes {
		if isMatching(name) {
			results = append(results, name)
		}
	}
	return results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: notes DIR")
		return
	}

	dir, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	notes := LoadNotes(dir)
	filtered := notes

	err = tty.Init()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	tty.SetInputMode(tty.InputEsc)

mainloop:
	for {
		_, height := tty.Size()

		paint(filtered, height)

		switch event := tty.PollEvent(); event.Type {
		case tty.EventKey:
			switch event.Key {
			case tty.KeyEsc:
				break mainloop

			case tty.KeyArrowUp:
				selection--
				if selection < 0 {
					selection = 0
				}

			case tty.KeyArrowDown:
				selection++
				if selection > len(filtered)-1 {
					selection = len(filtered) - 1
				}
				if selection > height-2 {
					selection = height - 2
				}

			case tty.KeyBackspace, tty.KeyBackspace2:
				input.Pop()
				filtered = FilterNotes(notes)
				selection = 0

			default:
				if event.Ch != 0 {
					input.Append(event.Ch)
					filtered = FilterNotes(notes)
					selection = 0
				}
			}

		case tty.EventError:
			panic(event.Err)
		}
	}
}
