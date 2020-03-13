package main

// Experiment! draw some crazy stuff!
// Gist it next week and I'll show it off on stream

// TODO
//  Freame rate indepedence
//  Score
//  Game Over State - win/lose
// 2 Player vs Playing computer
// AI needs to be more imperfect
// Handling resizing of the window
// Mouse? Joysticks

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	// YAGNI - Ya Aint Gonna Need it
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			// sqrt root is expensive
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(ball.x+x, ball.y+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv
	// handle collision
	if ball.y-ball.radius < 0 || ball.y+ball.radius > float32(winHeight) {
		ball.yv = -ball.yv
	}
	if ball.x < 0 || ball.x > float32(winWidth) {
		ball.x = 300
		ball.y = 300
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {
		if ball.y > leftPaddle.y-leftPaddle.h/2 &&
			ball.y < leftPaddle.y+leftPaddle.h/2 {
			ball.xv = -ball.xv
		}
	}
	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 &&
			ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xv = -ball.xv
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	color color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := paddle.x - paddle.w/2
	startY := paddle.y - paddle.h/2

	for y := float32(0); y < paddle.h; y++ {
		for x := float32(0); x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= 5
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 5
	}

}

func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y float32, c color, pixels []byte) {
	index := int(y)*winWidth + int(x)*4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}

}

func main() {

	// Added after EP06 to address macosx issues
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	// for y := 0; y < winHeight; y++ {
	// 	for x := 0; x < winWidth; x++ {
	// 		setPixel(x, y, color{byte(x % 255), byte(y % 255), 0}, pixels)
	// 	}
	// }

	player1 := paddle{pos{50, 100}, 20, 100, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 3, 3, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// type switch
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}

		}
		clear(pixels)
		player1.update(keyState)
		player2.aiUpdate(&ball)
		ball.update(&player1, &player2)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}

}
