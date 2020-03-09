package main

import "fmt"

type position struct {
	x float32
	y float32
}

type badGuy struct {
	name   string
	health int
	pos    position
}

func addOne(x *int) {
	*x = *x + 1
}

func whereIsBadGuy(b *badGuy) {
	x := b.pos.x
	y := b.pos.y
	fmt.Println(b.name, "is at (", x, ",", y, ")")
}

func main() {
	x := 5
	fmt.Println(x)
	// var xPtr *int = &x
	xPtr := &x
	fmt.Println(xPtr)
	addOne(xPtr)
	fmt.Println(x)

	p := position{4, 2}
	b := badGuy{"Jabba the Hut", 100, p}
	whereIsBadGuy(&b)
}
