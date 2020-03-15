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
	"math"
	"time"

	"github.com/JizongL/go_with_game/noise"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

// --
type gameState int

const (
	start gameState = iota
	play
)

var state = start

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

func lerp(b1 byte, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 color, pct float32) color {
	return color{lerp(c1.r, c2.r, pct), lerp(c1.g, c2.g, pct), lerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}
	return result
}

func getDualGradient(c1, c2, c3, c4 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		if pct < 0.5 {
			result[i] = colorLerp(c1, c2, pct*float32(2))
		} else {
			result[i] = colorLerp(c3, c4, pct*float32(1.5)-float32(0.5))
		}
	}
	return result
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	}
	if v > max {
		v = max
	}
	return v
}

func rescaleAndDraw(noise []float32, min, max float32, gradient []color, pixels []byte) {
	scale := 255.0 / (max - min)
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset
		c := gradient[clamp(0, 255, int(noise[i]))]
		// b := byte(noise[i])
		// pixels[i*4] = b
		// pixels[i*4+1] = b
		// pixels[i*4+2] = b
		pixels[i*4] = c.r
		pixels[i*4+1] = c.g
		pixels[i*4+2] = c.b
	}
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	nums := [][]byte{
		{
			1, 1, 1,
			1, 0, 1,
			1, 0, 1,
			1, 0, 1,
			1, 1, 1,
		},
		{
			1, 1, 0,
			1, 1, 0,
			1, 1, 0,
			1, 1, 0,
			1, 1, 1,
		},
		{
			1, 1, 1,
			0, 0, 1,
			1, 1, 1,
			1, 0, 0,
			1, 1, 1,
		},
		{
			1, 1, 1,
			0, 0, 1,
			0, 1, 1,
			0, 0, 1,
			1, 1, 1,
		},
	}

	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2
	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func (ball *ball) draw(pixels []byte) {
	// YAGNI - Ya Aint Gonna Need it
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			// sqrt root is expensive
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime
	// handle collision
	if ball.y-ball.radius < 0 || ball.y+ball.radius > float32(winHeight) {
		ball.yv = -ball.yv
	}

	if ball.x < 0 {
		rightPaddle.score++
		ball.pos = getCenter()
		state = start
	} else if int(ball.x) > winWidth {
		leftPaddle.score++
		ball.pos = getCenter()
		state = start
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {

		if ball.y > leftPaddle.y-leftPaddle.h/2 &&
			ball.y < leftPaddle.y+leftPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = leftPaddle.x + leftPaddle.w/2.0 + ball.radius
		}
	}
	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 &&
			ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = rightPaddle.x - rightPaddle.w/2.0 - ball.radius
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	score int
	color color
}

// linear interpolation
func flerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
	numX := flerp(paddle.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8, controllerAxis int16, elapsedTime float32) {

	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
		fmt.Println(paddle.y)
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
		fmt.Println(paddle.y)
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
		paddle.y += paddle.speed * pct * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	paddle.y = ball.y
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

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

	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	pixels := make([]byte, winWidth*winHeight*4)

	// for y := 0; y < winHeight; y++ {
	// 	for x := 0; x < winWidth; x++ {
	// 		setPixel(x, y, color{byte(x % 255), byte(y % 255), 0}, pixels)
	// 	}
	// }

	player1 := paddle{pos{50, 100}, 20, 100, 300, 0, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, 300, 0, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 200, 200, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()
	noise, min, max := noise.MakeNoise(noise.FBM, .001, 0.5, 2, 3, winWidth, winHeight)
	gradient := getGradient(color{255, 0, 0}, color{0, 0, 0})
	rescaleAndDraw(noise, min, max, gradient, pixels)
	var frameStart time.Time
	var elapsedTime float32
	var controllerAxis int16
	for {
		frameStart = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// type switch
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Joystick().Axis(sdl.CONTROLLER_AXIS_LEFTY)

			}
		}

		if state == play {
			// drawNumber(getCenter(), color{255, 255, 255}, 20, 2, pixels)
			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				state = play
			}
		}

		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds())
		fmt.Println(elapsedTime)
		if elapsedTime < .005 {
			sdl.Delay(uint32(5) - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}

}
