package main

import (
	runewidth "github.com/mattn/go-runewidth"
	tty "github.com/nsf/termbox-go"
)

func print(x int, y int, str string) int {
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
	pos := printRunes(x, y, input.text)
	tty.SetCursor(pos, y)
}

var input Input

func paint() {
	tty.Clear(tty.ColorDefault, tty.ColorDefault)
	defer tty.Flush()

	input.Paint(0, 0)
}

func main() {
	err := tty.Init()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	tty.SetInputMode(tty.InputEsc)

	paint()

mainloop:
	for {
		switch event := tty.PollEvent(); event.Type {
		case tty.EventKey:
			switch event.Key {
			case tty.KeyEsc:
				break mainloop

			case tty.KeyBackspace, tty.KeyBackspace2:
				input.Pop()

			default:
				if event.Ch != 0 {
					input.Append(event.Ch)
				}
			}

		case tty.EventError:
			panic(event.Err)
		}

		paint()
	}
}
